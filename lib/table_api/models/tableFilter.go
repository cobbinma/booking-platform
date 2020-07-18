package models

type TableFilter struct {
	Capacity Capacity
}

func NewTableFilter(capacity Capacity) *TableFilter {
	return &TableFilter{Capacity: capacity}
}
