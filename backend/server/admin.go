package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/backend/persistent"
)

type AdminServer struct {
}

func NewAdminServer() *AdminServer {
	return &AdminServer{}
}

func (a *AdminServer) Start(ctx context.Context, port string) (string, func(), error) {
	r := mux.NewRouter()
	r.PathPrefix("/mocks").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		db := persistent.GetDefault()
		mocks, err := db.GetConfigs(ctx)
		if err != nil {
			log.WithError(err).Error("get configs")
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response(w, http.StatusOK, mocks)
	}).Methods(http.MethodGet)

	r.PathPrefix("/mocks").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write([]byte("create new mock"))
	}).Methods(http.MethodPost)

	r.PathPrefix("/mocks/{mock_id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write([]byte("put mock"))
	}).Methods(http.MethodPut)

	addr := "0.0.0.0:" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: r, // Pass our instance of gorilla/mux in.
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
