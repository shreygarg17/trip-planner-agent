package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
)

var (
	itineraryToolOnce     sync.Once
	itineraryToolInstance interfaces.ToolExecutor
)

// GenerateItineraryTool wraps ItineraryGenerator as an LLM tool.
type GenerateItineraryTool struct {
	itineraryGen interfaces.ItineraryGenerator
}

// NewGenerateItineraryTool returns a new ToolExecutor instance (singleton).
func NewGenerateItineraryTool(itineraryGen interfaces.ItineraryGenerator) interfaces.ToolExecutor {
	itineraryToolOnce.Do(func() {
		itineraryToolInstance = &GenerateItineraryTool{
			itineraryGen: itineraryGen,
		}
	})
	return itineraryToolInstance
}

// GetDefinition returns the LLM tool definition metadata.
func (t *GenerateItineraryTool) GetDefinition() models.Tool {
	return models.Tool{
		Type: "function",
		Function: models.FunctionDefinition{
			Name:        "generate_itinerary",
			Description: "Generate a day-by-day travel itinerary for a specific destination and duration.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"destination": map[string]interface{}{
						"type":        "string",
						"description": "The destination name (e.g. Sikkim, Manali, Bali).",
					},
					"days": map[string]interface{}{
						"type":        "integer",
						"description": "Number of days for the itinerary.",
					},
				},
				"required": []string{"destination", "days"},
			},
		},
	}
}

// Execute parses function arguments and runs the itinerary generator service.
func (t *GenerateItineraryTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Destination string `json:"destination"`
		Days        int    `json:"days"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrUnmarshalToolArgs, err)
	}

	req := models.ItineraryRequest{
		Destination: args.Destination,
		Days:        args.Days,
	}

	itinerary, err := t.itineraryGen.GenerateItinerary(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrExecTool, err)
	}

	resBytes, err := json.Marshal(itinerary)
	if err != nil {
		return "", fmt.Errorf("failed to marshal itinerary: %w", err)
	}

	return string(resBytes), nil
}
