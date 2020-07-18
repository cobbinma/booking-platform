package models

import "context"

type Repository interface {
	CreateTable(ctx context.Context, newTable NewTable) error
	GetTables(ctx context.Context, options func(options *GetTablesOptions) *GetTablesOptions) ([]Table, error)
	DeleteTable(ctx context.Context, id TableID) error
	Migrate(ctx context.Context) error
}
