package storage

import (
	"crud_service/storage/db"
	"crud_service/storage/file"
	"crud_service/storage/mem"
)

func NewMemStorage() *mem.MemStorage {
	return mem.New()
}

func NewFileStorage(path string) *file.FileStorage {
	return file.New(path)
}

func NewDbStorage(connStr string) *db.DbStorage {
	return db.New(connStr)
}
