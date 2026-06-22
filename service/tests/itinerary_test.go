package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/repo"
	"github.com/shreygarg/trip-planner-agent/service"
	"github.com/shreygarg/trip-planner-agent/validations"
)

func TestGenerateItinerary_Success(t *testing.T) {
	validator := validations.NewValidator()
	itineraryRepo := repo.NewInMemoryItineraryRepository()
	gen := service.NewItineraryGenerator(validator, itineraryRepo)

	req := models.ItineraryRequest{
		Destination: "Sikkim",
		Days:        4,
	}

	plans, err := gen.GenerateItinerary(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(plans) != 4 {
		t.Fatalf("Expected exactly 4 day plans, got: %d", len(plans))
	}

	expectedActivities := []string{
		"Gangtok",
		"Tsomgo Lake",
		"Nathula Pass",
		"Local markets and return",
	}

	for i, plan := range plans {
		if plan.Day != i+1 {
			t.Errorf("Expected day %d, got %d", i+1, plan.Day)
		}
		if !strings.Contains(plan.Activity, expectedActivities[i]) {
			t.Errorf("Expected activity for day %d to contain %q, got %q", i+1, expectedActivities[i], plan.Activity)
		}
	}
}

func TestGenerateItinerary_ValidationError(t *testing.T) {
	validator := validations.NewValidator()
	itineraryRepo := repo.NewInMemoryItineraryRepository()
	gen := service.NewItineraryGenerator(validator, itineraryRepo)

	// Test 1: Empty destination
	_, err1 := gen.GenerateItinerary(context.Background(), models.ItineraryRequest{
		Destination: "",
		Days:        4,
	})
	if err1 == nil {
		t.Error("Expected error for empty destination, got nil")
	}

	// Test 2: Invalid days
	_, err2 := gen.GenerateItinerary(context.Background(), models.ItineraryRequest{
		Destination: "Sikkim",
		Days:        0,
	})
	if err2 == nil {
		t.Error("Expected error for 0 days, got nil")
	}
}
