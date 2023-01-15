package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Shop) CmpSearchField(data string) bool { return data == s.Name }
func (s *Shop) CheckRequired() error {
	if s.GetId() == "" {
		return &RequiredMissingError{"id"}
	}
	if s.GetName() == "" {
		return &RequiredMissingError{"name"}
	}
	if s.GetAddress() == "" {
		return &RequiredMissingError{"address"}
	}
	if s.IsClosed == nil {
		return &RequiredMissingError{"is_closed"}
	}
	return nil
}
func (s *Shop) GetTypeName() string        { return "shop" }
func (s *Shop) GetSearchField() string     { return s.Name }
func (s *Shop) GetSearchFieldName() string { return "name" }
func (s *Shop) SetDefaults()               { s.Id = uuid.New().String() }
func (s *Shop) SetFromMap(m map[string]interface{}) error {
	id, ok := m["id"].(string)
	if !ok {
		return errors.New("can't parse id from the map")
	}
	name, ok := m["name"].(string)
	if !ok {
		return errors.New("can't parse name from the map")
	}
	addr, ok := m["address"].(string)
	if !ok {
		return errors.New("can't parse address from the map")
	}
	isClosed, ok := m["is_closed"].(bool)
	if !ok {
		return errors.New("can't parse is_closed from the map")
	}
	owner, _ := m["owner"].(*string)
	*s = Shop{Id: id, Name: name, Address: addr, IsClosed: &isClosed, Owner: owner}
	return nil
}
func (s *Shop) Insert(ctx context.Context, connection *pgx.Conn) (pgx.Rows, error) {
	raw, err := connection.Query(ctx, `INSERT INTO shops (id, name, address, is_closed, owner) 
	values ($1, $2, $3, $4, $5) RETURNING id;`,
		s.Id, s.Name, s.Address, s.IsClosed, s.Owner)
	return raw, err
}

func (c *Shop) SetFieldsFromDbRow(ctx context.Context, row pgx.Rows) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	err := row.Scan(&c.Id, &c.Name, &c.Address, &c.IsClosed, &c.Owner)
	if err != nil {
		return err
	}
	return nil
}
