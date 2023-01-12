package mem

import (
	"context"
	"crud_service/app/repos"
	"crud_service/storage/common"
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

func (cs *customers) create(ctx context.Context, c repos.Customer) (*uuid.UUID, error) {
	cs.Lock()
	defer cs.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	c.Id = uuid.New()
	if c.RegDate.IsZero() {
		c.RegDate = time.Now()
	}

	if !common.CheckRequiredFields(&c) {
		return nil, errors.New("required field is missing")
	}
	cs.m[c.Id] = c
	return &c.Id, nil
}

func (cs *customers) readSurname(ctx context.Context, surname string) ([]repos.Data, error) {
	cs.RLock()
	defer cs.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data := []repos.Data{}
	for _, elem := range cs.m {
		if elem.Surname == surname {
			data = append(data, &elem)
		}
	}
	return data, nil
}

func (cs *customers) read(ctx context.Context, uid uuid.UUID) (repos.Data, error) {
	cs.RLock()
	defer cs.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, ok := cs.m[uid]

	if !ok {
		return nil, fmt.Errorf("no customer with such uuid")
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
