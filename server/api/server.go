package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (a *Server) Start(ctx context.Context, port string) (string, func(), error) {
	r := mux.NewRouter()
	r.PathPrefix("/mocks").HandlerFunc(GetMocksHandler).Methods(http.MethodGet)
	r.PathPrefix("/mocks").HandlerFunc(CreateMockHandler).Methods(http.MethodPost)

	addr := "0.0.0.0:" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
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
