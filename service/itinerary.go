package service

import (
	"context"
	"sync"

	"github.com/shreygarg/trip-planner-agent/models"
	repointerfaces "github.com/shreygarg/trip-planner-agent/repo/interfaces"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
	validationinterfaces "github.com/shreygarg/trip-planner-agent/validations/interfaces"
)

var (
	itineraryOnce     sync.Once
	itineraryInstance interfaces.ItineraryGenerator
)

// ItineraryGenerator implements the interfaces.ItineraryGenerator contract.
type ItineraryGenerator struct {
	validator     validationinterfaces.Validator
	itineraryRepo repointerfaces.ItineraryRepository
}

// NewItineraryGenerator returns a new ItineraryGenerator instance (singleton).
func NewItineraryGenerator(validator validationinterfaces.Validator, itineraryRepo repointerfaces.ItineraryRepository) interfaces.ItineraryGenerator {
	itineraryOnce.Do(func() {
		itineraryInstance = &ItineraryGenerator{
			validator:     validator,
			itineraryRepo: itineraryRepo,
		}
	})
	return itineraryInstance
}

// GenerateItinerary generates a day-by-day travel itinerary for a destination.
func (g *ItineraryGenerator) GenerateItinerary(ctx context.Context, req models.ItineraryRequest) ([]models.DayPlan, error) {
	if err := g.validator.ValidateItineraryRequest(ctx, req); err != nil {
		return nil, err
	}

	activities, err := g.itineraryRepo.GetActivities(ctx, req.Destination)
	if err != nil {
		activities = []string{"Local sightseeing and leisure"}
	}

	plans := make([]models.DayPlan, 0, req.Days)
	for i := 1; i <= req.Days; i++ {
		activity := g.getActivityForDay(activities, i)
		plans = append(plans, models.DayPlan{
			Day:      i,
			Activity: activity,
		})
	}

	return plans, nil
}

func (g *ItineraryGenerator) getActivityForDay(activities []string, day int) string {
	idx := day - 1
	if idx < len(activities) {
		return activities[idx]
	}
	return "Local sightseeing, relaxation, and leisure"
}
