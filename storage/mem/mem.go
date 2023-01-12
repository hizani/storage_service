package mem

import (
	"context"
	"crud_service/app/repos"
	"fmt"
	"reflect"
	"strings"

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
		if d.Id.String() == "00000000-0000-0000-0000-000000000000" || d.Id.String() == "" {
			return ms.customers.readSurname(ctx, d.Surname)
		}
		return ms.customers.read(ctx, d.Id)
	}

	if d, ok := data.(*repos.Shop); ok {
		if d.Id.String() == "00000000-0000-0000-0000-000000000000" || d.Id.String() == "" {
			return ms.shops.readName(ctx, d.Name)
		}
		return ms.shops.read(ctx, d.Id)
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

func checkRequiredFields(d repos.Data) bool {
	fields := reflect.ValueOf(d).Elem()
	for i := 0; i < fields.NumField(); i++ {
		tag := fields.Type().Field(i).Tag.Get("validate")
		if strings.Contains(tag, "required") && fields.Field(i).IsZero() {
			return false
		}
	}
	return true
}
