package file

import (
	"context"
	"crud_service/app/repos"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"

	"github.com/google/uuid"
)

var _ repos.Storage = &FileStorage{}

// File storage
type FileStorage struct {
	m       sync.Mutex // lock for map of mutexes
	mtxs    map[string]*sync.RWMutex
	dirPath string
}

func New(path string) (*FileStorage, error) {
	if path == "" {
		path, _ = os.Getwd()
	}
	folderInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, errors.New("folder does not exist")
	}
	if !folderInfo.IsDir() {
		return nil, errors.New("is not a directory")
	}
	if path[len(path)-1:] != string(os.PathSeparator) {
		path = fmt.Sprintf("%s%s", path, string(os.PathSeparator))
	}
	return &FileStorage{mtxs: make(map[string]*sync.RWMutex), dirPath: path}, nil
}

func (s *FileStorage) Create(ctx context.Context, data repos.Data) (*uuid.UUID, error) {
	s.m.Lock()
	if s.mtxs[data.GetTypeName()] == nil {
		s.mtxs[data.GetTypeName()] = &sync.RWMutex{}
	}
	s.m.Unlock()
	s.mtxs[data.GetTypeName()].Lock()
	defer s.mtxs[data.GetTypeName()].Unlock()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	data.SetDefaults()

	if err := data.CheckRequired(); err != nil {
		return nil, err
	}

	err := s.create(data)
	if err != nil {
		return nil, err
	}

	uid := data.GetId()

	return &uid, nil
}

func (s *FileStorage) Read(ctx context.Context, data repos.Data) (repos.Data, error) {
	s.m.Lock()
	if s.mtxs[data.GetTypeName()] == nil {
		s.mtxs[data.GetTypeName()] = &sync.RWMutex{}
	}
	s.m.Unlock()
	s.mtxs[data.GetTypeName()].RLock()
	defer s.mtxs[data.GetTypeName()].RUnlock()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.read(data)

}
func (s *FileStorage) Delete(ctx context.Context, data repos.Data) error {
	s.m.Lock()
	if s.mtxs[data.GetTypeName()] == nil {
		s.mtxs[data.GetTypeName()] = &sync.RWMutex{}
	}
	s.m.Unlock()
	s.mtxs[data.GetTypeName()].Lock()
	defer s.mtxs[data.GetTypeName()].Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return s.delete(data)
}
func (s *FileStorage) ReadBySearchField(ctx context.Context, data repos.Data) ([]repos.Data, error) {
	s.m.Lock()
	if s.mtxs[data.GetTypeName()] == nil {
		s.mtxs[data.GetTypeName()] = &sync.RWMutex{}
	}
	s.m.Unlock()
	s.mtxs[data.GetTypeName()].RLock()
	defer s.mtxs[data.GetTypeName()].RUnlock()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.readSearchField(data)
}

func (s *FileStorage) delete(data repos.Data) error {
	res, err := s.readSlice(data)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s%s%s.json", s.dirPath, string(os.PathSeparator), data.GetTypeName())
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	ok := false
	for idx, elem := range res {
		if elem.GetId() == data.GetId() {
			res = append(res[:idx], res[idx+1:]...)
			ok = true
			break
		}
	}
	if !ok {
		return nil
	}

	byteData, err := json.Marshal(res)
	if err != nil {
		return err
	}
	stringData := string(byteData)

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if _, err := file.WriteString(stringData); err != nil {
		return err
	}

	return file.Sync()

}

func (s *FileStorage) readSlice(data repos.Data) ([]repos.Data, error) {
	filename := fmt.Sprintf("%s%s%s.json", s.dirPath, string(os.PathSeparator), data.GetTypeName())
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonByte, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var ires = []map[string]interface{}{}
	if err := json.Unmarshal(jsonByte, &ires); err != nil {
		return nil, err
	}
	var res = []repos.Data{}
	for _, elem := range ires {
		newData, ok := reflect.New(reflect.ValueOf(data).Elem().Type()).Interface().(repos.Data)
		if !ok {
			return nil, errors.New("can't copy Data")
		}
		newData, err = newData.SetFromMap(elem)
		if err != nil {
			return nil, err
		}
		res = append(res, newData)
	}
	return res, nil
}

func (s *FileStorage) read(data repos.Data) (repos.Data, error) {
	res, err := s.readSlice(data)
	if err != nil {
		return nil, err
	}

	for _, elem := range res {
		if elem.GetId() == data.GetId() {
			return elem, nil
		}
	}

	return nil, nil

}

func (s *FileStorage) readSearchField(data repos.Data) ([]repos.Data, error) {
	res, err := s.readSlice(data)
	if err != nil {
		return nil, err
	}
	dataSlice := []repos.Data{}
	for _, elem := range res {
		if elem.CmpSearchField(data.GetSearchField()) {
			dataSlice = append(dataSlice, elem)
		}
	}

	return dataSlice, nil

}

func (s *FileStorage) create(data repos.Data) error {
	filename := fmt.Sprintf("%s%s%s.json", s.dirPath, string(os.PathSeparator), data.GetTypeName())
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonByte, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	stringData := string(byteData)
	if len(jsonByte) > 0 {
		stringData = fmt.Sprintf("%s,%s", jsonByte[1:len(jsonByte)-1], stringData)
	}
	stringData = fmt.Sprintf(`[%s]`, stringData)

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if _, err := file.WriteString(stringData); err != nil {
		return err
	}

	return file.Sync()
}
