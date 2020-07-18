package models

import "context"

type TableClient interface {
	GetTable(ctx context.Context, id TableID) (*Table, error)
	GetTablesWithCapacity(ctx context.Context, capacity int) ([]Table, error)
}
