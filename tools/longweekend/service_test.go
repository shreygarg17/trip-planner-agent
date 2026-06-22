package longweekend

import (
	"context"
	"sync"
	"testing"
	"time"
)

type mockHolidayRepository struct {
	holidays []Holiday
}

func (m *mockHolidayRepository) GetHolidays(ctx context.Context, year int) ([]Holiday, error) {
	var result []Holiday
	for _, h := range m.holidays {
		if h.Date.Year() == year {
			result = append(result, h)
		}
	}
	return result, nil
}

func TestGetLongWeekends_CombinedHolidaysAndSorting(t *testing.T) {
	// Reset service singleton for testing
	serviceOnce = sync.Once{}
	serviceInstance = nil

	hols := []Holiday{
		{Date: time.Date(2026, time.January, 26, 0, 0, 0, 0, time.UTC), Name: "Republic Day"}, // Monday
		{Date: time.Date(2026, time.April, 3, 0, 0, 0, 0, time.UTC), Name: "Good Friday"},   // Friday
	}
	repo := &mockHolidayRepository{holidays: hols}
	svc := NewLongWeekendService(repo)

	lws, err := svc.GetLongWeekends(context.Background(), 2026)
	if err != nil {
		t.Fatalf("Failed to get long weekends: %v", err)
	}

	if len(lws) < 2 {
		t.Fatalf("Expected at least 2 long weekends, got %d", len(lws))
	}

	// Verify Republic Day weekend sorting & details
	repDayLw := lws[0]
	if repDayLw.Reason != "Republic Day Weekend" {
		t.Errorf("Expected 'Republic Day Weekend', got %q", repDayLw.Reason)
	}
	if repDayLw.Days != 3 {
		t.Errorf("Expected 3 days, got %d", repDayLw.Days)
	}
	if repDayLw.StartDate.Format("2006-01-02") != "2026-01-24" {
		t.Errorf("Expected start date 2026-01-24, got %s", repDayLw.StartDate.Format("2006-01-02"))
	}
	if repDayLw.EndDate.Format("2006-01-02") != "2026-01-26" {
		t.Errorf("Expected end date 2026-01-26, got %s", repDayLw.EndDate.Format("2006-01-02"))
	}
}

func TestGetNextLongWeekend(t *testing.T) {
	// Reset service singleton for testing
	serviceOnce = sync.Once{}
	serviceInstance = nil

	now := time.Now()
	// Create a future holiday at least 1 day in the future, forced to be Friday or Monday
	futureHolidayDate := now.AddDate(0, 0, 2)
	for futureHolidayDate.Weekday() != time.Friday && futureHolidayDate.Weekday() != time.Monday {
		futureHolidayDate = futureHolidayDate.AddDate(0, 0, 1)
	}

	hols := []Holiday{
		{Date: futureHolidayDate, Name: "Future Holiday"},
	}
	repo := &mockHolidayRepository{holidays: hols}
	svc := NewLongWeekendService(repo)

	lw, err := svc.GetNextLongWeekend(context.Background())
	if err != nil {
		t.Fatalf("Failed to get next long weekend: %v", err)
	}

	if lw == nil {
		t.Fatal("Expected a long weekend but got nil")
	}

	if lw.Reason != "Future Holiday Weekend" {
		t.Errorf("Expected 'Future Holiday Weekend', got %q", lw.Reason)
	}
}
