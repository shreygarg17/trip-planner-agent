package interfaces

import (
	"context"

	"github.com/shreygarg/trip-planner-agent/models"
)

// Validator defines requests validation contracts.
type Validator interface {
	ValidatePlanRequest(ctx context.Context, req models.PlanRequest) error
	ValidateTripRequest(ctx context.Context, req models.TripRequest) error
	ValidateItineraryRequest(ctx context.Context, req models.ItineraryRequest) error
}
