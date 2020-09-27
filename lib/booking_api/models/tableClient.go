package models

import "context"

//go:generate mockgen -package=mock_models -destination=./mock/tableClient.go -source=tableClient.go
type TableClient interface {
	GetTable(ctx context.Context, id TableID) (*Table, error)
	GetTablesWithCapacity(ctx context.Context, capacity int) ([]Table, error)
}
