package models

import (
	"fmt"
	"strconv"
	"time"
)

const VenueCtxKey = "venue"

type VenueID int

func NewVenueID(id int) VenueID {
	return VenueID(id)
}

func (vid VenueID) String() string {
	return strconv.Itoa(int(vid))
}

type Venue struct {
	ID           VenueID        `json:"id" db:"id"`
	Name         string         `json:"name" db:"name"`
	OpeningHours []OpeningHours `json:"openingHours"`
}

func (v Venue) IsOpen(day int, startsAt time.Time, endsAt time.Time) bool {
	for i := range v.OpeningHours {
		if v.OpeningHours[i].DayOfWeek == day {
			openMins := v.OpeningHours[i].Opens.Time().Hour()*60 + v.OpeningHours[i].Opens.Time().Minute()
			closeMins := v.OpeningHours[i].Closes.Time().Hour()*60 + v.OpeningHours[i].Closes.Time().Minute()
			startsMins := startsAt.Hour()*60 + startsAt.Minute()
			endsMins := endsAt.Hour()*60 + endsAt.Minute()
			if openMins <= startsMins && closeMins >= endsMins {
				return true
			}
			break
		}
	}
	return false
}

type OpeningHours struct {
	DayOfWeek int       `json:"dayOfWeek" db:"day_of_week"`
	Opens     TimeOfDay `json:"opens" db:"opens"`
	Closes    TimeOfDay `json:"closes" db:"closes"`
}

type TimeOfDay time.Time

const timeLayout = "15:04"

func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(timeLayout) + `"`), nil
}

func (t *TimeOfDay) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) != 7 {
		return fmt.Errorf("time should be a string formatted as \"15:04\"")
	}
	ret, err := time.Parse(timeLayout, s[1:6])
	if err != nil {
		return err
	}
	tod := TimeOfDay(ret)
	*t = tod
	return nil
}

func (t TimeOfDay) Time() time.Time {
	return time.Time(t)
}
