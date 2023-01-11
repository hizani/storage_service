package mem

import (
	"context"
	"crud_service/app/repos"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Runtime storage of shops
type shops struct {
	sync.RWMutex
	m map[uuid.UUID]repos.Shop
}

func newShops() *shops {
	return &shops{
		m: make(map[uuid.UUID]repos.Shop),
	}
}

func (ss *shops) create(ctx context.Context, s repos.Shop) error {
	ss.Lock()
	defer ss.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.Id = uuid.New()
	if !checkRequiredFields(&s) {
		return errors.New("required field is missing")
	}
	ss.m[s.Id] = s
	return nil
}

func (ss *shops) readName(ctx context.Context, name string) (*repos.Shop, error) {
	for _, elem := range ss.m {
		if elem.Name == name {
			return &elem, nil
		}
	}
	return nil, fmt.Errorf("no customer with such surname")
}

func (ss *shops) read(ctx context.Context, uid uuid.UUID) (*repos.Shop, error) {
	ss.RLock()
	defer ss.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, ok := ss.m[uid]

	if !ok {
		return nil, fmt.Errorf("no such shop")
	}
	return &data, nil
}
func (cs *shops) delete(ctx context.Context, uid uuid.UUID) error {
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
