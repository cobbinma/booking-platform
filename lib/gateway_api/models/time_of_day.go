package models

import (
	"fmt"
	"io"
	"time"
)

const timeOfDayFormat = "15:04"

type TimeOfDay string

func NewTimeOfDay(t time.Time) TimeOfDay {
	return (TimeOfDay)(t.Format(timeOfDayFormat))
}

func (t *TimeOfDay) UnmarshalGQL(v interface{}) error {
	tod, ok := v.(string)
	if !ok {
		return fmt.Errorf("time of day must be a string")
	}

	if _, err := time.Parse(timeOfDayFormat, tod); err != nil {
		return fmt.Errorf("time of day must have format '%s'", timeOfDayFormat)
	}

	*t = (TimeOfDay)(tod)
	return nil
}

func (t TimeOfDay) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(fmt.Sprintf(`"%s"`, t)))
}

func (t TimeOfDay) Time() (time.Time, error) {
	return time.Parse(timeOfDayFormat, (string)(t))
}
