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

// DayPlan represents a single day activity in a travel itinerary.
type DayPlan struct {
	Day      int    `json:"day"`
	Activity string `json:"activity"`
}

// Coordinates represents geo-coordinates for a location.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// WeatherInfo represents aggregated forecast statistics for a destination.
type WeatherInfo struct {
	TemperatureMax  float64 `json:"temperature_max"`
	TemperatureMin  float64 `json:"temperature_min"`
	RainProbability float64 `json:"rain_probability"`
	Condition       string  `json:"condition"`
}

// ChatCompletionResponse represents the API response from an LLM.
type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice represents a single choice block returned by the LLM.
type Choice struct {
	Message Message `json:"message"`
}
