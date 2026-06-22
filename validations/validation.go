package validations

import (
	"errors"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/validations/interfaces"
)

// RequestValidator implements requests validator contracts.
type RequestValidator struct{}

// NewValidator returns an initialized request validator.
func NewValidator() interfaces.Validator {
	return &RequestValidator{}
}

// ValidatePlanRequest validates api prompt input requests.
func (v *RequestValidator) ValidatePlanRequest(req models.PlanRequest) error {
	if req.Prompt == "" {
		return errors.New(constants.ErrPromptRequired)
	}
	return nil
}

// ValidateTripRequest validates recommendations tool request parameters.
func (v *RequestValidator) ValidateTripRequest(req models.TripRequest) error {
	if req.Budget <= 0 {
		return errors.New(constants.ErrBudgetMustBePositive)
	}
	if req.Days <= 0 {
		return errors.New(constants.ErrDaysMustBePositive)
	}
	return nil
}
