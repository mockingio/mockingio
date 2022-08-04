package api

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/mockingio/engine/mock"
)

func (s *Server) GetMocksHandler(w http.ResponseWriter, r *http.Request) {
	mocks, err := s.db.GetMocks(r.Context())
	if err != nil {
		log.WithError(err).Error("get configs")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, mocks)
}

func (s *Server) GetMocksStatesHandler(w http.ResponseWriter, _ *http.Request) {
	response(w, http.StatusOK, s.mockServer.GetStates())
}

func (s *Server) CreateMockHandler(w http.ResponseWriter, r *http.Request) {
	mo := mock.New()

	if err := s.db.SetMock(r.Context(), mo); err != nil {
		log.WithError(err).Error("create new mock")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	url, err := s.mockServer.NewMockServer(r.Context(), mo)
	if err != nil {
		log.WithError(err).Error("start server")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusCreated, map[string]any{"id": mo.ID, "url": url})
}

func (s *Server) PatchRouteHandler(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
	routeID := mux.Vars(r)["route_id"]
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("read request body")
		responseError(w, http.StatusInternalServerError, err.Error())
	}

	if err = s.db.PatchRoute(r.Context(), mockID, routeID, string(data)); err != nil {
		log.WithError(err).Error("patch route")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, nil)
}

func (s *Server) PatchResponseHandler(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
	routeID := mux.Vars(r)["route_id"]
	responseID := mux.Vars(r)["response_id"]
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("read request body")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(data) == 0 {
		log.Error("request body empty")
		responseError(w, http.StatusBadRequest, "request body empty")
		return
	}

	if err = s.db.PatchResponse(r.Context(), mockID, routeID, responseID, string(data)); err != nil {
		log.WithError(err).Error("patch route")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, nil)
}

func (s *Server) StopMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	resp, err := s.mockServer.StopMockServer(id)
	if err != nil {
		log.WithError(err).Error("stop mock server")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, resp)
}

func (s *Server) StartMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	resp, err := s.mockServer.NewMockServerByID(r.Context(), id)
	if err != nil {
		log.WithError(err).Error("start server by id")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, resp)
}
