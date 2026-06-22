package service

import (
	"errors"
	"fmt"

	clientinterfaces "github.com/shreygarg/trip-planner-agent/clients/interfaces"
	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
)

// TripAgent coordinates prompt execution and tool calls.
type TripAgent struct {
	llmClient    clientinterfaces.LLMClient
	toolExecutor interfaces.ToolExecutor
	model        string
}

// NewTripAgent instantiates a new Agent service.
func NewTripAgent(llmClient clientinterfaces.LLMClient, toolExecutor interfaces.ToolExecutor, model string) interfaces.Agent {
	return &TripAgent{
		llmClient:    llmClient,
		toolExecutor: toolExecutor,
		model:        model,
	}
}

// PlanTrip processes a planning prompt, resolving requested tool calls iteratively.
func (a *TripAgent) PlanTrip(prompt string) (string, error) {
	messages := []models.Message{
		{Role: "user", Content: prompt},
	}

	availableTools := []models.Tool{
		a.toolExecutor.GetDefinition(),
	}

	for i := 0; i < 5; i++ {
		req := models.ChatCompletionRequest{
			MaxTokens: 500,
			Model:     a.model,
			Messages:  messages,
			Tools:     availableTools,
		}

		resp, err := a.llmClient.CreateChatCompletion(req)
		if err != nil {
			return "", fmt.Errorf("llm completion failure: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", errors.New(constants.ErrLLMEmptyChoices)
		}

		choiceMessage := resp.Choices[0].Message
		messages = append(messages, choiceMessage)

		if len(choiceMessage.ToolCalls) == 0 {
			return choiceMessage.Content, nil
		}

		toolMsgs, err := a.handleToolCalls(choiceMessage.ToolCalls)
		if err != nil {
			return "", err
		}
		messages = append(messages, toolMsgs...)
	}

	return "", errors.New(constants.ErrMaxIterationsExceeded)
}

func (a *TripAgent) handleToolCalls(toolCalls []models.ToolCall) ([]models.Message, error) {
	messages := make([]models.Message, 0, len(toolCalls))
	for _, tc := range toolCalls {
		if tc.Function.Name != "recommend_destinations" {
			messages = append(messages, models.Message{
				Role:       "tool",
				Content:    fmt.Sprintf(`{"error": "unsupported tool: %s"}`, tc.Function.Name),
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
			})
			continue
		}

		res, err := a.toolExecutor.Execute(tc.Function.Arguments)
		if err != nil {
			return nil, fmt.Errorf("tool execution failed: %w", err)
		}

		messages = append(messages, models.Message{
			Role:       "tool",
			Content:    res,
			ToolCallID: tc.ID,
			Name:       tc.Function.Name,
		})
	}
	return messages, nil
}
