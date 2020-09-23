package models

import "context"

type Repository interface {
	CreateTable(ctx context.Context, newTable NewTable) (*Table, error)
	GetTables(ctx context.Context, filter *TableFilter) ([]Table, error)
	DeleteTable(ctx context.Context, id TableID) error
	Migrate(ctx context.Context, sourceURL string) error
}
