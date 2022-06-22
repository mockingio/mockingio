package backend

import (
	"context"
	"fmt"
	"github.com/smockyio/smocky/engine/persistent"
	"net"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/smockyio/smocky/engine"
)

type Result struct {
	Error error
}

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context, mockID string) (string, func(), error) {
	db := persistent.GetDefault()
	mo, err := db.GetMock(ctx, mockID)
	if err != nil {
		return "", func() {}, err
	}

	eng := engine.New(mo.ID)
	srv := s.buildHTTPServer(eng)

	listener, err := net.Listen("tcp", ":"+mo.Port)
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

func (s *Server) buildHTTPServer(m *engine.Engine) *http.Server {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(m.Handler)

	return &http.Server{Handler: r}
}
