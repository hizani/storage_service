package storage_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hizani/crud_service/storage"
	"google.golang.org/protobuf/encoding/protojson"
)

func (ss *StorageService) CreateShop(ctx context.Context, req *storage.ShopRequest) (*storage.ShopResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.ShopResponce{}
	if len(req.GetShops()) < 1 {
		resp.Success = false
		resp.Message = "no shop provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.Create(ctx, req.GetShops()[0])
	if err, ok := err.(*storage.RequiredMissingError); ok {
		resp.Success = false
		resp.Message = "required field is missing: " + err.Error()
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	if err != nil {
		return nil, fmt.Errorf("create shop error: %v", err)
	}
	resp.Success = true
	resp.Message = "created successfully"
	resp.Status = http.StatusCreated
	s, ok := data.(*storage.Shop)
	if !ok {
		resp.Message = "return cutomer error"
		return resp, nil
	}
	resp.Shops = append(resp.Shops, s)

	return resp, nil
}
func (ss *StorageService) ReadShop(ctx context.Context, req *storage.ShopRequest) (*storage.ShopResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.ShopResponce{}
	if len(req.GetShops()) < 1 {
		resp.Success = false
		resp.Message = "no shops provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.Read(ctx, req.GetShops()[0])
	if err != nil {
		return nil, fmt.Errorf("read shops error: %v", err)
	}
	if data == nil {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}

	s, ok := data.(*storage.Shop)
	if !ok {
		return nil, errors.New("can't parse shop")
	}

	resp.Success = true
	resp.Message = "read successfully"
	resp.Status = http.StatusOK
	resp.Shops = append(resp.Shops, s)
	return resp, nil
}
func (ss *StorageService) DeleteShop(ctx context.Context, req *storage.ShopRequest) (*storage.ShopResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.ShopResponce{}
	if len(req.GetShops()) < 1 {
		resp.Success = false
		resp.Message = "no shop provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.Read(ctx, req.GetShops()[0])
	if err != nil {
		return nil, fmt.Errorf("delete shops error: %v", err)
	}
	if data == nil {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}
	s, ok := data.(*storage.Shop)
	if !ok {
		return nil, errors.New("can't parse shop")
	}
	resp.Success = true
	resp.Message = "deleted successfully"
	resp.Status = http.StatusOK
	resp.Shops = append(resp.Shops, s)
	return resp, ss.storage.Delete(ctx, data)
}
func (ss *StorageService) ReadShopBySearchField(ctx context.Context, req *storage.ShopRequest) (*storage.ShopResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.ShopResponce{}
	if len(req.GetShops()) < 1 {
		resp.Success = false
		resp.Message = "no shop provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.ReadBySearchField(ctx, req.GetShops()[0])
	if err != nil {
		return nil, fmt.Errorf("read shop error: %v", err)
	}
	if data == nil {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}

	result := make([]*storage.Shop, 0, len(data))
	for _, sh := range data {
		s, ok := sh.(*storage.Shop)
		if !ok {
			return nil, errors.New("can't parse shop")
		}
		result = append(result, s)
	}
	resp.Success = true
	resp.Message = "reading successful"
	resp.Status = http.StatusOK
	resp.Shops = result
	return resp, nil
}
func (ss *StorageService) ReadShopFieldById(ctx context.Context, req *storage.ShopRequest) (*storage.ShopResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.ShopResponce{}
	if len(req.GetShops()) < 1 {
		resp.Success = false
		resp.Message = "no shop provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.Read(ctx, req.GetShops()[0])
	if err != nil {
		return nil, fmt.Errorf("read shop error: %v", err)
	}
	if data == nil {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}

	s, ok := data.(*storage.Shop)
	if !ok {
		return nil, errors.New("can't parse shop")
	}

	jsoned, _ := protojson.Marshal(s)
	var fields map[string]interface{}
	_ = json.Unmarshal(jsoned, &fields)
	elem, ok := fields[req.GetFieldName()]
	if !ok {
		resp.Success = false
		resp.Message = "there is no such field"
		resp.Status = http.StatusNotFound
		return resp, nil
	}
	f, err := interfaceToFieldMessage(req.GetFieldName(), elem)
	if err != nil {
		return nil, fmt.Errorf("interface to field convertion error: %v", err)
	}

	resp.Success = true
	resp.Message = "read successfully"
	resp.Status = http.StatusOK
	resp.Shops = append(resp.Shops, s)
	resp.Fields = append(resp.Fields, f)
	return resp, nil
}
