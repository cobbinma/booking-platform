package models

import "fmt"

type Table struct {
	ID       TableID `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Capacity int     `json:"capacity" db:"capacity"`
}

type TableID int

type NewTable struct {
	Name     string `json:"name" db:"name"`
	Capacity int    `json:"capacity" db:"capacity"`
}

func (nt *NewTable) Valid() error {
	if nt.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if nt.Capacity == 0 {
		return fmt.Errorf("capacity must be greater than zero")
	}

	return nil
}
