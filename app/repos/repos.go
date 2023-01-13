package repos

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Interface of storage data
type Data interface {
	CheckRequired() bool
	CmpSearchField(data string) bool
	GetId() uuid.UUID
	GetTypeName() string
	GetSearchField() string
	GetSearchFieldName() string
	SetDefaults()
	SetFromMap(m map[string]interface{}) (Data, error)
	DbData() DbData
}

// Interface of database allowed data
type DbData interface {
	// Sets fields of DbData by reading database row
	SetFieldsFromDbRow(ctx context.Context, row pgx.Rows) error
	// Inserts DbData via database connection
	Insert(ctx context.Context, connection *pgx.Conn) (pgx.Rows, error)
}

// Inteface of a storage
type Storage interface {
	// Create Data in the storage and return its UUID
	Create(ctx context.Context, data Data) (*uuid.UUID, error)
	// Read Data from the storage by data UUID
	Read(ctx context.Context, data Data) (Data, error)
	// Delete Data from the storage by data UUID
	Delete(ctx context.Context, data Data) error
	// Return []Data with occurrences of Data search field
	ReadBySearchField(ctx context.Context, data Data) ([]Data, error)
}

func checkRequired(d Data) bool {
	fields := reflect.ValueOf(d).Elem()
	for i := 0; i < fields.NumField(); i++ {
		tag := fields.Type().Field(i).Tag.Get("validate")
		if strings.Contains(tag, "required") && fields.Field(i).IsZero() {
			fmt.Println(fields.Field(i))
			return false
		}
	}
	return true
}
