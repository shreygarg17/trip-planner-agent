package interfaces

import (
	"context"

	"github.com/shreygarg/trip-planner-agent/models"
)

// ItineraryGenerator defines the contract for generating a daily travel itinerary.
type ItineraryGenerator interface {
	// GenerateItinerary generates a day-by-day travel itinerary for a destination.
	GenerateItinerary(ctx context.Context, req models.ItineraryRequest) ([]models.DayPlan, error)
}
