package server

import (
	"crud_service/app/repos"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) getShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "not valid uuid", http.StatusBadRequest)
	}

	sh, err := s.shops.ReadId(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusFound)
	_ = json.NewEncoder(w).Encode(sh)
}

func (s *Server) getShopByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	sh, err := s.shops.ReadName(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusFound)
	_ = json.NewEncoder(w).Encode(sh)
}

func (s *Server) getShopFieldByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	field := vars["field"]

	uid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "not valid uuid", http.StatusBadRequest)
	}

	ss, err := s.shops.ReadId(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsoned, _ := json.Marshal(ss)
	var data map[string]interface{}
	_ = json.Unmarshal(jsoned, &data)

	elem, ok := data[field]
	if !ok {
		http.Error(w, "there is no such field", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusFound)
	msg := interfaceToJson(field, elem)
	w.Write(msg)
}

func (s *Server) deleteShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "not valid uuid", http.StatusBadRequest)
	}

	c, err := s.shops.Delete(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return

	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*c)
}

func (s *Server) createShopHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sh := repos.Shop{}
	if err := json.NewDecoder(r.Body).Decode(&sh); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	uid, err := s.shops.Create(r.Context(), sh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	msg := fmt.Sprintf(`{"id":"%v"}`, uid.String())
	w.Write([]byte(msg))
}
