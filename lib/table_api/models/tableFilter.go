package models

type TableFilter struct {
	Capacity int
}

func NewTableFilter(capacity int) *TableFilter {
	return &TableFilter{Capacity: capacity}
}
