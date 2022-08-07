package server

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mockingio/mockingio/engine"
	"github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
)

const (
	serverURLFormat    = "http://127.0.0.1:%v"
	serverURLTLSFormat = "https://127.0.0.1:%v"
)

var (
	//go:embed certs/default-cert.pem
	defaultTLSCert []byte

	//go:embed certs/default-key.pem
	defaultTLSCertKey []byte
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

	var listener net.Listener
	var err error

	lAddr := "0.0.0.0:" + mo.Port
	if mo.TLSEnabled() {
		cert, err := getTLSCert(mo)
		if err != nil {
			return nil, errors.Wrap(err, "get TLS cert")
		}
		if cert == nil {
			return nil, errors.New("no TLS cert found")
		}

		listener, err = tls.Listen("tcp", lAddr, &tls.Config{
			Certificates: []tls.Certificate{*cert},
		})
		if err != nil {
			return nil, errors.Wrap(err, "listen TLS TCP")
		}
	} else {
		listener, err = net.Listen("tcp", lAddr)
		if err != nil {
			return nil, errors.Wrap(err, "listen TCP")
		}
	}

	shutdownC := make(chan bool, 1)

	serverPort := listener.Addr().(*net.TCPAddr).Port
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Error(errors.Wrapf(err, "serving HTTP at %v", listener.Addr().String()))
		}
	}()

	go func() {
		<-shutdownC
		_ = srv.Shutdown(ctx)
	}()

	urlFormat := serverURLFormat
	if mo.TLSEnabled() {
		urlFormat = serverURLTLSFormat
	}
	serverURL := fmt.Sprintf(urlFormat, serverPort)

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

func getTLSCert(mo *mock.Mock) (*tls.Certificate, error) {
	if !mo.TLSEnabled() {
		return nil, nil
	}

	certPath := mo.TLS.PEMCertPath
	keyPath := mo.TLS.PEMKeyPath

	if certPath == "" || keyPath == "" {
		certificate, err := tls.X509KeyPair(defaultTLSCert, defaultTLSCertKey)
		if err != nil {
			return nil, err
		}
		return &certificate, nil
	}

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read TLS cert file: %v", certPath)
	}

	certKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read TLS key file: %v", keyPath)
	}

	certificate, err := tls.X509KeyPair(cert, certKey)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}
