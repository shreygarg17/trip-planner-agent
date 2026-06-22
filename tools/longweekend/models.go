package longweekend

import "time"

// LongWeekend represents a contiguous block of holiday and weekend days of 3 days or more.
type LongWeekend struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Days      int       `json:"days"`
	Reason    string    `json:"reason"`
}

// LongWeekendRequest is the input model for retrieving long weekends.
type LongWeekendRequest struct {
	Year int `json:"year"`
}
