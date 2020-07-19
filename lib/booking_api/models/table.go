package models

type Table struct {
	ID       TableID `json:"id"`
	Name     string  `json:"name"`
	Capacity int     `json:"capacity"`
}
