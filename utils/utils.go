package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
