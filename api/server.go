package api

import (
	"context"
	"github.com/tuongaz/smocky-engine/engine/persistent"
	"net/http"

	"github.com/gorilla/handlers"
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
	db := persistent.GetDefault()

	r.Path("/mocks").HandlerFunc(GetMocksHandler(db)).Methods(http.MethodGet)
	r.Path("/mocks/states").HandlerFunc(GetMocksStatesHandler).Methods(http.MethodGet)
	r.Path("/mocks").HandlerFunc(CreateMockHandler(db)).Methods(http.MethodPost)
	r.Path("/mocks/{mock_id}/stop").HandlerFunc(StopMockServerHandler).Methods(http.MethodDelete)
	r.Path("/mocks/{mock_id}/start").HandlerFunc(StartMockServerHandler).Methods(http.MethodPost)

	// routes
	r.Path("/mocks/{mock_id}/routes/{route_id}").HandlerFunc(PatchRouteHandler(db)).Methods(http.MethodPatch)
	r.Path("/mocks/{mock_id}/routes/{route_id}/responses/{response_id}").HandlerFunc(PatchResponseHandler(db)).Methods(http.MethodPatch)

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
