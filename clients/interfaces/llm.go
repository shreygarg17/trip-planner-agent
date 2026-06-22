package interfaces

import "github.com/shreygarg/trip-planner-agent/models"

// LLMClient defines the interface for communicating with an LLM.
type LLMClient interface {
	CreateChatCompletion(request models.ChatCompletionRequest) (*models.ChatCompletionResponse, error)
}
