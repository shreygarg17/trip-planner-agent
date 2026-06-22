package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
	"github.com/shreygarg/trip-planner-agent/tools/longweekend"
)

var (
	longWeekendToolOnce     sync.Once
	longWeekendToolInstance interfaces.ToolExecutor
)

// LongWeekendTool wraps LongWeekendService as an LLM tool.
type LongWeekendTool struct {
	service longweekend.LongWeekendService
}

// NewGetLongWeekendsTool returns a new ToolExecutor instance for long weekends (singleton).
func NewGetLongWeekendsTool(svc longweekend.LongWeekendService) interfaces.ToolExecutor {
	longWeekendToolOnce.Do(func() {
		longWeekendToolInstance = &LongWeekendTool{
			service: svc,
		}
	})
	return longWeekendToolInstance
}

// GetDefinition returns the LLM tool definition metadata.
func (t *LongWeekendTool) GetDefinition() models.Tool {
	return models.Tool{
		Type: "function",
		Function: models.FunctionDefinition{
			Name:        "get_long_weekends",
			Description: "Get upcoming long weekends for trip planning.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"year": map[string]interface{}{
						"type":        "integer",
						"description": "The year to fetch long weekends for (e.g. 2026). Defaults to current year.",
					},
				},
			},
		},
	}
}

// Execute parses function arguments and runs the long weekend service.
func (t *LongWeekendTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Year int `json:"year"`
	}

	// It's fine if the arguments are empty; we default to current year
	if argsJSON != "" && argsJSON != "{}" {
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("%s: %w", constants.ErrUnmarshalToolArgs, err)
		}
	}

	year := args.Year
	if year == 0 {
		year = time.Now().Year()
	}

	lws, err := t.service.GetLongWeekends(ctx, year)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrExecTool, err)
	}

	resBytes, err := json.Marshal(lws)
	if err != nil {
		return "", fmt.Errorf("failed to marshal long weekends: %w", err)
	}

	return string(resBytes), nil
}
