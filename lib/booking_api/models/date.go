package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Date time.Time

const DateFormat string = "2006-01-02"

func DateFromString(s string) (Date, error) {
	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return Date(time.Time{}), err
	}

	return Date(t), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(DateFormat))
}

func (d Date) Format(s string) string {
	t := time.Time(d)
	return t.Format(s)
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

func (d *Date) Scan(src interface{}) error {
	var value time.Time
	switch src.(type) {
	case time.Time:
		value = src.(time.Time)
	default:
		return fmt.Errorf("invalid type for Date")
	}
	*d = Date(value)
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return driver.Value(d.Time()), nil
}
