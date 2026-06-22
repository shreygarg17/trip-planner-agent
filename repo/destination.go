package repo

import (
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/repo/interfaces"
)

// InMemoryDestinationRepository implements DestinationRepository with in-memory data.
type InMemoryDestinationRepository struct{}

// NewDestinationRepository returns a new DestinationRepository.
func NewDestinationRepository() interfaces.DestinationRepository {
	return &InMemoryDestinationRepository{}
}

// GetAll returns all travel destinations.
func (r *InMemoryDestinationRepository) GetAll() []models.Destination {
	return []models.Destination{
		{
			Name:          "Manali",
			AverageCost:   15000,
			IdealTripDays: 5,
			Tags:          []string{"mountains", "adventure"},
		},
		{
			Name:          "Kasol",
			AverageCost:   8000,
			IdealTripDays: 3,
			Tags:          []string{"mountains", "adventure"},
		},
		{
			Name:          "Sikkim",
			AverageCost:   20000,
			IdealTripDays: 7,
			Tags:          []string{"mountains", "nature"},
		},
		{
			Name:          "Bali",
			AverageCost:   80000,
			IdealTripDays: 6,
			Tags:          []string{"beach", "international", "relaxation"},
		},
		{
			Name:          "Vietnam",
			AverageCost:   70000,
			IdealTripDays: 8,
			Tags:          []string{"beach", "international", "culture"},
		},
		{
			Name:          "Bhutan",
			AverageCost:   50000,
			IdealTripDays: 5,
			Tags:          []string{"mountains", "international", "culture"},
		},
	}
}
