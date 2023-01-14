package server

import (
	"crud_service/app/repos"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) getShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "bad uuid")
		return
	}

	sh, err := s.shops.ReadId(r.Context(), uid)
	if err != nil {
		internalError(w, err)
		return
	}

	if sh == nil {
		notFound(w)
		return
	}

	cj, _ := json.Marshal(sh)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusFound)
	w.Write(j)
}

func (s *Server) getShopByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	sh, err := s.shops.ReadName(r.Context(), name)
	if err != nil {
		internalError(w, err)
		return
	}

	if sh == nil || len(sh) < 1 {
		notFound(w)
		return
	}

	cj, _ := json.Marshal(sh)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusFound)
	w.Write(j)
}

func (s *Server) getShopFieldByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	field := vars["field"]

	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "bad uuid")
		return
	}

	ss, err := s.shops.ReadId(r.Context(), uid)
	if err != nil {
		internalError(w, err)
		return
	}

	if ss == nil {
		notFound(w)
		return
	}

	jsoned, _ := json.Marshal(ss)
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

func (s *Server) deleteShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "bad uuid")
		return
	}

	ss, err := s.shops.Delete(r.Context(), uid)
	if err != nil {
		internalError(w, err)
		return
	}

	if ss == nil {
		notFound(w)
		return
	}

	cj, _ := json.Marshal(ss)
	j, _ := json.Marshal(payload{true, string(cj)})
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *Server) createShopHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sh := repos.Shop{}
	if err := json.NewDecoder(r.Body).Decode(&sh); err != nil {
		badRequest(w, "bad json")
		return
	}

	uid, err := s.shops.Create(r.Context(), sh)
	if err != nil {
		if err, ok := err.(*repos.RequiredMissingError); ok {
			badRequest(w, err.Error())
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
