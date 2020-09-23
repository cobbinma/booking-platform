package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/cobbinma/booking/lib/table_api/models"
)

var _ models.Repository = (*postgres)(nil)

type postgres struct {
	dbClient DBClient
}

func (p *postgres) CreateTable(ctx context.Context, newTable models.NewTable) (*models.Table, error) {
	venue, ok := ctx.Value(models.VenueCtxKey).(models.Venue)
	if !ok {
		return nil, fmt.Errorf("venue was not in context")
	}

	args := map[string]interface{}{"venue_id": venue.ID, "name": newTable.Name, "capacity": newTable.Capacity}
	rows, err := p.dbClient.NamedQuery("INSERT INTO tables (venue_id, name, capacity) VALUES (:venue_id, :name, :capacity) RETURNING id", args)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not perform named query to insert into tables", err)
	}

	var tableID int
	if rows.Next() {
		if err := rows.Scan(&tableID); err != nil {
			return nil, fmt.Errorf("%s : %w", "could not scan row", err)
		}
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("%s : %w", "could not close rows", err)
	}

	return &models.Table{
		ID:       models.NewTableID(tableID),
		Name:     newTable.Name,
		Capacity: newTable.Capacity,
	}, nil
}

func (p *postgres) GetTables(ctx context.Context, filter *models.TableFilter) ([]models.Table, error) {
	venue, ok := ctx.Value(models.VenueCtxKey).(models.Venue)
	if !ok {
		return nil, fmt.Errorf("venue was not in context")
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "name", "capacity").From("tables")

	where := sq.And{}
	if filter != nil {
		if filter.Capacity != 0 {
			where = append(where, sq.GtOrEq{"capacity": filter.Capacity})
		}

		if filter.IDs != nil {
			where = append(where, sq.Eq{"id": filter.IDs})
		}
	}
	where = append(where, sq.Eq{"venue_id": venue.ID})
	builder = builder.Where(where)

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
