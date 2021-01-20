package models

import (
	"fmt"
	"io"
	"time"
)

const dateFormat = "02-01-2006"

type Date string

func (d *Date) UnmarshalGQL(v interface{}) error {
	dt, ok := v.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}

	if _, err := time.Parse(dateFormat, dt); err != nil {
		return fmt.Errorf("date must have format '%s'", dateFormat)
	}

	*d = (Date)(dt)
	return nil
}

func (d Date) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(d))
}
