package api_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/mockingio/engine/mock"
	"github.com/mockingio/engine/persistent"
	"github.com/mockingio/engine/persistent/memory"
	"github.com/mockingio/mockingio/api"
	"github.com/mockingio/mockingio/api/fixtures"
	"github.com/mockingio/mockingio/server"
)

func TestServer_CreateMockHandler(t *testing.T) {
	db := memory.New()
	mockServer := server.New(db)
	defer mockServer.StopAllServers()

	writer := httptest.NewRecorder()
	apiServer := api.NewServer(db, mockServer)
	apiServer.CreateMockHandler(writer, &http.Request{})

	assert.Equal(t, http.StatusCreated, writer.Code)
}

func TestServer_GetMocksHandler(t *testing.T) {
	db := newDB(fixtures.Mock1())

	writer := httptest.NewRecorder()
	apiServer := api.NewServer(db, nil)
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
                        "status": 201
                    }
                ]
            }
        ]
    }
]
`, writer.Body.String())
}

func TestServer_GetMocksStatesHandler(t *testing.T) {
	db := newDB(fixtures.Mock1())
	mockServer := server.New(db)
	state, _ := mockServer.NewMockServerByID(context.Background(), fixtures.Mock1().ID)
	defer mockServer.StopAllServers()

	writer := httptest.NewRecorder()
	apiServer := api.NewServer(db, mockServer)
	apiServer.GetMocksStatesHandler(writer, &http.Request{})

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, fmt.Sprintf(`{"mock1":{"mock_id":"mock1","url":"%v","status":"running"}}`, state.URL), writer.Body.String())
}

func TestServer_PatchRouteHandler(t *testing.T) {
	db := newDB(fixtures.Mock1())

	writer := httptest.NewRecorder()
	apiServer := api.NewServer(db, nil)

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
}

func TestServer_PatchResponseHandler(t *testing.T) {
	db := newDB(fixtures.Mock1())

	writer := httptest.NewRecorder()
	apiServer := api.NewServer(db, nil)

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
}

func TestServer_StartStopMockServerHandler(t *testing.T) {
	db := newDB(fixtures.Mock1())
	mockServer := server.New(db)
	defer mockServer.StopAllServers()

	writer := httptest.NewRecorder()
	apiServer := api.NewServer(db, mockServer)

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
}

func newDB(mocks ...*mock.Mock) persistent.Persistent {
	db := memory.New()
	for _, m := range mocks {
		_ = db.SetMock(context.Background(), m)
	}
	return db
}
