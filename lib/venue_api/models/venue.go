package models

import (
	"fmt"
	"time"
)

type Venue struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	OpeningHours []OpeningHours `json:"openingHours"`
}

type VenueInput struct {
	Name         string         `json:"name"`
	OpeningHours []OpeningHours `json:"openingHours"`
}

type OpeningHours struct {
	DayOfWeek int       `json:"dayOfWeek"`
	Opens     timeOfDay `json:"opens"`
	Closes    timeOfDay `json:"closes"`
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
