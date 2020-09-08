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
	venueID, ok := ctx.Value(models.VenueCtxKey).(models.VenueID)
	if !ok || venueID == "" {
		return fmt.Errorf("venue id was not in context")
	}

	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("tables").Columns("venue_id", "name", "capacity").
		Values(venueID, newTable.Name, newTable.Capacity).ToSql()
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
	venueID, ok := ctx.Value(models.VenueCtxKey).(models.VenueID)
	if !ok || venueID == "" {
		return nil, fmt.Errorf("venue id was not in context")
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
	where = append(where, sq.Eq{"venue_id": venueID})
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
