package interfaces

// Agent defines the contract for processing user planning prompts.
type Agent interface {
	PlanTrip(prompt string) (string, error)
}
