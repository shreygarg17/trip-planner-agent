package models

// PlanResponse represents the outgoing REST payload with the final trip recommendation.
type PlanResponse struct {
	Response string `json:"response"`
}

// DestinationRecommendation represents the recommendation recommendation engine outcome.
type DestinationRecommendation struct {
	Destination   string  `json:"destination"`
	EstimatedCost int     `json:"estimated_cost"`
	Score         float64 `json:"score"`
}

// ChatCompletionResponse represents the API response from an LLM.
type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice represents a single choice block returned by the LLM.
type Choice struct {
	Message Message `json:"message"`
}
