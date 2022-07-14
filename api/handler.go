package api

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/engine/mock"
	"github.com/smockyio/smocky/engine/persistent"
	"github.com/smockyio/smocky/server"
)

func GetMocksHandler(w http.ResponseWriter, r *http.Request) {
	db := persistent.GetDefault()
	mocks, err := db.GetMocks(r.Context())
	if err != nil {
		log.WithError(err).Error("get configs")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, mocks)
}

func GetMocksStatesHandler(w http.ResponseWriter, _ *http.Request) {
	response(w, http.StatusOK, server.GetStates())
}

func CreateMockHandler(w http.ResponseWriter, r *http.Request) {
	db := persistent.GetDefault()
	mo := mock.New()

	if err := db.SetMock(r.Context(), mo); err != nil {
		log.WithError(err).Error("create new mock")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	url, err := server.Start(r.Context(), mo)
	if err != nil {
		log.WithError(err).Error("start server")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusCreated, map[string]any{"id": mo.ID, "url": url})
}

func StopMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	server.RemoveServer(id)
	response(w, http.StatusOK, map[string]any{"id": id, "state": "stopped"})
}

func StartMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	if _, err := server.StartByID(r.Context(), id); err != nil {
		log.WithError(err).Error("start server by id")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, map[string]any{"id": id, "state": "running"})
}

func PauseMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	server.PauseServer(id)
	response(w, http.StatusOK, map[string]any{"id": id, "state": "paused"})
}

func ResumeMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	server.ResumeServer(id)
	response(w, http.StatusOK, map[string]any{"id": id, "state": "running"})
}
