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
	"github.com/shreygarg/trip-planner-agent/tools/longweekend"
	"github.com/shreygarg/trip-planner-agent/validations"
)

type mockLLMClient struct {
	callCount int
	t         *testing.T
}

func (m *mockLLMClient) CreateChatCompletion(ctx context.Context, req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	m.callCount++

	if m.callCount == 1 {
		// First turn: Verify LLM receives the user prompt and registers all 4 tools
		if len(req.Messages) != 1 || req.Messages[0].Role != "user" {
			m.t.Errorf("Expected first turn to have exactly 1 user message, got: %v", req.Messages)
		}

		if len(req.Tools) != 4 {
			m.t.Errorf("Expected first turn to register 4 tools, got: %d", len(req.Tools))
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
	service.ResetSingletons()
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
	holidayRepo := longweekend.NewStaticHolidayRepository()
	longWeekendSvc := longweekend.NewLongWeekendService(holidayRepo)
	longWeekendTool := service.NewGetLongWeekendsTool(longWeekendSvc)

	agent := service.NewTripAgent(mock, []interfaces.ToolExecutor{recommendTool, itineraryTool, weatherTool, longWeekendTool}, "gpt-4o")

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

type mockLLMLongWeekendClient struct {
	callCount int
	t         *testing.T
}

func (m *mockLLMLongWeekendClient) CreateChatCompletion(ctx context.Context, req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	m.callCount++
	m.t.Logf("[DEBUG] CreateChatCompletion callCount: %d, messages: %d, tools: %d", m.callCount, len(req.Messages), len(req.Tools))
	for i, msg := range req.Messages {
		m.t.Logf("  Message %d: Role=%s, Content=%s, ToolCalls=%d", i, msg.Role, msg.Content, len(msg.ToolCalls))
	}
	if m.callCount == 1 {
		// Verify tool count is 4
		if len(req.Tools) != 4 {
			m.t.Errorf("Expected 4 tools, got: %d", len(req.Tools))
		}
		// Simulate LLM calling get_long_weekends tool
		return &models.ChatCompletionResponse{
			Choices: []models.Choice{
				{
					Message: models.Message{
						Role: "assistant",
						ToolCalls: []models.ToolCall{
							{
								ID:   "call_lw1",
								Type: "function",
								Function: models.ToolFunction{
									Name:      "get_long_weekends",
									Arguments: `{"year": 2026}`,
								},
							},
						},
					},
				},
			},
		}, nil
	}

	if m.callCount == 2 {
		// Verify we got the tool response for get_long_weekends
		toolMsg := req.Messages[2]
		if !strings.Contains(toolMsg.Content, "Good Friday") {
			m.t.Errorf("Expected long weekends list to contain Good Friday, got: %s", toolMsg.Content)
		}
		// Simulate LLM calling recommend_destinations for the 3-day Good Friday weekend
		return &models.ChatCompletionResponse{
			Choices: []models.Choice{
				{
					Message: models.Message{
						Role: "assistant",
						ToolCalls: []models.ToolCall{
							{
								ID:   "call_rec1",
								Type: "function",
								Function: models.ToolFunction{
									Name:      "recommend_destinations",
									Arguments: `{"budget": 30000, "days": 3, "preferences": ["mountains"]}`,
								},
							},
						},
					},
				},
			},
		}, nil
	}

	if m.callCount == 3 {
		// Simulate LLM returning final trip suggestion
		return &models.ChatCompletionResponse{
			Choices: []models.Choice{
				{
					Message: models.Message{
						Role:    "assistant",
						Content: "For the upcoming Good Friday Long Weekend (2026-04-03 to 2026-04-05), I suggest Kasol.",
					},
				},
			},
		}, nil
	}

	return nil, nil
}

func TestTripAgent_PlanTrip_LongWeekend(t *testing.T) {
	service.ResetSingletons()
	mock := &mockLLMLongWeekendClient{t: t}
	destRepo := repo.NewInMemoryDestinationRepository()
	planner := service.NewDestinationPlanner(destRepo)
	validator := validations.NewValidator()
	itineraryRepo := repo.NewInMemoryItineraryRepository()
	itineraryGen := service.NewItineraryGenerator(validator, itineraryRepo)

	recommendTool := service.NewRecommendDestinationsTool(planner, validator)
	itineraryTool := service.NewGenerateItineraryTool(itineraryGen)
	weatherService := weather.NewOpenMeteoClient(nil)
	weatherTool := service.NewWeatherTool(weatherService)
	holidayRepo := longweekend.NewStaticHolidayRepository()
	longWeekendSvc := longweekend.NewLongWeekendService(holidayRepo)
	longWeekendTool := service.NewGetLongWeekendsTool(longWeekendSvc)

	agent := service.NewTripAgent(mock, []interfaces.ToolExecutor{recommendTool, itineraryTool, weatherTool, longWeekendTool}, "gpt-4o")

	prompt := "Suggest trips for upcoming long weekends in 2026"
	response, err := agent.PlanTrip(context.Background(), prompt)
	if err != nil {
		t.Fatalf("PlanTrip failed: %v", err)
	}

	if mock.callCount != 3 {
		t.Errorf("Expected exactly 3 LLM completion calls, got %d", mock.callCount)
	}

	if !strings.Contains(response, "Good Friday") || !strings.Contains(response, "Kasol") {
		t.Errorf("Expected response to mention Good Friday and Kasol, got %q", response)
	}
}
