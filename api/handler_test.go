package api_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smockyio/smocky/api"
	"github.com/tuongaz/smocky-engine/engine/mock"
	"github.com/tuongaz/smocky-engine/engine/persistent/memory"
)

func Test_GetMocksHandler(t *testing.T) {
	db := memory.New()
	_ = db.SetMock(context.Background(), &mock.Mock{
		ID: "123",
	})

	req := httptest.NewRequest(http.MethodGet, "/mocks", nil)
	w := httptest.NewRecorder()
	api.GetMocksHandler(db)(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `[{"id":"123"}]`, string(data))
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_CreateMockHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/mocks", nil)
	w := httptest.NewRecorder()
	api.CreateMockHandler(memory.New())(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	_, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}
