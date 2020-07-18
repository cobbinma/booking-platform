package models

type Table struct {
	id       string
	name     string
	capacity int
}

type NewTable struct {
	name     string
	capacity int
}
