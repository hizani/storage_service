package repos

import (
	"context"
)

type Data = interface{}

type Storage interface {
	Create(ctx context.Context, data Data) error
	Read(ctx context.Context, data Data) (Data, error)
	Delete(ctx context.Context, data Data) error
}
