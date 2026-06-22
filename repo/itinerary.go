package repo

import (
	"context"
	"log"
	"sync"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/repo/interfaces"
	"github.com/shreygarg/trip-planner-agent/utils"
	"gorm.io/gorm"
)

var (
	inMemItinOnce     sync.Once
	inMemItinInstance interfaces.ItineraryRepository
)

// InMemoryItineraryRepository implements ItineraryRepository for fallback/testing.
type InMemoryItineraryRepository struct{}

// NewInMemoryItineraryRepository returns a new InMemoryItineraryRepository (singleton).
func NewInMemoryItineraryRepository() interfaces.ItineraryRepository {
	inMemItinOnce.Do(func() {
		inMemItinInstance = &InMemoryItineraryRepository{}
	})
	return inMemItinInstance
}

// GetActivities returns in-memory daily activities.
func (r *InMemoryItineraryRepository) GetActivities(ctx context.Context, destination string) ([]string, error) {
	dest := utils.NormalizeString(destination)
	switch dest {
	case "manali":
		return []string{
			"Arrive and explore Hadimba Temple & Mall Road",
			"Solang Valley adventure sports",
			"Rohtang Pass snow point excursion",
			"Jogini Waterfall trek and Vashisht hot springs",
			"Rafting in Beas river and departure",
		}, nil
	case "kasol":
		return []string{
			"Arrive and relax by Parvati River",
			"Trek to Chalal village and local cafes",
			"Manikaran Sahib hot springs visit",
			"Trek to Tosh village and scenic exploration",
		}, nil
	case "sikkim":
		return []string{
			"Arrive and explore Gangtok",
			"Tsomgo Lake",
			"Nathula Pass",
			"Local markets and return",
		}, nil
	case "bali":
		return []string{
			"Arrive and transfer to hotel in Seminyak",
			"Ubud temple tour and sacred monkey forest",
			"Uluwatu cliff temple and Kecak fire dance",
			"Nusa Penida island day trip",
			"Beach club relaxation and departure",
		}, nil
	case "vietnam":
		return []string{
			"Arrive in Hanoi and explore Old Quarter",
			"Halong Bay overnight cruise",
			"Transfer to Hoi An ancient town",
			"My Son Sanctuary ruins tour",
			"Explore local street food and markets",
		}, nil
	case "bhutan":
		return []string{
			"Arrive in Paro and drive to Thimphu",
			"Thimphu city tour (Buddha Dordenma, Dzong)",
			"Hike to Tiger's Nest Monastery (Paro)",
			"Chele La Pass excursion and local exploration",
		}, nil
	default:
		return []string{"Local sightseeing and leisure"}, nil
	}
}

var (
	pgItinOnce     sync.Once
	pgItinInstance interfaces.ItineraryRepository
)

// PostgresItineraryRepository implements ItineraryRepository using Neon database.
type PostgresItineraryRepository struct {
	db *gorm.DB
}

// NewPostgresItineraryRepository returns a new PostgresItineraryRepository (singleton).
func NewPostgresItineraryRepository(db *gorm.DB) interfaces.ItineraryRepository {
	pgItinOnce.Do(func() {
		pgItinInstance = &PostgresItineraryRepository{
			db: db,
		}
	})
	return pgItinInstance
}

// GetActivities queries day activities for a specific destination.
func (r *PostgresItineraryRepository) GetActivities(ctx context.Context, destination string) ([]string, error) {
	reqID := utils.GetRequestID(ctx)

	dbCtx, cancel := context.WithTimeout(ctx, constants.DBQueryTimeout)
	defer cancel()

	var acts []ActivityDB
	err := r.db.WithContext(dbCtx).
		Where("LOWER(destination_name) = LOWER(?)", destination).
		Order("day_number ASC").
		Find(&acts).Error

	if err != nil {
		log.Printf("[%s] [ERROR] Failed to query activities for destination %s: %v", reqID, destination, err)
		return nil, err
	}

	if len(acts) == 0 {
		return []string{"Local sightseeing and leisure"}, nil
	}

	result := make([]string, 0, len(acts))
	for _, a := range acts {
		result = append(result, a.Activity)
	}

	return result, nil
}
