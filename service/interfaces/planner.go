package interfaces

import "github.com/shreygarg/trip-planner-agent/models"

// DestinationPlanner defines the contract for recommending travel destinations.
type DestinationPlanner interface {
	RecommendDestinations(request models.TripRequest) []models.DestinationRecommendation
}
