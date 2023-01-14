package repos

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ Data = &Shop{}
var _ DbData = &Shop{}

type Shop struct {
	Id       uuid.UUID `json:"id" validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Address  string    `json:"address" validate:"required"`
	IsClosed bool      `json:"is_closed"`
	Owner    *string   `json:"owner"`
}

func (s *Shop) CmpSearchField(data string) bool { return data == s.Name }
func (s *Shop) CheckRequired() bool             { return checkRequired(s) }
func (s *Shop) GetId() uuid.UUID                { return s.Id }
func (s *Shop) GetTypeName() string             { return "shop" }
func (s *Shop) GetSearchField() string          { return s.Name }
func (s *Shop) GetSearchFieldName() string      { return "name" }
func (s *Shop) DbData() DbData                  { return s }
func (s *Shop) SetDefaults()                    { s.Id = uuid.New() }
func (s *Shop) SetFromMap(m map[string]interface{}) (Data, error) {
	ids, ok := m["id"].(string)
	id, err := uuid.Parse(ids)
	if !ok || err != nil {
		return nil, errors.New("can't parse id from the map")
	}
	name, ok := m["name"].(string)
	if !ok {
		return nil, errors.New("can't parse name from the map")
	}
	addr, ok := m["address"].(string)
	if !ok {
		return nil, errors.New("can't parse address from the map")
	}
	isClosed, ok := m["is_closed"].(bool)
	if !ok {
		return nil, errors.New("can't parse is_closed from the map")
	}
	owner, _ := m["owner"].(*string)
	s = &Shop{id, name, addr, isClosed, owner}
	return s, nil
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

// storage wrapper
type Shops struct {
	storage Storage
}

func NewShops(storage Storage) *Shops {
	return &Shops{storage}
}

func (ss *Shops) Create(ctx context.Context, s Shop) (*uuid.UUID, error) {
	uid, err := ss.storage.Create(ctx, &s)
	if err != nil {
		return nil, fmt.Errorf("create shop error: %v", err)
	}
	return uid, nil
}

func (ss *Shops) ReadName(ctx context.Context, name string) ([]*Shop, error) {
	data, err := ss.storage.ReadBySearchField(ctx, &Shop{Name: name})
	if err != nil {
		return nil, fmt.Errorf("read shop error: %v", err)
	}
	result := []*Shop{}
	if len(data) < 1 {
		return result, nil
	}
	for _, foo := range data {
		c, ok := foo.(*Shop)
		if !ok {
			return nil, fmt.Errorf("read shop error: %v", err)
		}
		result = append(result, c)
	}

	return result, nil
}

func (ss *Shops) ReadId(ctx context.Context, uid uuid.UUID) (*Shop, error) {
	data, err := ss.storage.Read(ctx, &Shop{Id: uid})
	if err != nil {
		return nil, fmt.Errorf("read shop error: %v", err)
	}

	s, ok := data.(*Shop)
	if !ok || s == nil {
		return nil, nil
	}

	return s, nil
}

func (ss *Shops) Delete(ctx context.Context, uid uuid.UUID) (*Shop, error) {
	data, err := ss.storage.Read(ctx, &Shop{Id: uid})
	if err != nil {
		return nil, fmt.Errorf("read shop error: %v", err)
	}
	s, ok := data.(*Shop)
	if !ok || s == nil {
		return nil, nil
	}
	return s, ss.storage.Delete(ctx, s)
}
