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

func (cs *Customers) Create(ctx context.Context, c Customer) error {
	err := cs.storage.Create(ctx, &c)
	if err != nil {
		return fmt.Errorf("create customer error: %v", err)
	}
	return nil
}

func (cs *Customers) ReadSurname(ctx context.Context, surname string) (*Customer, error) {
	uinterface, err := cs.storage.Read(ctx, &Customer{Surname: surname})
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	u, ok := uinterface.(*Customer)
	if !ok {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	return u, nil
}

func (cs *Customers) ReadId(ctx context.Context, uid uuid.UUID) (*Customer, error) {
	uinterface, err := cs.storage.Read(ctx, &Customer{Id: uid})
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	u, ok := uinterface.(*Customer)
	if !ok {
		return nil, fmt.Errorf("read customer error: %v", err)
	}

	return u, nil
}

func (cs *Customers) Delete(ctx context.Context, surname string) (*Customer, error) {
	c := Customer{Surname: surname}
	u, err := cs.storage.Read(ctx, &c)
	if err != nil {
		return nil, fmt.Errorf("can not find customer: %v", err)
	}
	return u.(*Customer), cs.storage.Delete(ctx, &c)
}
