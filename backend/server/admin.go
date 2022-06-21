package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/backend/server/api"
)

type AdminServer struct {
}

func NewAdminServer() *AdminServer {
	return &AdminServer{}
}

func (a *AdminServer) Start(ctx context.Context, port string) (string, func(), error) {
	r := mux.NewRouter()
	r.PathPrefix("/mocks").HandlerFunc(api.GetMocksHandler).Methods(http.MethodGet)

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
