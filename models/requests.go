package models

// PlanRequest represents the incoming API payload for planning a trip.
type PlanRequest struct {
	Prompt string `json:"prompt"`
}

// TripRequest represents the structured request parameters for recommending destinations.
type TripRequest struct {
	SourceCity  string   `json:"source_city,omitempty"`
	Budget      int      `json:"budget"`
	Days        int      `json:"days"`
	Preferences []string `json:"preferences,omitempty"`
}

// ChatCompletionRequest is the payload sent to the LLM.
type ChatCompletionRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens,omitempty"`
	Messages  []Message `json:"messages"`
	Tools     []Tool    `json:"tools,omitempty"`
}

// Message represents a chat history message with roles, text, and tool calls.
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

// ToolCall represents a requested function call from the LLM.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents the arguments and name of the requested function call.
type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool represents a registered function capability of the agent.
type Tool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

// FunctionDefinition defines the description and schemas of a tool.
type FunctionDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
}
