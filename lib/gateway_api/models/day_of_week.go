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
	dow, err := strconv.Atoi(fmt.Sprintf("%v", v))
	if err != nil {
		return fmt.Errorf("could not parse day of week")
	}

	if dow < 1 || dow > 7 {
		return fmt.Errorf("day of week is not valid")
	}

	*i = (DayOfWeek)(dow)
	return nil
}

func (i DayOfWeek) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(fmt.Sprintf("%v", i)))
}
