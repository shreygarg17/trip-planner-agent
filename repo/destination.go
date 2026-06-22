package repo

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/repo/interfaces"
	"github.com/shreygarg/trip-planner-agent/utils"
	"gorm.io/gorm"
)

var (
	inMemDestOnce     sync.Once
	inMemDestInstance interfaces.DestinationRepository

	pgDestOnce     sync.Once
	pgDestInstance interfaces.DestinationRepository
)

// InMemoryDestinationRepository implements DestinationRepository with in-memory data.
type InMemoryDestinationRepository struct{}

// NewInMemoryDestinationRepository returns a new InMemoryDestinationRepository (singleton).
func NewInMemoryDestinationRepository() interfaces.DestinationRepository {
	inMemDestOnce.Do(func() {
		inMemDestInstance = &InMemoryDestinationRepository{}
	})
	return inMemDestInstance
}

// GetAll returns all in-memory travel destinations.
func (r *InMemoryDestinationRepository) GetAll(ctx context.Context) ([]models.Destination, error) {
	return []models.Destination{
		{Name: "Manali", AverageCost: 15000, IdealTripDays: 5, Tags: []string{"mountains", "adventure"}},
		{Name: "Kasol", AverageCost: 8000, IdealTripDays: 3, Tags: []string{"mountains", "adventure"}},
		{Name: "Sikkim", AverageCost: 20000, IdealTripDays: 7, Tags: []string{"mountains", "nature"}},
		{Name: "Bali", AverageCost: 80000, IdealTripDays: 6, Tags: []string{"beach", "international", "relaxation"}},
		{Name: "Vietnam", AverageCost: 70000, IdealTripDays: 8, Tags: []string{"beach", "international", "culture"}},
		{Name: "Bhutan", AverageCost: 50000, IdealTripDays: 5, Tags: []string{"mountains", "international", "culture"}},
	}, nil
}

// PostgresDestinationRepository implements DestinationRepository querying GORM database.
type PostgresDestinationRepository struct {
	db *gorm.DB
}

// NewPostgresDestinationRepository returns a new PostgresDestinationRepository (singleton).
func NewPostgresDestinationRepository(db *gorm.DB) interfaces.DestinationRepository {
	pgDestOnce.Do(func() {
		pgDestInstance = &PostgresDestinationRepository{
			db: db,
		}
	})
	return pgDestInstance
}

// GetAll queries all destinations from the Neon GORM database.
func (r *PostgresDestinationRepository) GetAll(ctx context.Context) ([]models.Destination, error) {
	reqID := utils.GetRequestID(ctx)

	dbCtx, cancel := context.WithTimeout(ctx, constants.DBQueryTimeout)
	defer cancel()

	var dbDests []DestinationDB
	if err := r.db.WithContext(dbCtx).Find(&dbDests).Error; err != nil {
		log.Printf("[%s] [ERROR] Failed to query destinations from DB: %v", reqID, err)
		return nil, err
	}

	dests := make([]models.Destination, 0, len(dbDests))
	for _, d := range dbDests {
		var tags []string
		if d.Tags != "" {
			tags = strings.Split(d.Tags, ",")
		}
		dests = append(dests, models.Destination{
			Name:          d.Name,
			AverageCost:   d.AverageCost,
			IdealTripDays: d.IdealTripDays,
			Tags:          tags,
		})
	}
	return dests, nil
}
