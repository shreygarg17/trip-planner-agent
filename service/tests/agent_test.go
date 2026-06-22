package tests

import (
	"strings"
	"testing"

	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/repo"
	"github.com/shreygarg/trip-planner-agent/service"
	"github.com/shreygarg/trip-planner-agent/validations"
)

type mockLLMClient struct {
	callCount int
	t         *testing.T
}

func (m *mockLLMClient) CreateChatCompletion(req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	m.callCount++

	if m.callCount == 1 {
		// First turn: Verify LLM receives the user prompt and registers the tool
		if len(req.Messages) != 1 || req.Messages[0].Role != "user" {
			m.t.Errorf("Expected first turn to have exactly 1 user message, got: %v", req.Messages)
		}

		if len(req.Tools) != 1 || req.Tools[0].Function.Name != "recommend_destinations" {
			m.t.Errorf("Expected first turn to register 'recommend_destinations' tool, got: %v", req.Tools)
		}

		// Simulate LLM calling the recommend_destinations tool
		return &models.ChatCompletionResponse{
			Choices: []models.Choice{
				{
					Message: models.Message{
						Role: "assistant",
						ToolCalls: []models.ToolCall{
							{
								ID:   "call_xyz987",
								Type: "function",
								Function: models.ToolFunction{
									Name:      "recommend_destinations",
									Arguments: `{"budget": 50000, "days": 4, "preferences": ["mountains"]}`,
								},
							},
						},
					},
				},
			},
		}, nil
	}

	if m.callCount == 2 {
		// Second turn: Verify LLM receives the tool response along with previous context
		if len(req.Messages) != 3 {
			m.t.Errorf("Expected second turn to have 3 messages (user, assistant tool-call, tool response), got: %d", len(req.Messages))
		}

		toolMsg := req.Messages[2]
		if toolMsg.Role != "tool" || toolMsg.Name != "recommend_destinations" || toolMsg.ToolCallID != "call_xyz987" {
			m.t.Errorf("Expected third message to be the tool response from 'recommend_destinations', got: %v", toolMsg)
		}

		// The content returned by planner should contain recommended destinations from the mock arguments (Bhutan, Kasol, Manali)
		if !strings.Contains(toolMsg.Content, "Bhutan") && !strings.Contains(toolMsg.Content, "Kasol") {
			m.t.Errorf("Expected tool response to contain destinations like Bhutan or Kasol, got: %s", toolMsg.Content)
		}

		// Return final LLM formatted recommendation response
		return &models.ChatCompletionResponse{
			Choices: []models.Choice{
				{
					Message: models.Message{
						Role:    "assistant",
						Content: "I recommend Kasol (score: 0.90) and Bhutan (score: 0.80) for mountains under ₹50k.",
					},
				},
			},
		}, nil
	}

	return nil, nil
}

func TestTripAgent_PlanTrip(t *testing.T) {
	mock := &mockLLMClient{t: t}
	destRepo := repo.NewDestinationRepository()
	planner := service.NewDestinationPlanner(destRepo)
	validator := validations.NewValidator()
	toolExecutor := service.NewRecommendDestinationsTool(planner, validator)
	agent := service.NewTripAgent(mock, toolExecutor, "gpt-4o")

	prompt := "I have ₹50k budget, 4 days, prefer mountains"
	response, err := agent.PlanTrip(prompt)
	if err != nil {
		t.Fatalf("PlanTrip failed: %v", err)
	}

	if mock.callCount != 2 {
		t.Errorf("Expected exactly 2 LLM completion calls, got %d", mock.callCount)
	}

	expectedResult := "I recommend Kasol"
	if !strings.Contains(response, expectedResult) {
		t.Errorf("Expected response to contain %q, got %q", expectedResult, response)
	}
}
