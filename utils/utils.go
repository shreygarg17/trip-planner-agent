package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

type contextKey string

const requestIDKey contextKey = "request_id"

var requestIDCounter uint64

// NormalizeString trims and lowercases a string.
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// WriteJSON writes a JSON response with status code and headers.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// ReadJSON reads and parses JSON from request body.
func ReadJSON(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

// GenerateRequestID creates a unique, thread-safe request ID.
func GenerateRequestID() string {
	return fmt.Sprintf("req-%d", atomic.AddUint64(&requestIDCounter, 1))
}

// WithRequestID binds a request ID to the given context.
func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, reqID)
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		return v.(string)
	}
	return "req-system"
}
