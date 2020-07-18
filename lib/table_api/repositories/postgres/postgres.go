package postgres

import (
	"context"
	"github.com/cobbinma/booking/lib/table_api/models"
)

type postgres struct {
	dbClient DBClient
}

func (p *postgres) CreateTable(ctx context.Context, newTable models.NewTable) error {
	panic("implement me")
}

func (p *postgres) GetTables(ctx context.Context, options func(options *models.GetTablesOptions) *models.GetTablesOptions) ([]models.Table, error) {
	panic("implement me")
}

func (p *postgres) DeleteTable(ctx context.Context, id models.TableID) error {
	panic("implement me")
}

func NewPostgres(client DBClient) models.Repository {
	return &postgres{dbClient: client}
}
