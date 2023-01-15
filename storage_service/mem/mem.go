package mem

import (
	"context"
	"sync"

	"github.com/hizani/crud_service/storage_service/model"
)

var _ model.Storage = &MemStorage{}

// Runtime storage
type MemStorage struct {
	sync.RWMutex
	m map[string]map[string]model.Data
}

func New() *MemStorage {
	return &MemStorage{
		m: make(map[string]map[string]model.Data),
	}

}

func (s *MemStorage) Create(ctx context.Context, d model.Data) (model.Data, error) {
	s.Lock()
	defer s.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	d.SetDefaults()
	if err := d.CheckRequired(); err != nil {
		return nil, err
	}
	id := d.GetId()
	if s.m[d.GetTypeName()] == nil {
		s.m[d.GetTypeName()] = make(map[string]model.Data)
	}
	s.m[d.GetTypeName()][id] = d
	return d, nil
}
func (s *MemStorage) Read(ctx context.Context, d model.Data) (model.Data, error) {
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
func (s *MemStorage) Delete(ctx context.Context, d model.Data) error {
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

func (s *MemStorage) ReadBySearchField(ctx context.Context, d model.Data) ([]model.Data, error) {
	s.RLock()
	defer s.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data := []model.Data{}
	for _, elem := range s.m[d.GetTypeName()] {
		if elem.CmpSearchField(d.GetSearchField()) {
			data = append(data, elem)
		}
	}
	return data, nil
}
