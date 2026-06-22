package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	clientinterfaces "github.com/shreygarg/trip-planner-agent/clients/interfaces"
	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
)

var (
	agentOnce     sync.Once
	agentInstance interfaces.Agent
)

// TripAgent coordinates prompt execution and tool calls.
type TripAgent struct {
	llmClient     clientinterfaces.LLMClient
	toolExecutors []interfaces.ToolExecutor
	model         string
}

// NewTripAgent instantiates a new Agent service (singleton).
func NewTripAgent(llmClient clientinterfaces.LLMClient, toolExecutors []interfaces.ToolExecutor, model string) interfaces.Agent {
	agentOnce.Do(func() {
		agentInstance = &TripAgent{
			llmClient:     llmClient,
			toolExecutors: toolExecutors,
			model:         model,
		}
	})
	return agentInstance
}

// PlanTrip processes a planning prompt, resolving requested tool calls iteratively.
func (a *TripAgent) PlanTrip(ctx context.Context, prompt string) (string, error) {
	messages := []models.Message{
		{Role: "user", Content: prompt},
	}

	availableTools := make([]models.Tool, 0, len(a.toolExecutors))
	for _, te := range a.toolExecutors {
		availableTools = append(availableTools, te.GetDefinition())
	}

	for i := 0; i < 5; i++ {
		req := models.ChatCompletionRequest{
			Model:     a.model,
			MaxTokens: 500,
			Messages:  messages,
			Tools:     availableTools,
		}

		resp, err := a.llmClient.CreateChatCompletion(ctx, req)
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

		toolMsgs, err := a.handleToolCalls(ctx, choiceMessage.ToolCalls)
		if err != nil {
			return "", err
		}
		messages = append(messages, toolMsgs...)
	}

	return "", errors.New(constants.ErrMaxIterationsExceeded)
}

func (a *TripAgent) handleToolCalls(ctx context.Context, toolCalls []models.ToolCall) ([]models.Message, error) {
	messages := make([]models.Message, 0, len(toolCalls))
	executorMap := make(map[string]interfaces.ToolExecutor)
	for _, te := range a.toolExecutors {
		executorMap[te.GetDefinition().Function.Name] = te
	}

	for _, tc := range toolCalls {
		te, exists := executorMap[tc.Function.Name]
		if !exists {
			messages = append(messages, models.Message{
				Role:       "tool",
				Content:    fmt.Sprintf(`{"error": "unsupported tool: %s"}`, tc.Function.Name),
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
			})
			continue
		}

		res, err := te.Execute(ctx, tc.Function.Arguments)
		if err != nil {
			return nil, fmt.Errorf("tool %s execution failed: %w", tc.Function.Name, err)
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

// ResetSingletons resets all singletons in the service package for testing.
func ResetSingletons() {
	agentOnce = sync.Once{}
	agentInstance = nil
	plannerOnce = sync.Once{}
	plannerInstance = nil
	itineraryOnce = sync.Once{}
	itineraryInstance = nil
	recommendToolOnce = sync.Once{}
	recommendToolInstance = nil
	itineraryToolOnce = sync.Once{}
	itineraryToolInstance = nil
	weatherToolOnce = sync.Once{}
	weatherToolInstance = nil
	longWeekendToolOnce = sync.Once{}
	longWeekendToolInstance = nil
}
