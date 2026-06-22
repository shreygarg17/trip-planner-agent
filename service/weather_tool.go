package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/internal/tools/weather"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
)

var (
	weatherToolOnce     sync.Once
	weatherToolInstance interfaces.ToolExecutor
)

// WeatherTool implements interfaces.ToolExecutor for the weather service.
type WeatherTool struct {
	weatherService weather.WeatherService
}

// NewWeatherTool returns a new ToolExecutor instance for weather (singleton).
func NewWeatherTool(ws weather.WeatherService) interfaces.ToolExecutor {
	weatherToolOnce.Do(func() {
		weatherToolInstance = &WeatherTool{
			weatherService: ws,
		}
	})
	return weatherToolInstance
}

// GetDefinition returns the LLM tool definition metadata.
func (t *WeatherTool) GetDefinition() models.Tool {
	return models.Tool{
		Type: "function",
		Function: models.FunctionDefinition{
			Name:        "search_weather",
			Description: "Get weather forecast information for a travel destination.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"destination": map[string]interface{}{
						"type":        "string",
						"description": "The destination name (e.g. Sikkim, Bali, Manali).",
					},
				},
				"required": []string{"destination"},
			},
		},
	}
}

// Execute parses function arguments and runs the weather search service.
func (t *WeatherTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Destination string `json:"destination"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrUnmarshalToolArgs, err)
	}

	if args.Destination == "" {
		return "", fmt.Errorf("%s: destination is required", constants.ErrExecTool)
	}

	info, err := t.weatherService.SearchWeather(ctx, args.Destination)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrExecTool, err)
	}

	resBytes, err := json.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("failed to marshal weather info: %w", err)
	}

	return string(resBytes), nil
}
