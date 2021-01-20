package models

import (
	"fmt"
	"io"
	"time"
)

const timeOfDayFormat = "15:04"

type TimeOfDay string

func (t *TimeOfDay) UnmarshalGQL(v interface{}) error {
	tod, ok := v.(string)
	if !ok {
		return fmt.Errorf("time of day must be a string")
	}

	if _, err := time.Parse(timeOfDayFormat, tod); err != nil {
		return fmt.Errorf("could not parse time of day : %w", err)
	}

	*t = (TimeOfDay)(tod)
	return nil
}

func (t TimeOfDay) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(t))
}
