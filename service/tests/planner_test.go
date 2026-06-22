package tests

import (
	"context"
	"testing"

	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/repo"
	"github.com/shreygarg/trip-planner-agent/service"
)

func TestRecommendDestinations_MountainPreference(t *testing.T) {
	service.ResetSingletons()
	destRepo := repo.NewInMemoryDestinationRepository()
	planner := service.NewDestinationPlanner(destRepo)

	request := models.TripRequest{
		SourceCity:  "Delhi",
		Budget:      100000,
		Days:        5,
		Preferences: []string{"mountains"},
	}

	recs, err := planner.RecommendDestinations(context.Background(), request)
	if err != nil {
		t.Fatalf("RecommendDestinations failed: %v", err)
	}

	if len(recs) == 0 {
		t.Fatal("Expected recommendations, got none")
	}

	mountainDestinations := map[string]bool{
		"Manali": true, "Kasol": true, "Sikkim": true, "Bhutan": true,
	}

	topRec := recs[0]
	if !mountainDestinations[topRec.Destination] {
		t.Errorf("Expected top recommendation to be a mountain destination, but got %s", topRec.Destination)
	}

	var manaliScore, baliScore float64
	for _, rec := range recs {
		if rec.Destination == "Manali" {
			manaliScore = rec.Score
		} else if rec.Destination == "Bali" {
			baliScore = rec.Score
		}
	}

	if manaliScore <= baliScore {
		t.Errorf("Expected Manali score (%f) to be strictly greater than Bali score (%f)", manaliScore, baliScore)
	}
}

func TestRecommendDestinations_LowBudget(t *testing.T) {
	service.ResetSingletons()
	destRepo := repo.NewInMemoryDestinationRepository()
	planner := service.NewDestinationPlanner(destRepo)

	request := models.TripRequest{
		SourceCity:  "Delhi",
		Budget:      10000,
		Days:        5,
		Preferences: []string{},
	}

	recs, err := planner.RecommendDestinations(context.Background(), request)
	if err != nil {
		t.Fatalf("RecommendDestinations failed: %v", err)
	}

	if len(recs) == 0 {
		t.Fatal("Expected recommendations, got none")
	}

	var kasolScore, baliScore float64
	for _, rec := range recs {
		if rec.Destination == "Kasol" {
			kasolScore = rec.Score
		} else if rec.Destination == "Bali" {
			baliScore = rec.Score
		}
	}

	if kasolScore <= baliScore {
		t.Errorf("Expected Kasol score (%f) to be strictly greater than Bali score (%f) for low budget", kasolScore, baliScore)
	}
}

func TestRecommendDestinations_InternationalPreference(t *testing.T) {
	service.ResetSingletons()
	destRepo := repo.NewInMemoryDestinationRepository()
	planner := service.NewDestinationPlanner(destRepo)

	request := models.TripRequest{
		SourceCity:  "Delhi",
		Budget:      100000,
		Days:        6,
		Preferences: []string{"international"},
	}

	recs, err := planner.RecommendDestinations(context.Background(), request)
	if err != nil {
		t.Fatalf("RecommendDestinations failed: %v", err)
	}

	if len(recs) == 0 {
		t.Fatal("Expected recommendations, got none")
	}

	internationalDestinations := map[string]bool{
		"Bali": true, "Vietnam": true, "Bhutan": true,
	}

	topRec := recs[0]
	if !internationalDestinations[topRec.Destination] {
		t.Errorf("Expected top recommendation to be an international destination, but got %s", topRec.Destination)
	}

	var baliScore, manaliScore float64
	for _, rec := range recs {
		if rec.Destination == "Bali" {
			baliScore = rec.Score
		} else if rec.Destination == "Manali" {
			manaliScore = rec.Score
		}
	}

	if baliScore <= manaliScore {
		t.Errorf("Expected Bali score (%f) to be strictly greater than Manali score (%f)", baliScore, manaliScore)
	}
}
