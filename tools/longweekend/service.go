package longweekend

import (
	"context"
	"sort"
	"sync"
	"time"
)

// LongWeekendService defines contracts for identifying long weekends.
type LongWeekendService interface {
	GetLongWeekends(ctx context.Context, year int) ([]LongWeekend, error)
	GetNextLongWeekend(ctx context.Context) (*LongWeekend, error)
}

type offDayInfo struct {
	IsWeekend   bool
	HolidayName string
}

var (
	serviceOnce     sync.Once
	serviceInstance LongWeekendService
)

// LongWeekendServiceImpl implements LongWeekendService.
type LongWeekendServiceImpl struct {
	holidayRepo HolidayRepository
}

// NewLongWeekendService creates an initialized LongWeekendService (singleton).
func NewLongWeekendService(holidayRepo HolidayRepository) LongWeekendService {
	serviceOnce.Do(func() {
		serviceInstance = &LongWeekendServiceImpl{
			holidayRepo: holidayRepo,
		}
	})
	return serviceInstance
}

// GetLongWeekends retrieves all contiguous long weekends of 3 days or more.
func (s *LongWeekendServiceImpl) GetLongWeekends(ctx context.Context, year int) ([]LongWeekend, error) {
	holidays, err := s.holidayRepo.GetHolidays(ctx, year)
	if err != nil {
		return nil, err
	}

	weekends := s.getWeekendsForYear(year)
	offDays := s.buildOffDaysMap(weekends, holidays)

	var dates []time.Time
	for dtStr := range offDays {
		t, _ := time.Parse("2006-01-02", dtStr)
		dates = append(dates, t)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	groups := s.groupContiguousDates(dates)
	return s.buildLongWeekends(groups, offDays), nil
}

// GetNextLongWeekend finds the nearest upcoming long weekend relative to time.Now().
func (s *LongWeekendServiceImpl) GetNextLongWeekend(ctx context.Context) (*LongWeekend, error) {
	now := time.Now()
	year := now.Year()

	currentYearWeekends, err := s.GetLongWeekends(ctx, year)
	if err != nil {
		return nil, err
	}

	// Filter for future long weekends in the current year
	for _, lw := range currentYearWeekends {
		if lw.StartDate.After(now) || lw.StartDate.Format("2006-01-02") == now.Format("2006-01-02") {
			return &lw, nil
		}
	}

	// If none found, check the next year
	nextYearWeekends, err := s.GetLongWeekends(ctx, year+1)
	if err != nil {
		return nil, err
	}

	if len(nextYearWeekends) > 0 {
		return &nextYearWeekends[0], nil
	}

	return nil, nil
}

func (s *LongWeekendServiceImpl) getWeekendsForYear(year int) map[string]bool {
	weekends := make(map[string]bool)
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	for start.Year() == year {
		if start.Weekday() == time.Saturday || start.Weekday() == time.Sunday {
			weekends[start.Format("2006-01-02")] = true
		}
		start = start.AddDate(0, 0, 1)
	}
	return weekends
}

func (s *LongWeekendServiceImpl) buildOffDaysMap(weekends map[string]bool, holidays []Holiday) map[string]*offDayInfo {
	offDays := make(map[string]*offDayInfo)
	for dtStr := range weekends {
		offDays[dtStr] = &offDayInfo{IsWeekend: true}
	}
	for _, hol := range holidays {
		dtStr := hol.Date.Format("2006-01-02")
		if info, exists := offDays[dtStr]; exists {
			info.HolidayName = hol.Name
		} else {
			offDays[dtStr] = &offDayInfo{HolidayName: hol.Name}
		}
	}
	return offDays
}

func (s *LongWeekendServiceImpl) groupContiguousDates(dates []time.Time) [][]time.Time {
	var groups [][]time.Time
	if len(dates) == 0 {
		return groups
	}
	currentGroup := []time.Time{dates[0]}
	for i := 1; i < len(dates); i++ {
		prev := dates[i-1]
		curr := dates[i]
		if curr.Sub(prev) <= 24*time.Hour {
			currentGroup = append(currentGroup, curr)
		} else {
			groups = append(groups, currentGroup)
			currentGroup = []time.Time{curr}
		}
	}
	groups = append(groups, currentGroup)
	return groups
}

func (s *LongWeekendServiceImpl) buildLongWeekends(groups [][]time.Time, offDays map[string]*offDayInfo) []LongWeekend {
	var longWeekends []LongWeekend
	for _, group := range groups {
		if len(group) < 3 {
			continue
		}

		var holidayNames []string
		for _, dt := range group {
			dtStr := dt.Format("2006-01-02")
			if info, exists := offDays[dtStr]; exists && info.HolidayName != "" {
				holidayNames = append(holidayNames, info.HolidayName)
			}
		}

		reason := "Long Weekend"
		if len(holidayNames) > 0 {
			reason = holidayNames[0] + " Weekend"
		}

		longWeekends = append(longWeekends, LongWeekend{
			StartDate: group[0],
			EndDate:   group[len(group)-1],
			Days:      len(group),
			Reason:    reason,
		})
	}
	return longWeekends
}
