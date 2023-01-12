package mem

import (
	"context"
	"crud_service/app/repos"
	"fmt"

	"github.com/google/uuid"
)

var _ repos.Storage = &MemStorage{}

// Runtime storage
type MemStorage struct {
	customers *customers
	shops     *shops
}

func New() *MemStorage {
	c := newCustomers()
	s := newShops()

	return &MemStorage{c, s}
}

func (ms *MemStorage) Create(ctx context.Context, data repos.Data) (*uuid.UUID, error) {
	if d, ok := data.(*repos.Customer); ok {
		return ms.customers.create(ctx, *d)
	}
	if d, ok := data.(*repos.Shop); ok {
		return ms.shops.create(ctx, *d)
	}

	return nil, fmt.Errorf("there is no storage for this type of data")
}
func (ms *MemStorage) Read(ctx context.Context, data repos.Data) ([]repos.Data, error) {
	if d, ok := data.(*repos.Customer); ok {
		if d.Id == uuid.Nil || d.Id.String() == "" {
			return ms.customers.readSurname(ctx, d.Surname)
		}
		val, err := ms.customers.read(ctx, d.Id)
		return []interface{}{val}, err
	}

	if d, ok := data.(*repos.Shop); ok {
		if d.Id == uuid.Nil || d.Id.String() == "" {
			return ms.shops.readName(ctx, d.Name)
		}
		val, err := ms.shops.read(ctx, d.Id)
		return []interface{}{val}, err
	}

	return nil, fmt.Errorf("there is no storage for this type of data")
}
func (ms *MemStorage) Delete(ctx context.Context, data repos.Data) error {
	if d, ok := data.(*repos.Customer); ok {
		return ms.customers.delete(ctx, d.Id)
	}

	if d, ok := data.(*repos.Shop); ok {
		return ms.shops.delete(ctx, d.Id)
	}

	return fmt.Errorf("there is no storage for this type of data")
}
