package main

import (
	"log"
	"net/http"

	"github.com/shreygarg/trip-planner-agent/clients"
	"github.com/shreygarg/trip-planner-agent/cmd/api"
	"github.com/shreygarg/trip-planner-agent/config"
	"github.com/shreygarg/trip-planner-agent/repo"
	"github.com/shreygarg/trip-planner-agent/service"
	"github.com/shreygarg/trip-planner-agent/validations"
)

func main() {
	cfg := config.NewConfigProvider()
	validator := validations.NewValidator()

	destRepo := repo.NewDestinationRepository()
	planner := service.NewDestinationPlanner(destRepo)
	toolExecutor := service.NewRecommendDestinationsTool(planner, validator)

	llmClient := clients.NewOpenAIClient(cfg)
	agent := service.NewTripAgent(llmClient, toolExecutor, cfg.GetModel())

	handler := api.NewAPIHandler(agent, validator)

	http.HandleFunc("/api/v1/trips/plan", handler.PlanTripHandler)

	port := cfg.GetPort()
	log.Printf("Starting trip planner agent server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
