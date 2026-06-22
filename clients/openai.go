package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/shreygarg/trip-planner-agent/clients/interfaces"
	"github.com/shreygarg/trip-planner-agent/config"
	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/utils"
)

var (
	openaiOnce     sync.Once
	openaiInstance interfaces.LLMClient
)

// OpenAIClient implements LLMClient using OpenRouter/OpenAI endpoints.
type OpenAIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewOpenAIClient returns an initialized LLMClient (singleton).
func NewOpenAIClient(cfg config.ConfigProvider) interfaces.LLMClient {
	openaiOnce.Do(func() {
		openaiInstance = &OpenAIClient{
			apiKey:  cfg.GetAPIKey(),
			baseURL: cfg.GetBaseURL(),
			client: &http.Client{
				Timeout: constants.HTTPTimeout,
			},
		}
	})
	return openaiInstance
}

// CreateChatCompletion initiates a chat completion workflow.
func (c *OpenAIClient) CreateChatCompletion(ctx context.Context, req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	reqID := utils.GetRequestID(ctx)

	httpCtx, cancel := context.WithTimeout(ctx, constants.HTTPTimeout)
	defer cancel()

	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(httpCtx, http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	c.setHeaders(httpReq)

	log.Printf("[%s] [INFO] Sending ChatCompletion request to LLM (Model: %s)", reqID, req.Model)
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request execution failed: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(reqID, resp)
}

func (c *OpenAIClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
}

func (c *OpenAIClient) parseResponse(reqID string, resp *http.Response) (*models.ChatCompletionResponse, error) {
	if resp.StatusCode != http.StatusOK {
		var errData map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errData)
		return nil, fmt.Errorf("unexpected status: %d, response: %v", resp.StatusCode, errData)
	}

	var completionResponse models.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("[%s] [INFO] Received successful ChatCompletion response from LLM", reqID)
	return &completionResponse, nil
}
