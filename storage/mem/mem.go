package mem

import (
	"context"
	"crud_service/app/repos"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var _ repos.Storage = &MemStorage{}

// Runtime storage
type MemStorage struct {
	sync.RWMutex
	m map[string]map[uuid.UUID]repos.Data
}

func New() *MemStorage {
	return &MemStorage{
		m: make(map[string]map[uuid.UUID]repos.Data),
	}

}

func (s *MemStorage) Create(ctx context.Context, d repos.Data) (*uuid.UUID, error) {
	s.Lock()
	defer s.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	d.SetDefaults()
	if !d.CheckRequired() {
		return nil, errors.New("required field is missing")
	}
	id := d.GetId()
	if s.m[d.GetTypeName()] == nil {
		s.m[d.GetTypeName()] = make(map[uuid.UUID]repos.Data)
	}
	s.m[d.GetTypeName()][id] = d
	return &id, nil
}
func (s *MemStorage) Read(ctx context.Context, d repos.Data) (repos.Data, error) {
	s.RLock()
	defer s.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if s.m[d.GetTypeName()] == nil {
		return nil, nil
	}

	data, ok := s.m[d.GetTypeName()][d.GetId()]

	if !ok {
		return nil, nil
	}
	return data, nil
}
func (s *MemStorage) Delete(ctx context.Context, d repos.Data) error {
	s.Lock()
	defer s.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if s.m[d.GetTypeName()] == nil {
		return nil
	}

	delete(s.m[d.GetTypeName()], d.GetId())

	return nil
}

func (s *MemStorage) ReadBySearchField(ctx context.Context, d repos.Data) ([]repos.Data, error) {
	s.RLock()
	defer s.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data := []repos.Data{}
	for _, elem := range s.m[d.GetTypeName()] {
		if elem.CmpSearchField(d.GetSearchField()) {
			data = append(data, elem)
		}
	}
	return data, nil
}
