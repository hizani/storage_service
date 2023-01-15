package storage_service

import (
	"log"
	"net"
	"sync"

	"github.com/hizani/crud_service/storage"
	"github.com/hizani/crud_service/storage_service/model"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type StorageService struct {
	storage.UnimplementedStorageServiceServer
	Wg      *sync.WaitGroup
	storage model.Storage
}

func New(st model.Storage) *StorageService {
	return &StorageService{storage: st, Wg: &sync.WaitGroup{}}
}

func (s *StorageService) Start(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	storage.RegisterStorageServiceServer(srv, s)
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Printf("server crushed %v:", err)
		}
	}()
}

func interfaceToFieldMessage(name string, value interface{}) (*storage.Field, error) {
	v, err := structpb.NewValue(value)
	if err != nil {
		return nil, err
	}
	return &storage.Field{Name: name, Value: v}, nil
}
