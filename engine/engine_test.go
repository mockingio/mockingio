package engine_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smockyio/smocky/engine"
	"github.com/smockyio/smocky/engine/mock"
	"github.com/smockyio/smocky/engine/persistent"
	"github.com/smockyio/smocky/engine/persistent/memory"
)

func TestEngine_Pause(t *testing.T) {
	eng := engine.New("mock-id")
	eng.Pause()

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	eng.Handler(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func TestEngine_PauseResume(t *testing.T) {
	setupMock()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	eng := engine.New("mock-id")

	eng.Pause()
	w := httptest.NewRecorder()
	eng.Handler(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	eng.Resume()
	w = httptest.NewRecorder()
	eng.Handler(w, req)
	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestEngine_Match(t *testing.T) {
	setupMock()
	eng := engine.New("mock-id")

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	eng.Handler(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	bod, err := ioutil.ReadAll(res.Body)

	require.NoError(t, err)
	assert.Equal(t, "Hello World", string(bod))
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func setupMock() {
	mok := &mock.Mock{
		ID: "mock-id",
		Routes: []*mock.Route{
			{
				Request: "GET /hello",
				Responses: []mock.Response{
					{
						Status: 200,
						Body:   "Hello World",
					},
				},
			},
		},
	}
	mem := memory.New()
	persistent.New(mem)
	_ = mem.SetMock(context.Background(), mok)
}
