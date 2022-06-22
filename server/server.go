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

type Controller struct {
	Pause    func()
	Resume   func()
	Shutdown func()
	URL      string
}

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

	serverURL := fmt.Sprintf("http://0.0.0.0:%v", listener.Addr().(*net.TCPAddr).Port)
	addServer(mo.ID, &Controller{
		Pause:  eng.Pause,
		Resume: eng.Resume,
		Shutdown: func() {
			fmt.Printf("shutting down server: %v\n", serverURL)
			done <- true
		},
		URL: serverURL,
	})

	return serverURL, nil
}

func buildHTTPServer(e *engine.Engine) *http.Server {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(e.Handler)

	return &http.Server{Handler: r}
}
