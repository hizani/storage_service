package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var _ Data = &Customer{}

type Customer struct {
	Id         uuid.UUID `json:"id" validate:"required"`
	Surname    string    `json:"surname" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Patronymic string    `json:"patronymic" validate:"required"`
	Age        uint      `json:"age"`
	RegDate    time.Time `json:"reg_date" validate:"required"`
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
	data, err := cs.storage.Read(ctx, &Customer{Surname: surname})
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	customers := []*Customer{}
	for _, elem := range data {
		u, ok := elem.(*Customer)
		if !ok {
			return nil, fmt.Errorf("read customer error: %v", err)
		}
		customers = append(customers, u)
	}

	return customers, nil
}

func (cs *Customers) ReadId(ctx context.Context, uid uuid.UUID) ([]*Customer, error) {
	data, err := cs.storage.Read(ctx, &Customer{Id: uid})
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	customers := []*Customer{}
	for _, elem := range data {
		u, ok := elem.(*Customer)
		if !ok {
			return nil, fmt.Errorf("read customer error: %v", err)
		}
		customers = append(customers, u)
	}

	return customers, nil
}

func (cs *Customers) Delete(ctx context.Context, uid uuid.UUID) (*Customer, error) {
	c, err := cs.storage.Read(ctx, &Customer{Id: uid})
	if err != nil {
		return nil, fmt.Errorf("can not find customer: %v", err)
	}
	return c[0].(*Customer), cs.storage.Delete(ctx, &c[0])
}
