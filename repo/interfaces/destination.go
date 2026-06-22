package interfaces

import "github.com/shreygarg/trip-planner-agent/models"

// DestinationRepository defines the contract for loading destinations data.
type DestinationRepository interface {
	GetAll() []models.Destination
}
