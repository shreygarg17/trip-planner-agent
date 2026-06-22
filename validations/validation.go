package validations

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/utils"
	"github.com/shreygarg/trip-planner-agent/validations/interfaces"
)

var (
	validatorOnce     sync.Once
	validatorInstance interfaces.Validator
)

// RequestValidator implements requests validator contracts.
type RequestValidator struct{}

// NewValidator returns an initialized request validator (singleton).
func NewValidator() interfaces.Validator {
	validatorOnce.Do(func() {
		validatorInstance = &RequestValidator{}
	})
	return validatorInstance
}

// ValidatePlanRequest validates api prompt input requests.
func (v *RequestValidator) ValidatePlanRequest(ctx context.Context, req models.PlanRequest) error {
	reqID := utils.GetRequestID(ctx)
	if req.Prompt == "" {
		log.Printf("[%s] [WARN] Validation failed: prompt is empty", reqID)
		return errors.New(constants.ErrPromptRequired)
	}
	return nil
}

// ValidateTripRequest validates recommendations tool request parameters.
func (v *RequestValidator) ValidateTripRequest(ctx context.Context, req models.TripRequest) error {
	reqID := utils.GetRequestID(ctx)
	if req.Budget <= 0 {
		log.Printf("[%s] [WARN] Validation failed: budget <= 0", reqID)
		return errors.New(constants.ErrBudgetMustBePositive)
	}
	if req.Days <= 0 {
		log.Printf("[%s] [WARN] Validation failed: days <= 0", reqID)
		return errors.New(constants.ErrDaysMustBePositive)
	}
	return nil
}

// ValidateItineraryRequest validates itinerary generator request parameters.
func (v *RequestValidator) ValidateItineraryRequest(ctx context.Context, req models.ItineraryRequest) error {
	reqID := utils.GetRequestID(ctx)
	if req.Destination == "" {
		log.Printf("[%s] [WARN] Validation failed: destination is empty", reqID)
		return errors.New(constants.ErrDestinationRequired)
	}
	if req.Days <= 0 {
		log.Printf("[%s] [WARN] Validation failed: days <= 0", reqID)
		return errors.New(constants.ErrDaysMustBePositive)
	}
	return nil
}
