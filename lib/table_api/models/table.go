package models

type Table struct {
	ID       TableID `json:"id"`
	Name     string  `json:"name"`
	Capacity int     `json:"capacity"`
}

type TableID string

type NewTable struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type GetTablesOptions struct {
	capacity int
}

func NewGetTablesOptions() *GetTablesOptions {
	return &GetTablesOptions{}
}

func WithCapacity(capacity int) func(options *GetTablesOptions) *GetTablesOptions {
	return func(options *GetTablesOptions) *GetTablesOptions {
		options.capacity = capacity
		return options
	}
}
