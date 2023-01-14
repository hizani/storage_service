package server

import (
	"crud_service/app/repos"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) getCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "bad uuid")
		return
	}

	c, err := s.customers.ReadId(r.Context(), uid)
	if err != nil {
		internalError(w, err)
		return
	}

	if c == nil {
		notFound(w)
		return
	}
	cj, _ := json.Marshal(c)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusFound)
	w.Write(j)
}

func (s *Server) getCustomerBySurnameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	surname := vars["surname"]

	c, err := s.customers.ReadSurname(r.Context(), surname)
	if err != nil {
		internalError(w, err)
		return
	}

	if c == nil || len(c) < 1 {
		notFound(w)
		return
	}

	cj, _ := json.Marshal(c)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusFound)
	w.Write(j)
}

func (s *Server) getCustomerFieldByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	field := vars["field"]

	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "bad uuid")
		return
	}

	c, err := s.customers.ReadId(r.Context(), uid)
	if err != nil {
		internalError(w, err)
		return
	}
	if c == nil {
		notFound(w)
		return
	}

	jsoned, _ := json.Marshal(c)
	var data map[string]interface{}
	_ = json.Unmarshal(jsoned, &data)

	elem, ok := data[field]
	if !ok {
		notFound(w)
		return
	}

	cj, _ := json.Marshal(elem)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusFound)
	w.Write(j)
}

func (s *Server) deleteCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "bad uuid")
		return
	}

	c, err := s.customers.Delete(r.Context(), uid)
	if err != nil {
		internalError(w, err)
		return
	}

	if c == nil {
		notFound(w)
		return
	}

	cj, _ := json.Marshal(c)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *Server) createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	c := &repos.Customer{}
	if err := json.NewDecoder(r.Body).Decode(c); err != nil {
		badRequest(w, "bad json")
		return
	}
	uid, err := s.customers.Create(r.Context(), *c)
	if err != nil {
		if err, ok := err.(*repos.RequiredMissingError); ok {
			badRequest(w, fmt.Sprint("required field is missing: ", err))
			return
		}
		internalError(w, err)
		return
	}

	cj, _ := json.Marshal(uid)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
