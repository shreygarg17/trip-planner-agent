package service

import (
	"context"
	"math"
	"sort"
	"sync"

	"github.com/shreygarg/trip-planner-agent/models"
	repointerfaces "github.com/shreygarg/trip-planner-agent/repo/interfaces"
	"github.com/shreygarg/trip-planner-agent/service/interfaces"
	"github.com/shreygarg/trip-planner-agent/utils"
)

var (
	plannerOnce     sync.Once
	plannerInstance interfaces.DestinationPlanner
)

// InMemoryPlanner implements DestinationPlanner logic.
type InMemoryPlanner struct {
	destRepo repointerfaces.DestinationRepository
}

// NewDestinationPlanner instantiates a new DestinationPlanner service (singleton).
func NewDestinationPlanner(destRepo repointerfaces.DestinationRepository) interfaces.DestinationPlanner {
	plannerOnce.Do(func() {
		plannerInstance = &InMemoryPlanner{
			destRepo: destRepo,
		}
	})
	return plannerInstance
}

// RecommendDestinations recommends travel destinations based on a trip request.
func (p *InMemoryPlanner) RecommendDestinations(ctx context.Context, request models.TripRequest) ([]models.DestinationRecommendation, error) {
	dests, err := p.destRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	recommendations := make([]models.DestinationRecommendation, 0, len(dests))
	for _, d := range dests {
		bScore := p.scoreBudget(request.Budget, d.AverageCost)
		pScore := p.scorePreferences(request.Preferences, d.Tags)
		dScore := p.scoreDuration(request.Days, d.IdealTripDays)

		score := (bScore * 0.40) + (pScore * 0.40) + (dScore * 0.20)

		recommendations = append(recommendations, models.DestinationRecommendation{
			Destination:   d.Name,
			EstimatedCost: d.AverageCost,
			Score:         score,
		})
	}

	p.sortRecommendations(recommendations)
	return recommendations, nil
}

func (p *InMemoryPlanner) scoreBudget(budget, averageCost int) float64 {
	if budget <= 0 {
		return 0.0
	}
	if averageCost <= budget {
		return 1.0
	}
	return float64(budget) / float64(averageCost)
}

func (p *InMemoryPlanner) scorePreferences(reqPrefs, tags []string) float64 {
	if len(reqPrefs) == 0 {
		return 1.0
	}

	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[utils.NormalizeString(tag)] = true
	}

	matched := 0
	for _, pref := range reqPrefs {
		if tagMap[utils.NormalizeString(pref)] {
			matched++
		}
	}
	return float64(matched) / float64(len(reqPrefs))
}

func (p *InMemoryPlanner) scoreDuration(reqDays, idealDays int) float64 {
	if idealDays <= 0 {
		return 0.0
	}
	diff := math.Abs(float64(reqDays - idealDays))
	score := 1.0 - (diff / float64(idealDays))
	if score < 0 {
		return 0.0
	}
	return score
}

func (p *InMemoryPlanner) sortRecommendations(recs []models.DestinationRecommendation) {
	sort.Slice(recs, func(i, j int) bool {
		if recs[i].Score == recs[j].Score {
			return recs[i].Destination < recs[j].Destination
		}
		return recs[i].Score > recs[j].Score
	})
}
