package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/cobbinma/booking/lib/table_api/models"
)

type postgres struct {
	dbClient DBClient
}

func (p *postgres) CreateTable(ctx context.Context, newTable models.NewTable) error {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("tables").Columns("name", "capacity").
		Values(newTable.Name, newTable.Capacity).ToSql()
	if err != nil {
		return fmt.Errorf("%s : %w", "could not build statement", err)
	}

	_, err = p.dbClient.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("%s : %w", "could not execute", err)
	}

	return nil
}

func (p *postgres) GetTables(ctx context.Context, filter *models.TableFilter) ([]models.Table, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "name", "capacity").From("tables")

	if filter != nil {
		where := sq.And{}
		if filter.Capacity != 0 {
			where = append(where, sq.GtOrEq{"capacity": filter.Capacity})
		}

		if filter.IDs != nil {
			where = append(where, sq.Eq{"id": filter.IDs})
		}

		builder = builder.Where(where)
	}

	sql, args, err := builder.OrderBy("capacity ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not build statement", err)
	}

	tables, err := p.dbClient.GetTables(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not get tables from db client", err)
	}

	return tables, nil
}

func (p *postgres) DeleteTable(ctx context.Context, id models.TableID) error {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("tables").
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("%s : %w", "could not build statement", err)
	}

	_, err = p.dbClient.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("%s : %w", "could not execute", err)
	}

	return nil
}

func NewPostgres(client DBClient) models.Repository {
	return &postgres{dbClient: client}
}
