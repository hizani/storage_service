package repos

import (
	"context"

	"github.com/google/uuid"
)

type Data = interface{}

type Storage interface {
	Create(ctx context.Context, data Data) (*uuid.UUID, error)
	Read(ctx context.Context, data Data) ([]Data, error)
	Delete(ctx context.Context, data Data) error
}
