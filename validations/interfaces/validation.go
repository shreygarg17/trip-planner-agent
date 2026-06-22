package interfaces

import "github.com/shreygarg/trip-planner-agent/models"

// Validator defines requests validation contracts.
type Validator interface {
	ValidatePlanRequest(req models.PlanRequest) error
	ValidateTripRequest(req models.TripRequest) error
}
