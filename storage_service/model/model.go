package model

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Interface of storage data
type Data interface {
	CheckRequired() error
	CmpSearchField(data string) bool
	GetId() string
	GetTypeName() string
	GetSearchField() string
	GetSearchFieldName() string
	SetDefaults()
	SetFromMap(m map[string]interface{}) error
	SetFieldsFromDbRow(ctx context.Context, row pgx.Rows) error
	Insert(ctx context.Context, connection *pgx.Conn) (pgx.Rows, error)
}

// Inteface of a storage
type Storage interface {
	// Create Data in the storage and return its UUID
	Create(ctx context.Context, data Data) (Data, error)
	// Read Data from the storage by data UUID
	Read(ctx context.Context, data Data) (Data, error)
	// Delete Data from the storage by data UUID
	Delete(ctx context.Context, data Data) error
	// Return []Data with occurrences of Data search field
	ReadBySearchField(ctx context.Context, data Data) ([]Data, error)
}
