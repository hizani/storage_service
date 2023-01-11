package server

import (
	"crud_service/app/repos"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) getCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surname := vars["surname"]

	c, err := s.customers.Read(r.Context(), surname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusFound)
	_ = json.NewEncoder(w).Encode(*c)
}

func (s *Server) getCustomerFieldByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surname := vars["surname"]
	field := vars["field"]

	c, err := s.customers.Read(r.Context(), surname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsoned, _ := json.Marshal(&c)
	var data map[string]interface{}
	_ = json.Unmarshal(jsoned, &data)

	elem, ok := data[field]
	if !ok {
		http.Error(w, "there is no such field", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusFound)
	msg := fmt.Sprintf(`{"%v":"%v"}`, field, elem)
	if elem, ok := elem.(float64); ok {
		msg = fmt.Sprintf(`{"%v":%v}`, field, elem)
	}
	w.Write([]byte(fmt.Sprintln(msg)))
}

func (s *Server) deleteCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surname := vars["surname"]

	c, err := s.customers.Delete(r.Context(), surname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return

	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*c)
}

func (s *Server) createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	c := &repos.Customer{}
	if err := json.NewDecoder(r.Body).Decode(c); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err := s.customers.Create(r.Context(), *c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
