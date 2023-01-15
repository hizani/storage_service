package http_server

import (
	"log"
	"net/http"

	"github.com/hizani/crud_service/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	srv           http.Server
	storageClient storage.StorageServiceClient
}

func New(serverAddress string, storageAddress string) *HTTPServer {
	conn, err := grpc.Dial(storageAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	c := storage.NewStorageServiceClient(conn)
	var s *HTTPServer = &HTTPServer{srv: http.Server{Addr: serverAddress}, storageClient: c}
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

func (s *HTTPServer) Start() {
	err := s.srv.ListenAndServe()
	if err != nil {
		log.Fatalf("serve error %v:", err)
	}

}

func internalError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func respond(w http.ResponseWriter, status int, msg []byte) {
	w.WriteHeader(int(status))
	w.Write(msg)
}

func badRequest(w http.ResponseWriter) {
	http.Error(w, "bad request", http.StatusBadRequest)
}
