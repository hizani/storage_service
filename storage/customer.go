package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hizani/crud_service/storage_service/model"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ model.Data = &Customer{}

func (c *Customer) CmpSearchField(data string) bool { return data == c.Surname }
func (c *Customer) GetTypeName() string             { return "customer" }
func (c *Customer) GetSearchField() string          { return c.Surname }
func (c *Customer) GetSearchFieldName() string      { return "surname" }
func (c *Customer) CheckRequired() error {
	if c.GetId() == "" {
		return &RequiredMissingError{"id"}
	}
	if c.GetName() == "" {
		return &RequiredMissingError{"name"}
	}
	if c.GetSurname() == "" {
		return &RequiredMissingError{"surname"}
	}
	if c.GetPatronymic() == "" {
		return &RequiredMissingError{"patronymic"}
	}
	if c.GetRegDate().AsTime().IsZero() {
		return &RequiredMissingError{"reg_date"}
	}
	return nil
}
func (c *Customer) SetDefaults() {
	c.Id = uuid.New().String()
	c.RegDate = timestamppb.Now()
}
func (c *Customer) SetFromMap(m map[string]interface{}) error {
	id, ok := m["id"].(string)
	if !ok {
		return errors.New("can't parse id from the map")
	}
	surname, ok := m["surname"].(string)
	if !ok {
		return errors.New("can't parse surname from the map")
	}
	name, ok := m["name"].(string)
	if !ok {
		return errors.New("can't parse name from the map")
	}
	patr, ok := m["patronymic"].(string)
	if !ok {
		return errors.New("can't parse patronymic from the map")
	}
	age, _ := m["age"].(*uint32)

	st, ok := m["reg_date"].(map[string]interface{})
	if !ok {
		return errors.New("can't parse time from the map")
	}
	nanos := st["nanos"].(float64)
	seconds := st["seconds"].(float64)
	t := &timestamppb.Timestamp{Seconds: int64(seconds), Nanos: int32(nanos)}
	*c = Customer{Id: id, Surname: surname, Name: name, Patronymic: patr, Age: age, RegDate: t}
	return nil
}

func (c *Customer) Insert(ctx context.Context, connection *pgx.Conn) (pgx.Rows, error) {
	raw, err := connection.Query(ctx, `INSERT INTO customers (id, surname, name, patronymic, age, reg_date) 
	values ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		c.Id, c.Surname, c.Name, c.Patronymic, c.Age, c.RegDate.AsTime())
	return raw, err
}
func (c *Customer) SetFieldsFromDbRow(ctx context.Context, row pgx.Rows) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	var time time.Time
	err := row.Scan(&c.Id, &c.Surname, &c.Name, &c.Patronymic, &c.Age, &time)
	c.RegDate = timestamppb.New(time)
	if err != nil {
		return err
	}

	return nil
}
