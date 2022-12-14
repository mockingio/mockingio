package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/mockingio/mockingio/engine/mock"
)

func (s *Server) GetMocksHandler(w http.ResponseWriter, r *http.Request) {
	mocks, err := s.db.GetMocks(r.Context())
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, mocks)
}

func (s *Server) GetMocksStatesHandler(w http.ResponseWriter, _ *http.Request) {
	response(w, http.StatusOK, s.mockServer.GetMockServerStates())
}

func (s *Server) CreateMockHandler(w http.ResponseWriter, r *http.Request) {
	mo := mock.New()

	if err := s.db.SetMock(r.Context(), mo); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	url, err := s.mockServer.NewMockServer(r.Context(), mo)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusCreated, map[string]any{"id": mo.ID, "url": url})
}

func (s *Server) PatchRouteHandler(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
	routeID := mux.Vars(r)["route_id"]
	data, err := io.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
	}

	if err = s.db.PatchRoute(r.Context(), mockID, routeID, string(data)); err != nil {
		responseError(w, http.StatusInternalServerError, err)
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
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	if len(data) == 0 {
		log.Error("request body empty")
		responseError(w, http.StatusBadRequest, errors.New("request body empty"))
		return
	}

	if err = s.db.PatchResponse(r.Context(), mockID, routeID, responseID, string(data)); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, nil)
}

func (s *Server) StopMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	resp, err := s.mockServer.StopMockServer(id)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, resp)
}

func (s *Server) StartMockServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["mock_id"]
	resp, err := s.mockServer.NewMockServerByID(r.Context(), id)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, resp)
}
