package constants

import "time"

const (
	// APITimeout sets the maximum execution duration for the plan trip HTTP handler.
	APITimeout = 15 * time.Second

	// HTTPTimeout sets the maximum duration for external API requests (OpenRouter, Open-Meteo).
	HTTPTimeout = 5 * time.Second

	// DBQueryTimeout sets the maximum duration for database operations.
	DBQueryTimeout = 3 * time.Second
)
