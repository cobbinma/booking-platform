package models

import "context"

type TableClient interface {
	GetTable(ctx context.Context, id TableID) (*Table, error)
}
