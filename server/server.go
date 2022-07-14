package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/smockyio/smocky/engine"
	"github.com/smockyio/smocky/engine/mock"
)

func Start(ctx context.Context, mo *mock.Mock) (string, error) {
	eng := engine.New(mo.ID)
	srv := buildHTTPServer(eng)

	listener, err := net.Listen("tcp", ":"+mo.Port)
	if err != nil {
		return "", err
	}

	done := make(chan bool, 1)

	go func() {
		_ = srv.Serve(listener)
	}()

	go func() {
		<-done
		_ = srv.Shutdown(ctx)
	}()

	serverPort := listener.Addr().(*net.TCPAddr).Port
	serverURL := fmt.Sprintf("http://127.0.0.1:%v", serverPort)
	addServer(mo.ID, &Controller{
		Pause:  eng.Pause,
		Resume: eng.Resume,
		Shutdown: func() {
			fmt.Printf("shutting down server: %v\n", serverURL)
			done <- true
		},
		State: State{
			MockID: mo.ID,
			URL:    serverURL,
			Status: Running,
		},
	})

	return serverURL, nil
}

func buildHTTPServer(e *engine.Engine) *http.Server {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(e.Handler)

	return &http.Server{Handler: r}
}
