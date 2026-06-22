package interfaces

import "context"

// Agent defines the contract for processing user planning prompts.
type Agent interface {
	// PlanTrip processes a planning prompt, resolving requested tool calls iteratively.
	PlanTrip(ctx context.Context, prompt string) (string, error)
}
