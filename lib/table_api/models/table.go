package models

type Table struct {
	ID       TableID `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Capacity int     `json:"capacity" db:"capacity"`
}

type TableID string

type NewTable struct {
	Name     string `json:"name" db:"name"`
	Capacity int    `json:"capacity" db:"capacity"`
}
