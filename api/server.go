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
	r.Path("/mocks").HandlerFunc(GetMocksHandler).Methods(http.MethodGet)
	r.Path("/mocks").HandlerFunc(CreateMockHandler).Methods(http.MethodPost)
	r.Path("/mocks/{mock_id}/pause").HandlerFunc(PauseMockServerHandler).Methods(http.MethodPost)
	r.Path("/mocks/{mock_id}/resume").HandlerFunc(ResumeMockServerHandler).Methods(http.MethodPost)

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
