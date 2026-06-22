package repo

import (
	"fmt"
	"log"
	"strings"

	"github.com/shreygarg/trip-planner-agent/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DestinationDB represents the GORM schema for travel destinations.
type DestinationDB struct {
	gorm.Model
	Name          string `gorm:"uniqueIndex"`
	AverageCost   int
	IdealTripDays int
	Tags          string // Comma-separated tags
}

// ActivityDB represents the GORM schema for daily activities.
type ActivityDB struct {
	gorm.Model
	DestinationName string `gorm:"index"`
	DayNumber       int
	Activity        string
}

// ConnectNeon initializes a GORM connection to the Neon Postgres database.
func ConnectNeon(cfg config.ConfigProvider) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseURL()
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is empty")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	log.Println("[INFO] Connected to Neon database successfully via GORM.")
	return db, nil
}

// AutoMigrateAndSeed performs database migrations and seeds static values if empty.
func AutoMigrateAndSeed(db *gorm.DB) error {
	log.Println("[INFO] Running database migrations...")
	if err := db.AutoMigrate(&DestinationDB{}, &ActivityDB{}); err != nil {
		return fmt.Errorf("failed to run GORM AutoMigrate: %w", err)
	}

	var count int64
	if err := db.Model(&DestinationDB{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count destinations: %w", err)
	}

	if count == 0 {
		log.Println("[INFO] Database is empty. Seeding destinations and activities...")
		if err := seedDatabase(db); err != nil {
			return fmt.Errorf("seeding failed: %w", err)
		}
	}
	return nil
}

func seedDatabase(db *gorm.DB) error {
	dests := getSeedDestinations()
	for _, d := range dests {
		destDB := DestinationDB{
			Name:          d.name,
			AverageCost:   d.cost,
			IdealTripDays: d.days,
			Tags:          strings.Join(d.tags, ","),
		}
		if err := db.Create(&destDB).Error; err != nil {
			return err
		}
		for i, act := range d.activities {
			actDB := ActivityDB{
				DestinationName: d.name,
				DayNumber:       i + 1,
				Activity:        act,
			}
			if err := db.Create(&actDB).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

type seedDest struct {
	name       string
	cost       int
	days       int
	tags       []string
	activities []string
}

func getSeedDestinations() []seedDest {
	return []seedDest{
		{
			name: "Manali", cost: 15000, days: 5, tags: []string{"mountains", "adventure"},
			activities: []string{
				"Arrive and explore Hadimba Temple & Mall Road",
				"Solang Valley adventure sports",
				"Rohtang Pass snow point excursion",
				"Jogini Waterfall trek and Vashisht hot springs",
				"Rafting in Beas river and departure",
			},
		},
		{
			name: "Kasol", cost: 8000, days: 3, tags: []string{"mountains", "adventure"},
			activities: []string{
				"Arrive and relax by Parvati River",
				"Trek to Chalal village and local cafes",
				"Manikaran Sahib hot springs visit",
				"Trek to Tosh village and scenic exploration",
			},
		},
		{
			name: "Sikkim", cost: 20000, days: 7, tags: []string{"mountains", "nature"},
			activities: []string{
				"Arrive and explore Gangtok",
				"Tsomgo Lake",
				"Nathula Pass",
				"Local markets and return",
			},
		},
		{
			name: "Bali", cost: 80000, days: 6, tags: []string{"beach", "international", "relaxation"},
			activities: []string{
				"Arrive and transfer to hotel in Seminyak",
				"Ubud temple tour and sacred monkey forest",
				"Uluwatu cliff temple and Kecak fire dance",
				"Nusa Penida island day trip",
				"Beach club relaxation and departure",
			},
		},
		{
			name: "Vietnam", cost: 70000, days: 8, tags: []string{"beach", "international", "culture"},
			activities: []string{
				"Arrive in Hanoi and explore Old Quarter",
				"Halong Bay overnight cruise",
				"Transfer to Hoi An ancient town",
				"My Son Sanctuary ruins tour",
				"Explore local street food and markets",
			},
		},
		{
			name: "Bhutan", cost: 50000, days: 5, tags: []string{"mountains", "international", "culture"},
			activities: []string{
				"Arrive in Paro and drive to Thimphu",
				"Thimphu city tour (Buddha Dordenma, Dzong)",
				"Hike to Tiger's Nest Monastery (Paro)",
				"Chele La Pass excursion and local exploration",
			},
		},
	}
}
