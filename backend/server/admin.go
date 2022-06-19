package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AdminServer struct{}

func NewAdminServer() *AdminServer {
	return &AdminServer{}
}

func (*AdminServer) Start(ctx context.Context, port int32) (string, func(), error) {
	r := mux.NewRouter()
	r.PathPrefix("/mocks").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write([]byte("get mocks"))
	}).Methods(http.MethodGet)

	r.PathPrefix("/mocks/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write([]byte("put mock"))
	}).Methods(http.MethodPut)

	r.PathPrefix("/mocks").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write([]byte("create new mock"))
	}).Methods(http.MethodPost)

	addr := "0.0.0.0:2601"
	srv := &http.Server{
		Addr:    addr,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return addr, func() {
		_ = srv.Shutdown(context.Background())
	}, nil
}
