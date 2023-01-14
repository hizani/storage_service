package db

import (
	"context"
	"crud_service/app/repos"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ repos.Storage = &DbStorage{}

// Database storage
type DbStorage struct {
	connection *pgx.Conn
}

func New(connection *pgx.Conn) *DbStorage {
	return &DbStorage{connection}
}

func (s *DbStorage) Create(ctx context.Context, d repos.Data) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	d.SetDefaults()

	if err := d.CheckRequired(); err != nil {
		return nil, err
	}
	raw, err := d.DbData().Insert(ctx, s.connection)

	if err != nil {
		return nil, err
	}
	defer raw.Close()

	uid := &uuid.UUID{}

	raw.Next()
	err = raw.Scan(uid)
	if err != nil {
		return nil, err
	}
	return uid, nil
}
func (s *DbStorage) Read(ctx context.Context, d repos.Data) (repos.Data, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	query := fmt.Sprintf(`SELECT * FROM %ss WHERE id = $1`, d.GetTypeName())

	row, err := s.connection.Query(ctx, query, d.GetId())
	if err != nil {
		return nil, err
	}
	defer row.Close()

	ok := row.Next()
	if !ok {
		return nil, nil
	}
	err = d.DbData().SetFieldsFromDbRow(ctx, row)
	if err != nil {
		return nil, err
	}
	return d, nil
}
func (s *DbStorage) Delete(ctx context.Context, d repos.Data) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	query := fmt.Sprintf(`DELETE FROM %ss WHERE id = $1`, d.GetTypeName())

	_, err := s.connection.Exec(ctx, query, d.GetId())
	if err != nil {
		return err
	}

	return nil
}

func (s *DbStorage) ReadBySearchField(ctx context.Context, d repos.Data) ([]repos.Data, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	query := fmt.Sprintf(`SELECT * FROM %ss WHERE %s = $1`, d.GetTypeName(), d.GetSearchFieldName())

	rows, err := s.connection.Query(ctx, query, d.GetSearchField())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := []repos.Data{}
	for rows.Next() {
		// Copy d variable
		newData, ok := reflect.New(reflect.ValueOf(d).Elem().Type()).Interface().(repos.Data)
		if !ok {
			return nil, errors.New("can't copy Data")
		}
		err = newData.DbData().SetFieldsFromDbRow(ctx, rows)
		if err != nil {
			return nil, err
		}
		data = append(data, newData)
	}
	return data, nil
}
