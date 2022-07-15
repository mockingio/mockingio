package api

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/server"
	"github.com/tuongaz/smocky-engine/engine/mock"
	"github.com/tuongaz/smocky-engine/engine/persistent"
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
	resp, err := server.RemoveServer(id)
	if err != nil {
		log.WithError(err).Error("stop mock server")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, resp)
}

func StartMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	resp, err := server.StartByID(r.Context(), id)
	if err != nil {
		log.WithError(err).Error("start server by id")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, resp)
}
