package mock_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/mockingio/mockingio/pkg/mock"
)

func TestBuilder_Validation(t *testing.T) {
	_, err := New().Start()
	assert.Error(t, err)
}

func TestBuilder_MatchedRoute(t *testing.T) {
	t.Run("simple get", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			Start()
		defer srv.Close()

		require.NoError(t, err)
		assertHTTPGETRequest(t, url(srv, "/hello"), 200, "world")
	})

	t.Run("simple post", func(t *testing.T) {
		srv, err := New().
			Post("/hello").
			Response(http.StatusOK, "world").
			Start()
		defer srv.Close()

		require.NoError(t, err)
		assertHTTPPOSTRequest(t, url(srv, "/hello"), "", 200, "world")
	})

	t.Run("simple put", func(t *testing.T) {
		srv, err := New().
			Put("/hello").
			Response(http.StatusOK, "world").
			Start()
		defer srv.Close()

		require.NoError(t, err)
		assertHTTPPUTRequest(t, url(srv, "/hello"), "", 200, "world")
	})

	t.Run("simple options", func(t *testing.T) {
		srv, err := New().
			Option("/hello").
			Response(http.StatusOK, "world").
			Start()
		defer srv.Close()

		require.NoError(t, err)
		assertHTTPOPTIONSRequest(t, url(srv, "/hello"), "", 200, "world")
	})

	t.Run("simple put", func(t *testing.T) {
		srv, err := New().
			Delete("/hello").
			Response(http.StatusOK, "world").
			Start()
		defer srv.Close()

		require.NoError(t, err)
		assertHTTPDELETETRequest(t, url(srv, "/hello"), "", 200, "world")
	})

	t.Run("simple get with builder", func(t *testing.T) {
		builder := New()
		builder.Get("/hello").
			Response(http.StatusOK, "world")
		srv, err := builder.Start()
		defer srv.Close()

		require.NoError(t, err)
		assertHTTPGETRequest(t, url(srv, "/hello"), 200, "world")
	})

	t.Run("with a condition", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			When(Cookie, "name", Equal, "Jack").
			Start()

		require.NoError(t, err)
		assertHTTPGETRequest(t, url(srv, "/hello"), 200, "world")
	})

	t.Run("with AND condition", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			When(Cookie, "name", Equal, "Jack").
			And(Cookie, "name", Regex, "[a-zA-Z]+").
			And(Header, "x-type", Regex, "x-men").
			Start()

		require.NoError(t, err)
		assertHTTPGETRequest(t, url(srv, "/hello"), 200, "world")
	})

	t.Run("with OR condition", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			When(Cookie, "name", Equal, "Jack").
			Or(Header, "name", Regex, "non exist value").
			Or(Header, "age", Equal, "18").
			Start()

		require.NoError(t, err)
		assertHTTPGETRequest(t, url(srv, "/hello"), 200, "world")
	})

	t.Run("with multiple routes", func(t *testing.T) {
		builder := New()
		builder.Get("/hello1").
			Response(http.StatusOK, "world1")

		builder.Get("/hello2").
			Response(http.StatusOK, "world2")

		builder.Post("/hello3").
			Response(http.StatusOK, "world3")

		srv, err := builder.Start()
		defer srv.Close()

		require.NoError(t, err)
		assertHTTPGETRequest(t, url(srv, "/hello1"), 200, "world1")
		assertHTTPGETRequest(t, url(srv, "/hello2"), 200, "world2")
		assertHTTPPOSTRequest(t, url(srv, "/hello3"), "", 200, "world3")
	})

	t.Run("with delay", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			Delay(50, 100).
			Start()
		defer srv.Close()

		timer := time.Now()

		require.NoError(t, err)
		assertHTTPGETRequest(t, url(srv, "/hello"), 200, "world")
		assert.Greater(t, time.Since(timer), 50*time.Millisecond)
	})
}

func TestBuilder_NoMatchedRoute(t *testing.T) {
	t.Run("simple get", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			Start()
		defer srv.Close()

		require.NoError(t, err)
		assertNoMatchHTTPGETRequest(t, url(srv, "/hellox"), 200)
	})

	t.Run("with a condition", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			When(Cookie, "name", Equal, "Jack not found").
			Start()

		require.NoError(t, err)
		assertNoMatchHTTPGETRequest(t, url(srv, "/hello"), 200)
	})

	t.Run("with AND condition", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			When(Cookie, "name", Equal, "Jack").
			And(Cookie, "name", Regex, "[a-zA-Z]+").
			And(Header, "x-type", Regex, "x-women").
			Start()

		require.NoError(t, err)
		assertNoMatchHTTPGETRequest(t, url(srv, "/hello"), 200)
	})

	t.Run("with OR condition", func(t *testing.T) {
		srv, err := New().
			Get("/hello").
			Response(http.StatusOK, "world").
			When(Cookie, "name", Equal, "Jack not found").
			Or(Header, "name", Regex, "non exist value").
			Start()

		require.NoError(t, err)
		assertNoMatchHTTPGETRequest(t, url(srv, "/hello"), 200)
	})
}

func url(srv *httptest.Server, path string) string {
	return fmt.Sprintf("%v%v", srv.URL, path)
}

func assertNoMatchHTTPGETRequest(t *testing.T, url string, responseStatus int) {
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	req.Header.Set("x-type", "x-men")
	req.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "Jack",
	})

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.NotEqual(t, responseStatus, resp.StatusCode)
}

func assertHTTPGETRequest(t *testing.T, url string, expectedStatusCode int, expectedResponseBody string) {
	t.Helper()
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	req.Header.Set("x-type", "x-men")
	req.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "Jack",
	})

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expectedResponseBody, string(body))
}

func assertHTTPPOSTRequest(t *testing.T, url, payload string, expectedStatusCode int, expectedResponseBody string) {
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	require.NoError(t, err)

	req.Header.Set("x-type", "x-men")
	req.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "Jack",
	})

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expectedResponseBody, string(body))
}

func assertHTTPOPTIONSRequest(t *testing.T, url, payload string, expectedStatusCode int, expectedResponseBody string) {
	req, err := http.NewRequest("OPTIONS", url, strings.NewReader(payload))
	require.NoError(t, err)

	req.Header.Set("x-type", "x-men")
	req.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "Jack",
	})

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expectedResponseBody, string(body))
}

func assertHTTPPUTRequest(t *testing.T, url, payload string, expectedStatusCode int, expectedResponseBody string) {
	req, err := http.NewRequest("PUT", url, strings.NewReader(payload))
	require.NoError(t, err)

	req.Header.Set("x-type", "x-men")
	req.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "Jack",
	})

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expectedResponseBody, string(body))
}

func assertHTTPDELETETRequest(t *testing.T, url, payload string, expectedStatusCode int, expectedResponseBody string) {
	req, err := http.NewRequest("DELETE", url, strings.NewReader(payload))
	require.NoError(t, err)

	req.Header.Set("x-type", "x-men")
	req.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "Jack",
	})

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expectedResponseBody, string(body))
}
