package mem

import (
	"context"
	"crud_service/app/repos"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Runtime storage of customers
type customers struct {
	sync.RWMutex
	m map[uuid.UUID]repos.Customer
}

func newCustomers() *customers {
	return &customers{
		m: make(map[uuid.UUID]repos.Customer),
	}
}

func (cs *customers) create(ctx context.Context, c repos.Customer) error {
	cs.Lock()
	defer cs.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	c.Id = uuid.New()
	if c.RegDate.IsZero() {
		c.RegDate = time.Now()
	}

	if !checkRequiredFields(&c) {
		return errors.New("required field is missing")
	}
	cs.m[c.Id] = c
	return nil
}

func (cs *customers) readSurname(ctx context.Context, surname string) (*repos.Customer, error) {
	for _, elem := range cs.m {
		if elem.Surname == surname {
			return &elem, nil
		}
	}
	return nil, fmt.Errorf("no customer with such surname")
}

func (cs *customers) read(ctx context.Context, uid uuid.UUID) (*repos.Customer, error) {
	cs.RLock()
	defer cs.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, ok := cs.m[uid]

	if !ok {
		return nil, fmt.Errorf("no customer with such uuid: %v", uid)
	}
	return &data, nil

}

func (cs *customers) delete(ctx context.Context, uid uuid.UUID) error {
	cs.Lock()
	defer cs.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	delete(cs.m, uid)
	return nil
}
