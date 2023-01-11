package repos

import (
	"context"
	"fmt"
)

var _ Data = &Shop{}

type Shop struct {
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address" validate:"required"`
	IsClosed bool   `json:"is_closed" validate:"required"`
	Owner    string `json:"owner"`
}

// storage wrapper
type Shops struct {
	storage Storage
}

func NewShops(storage Storage) *Shops {
	return &Shops{storage}
}

func (ss *Shops) Create(ctx context.Context, s Shop) error {
	err := ss.storage.Create(ctx, &s)
	if err != nil {
		return fmt.Errorf("create user error: %v", err)
	}
	return nil
}

func (ss *Shops) Read(ctx context.Context, name string) (*Shop, error) {
	data, err := ss.storage.Read(ctx, &Shop{Name: name})
	if err != nil {
		return nil, fmt.Errorf("read user error: %v", err)
	}

	s, ok := data.(*Shop)
	if !ok {
		return nil, fmt.Errorf("read user error: %v", err)
	}

	return s, nil
}

func (ss *Shops) Delete(ctx context.Context, name string) (*Shop, error) {
	s := Shop{Name: name}
	u, err := ss.storage.Read(ctx, &s)
	if err != nil {
		return nil, fmt.Errorf("read user error: %v", err)
	}
	return u.(*Shop), ss.storage.Delete(ctx, &s)
}
