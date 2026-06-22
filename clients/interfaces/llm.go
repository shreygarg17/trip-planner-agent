package interfaces

import (
	"context"

	"github.com/shreygarg/trip-planner-agent/models"
)

// LLMClient defines the interface for communicating with an LLM.
type LLMClient interface {
	CreateChatCompletion(ctx context.Context, request models.ChatCompletionRequest) (*models.ChatCompletionResponse, error)
}
