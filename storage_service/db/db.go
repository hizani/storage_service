package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hizani/crud_service/storage_service/model"

	"github.com/jackc/pgx/v5"
)

var _ model.Storage = &DbStorage{}

// Database storage
type DbStorage struct {
	connection *pgx.Conn
}

func New(connection *pgx.Conn) *DbStorage {
	return &DbStorage{connection}
}

func (s *DbStorage) Create(ctx context.Context, d model.Data) (model.Data, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	d.SetDefaults()

	if err := d.CheckRequired(); err != nil {
		return nil, err
	}
	raw, err := d.Insert(ctx, s.connection)

	if err != nil {
		return nil, err
	}
	defer raw.Close()

	return d, nil
}
func (s *DbStorage) Read(ctx context.Context, d model.Data) (model.Data, error) {
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
	err = d.SetFieldsFromDbRow(ctx, row)
	if err != nil {
		return nil, err
	}
	return d, nil
}
func (s *DbStorage) Delete(ctx context.Context, d model.Data) error {
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

func (s *DbStorage) ReadBySearchField(ctx context.Context, d model.Data) ([]model.Data, error) {
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
	data := []model.Data{}
	for rows.Next() {
		// Copy d variable
		newData, ok := reflect.New(reflect.ValueOf(d).Elem().Type()).Interface().(model.Data)
		if !ok {
			return nil, errors.New("can't copy Data")
		}
		err = newData.SetFieldsFromDbRow(ctx, rows)
		if err != nil {
			return nil, err
		}
		data = append(data, newData)
	}
	return data, nil
}
