package http_server

import (
	"io"
	"net/http"

	"github.com/hizani/crud_service/storage"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/gorilla/mux"
)

func (s *HTTPServer) getShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	req := []*storage.Shop{{Id: id}}
	resp, err := s.storageClient.ReadShop(r.Context(), &storage.ShopRequest{Shops: req})
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)
}

func (s *HTTPServer) getShopByNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	name := vars["name"]

	req := []*storage.Shop{{Name: name}}
	resp, err := s.storageClient.ReadShopBySearchField(r.Context(), &storage.ShopRequest{Shops: req})
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)
}

func (s *HTTPServer) getShopFieldByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	field := vars["field"]

	req := []*storage.Shop{{Id: id}}
	resp, err := s.storageClient.ReadShopFieldById(
		r.Context(),
		&storage.ShopRequest{Shops: req, FieldName: &field},
	)
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)
}

func (s *HTTPServer) deleteShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	req := []*storage.Shop{{Id: id}}
	resp, err := s.storageClient.DeleteShop(
		r.Context(),
		&storage.ShopRequest{Shops: req},
	)
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)
}

func (s *HTTPServer) createShopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	sh := &storage.Shop{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w)
		return
	}
	if err := protojson.Unmarshal(b, sh); err != nil {
		badRequest(w)
		return
	}

	req := []*storage.Shop{sh}
	resp, err := s.storageClient.CreateShop(r.Context(), &storage.ShopRequest{Shops: req})
	if err != nil {
		internalError(w, err)
		return
	}

	sj, _ := protojson.Marshal(resp)
	var status int32 = resp.GetStatus()
	respond(w, int(status), sj)

}
