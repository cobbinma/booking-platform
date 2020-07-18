package models

import (
	"fmt"
)

type CustomerID string

func (c CustomerID) Valid() error {
	if c == "" {
		return fmt.Errorf("customer id cannot be empty")
	}

	return nil
}

type TableID int

func (t TableID) Valid() error {
	if t < 0 {
		return fmt.Errorf("table id must be positive")
	}

	return nil
}
