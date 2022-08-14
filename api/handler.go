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

func (s *Server) DeleteRouteHandler(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
	routeID := mux.Vars(r)["route_id"]

	if err := s.db.DeleteRoute(r.Context(), mockID, routeID); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, nil)
}

func (s *Server) CreateRouteHandle(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
	data, err := io.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
	}

	newRoute, err := mock.Route{}.PatchString(string(data))

	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	if err = s.db.CreateRoute(r.Context(), mockID, *newRoute); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, *newRoute)
}

func (s *Server) DuplicateRouteHandle(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
	routeID := mux.Vars(r)["route_id"]

	route, err := s.db.GetRoute(r.Context(), mockID, routeID)

	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	newRoute := route.Clone()

	if err = s.db.CreateRoute(r.Context(), mockID, *newRoute); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, *newRoute)
}

func (s *Server) PatchResponseHandler(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
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

	if err = s.db.PatchResponse(r.Context(), mockID, responseID, string(data)); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, nil)
}

func (s *Server) CreateRuleHandler(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
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

	newRule, err := mock.Rule{}.PatchString(string(data))

	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	if err = s.db.CreateRule(r.Context(), mockID, responseID, *newRule); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, newRule)
}
func (s *Server) DeleteRuleHandler(w http.ResponseWriter, r *http.Request) {
	mockID := mux.Vars(r)["mock_id"]
	ruleID := mux.Vars(r)["rule_id"]

	if err := s.db.DeleteRule(r.Context(), mockID, ruleID); err != nil {
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
