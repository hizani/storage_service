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

func (ss *StorageService) CreateCustomer(ctx context.Context, req *storage.CustomerRequest) (*storage.CustomerResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.CustomerResponce{}
	if len(req.GetCustomers()) < 1 {
		resp.Success = false
		resp.Message = "no cutomer provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}

	data, err := ss.storage.Create(ctx, req.GetCustomers()[0])
	if err, ok := err.(*storage.RequiredMissingError); ok {
		resp.Success = false
		resp.Message = "required field is missing: " + err.Error()
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	if err != nil {
		return nil, fmt.Errorf("create customer error: %v", err)
	}
	resp.Success = true
	resp.Message = "created successfully"
	resp.Status = http.StatusCreated
	c, ok := data.(*storage.Customer)
	if !ok {
		resp.Message = "return cutomer error"
		return resp, nil
	}
	resp.Customers = append(resp.Customers, c)

	return resp, nil
}
func (ss *StorageService) ReadCustomer(ctx context.Context, req *storage.CustomerRequest) (*storage.CustomerResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.CustomerResponce{}
	if len(req.GetCustomers()) < 1 {
		resp.Success = false
		resp.Message = "no cutomer provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.Read(ctx, req.GetCustomers()[0])
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}
	if data == nil {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}

	c, ok := data.(*storage.Customer)
	if !ok {
		return nil, errors.New("can't parse customer")
	}

	resp.Success = true
	resp.Message = "read successfully"
	resp.Status = http.StatusOK
	resp.Customers = append(resp.Customers, c)
	return resp, nil
}
func (ss *StorageService) DeleteCustomer(ctx context.Context, req *storage.CustomerRequest) (*storage.CustomerResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.CustomerResponce{}
	if len(req.GetCustomers()) < 1 {
		resp.Success = false
		resp.Message = "no cutomer provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.Read(ctx, req.GetCustomers()[0])
	if err != nil {
		return nil, fmt.Errorf("delete customer error: %v", err)
	}
	if data == nil {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}
	c, ok := data.(*storage.Customer)
	if !ok {
		return nil, errors.New("can't parse customer")
	}
	resp.Success = true
	resp.Message = "deleted successfully"
	resp.Status = http.StatusOK
	resp.Customers = append(resp.Customers, c)
	return resp, ss.storage.Delete(ctx, data)
}
func (ss *StorageService) ReadCustomerBySearchField(ctx context.Context, req *storage.CustomerRequest) (*storage.CustomerResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.CustomerResponce{}
	if len(req.GetCustomers()) < 1 {
		resp.Success = false
		resp.Message = "no cutomer provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.ReadBySearchField(ctx, req.GetCustomers()[0])
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}
	if len(data) < 1 {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}

	result := make([]*storage.Customer, 0, len(data))
	for _, cs := range data {
		c, ok := cs.(*storage.Customer)
		if !ok {
			return nil, errors.New("can't parse customer")
		}
		result = append(result, c)
	}
	resp.Success = true
	resp.Message = "reading successful"
	resp.Status = http.StatusOK
	resp.Customers = result
	return resp, nil
}
func (ss *StorageService) ReadCustomerFieldById(ctx context.Context, req *storage.CustomerRequest) (*storage.CustomerResponce, error) {
	ss.Wg.Add(1)
	defer ss.Wg.Done()
	resp := &storage.CustomerResponce{}
	if len(req.GetCustomers()) < 1 {
		resp.Success = false
		resp.Message = "no cutomer provided"
		resp.Status = http.StatusBadRequest
		return resp, nil
	}
	data, err := ss.storage.Read(ctx, req.GetCustomers()[0])
	if err != nil {
		return nil, fmt.Errorf("read customer error: %v", err)
	}
	if data == nil {
		resp.Success = false
		resp.Message = "not found"
		resp.Status = http.StatusOK
		return resp, nil
	}

	c, ok := data.(*storage.Customer)
	if !ok {
		return nil, errors.New("can't parse customer")
	}

	jsoned, _ := protojson.Marshal(c)
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
	resp.Customers = append(resp.Customers, c)
	resp.Fields = append(resp.Fields, f)
	return resp, nil
}
