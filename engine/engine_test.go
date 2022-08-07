package engine_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mockingio/mockingio/engine"
	"github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
	"github.com/mockingio/mockingio/engine/persistent/memory"
)

func TestEngine_Pause(t *testing.T) {
	eng := engine.New("mock-id", memory.New())
	eng.Pause()

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	eng.Handler(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func TestEngine_PauseResume(t *testing.T) {
	mem := setupMock()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	eng := engine.New("mock-id", mem)

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
	mem := setupMock()
	eng := engine.New("mock-id", mem)

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
	assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
	assert.Equal(t, "test", res.Header.Get("X-Test"))
}

func TestEngine_Match_With_Delay_Response(t *testing.T) {
	mem := memory.New()
	_ = mem.SetMock(context.Background(), &mock.Mock{
		ID: "mock-id",
		Routes: []*mock.Route{
			{
				Method: "GET",
				Path:   "/hello",
				Responses: []mock.Response{
					{
						Status: 200,
						Delay:  50,
					},
				},
			},
		},
	})

	eng := engine.New("mock-id", mem)

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	timer := time.Now()
	eng.Handler(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.True(t, time.Since(timer) > 50*time.Millisecond)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestEngine_NoResponses(t *testing.T) {
	mem := memory.New()
	_ = mem.SetMock(context.Background(), &mock.Mock{
		ID: "mock-id",
		Routes: []*mock.Route{
			{
				Method:    "GET",
				Path:      "/hello",
				Responses: []mock.Response{},
			},
		},
	})

	eng := engine.New("mock-id", mem)

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	eng.Handler(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestEngine_ProxyHandler(t *testing.T) {
	proxyServer := httptest.NewServer(proxyHandler(t))
	defer proxyServer.Close()

	httpsProxyServer := httptest.NewTLSServer(proxyHandler(t))
	defer httpsProxyServer.Close()

	tests := []struct {
		name     string
		mock     *mock.Mock
		assertFn func(t *testing.T, res *http.Response)
	}{
		{
			name: "proxy to http host",
			mock: proxyMock(proxyServer.URL, false),
			assertFn: func(t *testing.T, res *http.Response) {
				body, _ := io.ReadAll(res.Body)
				defer func() {
					_ = res.Body.Close()
				}()

				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "From Proxy", string(body))
				assert.Equal(t, "html/text", res.Header.Get("Content-Type"))
				assert.Equal(t, "from response", res.Header.Get("X-Response"))
			},
		},
		{
			name: "proxy to https host, skip TLS check",
			mock: proxyMock(httpsProxyServer.URL, true),
			assertFn: func(t *testing.T, res *http.Response) {
				body, _ := io.ReadAll(res.Body)
				defer func() {
					_ = res.Body.Close()
				}()

				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "From Proxy", string(body))
				assert.Equal(t, "html/text", res.Header.Get("Content-Type"))
				assert.Equal(t, "from response", res.Header.Get("X-Response"))
			},
		},
		{
			name: "proxy to https host, TLS check, expect fall back to not found",
			mock: proxyMock(httpsProxyServer.URL, false),
			assertFn: func(t *testing.T, res *http.Response) {
				assert.Equal(t, http.StatusNotFound, res.StatusCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.New()
			_ = mem.SetMock(context.Background(), tt.mock)
			_ = mem.SetActiveSession(context.Background(), tt.mock.ID, "session-id")
			eng := engine.New("mock-id", mem)

			w := httptest.NewRecorder()
			eng.Handler(w, httptest.NewRequest(http.MethodGet, "/hello", nil))
			res := w.Result()
			tt.assertFn(t, res)
		})
	}
}

func TestEngine_CORS_Request(t *testing.T) {
	tests := []struct {
		name               string
		mok                *mock.Mock
		expectedStatusCode int
	}{
		{
			"Auto CORS option is disabled, and no matching response, expect default 404",
			&mock.Mock{
				ID: "mock-id",
			},
			http.StatusNotFound,
		}, {
			"Auto CORS option is enabled, and no matching response, expect 200 OK",
			&mock.Mock{
				ID:       "mock-id",
				AutoCORS: true,
			},
			http.StatusOK,
		},
		{
			"Auto CORS option is enabled, and has matching response, use the matched response",
			&mock.Mock{
				ID:       "mock-id",
				AutoCORS: true,
				Routes: []*mock.Route{
					{
						Method: "OPTIONS",
						Path:   "/hello",
						Responses: []mock.Response{
							{
								Status: 201,
							},
						},
					},
				},
			},
			http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.New()
			_ = mem.SetMock(context.Background(), tt.mok)
			eng := engine.New("mock-id", mem)
			req := httptest.NewRequest(http.MethodOptions, "/hello", nil)
			w := httptest.NewRecorder()
			eng.Handler(w, req)
			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestEngine_MockNotFound(t *testing.T) {
	mem := memory.New()
	_ = mem.SetMock(context.Background(), &mock.Mock{
		ID: "mock-2",
	})
	eng := engine.New("mock-1", mem)
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)

	assert.Nil(t, eng.Match(req))
}

func setupMock() persistent.Persistent {
	mok := &mock.Mock{
		ID:       "mock-id",
		AutoCORS: true,
		Routes: []*mock.Route{
			{
				Method: "GET",
				Path:   "/hello",
				Responses: []mock.Response{
					{
						Status: 200,
						Body:   "Hello World",
						Headers: map[string]string{
							"Content-Type": "text/plain",
							"X-Test":       "test",
						},
					},
				},
			},
		},
	}
	mem := memory.New()
	_ = mem.SetMock(context.Background(), mok)

	return mem
}

func proxyHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hello" {
			t.Errorf("request path is %s, not /hello", r.URL.Path)
		}

		if r.Header.Get("X-Request") != "from request" {
			t.Error("header is not append to the request")
		}

		w.Header().Set("Content-Type", "html/text")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("From Proxy"))
	}
}

func proxyMock(proxyHost string, skipTLS bool) *mock.Mock {
	return &mock.Mock{
		ID: "mock-id",
		Proxy: &mock.Proxy{
			Enabled:            true,
			Host:               proxyHost,
			InsecureSkipVerify: skipTLS,
			RequestHeaders: map[string]string{
				"X-Request": "from request",
			},
			ResponseHeaders: map[string]string{
				"X-Response": "from response",
			},
		},
	}
}
