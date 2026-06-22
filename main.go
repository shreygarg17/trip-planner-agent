package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/shreygarg/trip-planner-agent/clients"
	"github.com/shreygarg/trip-planner-agent/cmd/api"
	"github.com/shreygarg/trip-planner-agent/config"
	"github.com/shreygarg/trip-planner-agent/internal/tools/weather"
	"github.com/shreygarg/trip-planner-agent/repo"
	repointerfaces "github.com/shreygarg/trip-planner-agent/repo/interfaces"
	"github.com/shreygarg/trip-planner-agent/service"
	serviceinterfaces "github.com/shreygarg/trip-planner-agent/service/interfaces"
	"github.com/shreygarg/trip-planner-agent/validations"
)

func main() {
	_ = godotenv.Load(".env.local")
	cfg := config.NewConfigProvider()
	validator := validations.NewValidator()

	var destRepo repointerfaces.DestinationRepository
	var itineraryRepo repointerfaces.ItineraryRepository

	log.Println("[INFO] Initializing repositories...")
	db, err := repo.ConnectNeon(cfg)
	if err != nil || db == nil {
		log.Printf("[WARN] Failed to connect to Neon database (falling back to in-memory datasets): %v", err)
		destRepo = repo.NewInMemoryDestinationRepository()
		itineraryRepo = repo.NewInMemoryItineraryRepository()
	} else {
		if err := repo.AutoMigrateAndSeed(db); err != nil {
			log.Fatalf("[FATAL] Database auto-migration or seeding failed: %v", err)
		}
		destRepo = repo.NewPostgresDestinationRepository(db)
		itineraryRepo = repo.NewPostgresItineraryRepository(db)
	}

	planner := service.NewDestinationPlanner(destRepo)
	itineraryGen := service.NewItineraryGenerator(validator, itineraryRepo)
	weatherService := weather.NewOpenMeteoClient(nil)

	recommendTool := service.NewRecommendDestinationsTool(planner, validator)
	itineraryTool := service.NewGenerateItineraryTool(itineraryGen)
	weatherTool := service.NewWeatherTool(weatherService)

	llmClient := clients.NewOpenAIClient(cfg)
	agent := service.NewTripAgent(llmClient, []serviceinterfaces.ToolExecutor{recommendTool, itineraryTool, weatherTool}, cfg.GetModel())

	handler := api.NewAPIHandler(agent, validator)

	http.HandleFunc("/api/v1/trips/plan", handler.PlanTripHandler)

	port := cfg.GetPort()
	log.Printf("[INFO] Starting trip planner server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("[FATAL] Server failed to start: %v", err)
	}
}
