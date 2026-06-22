package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// Cors configures and returns a CORS middleware handler for browser-based clients.
// This is required to allow the Next.js frontend (running on http://localhost:3000) to
// communicate with the Go backend (running on http://localhost:8080) without encountering
// browser-enforced Cross-Origin Resource Sharing (CORS) security restrictions.
func Cors() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"*",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
