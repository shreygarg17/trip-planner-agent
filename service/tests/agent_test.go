package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/shreygarg/trip-planner-agent/internal/tools/weather"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/repo"
	"github.com/shreygarg/trip-planner-agent/service"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
	"github.com/shreygarg/trip-planner-agent/validations"
)

type mockLLMClient struct {
	callCount int
	t         *testing.T
}

func (m *mockLLMClient) CreateChatCompletion(ctx context.Context, req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	m.callCount++

	if m.callCount == 1 {
		// First turn: Verify LLM receives the user prompt and registers all 3 tools
		if len(req.Messages) != 1 || req.Messages[0].Role != "user" {
			m.t.Errorf("Expected first turn to have exactly 1 user message, got: %v", req.Messages)
		}

		if len(req.Tools) != 3 {
			m.t.Errorf("Expected first turn to register 3 tools, got: %d", len(req.Tools))
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
		// Second turn: Verify LLM receives the recommend_destinations response and requests the generate_itinerary tool
		if len(req.Messages) != 3 {
			m.t.Errorf("Expected second turn to have 3 messages (user, assistant tool-call, tool response), got: %d", len(req.Messages))
		}

		toolMsg := req.Messages[2]
		if !strings.Contains(toolMsg.Content, "Sikkim") {
			m.t.Errorf("Expected first tool response to contain Sikkim, got: %s", toolMsg.Content)
		}

		// Simulate LLM calling the generate_itinerary tool
		return &models.ChatCompletionResponse{
			Choices: []models.Choice{
				{
					Message: models.Message{
						Role: "assistant",
						ToolCalls: []models.ToolCall{
							{
								ID:   "call_abc123",
								Type: "function",
								Function: models.ToolFunction{
									Name:      "generate_itinerary",
									Arguments: `{"destination": "Sikkim", "days": 4}`,
								},
							},
						},
					},
				},
			},
		}, nil
	}

	if m.callCount == 3 {
		// Third turn: Verify LLM receives the generate_itinerary response and outputs final response
		if len(req.Messages) != 5 {
			m.t.Errorf("Expected third turn to have 5 messages, got: %d", len(req.Messages))
		}

		itineraryMsg := req.Messages[4]
		if !strings.Contains(itineraryMsg.Content, "Gangtok") {
			m.t.Errorf("Expected second tool response to contain Gangtok activity, got: %s", itineraryMsg.Content)
		}

		return &models.ChatCompletionResponse{
			Choices: []models.Choice{
				{
					Message: models.Message{
						Role:    "assistant",
						Content: "Recommended Destination: Sikkim\nEstimated Cost: ₹20,000\nDay 1: Arrive and explore Gangtok\nDay 2: Tsomgo Lake\nDay 3: Nathula Pass\nDay 4: Local markets and return",
					},
				},
			},
		}, nil
	}

	return nil, nil
}

func TestTripAgent_PlanTrip(t *testing.T) {
	mock := &mockLLMClient{t: t}
	destRepo := repo.NewInMemoryDestinationRepository()
	planner := service.NewDestinationPlanner(destRepo)
	validator := validations.NewValidator()
	itineraryRepo := repo.NewInMemoryItineraryRepository()
	itineraryGen := service.NewItineraryGenerator(validator, itineraryRepo)

	recommendTool := service.NewRecommendDestinationsTool(planner, validator)
	itineraryTool := service.NewGenerateItineraryTool(itineraryGen)
	weatherService := weather.NewOpenMeteoClient(nil)
	weatherTool := service.NewWeatherTool(weatherService)

	agent := service.NewTripAgent(mock, []interfaces.ToolExecutor{recommendTool, itineraryTool, weatherTool}, "gpt-4o")

	prompt := "Plan a 4-day mountain trip under ₹50k"
	response, err := agent.PlanTrip(context.Background(), prompt)
	if err != nil {
		t.Fatalf("PlanTrip failed: %v", err)
	}

	if mock.callCount != 3 {
		t.Errorf("Expected exactly 3 LLM completion calls, got %d", mock.callCount)
	}

	expectedResult := "Recommended Destination: Sikkim"
	if !strings.Contains(response, expectedResult) {
		t.Errorf("Expected response to contain %q, got %q", expectedResult, response)
	}
}
