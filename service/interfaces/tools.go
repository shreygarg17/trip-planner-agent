package interfaces

import "github.com/shreygarg/trip-planner-agent/models"

// ToolExecutor defines the contract for registering and executing a planner tool.
type ToolExecutor interface {
	GetDefinition() models.Tool
	Execute(argsJSON string) (string, error)
}
