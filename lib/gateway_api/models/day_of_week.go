package models

import (
	"fmt"
	"io"
	"strconv"
)

type DayOfWeek int

const (
	Monday    DayOfWeek = iota + 1
	Tuesday   DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday  DayOfWeek = 4
	Friday    DayOfWeek = 5
	Saturday  DayOfWeek = 6
	Sunday    DayOfWeek = 7
)

func (i *DayOfWeek) UnmarshalGQL(v interface{}) error {
	dow, ok := v.(int)
	if !ok {
		return fmt.Errorf("day of week must be an integer")
	}

	if dow < 1 || dow > 7 {
		return fmt.Errorf("day of week is not valid")
	}

	*i = (DayOfWeek)(dow)
	return nil
}

func (i DayOfWeek) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Itoa((int)(i))))
}
