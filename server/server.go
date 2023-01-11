package server

import (
	"context"
	"crud_service/app/repos"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	srv       http.Server
	customers *repos.Customers
	shops     *repos.Shops
}

func New(address string) *Server {
	var s *Server = &Server{srv: http.Server{Addr: address}}
	router := mux.NewRouter()
	router.HandleFunc("/customers/search/{surname}", s.getCustomerBySurnameHandler).Methods("GET")
	router.HandleFunc("/customers/{id}", s.getCustomerByIdHandler).Methods("GET")
	router.HandleFunc("/customers/{id}/{field}", s.getCustomerFieldByIdHandler).Methods("GET")
	router.HandleFunc("/customers/delete/{id}", s.deleteCustomerByIdHandler).Methods("DELETE")
	router.HandleFunc("/customers/create", s.createCustomerHandler).Methods("POST")

	router.HandleFunc("/shops/search/{name}", s.getShopByNameHandler).Methods("GET")
	router.HandleFunc("/shops/{id}", s.getShopByIdHandler).Methods("GET")
	router.HandleFunc("/shops/{id}/{field}", s.getShopFieldByIdHandler).Methods("GET")
	router.HandleFunc("/shops/delete/{id}", s.deleteShopByIdHandler).Methods("DELETE")
	router.HandleFunc("/shops/create", s.createShopHandler).Methods("POST")
	s.srv.Handler = router
	return s
}

func (s *Server) Start(cs *repos.Customers, ss *repos.Shops) {
	s.customers = cs
	s.shops = ss
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil {
			log.Printf("serve error %v:", err)
		}
	}()
}

// Stop метод для остановки сервера, для этого у http сервера есть Shutdown(), который принимает контекст.
// Этот контекст сделаем с таймаутом
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_ = s.srv.Shutdown(ctx)
	cancel()
}
