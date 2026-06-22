package constants

const (
	ErrMethodNotAllowed      = "Method Not Allowed"
	ErrInvalidJSON           = "Invalid request JSON"
	ErrPromptRequired        = "Prompt field is required"
	ErrBudgetMustBePositive  = "budget must be greater than zero"
	ErrDaysMustBePositive    = "duration must be greater than zero"
	ErrMaxIterationsExceeded = "agent exceeded maximum iterations"
	ErrUnmarshalToolArgs     = "failed to unmarshal recommend_destinations args"
	ErrExecTool              = "failed executing recommend_destinations"
	ErrLLMEmptyChoices       = "llm returned an empty choice set"
)
