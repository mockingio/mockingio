package api

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/server"
	"github.com/tuongaz/smocky-engine/engine/mock"
	"github.com/tuongaz/smocky-engine/engine/persistent"
)

func GetMocksHandler(db persistent.Persistent) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mocks, err := db.GetMocks(r.Context())
		if err != nil {
			log.WithError(err).Error("get configs")
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response(w, http.StatusOK, mocks)
	}
}

func GetMocksStatesHandler(w http.ResponseWriter, _ *http.Request) {
	response(w, http.StatusOK, server.GetStates())
}

func CreateMockHandler(db persistent.Persistent) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mo := mock.New()

		if err := db.SetMock(r.Context(), mo); err != nil {
			log.WithError(err).Error("create new mock")
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		url, err := server.Start(r.Context(), mo, db)
		if err != nil {
			log.WithError(err).Error("start server")
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response(w, http.StatusCreated, map[string]any{"id": mo.ID, "url": url})
	}
}

func PatchRouteHandler(db persistent.Persistent) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mockID := mux.Vars(r)["mock_id"]
		routeID := mux.Vars(r)["route_id"]
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("read request body")
			responseError(w, http.StatusInternalServerError, err.Error())
		}

		if err := db.PatchRoute(r.Context(), mockID, routeID, string(data)); err != nil {
			log.WithError(err).Error("patch route")
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response(w, http.StatusOK, nil)
	}
}

func PatchResponseHandler(db persistent.Persistent) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mockID := mux.Vars(r)["mock_id"]
		routeID := mux.Vars(r)["route_id"]
		responseID := mux.Vars(r)["response_id"]
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("read request body")
			responseError(w, http.StatusInternalServerError, err.Error())
		}

		if err := db.PatchResponse(r.Context(), mockID, routeID, responseID, string(data)); err != nil {
			log.WithError(err).Error("patch route")
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response(w, http.StatusOK, nil)
	}
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

func StartMockServerHandler(db persistent.Persistent) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["mock_id"]
		resp, err := server.StartByID(r.Context(), id, db)
		if err != nil {
			log.WithError(err).Error("start server by id")
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response(w, http.StatusOK, resp)
	}
}
