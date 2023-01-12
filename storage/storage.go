package storage

import (
	"crud_service/storage/db"
	"crud_service/storage/file"
	"crud_service/storage/mem"

	"github.com/jackc/pgx/v5"
)

func NewMemStorage() *mem.MemStorage {
	return mem.New()
}

func NewFileStorage(path string) *file.FileStorage {
	return file.New(path)
}

func NewDbStorage(conn *pgx.Conn) *db.DbStorage {
	return db.New(conn)
}
