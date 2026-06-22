package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/validations"
)

type mockAgent struct{}

func (a *mockAgent) PlanTrip(ctx context.Context, prompt string) (string, error) {
	return "I recommend Sikkim as it matches your mountain preference and fits under ₹50k.", nil
}

func TestPlanTripHandler(t *testing.T) {
	agent := &mockAgent{}
	validator := validations.NewValidator()
	handler := NewAPIHandler(agent, validator)

	// Test 1: Successful POST Request
	reqBody, _ := json.Marshal(models.PlanRequest{Prompt: "I have 4 days and ₹50k budget. I like mountains."})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/trips/plan", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()

	handler.PlanTripHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var resp models.PlanResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response payload: %v", err)
	}

	expectedSnippet := "Sikkim"
	if !strings.Contains(resp.Response, expectedSnippet) {
		t.Errorf("Expected response content to mention %q, got: %q", expectedSnippet, resp.Response)
	}

	// Test 2: Method Not Allowed
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/trips/plan", nil)
	recGet := httptest.NewRecorder()

	handler.PlanTripHandler(recGet, reqGet)

	if recGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, recGet.Code)
	}
}
