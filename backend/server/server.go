package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/smockyio/smocky/backend/mock"
	"github.com/smockyio/smocky/backend/mock/config"
	"github.com/smockyio/smocky/backend/persistent"
)

type Result struct {
	Error error
}

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) StartFromFile(ctx context.Context, file string) (string, func(), error) {
	cfg, err := config.FromYamlFile(file)
	if err != nil {
		return "", nil, err
	}
	mockID := uuid.NewString()
	cfg.ID = mockID

	sessionID := uuid.NewString()

	db := persistent.GetDefault()
	if err := db.SetActiveSession(ctx, mockID, sessionID); err != nil {
		return "", nil, err
	}

	_ = db.SetConfig(ctx, cfg)

	m, err := mock.New(mockID)
	if err != nil {
		return "", nil, err
	}

	srv := s.buildHTTPServer(m)

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return "", nil, err
	}

	done := make(chan bool, 1)
	serverStopped := make(chan bool, 1)

	go func() {
		_ = srv.Serve(listener)
	}()

	go func() {
		<-done
		fmt.Println("shutting down server")
		_ = srv.Shutdown(ctx)
		serverStopped <- true
	}()

	serverURL := fmt.Sprintf("http://0.0.0.0:%v", listener.Addr().(*net.TCPAddr).Port)
	return serverURL, func() {
		fmt.Printf("shutting down server: %v\n", serverURL)
		done <- true
	}, nil
}

func (s *Server) buildHTTPServer(m *mock.Mock) *http.Server {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(m.Handler)

	return &http.Server{Handler: r}
}
