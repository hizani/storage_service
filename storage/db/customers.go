package db

import (
	"context"
	"crud_service/app/repos"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Database storage of customers
type customers struct {
	connection *pgx.Conn
}

func newCustomers(connection *pgx.Conn) *customers {
	return &customers{connection}
}

func (cs *customers) create(ctx context.Context, c repos.Customer) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if c.RegDate.IsZero() {
		c.RegDate = time.Now()
	}
	var err error
	var raw pgx.Rows
	if c.Id == uuid.Nil {
		raw, err = cs.connection.Query(ctx, `INSERT INTO customers (id, surname, name, patronymic, age, reg_date) 
		values (DEFAULT, $1, $2, $3, $4, $5) RETURNING id;`,
			c.Surname, c.Name, c.Patronymic, c.Age, c.RegDate)
	} else {
		raw, err = cs.connection.Query(ctx, `INSERT INTO customers (id, surname, name, patronymic, age, reg_date) 
			values ($1, $2, $3, $4, $5, $6) RETURNING id;`,
			c.Id, c.Surname, c.Name, c.Patronymic, c.Age, c.RegDate)
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

func (cs *customers) readSurname(ctx context.Context, surname string) ([]repos.Data, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	raw, err := cs.connection.Query(ctx, `SELECT * FROM customers WHERE surname = $1`, surname)
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	data := []repos.Data{}
	for raw.Next() {
		row := &repos.Customer{}
		err := raw.Scan(&row.Id, &row.Surname, &row.Name, &row.Patronymic, &row.Age, &row.RegDate)
		if err != nil {
			return data, err
		}
		data = append(data, row)
	}
	if len(data) < 1 {
		return nil, fmt.Errorf("no customer with such surname")
	}
	return data, nil
}

func (cs *customers) read(ctx context.Context, uid uuid.UUID) ([]repos.Data, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	raw, err := cs.connection.Query(ctx, `SELECT * FROM customers WHERE id = $1`, uid.String())
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	data := []repos.Data{}
	raw.Next()
	row := &repos.Customer{}
	err = raw.Scan(&row.Id, &row.Surname, &row.Name, &row.Patronymic, &row.Age, &row.RegDate)
	if err != nil {
		return data, fmt.Errorf("no customer with such uuid")
	}
	data = append(data, row)

	if len(data) < 1 {
		return nil, fmt.Errorf("no customer with such uuid")
	}
	return data, nil

}

func (cs *customers) delete(ctx context.Context, uid uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	q, err := cs.connection.Exec(ctx, `DELETE FROM customers WHERE id = $1`, uid)
	if err != nil {
		return fmt.Errorf("no customer with such uuid")
	}
	if q.RowsAffected() < 1 {
		return fmt.Errorf("no customer with such uuid")
	}

	return nil
}
