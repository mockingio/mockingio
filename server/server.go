package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/mockingio/engine"
	"github.com/mockingio/engine/mock"
	"github.com/mockingio/engine/persistent"
)

type Server struct {
	mu               sync.Mutex
	db               persistent.Persistent
	mockServerStates map[string]*MockServerState
}

func New(db persistent.Persistent) *Server {
	return &Server{db: db, mockServerStates: make(map[string]*MockServerState)}
}

func (s *Server) NewMockServerByID(ctx context.Context, id string) (*MockServerState, error) {
	mo, err := s.db.GetMock(ctx, id)
	if err != nil {
		return nil, err
	}

	if mo == nil {
		return nil, fmt.Errorf("mock with ID: %s not found", id)
	}

	return s.NewMockServer(ctx, mo)
}

func (s *Server) NewMockServer(ctx context.Context, mo *mock.Mock) (*MockServerState, error) {
	eng := engine.New(mo.ID, s.db)
	srv := buildHTTPServer(eng)

	listener, err := net.Listen("tcp", ":"+mo.Port)
	if err != nil {
		return nil, err
	}

	shutdownC := make(chan bool, 1)

	go func() {
		_ = srv.Serve(listener)
	}()

	go func() {
		<-shutdownC
		_ = srv.Shutdown(ctx)
	}()

	serverPort := listener.Addr().(*net.TCPAddr).Port
	serverURL := fmt.Sprintf("http://127.0.0.1:%v", serverPort)

	state := s.addNewMockServerState(mo.ID, serverURL, func() {
		fmt.Printf("shutting down server: %v\n", serverURL)
		shutdownC <- true
	})

	return state, nil
}

func (s *Server) GetMockServerURLs() []string {
	var urls []string
	for _, state := range s.mockServerStates {
		if state.Status == Running {
			urls = append(urls, state.URL)
		}
	}

	return urls
}

func (s *Server) StopMockServer(mockID string) (*MockServerState, error) {
	state, err := s.getMockServerState(mockID)
	if err != nil {
		return nil, err
	}

	state.shutdownServer()

	return state, nil
}

func (s *Server) GetMockServerStates() map[string]*MockServerState {
	return s.mockServerStates
}

func (s *Server) StopAllServers() {
	for _, state := range s.mockServerStates {
		state.shutdownServer()
	}
}

func (s *Server) addNewMockServerState(mockID string, url string, shutdownServer func()) *MockServerState {
	state := &MockServerState{
		MockID:           mockID,
		URL:              url,
		Status:           Running,
		shutdownServerFn: shutdownServer,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.mockServerStates[mockID] = state

	return state
}

func (s *Server) getMockServerState(mockID string) (*MockServerState, error) {
	if state, ok := s.mockServerStates[mockID]; ok {
		return state, nil
	}

	return nil, fmt.Errorf("mock server: %v not found", mockID)
}

func buildHTTPServer(e *engine.Engine) *http.Server {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(e.Handler)

	return &http.Server{Handler: r}
}
