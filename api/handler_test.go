package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	engine "github.com/mockingio/engine/mock"
	"github.com/mockingio/mockingio/api"
	"github.com/mockingio/mockingio/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_Handlers(t *testing.T) {
	tests := []struct {
		name                 string
		handler              http.HandlerFunc
		request              *http.Request
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name: "get mocks",
			handler: func() http.HandlerFunc {
				db := mocks.NewDBMock(t)
				db.On("GetMocks", mock.Anything).Return([]*engine.Mock{}, nil)
				serv := api.NewServer(db, nil)
				return serv.GetMocksHandler
			}(),
			request:              &http.Request{},
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			tt.handler(writer, tt.request)

			assert.Equal(t, tt.expectedResponseCode, writer.Code)
		})
	}
}
