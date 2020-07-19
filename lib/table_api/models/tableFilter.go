package models

type TableFilter struct {
	Capacity Capacity
	IDs      []TableID
}

func NewTableFilter(capacity Capacity, IDs []TableID) *TableFilter {
	return &TableFilter{Capacity: capacity, IDs: IDs}
}
