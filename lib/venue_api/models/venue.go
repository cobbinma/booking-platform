package models

import (
	"fmt"
	"time"
)

type VenueID int

func NewVenueID(id int) VenueID {
	return VenueID(id)
}

type Venue struct {
	ID           VenueID        `json:"id" db:"id"`
	Name         string         `json:"name" db:"name"`
	OpeningHours []OpeningHours `json:"openingHours"`
}

type VenueInput struct {
	Name         string         `json:"name"`
	OpeningHours []OpeningHours `json:"openingHours"`
}

func (vi VenueInput) Valid() error {
	if vi.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	days := map[int]bool{}
	for i := range vi.OpeningHours {
		if err := vi.OpeningHours[i].Valid(); err != nil {
			return fmt.Errorf("opening hours was not valid : %w", err)
		}
		if days[vi.OpeningHours[i].DayOfWeek] {
			return fmt.Errorf("cannot have duplicate day opening hours")
		}
		days[vi.OpeningHours[i].DayOfWeek] = true
	}

	return nil
}

type OpeningHours struct {
	DayOfWeek int       `json:"dayOfWeek" db:"day_of_week"`
	Opens     timeOfDay `json:"opens" db:"opens"`
	Closes    timeOfDay `json:"closes" db:"closes"`
}

func (oh OpeningHours) Valid() error {
	if oh.DayOfWeek > 6 {
		return fmt.Errorf("day of week cannot be greater than 6")
	}
	if oh.Opens.Time().IsZero() || oh.Closes.Time().IsZero() {
		return fmt.Errorf("times cannot be zero")
	}

	return nil
}

type timeOfDay time.Time

const timeLayout = "15:04"

func (t timeOfDay) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(timeLayout) + `"`), nil
}

func (t *timeOfDay) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) != 7 {
		return fmt.Errorf("time should be a string formatted as \"15:04\"")
	}
	ret, err := time.Parse(timeLayout, s[1:6])
	if err != nil {
		return err
	}
	tod := timeOfDay(ret)
	*t = tod
	return nil
}

func (t timeOfDay) Time() time.Time {
	return time.Time(t)
}
