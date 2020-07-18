package models

import "fmt"

type Capacity int

func NewCapacity(capacity int) Capacity {
	return Capacity(capacity)
}

func (c Capacity) Valid() error {
	if c < 1 {
		return fmt.Errorf("capacity cannot be less than 1")
	}

	return nil
}
