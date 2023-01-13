package repos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ Data = &Customer{}
var _ DbData = &Customer{}

type Customer struct {
	Id         uuid.UUID `json:"id" validate:"required"`
	Surname    string    `json:"surname" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Patronymic string    `json:"patronymic" validate:"required"`
	Age        *uint     `json:"age"`
	RegDate    time.Time `json:"reg_date" validate:"required"`
}

func (c *Customer) CmpSearchField(data string) bool { return data == c.Surname }
func (c *Customer) CheckRequired() bool             { return checkRequired(c) }
func (c *Customer) GetId() uuid.UUID                { return c.Id }
func (c *Customer) GetTypeName() string             { return "customer" }
func (c *Customer) GetSearchField() string          { return c.Surname }
func (c *Customer) GetSearchFieldName() string      { return "surname" }
func (c *Customer) DbData() DbData                  { return c }
func (c *Customer) SetDefaults() {
	if c.Id == uuid.Nil || c.Id.String() == "" {
		c.Id = uuid.New()
	}
	if c.RegDate.IsZero() {
		c.RegDate = time.Now()
	}
}
func (c *Customer) SetFromMap(m map[string]interface{}) (Data, error) {
	ids, ok := m["id"].(string)
	id, err := uuid.Parse(ids)
	if !ok || err != nil {
		return nil, errors.New("can't parse id from the map")
	}
	surname, ok := m["surname"].(string)
	if !ok {
		return nil, errors.New("can't parse surname from the map")
	}
	name, ok := m["name"].(string)
	if !ok {
		return nil, errors.New("can't parse name from the map")
	}
	patr, ok := m["patronymic"].(string)
	if !ok {
		return nil, errors.New("can't parse patronymic from the map")
	}
	age, _ := m["age"].(*uint)
	st, ok := m["reg_date"].(string)
	if !ok {
		return nil, errors.New("can't parse time from the map")
	}
	t, err := time.Parse(time.RFC3339, st)
	if err != nil {
		return nil, errors.New("can't parse time")
	}
	c = &Customer{id, surname, name, patr, age, t}
	return c, nil
}

func (c *Customer) Insert(ctx context.Context, connection *pgx.Conn) (pgx.Rows, error) {
	raw, err := connection.Query(ctx, `INSERT INTO customers (id, surname, name, patronymic, age, reg_date) 
	values ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		c.Id, c.Surname, c.Name, c.Patronymic, c.Age, c.RegDate)
	return raw, err
}
func (c *Customer) SetFieldsFromDbRow(ctx context.Context, row pgx.Rows) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	err := row.Scan(&c.Id, &c.Surname, &c.Name, &c.Patronymic, &c.Age, &c.RegDate)
	if err != nil {
		return err
	}

	return nil
}

// storage wrapper
type Customers struct {
	storage Storage
}

func NewCustomers(storage Storage) *Customers {
	return &Customers{storage}
}

func (cs *Customers) Create(ctx context.Context, c Customer) (*uuid.UUID, error) {
	uid, err := cs.storage.Create(ctx, &c)
	if err != nil {
		return nil, fmt.Errorf("create customer error: %v", err)
	}
	return uid, nil
}

func (cs *Customers) ReadSurname(ctx context.Context, surname string) ([]*Customer, error) {
	data, err := cs.storage.ReadBySearchField(ctx, &Customer{Surname: surname})
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	var result []*Customer
	for _, foo := range data {
		c, ok := foo.(*Customer)
		if !ok {
			return nil, fmt.Errorf("read customer error: %v", err)
		}
		result = append(result, c)
	}

	return result, nil
}

func (cs *Customers) ReadId(ctx context.Context, uid uuid.UUID) (*Customer, error) {
	data, err := cs.storage.Read(ctx, &Customer{Id: uid})
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	return data.(*Customer), nil
}

func (cs *Customers) Delete(ctx context.Context, uid uuid.UUID) (*Customer, error) {
	data, err := cs.storage.Read(ctx, &Customer{Id: uid})
	if err != nil {
		return nil, fmt.Errorf("can not find customer: %v", err)
	}
	return data.(*Customer), cs.storage.Delete(ctx, data)
}
