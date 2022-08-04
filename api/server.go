package api

import (
	"context"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/mockingio/engine/persistent"
	"github.com/mockingio/mockingio/server"
)

type Server struct {
	db         persistent.Persistent
	mockServer *server.Server
}

func NewServer(db persistent.Persistent, mockServer *server.Server) *Server {
	return &Server{
		db:         db,
		mockServer: mockServer,
	}
}

func (s *Server) Start(_ context.Context, port string) (string, func(), error) {
	r := mux.NewRouter()

	r.Path("/mocks").HandlerFunc(s.GetMocksHandler).Methods(http.MethodGet)
	r.Path("/mocks/states").HandlerFunc(s.GetMocksStatesHandler).Methods(http.MethodGet)
	r.Path("/mocks").HandlerFunc(s.CreateMockHandler).Methods(http.MethodPost)
	r.Path("/mocks/{mock_id}/stop").HandlerFunc(s.StopMockServerHandler).Methods(http.MethodDelete)
	r.Path("/mocks/{mock_id}/start").HandlerFunc(s.StartMockServerHandler).Methods(http.MethodPost)

	// routes
	r.Path("/mocks/{mock_id}/routes/{route_id}").HandlerFunc(s.PatchRouteHandler).Methods(http.MethodPatch)
	r.Path("/mocks/{mock_id}/routes/{route_id}/responses/{response_id}").HandlerFunc(s.PatchResponseHandler).Methods(http.MethodPatch)

	addr := "0.0.0.0:" + port

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPut,
		http.MethodPatch,
		http.MethodOptions,
	}
	srv := &http.Server{
		Addr: addr,
		Handler: handlers.CORS(
			handlers.AllowedMethods(methods),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r),
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
