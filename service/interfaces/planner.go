package interfaces

import (
	"context"

	"github.com/shreygarg/trip-planner-agent/models"
)

// DestinationPlanner defines the contract for recommending travel destinations.
type DestinationPlanner interface {
	// RecommendDestinations recommends travel destinations based on a trip request.
	RecommendDestinations(ctx context.Context, request models.TripRequest) ([]models.DestinationRecommendation, error)
}
