package service

import (
	"encoding/json"
	"fmt"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
	validationinterfaces "github.com/shreygarg/trip-planner-agent/validations/interfaces"
)

// RecommendDestinationsTool wraps DestinationPlanner as an LLM tool.
type RecommendDestinationsTool struct {
	planner   interfaces.DestinationPlanner
	validator validationinterfaces.Validator
}

// NewRecommendDestinationsTool returns a new ToolExecutor instance.
func NewRecommendDestinationsTool(planner interfaces.DestinationPlanner, validator validationinterfaces.Validator) interfaces.ToolExecutor {
	return &RecommendDestinationsTool{
		planner:   planner,
		validator: validator,
	}
}

// GetDefinition returns the LLM tool definition metadata.
func (t *RecommendDestinationsTool) GetDefinition() models.Tool {
	return models.Tool{
		Type: "function",
		Function: models.FunctionDefinition{
			Name:        "recommend_destinations",
			Description: "Recommend travel destinations based on budget, ideal duration, and tags/preferences.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"source_city": map[string]interface{}{
						"type":        "string",
						"description": "The city the trip starts from.",
					},
					"budget": map[string]interface{}{
						"type":        "integer",
						"description": "Max budget for the trip.",
					},
					"days": map[string]interface{}{
						"type":        "integer",
						"description": "Number of days for the trip.",
					},
					"preferences": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
						"description": "List of tags/preferences like mountains, beach, international, adventure.",
					},
				},
				"required": []string{"budget", "days"},
			},
		},
	}
}

// Execute parses function arguments and runs the planner service.
func (t *RecommendDestinationsTool) Execute(argsJSON string) (string, error) {
	var args struct {
		SourceCity  string   `json:"source_city,omitempty"`
		Budget      int      `json:"budget"`
		Days        int      `json:"days"`
		Preferences []string `json:"preferences,omitempty"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrUnmarshalToolArgs, err)
	}

	req := models.TripRequest{
		SourceCity:  args.SourceCity,
		Budget:      args.Budget,
		Days:        args.Days,
		Preferences: args.Preferences,
	}

	if err := t.validator.ValidateTripRequest(req); err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrExecTool, err)
	}

	recommendations := t.planner.RecommendDestinations(req)
	resBytes, err := json.Marshal(recommendations)
	if err != nil {
		return "", fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	return string(resBytes), nil
}
