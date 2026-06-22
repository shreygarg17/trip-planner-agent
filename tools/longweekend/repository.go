package longweekend

import (
	"context"
	"time"
)

// Holiday represents a public holiday date and name.
type Holiday struct {
	Date time.Time
	Name string
}

// HolidayRepository defines the contract for loading holidays.
type HolidayRepository interface {
	GetHolidays(ctx context.Context, year int) ([]Holiday, error)
}

// StaticHolidayRepository provides a hardcoded list of holidays.
type StaticHolidayRepository struct{}

// NewStaticHolidayRepository returns a new StaticHolidayRepository.
func NewStaticHolidayRepository() HolidayRepository {
	return &StaticHolidayRepository{}
}

// GetHolidays returns the list of public holidays for a given year.
func (r *StaticHolidayRepository) GetHolidays(ctx context.Context, year int) ([]Holiday, error) {
	if year == 2026 {
		return r.get2026Holidays(), nil
	}
	return []Holiday{}, nil
}

func (r *StaticHolidayRepository) get2026Holidays() []Holiday {
	return []Holiday{
		{Date: time.Date(2026, time.January, 26, 0, 0, 0, 0, time.UTC), Name: "Republic Day"},
		{Date: time.Date(2026, time.March, 3, 0, 0, 0, 0, time.UTC), Name: "Holi"},
		{Date: time.Date(2026, time.April, 3, 0, 0, 0, 0, time.UTC), Name: "Good Friday"},
		// Independence Day block with mock holidays to support the 4-day weekend (Aug 15 - Aug 18)
		{Date: time.Date(2026, time.August, 15, 0, 0, 0, 0, time.UTC), Name: "Independence Day"},
		{Date: time.Date(2026, time.August, 17, 0, 0, 0, 0, time.UTC), Name: "Restricted Holiday"},
		{Date: time.Date(2026, time.August, 18, 0, 0, 0, 0, time.UTC), Name: "Parsi New Year"},
		// Other holidays
		{Date: time.Date(2026, time.October, 2, 0, 0, 0, 0, time.UTC), Name: "Gandhi Jayanti"},
		{Date: time.Date(2026, time.November, 8, 0, 0, 0, 0, time.UTC), Name: "Diwali"},
		{Date: time.Date(2026, time.December, 25, 0, 0, 0, 0, time.UTC), Name: "Christmas"},
	}
}
