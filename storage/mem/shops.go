package mem

import (
	"context"
	"crud_service/app/repos"
	"crud_service/storage/common"
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

func (ss *shops) create(ctx context.Context, s repos.Shop) (*uuid.UUID, error) {
	ss.Lock()
	defer ss.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.Id = uuid.New()
	if !common.CheckRequiredFields(&s) {
		return nil, errors.New("required field is missing")
	}
	ss.m[s.Id] = s
	return &s.Id, nil
}

func (ss *shops) readName(ctx context.Context, name string) ([]repos.Data, error) {
	ss.RLock()
	defer ss.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data := []repos.Data{}
	for _, elem := range ss.m {
		if elem.Name == name {
			data = append(data, &elem)
		}
	}
	return data, nil
}

func (ss *shops) read(ctx context.Context, uid uuid.UUID) (repos.Data, error) {
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
