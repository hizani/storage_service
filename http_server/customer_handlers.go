package http_server

import (
	"io"
	"net/http"

	"github.com/hizani/crud_service/storage"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/gorilla/mux"
)

func (s *HTTPServer) getCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	req := []*storage.Customer{{Id: id}}
	c, err := s.storageClient.ReadCustomer(r.Context(), &storage.CustomerRequest{Customers: req})
	if err != nil {
		internalError(w, err)
		return
	}

	cj, _ := protojson.Marshal(c)
	var resp int32 = c.GetStatus()
	w.WriteHeader(int(resp))
	w.Write(cj)
}

func (s *HTTPServer) getCustomerBySurnameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	surname := vars["surname"]

	req := []*storage.Customer{{Surname: surname}}
	resp, err := s.storageClient.ReadCustomerBySearchField(r.Context(), &storage.CustomerRequest{Customers: req})
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)
}

func (s *HTTPServer) getCustomerFieldByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	field := vars["field"]

	req := []*storage.Customer{{Id: id}}
	resp, err := s.storageClient.ReadCustomerFieldById(
		r.Context(),
		&storage.CustomerRequest{Customers: req, FieldName: &field},
	)
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)
}

func (s *HTTPServer) deleteCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	req := []*storage.Customer{{Id: id}}
	resp, err := s.storageClient.DeleteCustomer(
		r.Context(),
		&storage.CustomerRequest{Customers: req},
	)
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)
}

func (s *HTTPServer) createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	c := &storage.Customer{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w)
		return
	}
	protojson.Unmarshal(b, c)

	req := []*storage.Customer{c}
	resp, err := s.storageClient.CreateCustomer(r.Context(), &storage.CustomerRequest{Customers: req})
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)

}
