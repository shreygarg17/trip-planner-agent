package interfaces

import (
	"context"

	"github.com/shreygarg/trip-planner-agent/models"
)

// DestinationRepository defines the contract for loading destinations data.
type DestinationRepository interface {
	GetAll(ctx context.Context) ([]models.Destination, error)
}
