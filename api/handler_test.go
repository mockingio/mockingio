package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/mockingio/mockingio/api/fixtures"
	"github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
	"github.com/mockingio/mockingio/engine/persistent/memory"
	"github.com/mockingio/mockingio/engine/server"
)

func TestServer_CreateMockHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := memory.New()
		mockServer := server.New(db)
		defer mockServer.StopAllServers()

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, mockServer)
		apiServer.CreateMockHandler(writer, &http.Request{})

		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("db error", func(t *testing.T) {
		db := &mockDB{}

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, nil)
		apiServer.CreateMockHandler(writer, &http.Request{})

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestServer_GetMocksHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := newDB(fixtures.Mock1())

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, nil)
		apiServer.GetMocksHandler(writer, &http.Request{})

		assert.Equal(t, http.StatusOK, writer.Code)
		assert.JSONEq(t, `
[
    {
        "id": "mock1",
        "routes": [
            {
                "id": "route1",
                "method": "GET",
                "path": "",
                "description": "",
                "responses": [
                    {
                        "id": "response1",
                        "status": 201,
						"delay": {
							"min": 0,
							"max": 0
						}
                    }
                ]
            }
        ]
    }
]
`, writer.Body.String())
	})

	t.Run("db error", func(t *testing.T) {
		db := &mockDB{}

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, nil)
		apiServer.GetMocksHandler(writer, &http.Request{})

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestServer_GetMocksStatesHandler(t *testing.T) {
	db := newDB(fixtures.Mock1())
	mockServer := server.New(db)
	state, _ := mockServer.NewMockServerByID(context.Background(), fixtures.Mock1().ID)
	defer mockServer.StopAllServers()

	writer := httptest.NewRecorder()
	apiServer := NewServer(db, mockServer)
	apiServer.GetMocksStatesHandler(writer, &http.Request{})

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, fmt.Sprintf(`{"mock1":{"mock_id":"mock1","url":"%v","status":"running"}}`, state.URL), writer.Body.String())
}

func TestServer_PatchRouteHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := newDB(fixtures.Mock1())

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, nil)

		req := &http.Request{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"method":"OPTIONS"}`)),
		}
		req = mux.SetURLVars(req, map[string]string{
			"mock_id":  fixtures.Mock1().ID,
			"route_id": "route1",
		})
		apiServer.PatchRouteHandler(writer, req)

		assert.Equal(t, http.StatusOK, writer.Code)

		mok, _ := db.GetMock(context.Background(), fixtures.Mock1().ID)
		assert.Equal(t, "OPTIONS", mok.Routes[0].Method)
	})

	t.Run("db error", func(t *testing.T) {
		db := &mockDB{}

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, nil)

		req := &http.Request{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"method":"OPTIONS"}`)),
		}
		req = mux.SetURLVars(req, map[string]string{
			"mock_id":  fixtures.Mock1().ID,
			"route_id": "route1",
		})
		apiServer.PatchRouteHandler(writer, req)

		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func TestServer_PatchResponseHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := newDB(fixtures.Mock1())

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, nil)

		req := &http.Request{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"status":407}`)),
		}
		req = mux.SetURLVars(req, map[string]string{
			"mock_id":     fixtures.Mock1().ID,
			"route_id":    "route1",
			"response_id": "response1",
		})
		apiServer.PatchResponseHandler(writer, req)

		mok, _ := db.GetMock(context.Background(), fixtures.Mock1().ID)
		assert.Equal(t, 407, mok.Routes[0].Responses[0].Status)
	})

	t.Run("db error", func(t *testing.T) {
		db := &mockDB{}

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, nil)

		req := &http.Request{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"status":407}`)),
		}
		req = mux.SetURLVars(req, map[string]string{
			"mock_id":     fixtures.Mock1().ID,
			"route_id":    "route1",
			"response_id": "response1",
		})
		apiServer.PatchResponseHandler(writer, req)
		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})

	t.Run("empty body request", func(t *testing.T) {
		apiServer := NewServer(nil, nil)

		req := &http.Request{
			Body: ioutil.NopCloser(bytes.NewBufferString(``)),
		}
		writer := httptest.NewRecorder()
		apiServer.PatchResponseHandler(writer, req)

		assert.Equal(t, 400, writer.Code)
	})
}

func TestServer_StartStopMockServerHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := newDB(fixtures.Mock1())
		mockServer := server.New(db)
		defer mockServer.StopAllServers()

		writer := httptest.NewRecorder()
		apiServer := NewServer(db, mockServer)

		req := &http.Request{}
		req = mux.SetURLVars(req, map[string]string{
			"mock_id": fixtures.Mock1().ID,
		})

		// Start mock server
		apiServer.StartMockServerHandler(writer, req)
		state := mockServer.GetMockServerStates()
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(
			t,
			fmt.Sprintf(`{"mock_id":"mock1","url":"%s","status":"running"}`, state["mock1"].URL),
			writer.Body.String(),
		)

		// Stop mock server
		writer = httptest.NewRecorder()
		apiServer.StopMockServerHandler(writer, req)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Equal(
			t,
			`{"mock_id":"mock1","url":"","status":"stopped"}`,
			writer.Body.String(),
		)
	})

	t.Run("start mock server, mock server error", func(t *testing.T) {
		writer := httptest.NewRecorder()
		apiServer := NewServer(nil, &mockMockServer{})

		req := &http.Request{}
		req = mux.SetURLVars(req, map[string]string{
			"mock_id": fixtures.Mock1().ID,
		})

		// Start mock server
		apiServer.StartMockServerHandler(writer, req)
		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})

	t.Run("stop mock server, mock server error", func(t *testing.T) {
		writer := httptest.NewRecorder()
		apiServer := NewServer(nil, &mockMockServer{})

		req := &http.Request{}
		req = mux.SetURLVars(req, map[string]string{
			"mock_id": fixtures.Mock1().ID,
		})

		// Start mock server
		apiServer.StopMockServerHandler(writer, req)
		assert.Equal(t, http.StatusInternalServerError, writer.Code)
	})
}

func newDB(mocks ...*mock.Mock) persistent.Database {
	db := memory.New()
	for _, m := range mocks {
		_ = db.SetMock(context.Background(), m)
	}
	return db
}

type mockDB struct {
	persistent.CRUD
}

func (m *mockDB) GetMocks(_ context.Context) ([]*mock.Mock, error) {
	return nil, errors.New("something is not right")
}

func (m *mockDB) SetMock(_ context.Context, _ *mock.Mock) error {
	return errors.New("something is not right")
}

func (m *mockDB) PatchRoute(_ context.Context, _ string, _ string, _ string) error {
	return errors.New("something is not right")
}

func (m *mockDB) PatchResponse(_ context.Context, _, _, _, _ string) error {
	return errors.New("something is not right")
}

type mockMockServer struct {
	mockServer
}

func (m *mockMockServer) NewMockServerByID(_ context.Context, _ string) (*server.MockServerState, error) {
	return nil, errors.New("something is not right")
}

func (m *mockMockServer) NewMockServer(_ context.Context, _ *mock.Mock) (*server.MockServerState, error) {
	return nil, errors.New("something is not right")
}

func (m *mockMockServer) StopMockServer(_ string) (*server.MockServerState, error) {
	return nil, errors.New("something is not right")
}
