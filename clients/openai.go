package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shreygarg/trip-planner-agent/clients/interfaces"
	"github.com/shreygarg/trip-planner-agent/config"
	"github.com/shreygarg/trip-planner-agent/models"
)

// OpenAIClient implements LLMClient using OpenRouter/OpenAI endpoints.
type OpenAIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewOpenAIClient returns an initialized LLMClient.
func NewOpenAIClient(cfg config.ConfigProvider) interfaces.LLMClient {
	return &OpenAIClient{
		apiKey:  cfg.GetAPIKey(),
		baseURL: cfg.GetBaseURL(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateChatCompletion initiates a chat completion workflow.
func (c *OpenAIClient) CreateChatCompletion(req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request execution failed: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *OpenAIClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
}

func (c *OpenAIClient) parseResponse(resp *http.Response) (*models.ChatCompletionResponse, error) {
	if resp.StatusCode != http.StatusOK {
		var errData map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errData)
		return nil, fmt.Errorf("unexpected status: %d, response: %v", resp.StatusCode, errData)
	}

	var completionResponse models.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &completionResponse, nil
}
