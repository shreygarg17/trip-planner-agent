package interfaces

import "context"

// ItineraryRepository defines the contract for fetching daily activities.
type ItineraryRepository interface {
	GetActivities(ctx context.Context, destination string) ([]string, error)
}
