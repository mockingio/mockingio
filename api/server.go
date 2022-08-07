package api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	mockEngine "github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
	"github.com/mockingio/mockingio/server"
)

type Server struct {
	db         persistent.Persistent
	mockServer mockServer
}

func NewServer(db persistent.Persistent, mockServer mockServer) *Server {
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

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return "", nil, errors.Wrapf(err, "listen to tcp port: %s", port)
	}
	addr := fmt.Sprintf("0.0.0.0:%v", listener.Addr().(*net.TCPAddr).Port)

	srv := &http.Server{
		Addr: addr,
		Handler: handlers.CORS(
			handlers.AllowedMethods([]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodPut,
				http.MethodPatch,
				http.MethodOptions,
			}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r),
	}

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Println(err)
		}
	}()

	return addr, func() {
		_ = srv.Shutdown(context.Background())
	}, nil
}

type mockServer interface {
	NewMockServerByID(ctx context.Context, id string) (*server.MockServerState, error)
	NewMockServer(ctx context.Context, mo *mockEngine.Mock) (*server.MockServerState, error)
	GetMockServerURLs() []string
	StopMockServer(mockID string) (*server.MockServerState, error)
	GetMockServerStates() map[string]*server.MockServerState
	StopAllServers()
}

var _ mockServer = server.New(nil)
