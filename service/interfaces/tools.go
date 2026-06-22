package interfaces

import (
	"context"

	"github.com/shreygarg/trip-planner-agent/models"
)

// ToolExecutor defines the contract for registering and executing a planner tool.
type ToolExecutor interface {
	// GetDefinition returns the LLM tool definition metadata.
	GetDefinition() models.Tool
	// Execute parses function arguments and runs the tool executor service.
	Execute(ctx context.Context, argsJSON string) (string, error)
}
