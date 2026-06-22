package models

// Destination represents a location entry in our dataset.
type Destination struct {
	Name          string
	AverageCost   int
	IdealTripDays int
	Tags          []string
}
