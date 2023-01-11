package server

import (
	"crud_service/app/repos"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) getShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	sh, err := s.shops.Read(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusFound)
	_ = json.NewEncoder(w).Encode(*sh)
}

func (s *Server) getShopFieldByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	field := vars["field"]

	c, err := s.shops.Read(r.Context(), name)
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

func (s *Server) deleteShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	c, err := s.shops.Delete(r.Context(), name)
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

	err := s.shops.Create(r.Context(), sh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
