package db

import (
	"context"
	"crud_service/app/repos"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Database storage of shops
type shops struct {
	connection *pgx.Conn
}

func newShops(connection *pgx.Conn) *shops {
	return &shops{connection}
}

func (cs *shops) create(ctx context.Context, c repos.Shop) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var err error
	var raw pgx.Rows
	if c.Id == uuid.Nil {
		raw, err = cs.connection.Query(ctx, `INSERT INTO shops (id, name, address, is_closed, owner) 
	values (DEFAULT, $1, $2, $3, $4) RETURNING id;`,
			c.Name, c.Address, c.IsClosed, c.Owner)
	} else {
		raw, err = cs.connection.Query(ctx, `INSERT INTO shops (id, name, address, is_closed, owner) 
	values ($1, $2, $3, $4, $5) RETURNING id;`,
			c.Name, c.Address, c.IsClosed, c.Owner)
	}

	if err != nil {
		return nil, err
	}
	defer raw.Close()

	raw.Next()
	err = raw.Scan(&c.Id)
	if err != nil {
		return nil, err
	}
	return &c.Id, nil
}

func (cs *shops) readName(ctx context.Context, name string) ([]repos.Data, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	raw, err := cs.connection.Query(ctx, `SELECT * FROM shops WHERE name = $1`, name)
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	data := []repos.Data{}
	for raw.Next() {
		row := &repos.Shop{}
		err := raw.Scan(&row.Id, &row.Name, &row.Address, &row.IsClosed, &row.Owner)
		if err != nil {
			return data, err
		}
		data = append(data, row)
	}
	if len(data) < 1 {
		return nil, fmt.Errorf("no shops with such name")
	}
	return data, nil
}

func (ss *shops) read(ctx context.Context, uid uuid.UUID) ([]repos.Data, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	raw, err := ss.connection.Query(ctx, `SELECT * FROM shops WHERE id = $1`, uid)
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	data := []repos.Data{}
	raw.Next()
	row := &repos.Shop{}
	err = raw.Scan(&row.Id, &row.Name, &row.Address, &row.IsClosed, &row.Owner)
	if err != nil {
		return data, fmt.Errorf("no shop with such uuid")
	}
	data = append(data, row)

	if len(data) < 1 {
		return nil, fmt.Errorf("no shop with such uuid")
	}
	return data, nil

}

func (ss *shops) delete(ctx context.Context, uid uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	q, err := ss.connection.Exec(ctx, `DELETE FROM shops WHERE id = $1`, uid)
	if err != nil {
		return err
	}
	if q.RowsAffected() < 1 {
		return fmt.Errorf("no shop with such uuid")
	}

	return nil
}
