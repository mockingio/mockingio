package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smockyio/smocky/backend/server"
	"github.com/smockyio/smocky/engine/persistent"
	"github.com/smockyio/smocky/engine/persistent/memory"
)

func TestIntegration(t *testing.T) {
	endpoint, stop := mustStartServer(t)
	t.Logf("server: %v", endpoint)
	defer stop()

	url := func(path string) string {
		return fmt.Sprintf("%v%v", endpoint, path)
	}

	t.Run("match url", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/hello"), 200, "world")
	})

	t.Run("response in sequence", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/hello/response/sequence"), 201, "one")
		assertHTTPGETRequest(t, url("/hello/response/sequence"), 202, "two")
		assertHTTPGETRequest(t, url("/hello/response/sequence"), 203, "three")
	})

	t.Run("default response", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/hello/response/default"), 202, "two")
	})

	t.Run("success after 3 requests", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/request/number/rule"), 404, "")
		assertHTTPGETRequest(t, url("/request/number/rule"), 404, "")
		assertHTTPGETRequest(t, url("/request/number/rule"), 200, "success")
	})

	t.Run("match wildcard", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/match/wildcard/one"), 200, "success")
		assertHTTPGETRequest(t, url("/match/wildcard/two"), 200, "success")
	})

	t.Run("match request body", func(t *testing.T) {
		assertHTTPPOSTRequest(t, url("/match/request/body"), `{"address": {"suburb": "Melbourne", "state": "Victoria"}}`, 201, "success")
	})

	t.Run("match route param", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/match/route/param/joe"), 200, "success")
	})

	t.Run("match cookie", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/match/cookie"), 200, "success")
	})

	t.Run("match header", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/match/header"), 200, "success")
	})

	t.Run("match query string", func(t *testing.T) {
		assertHTTPGETRequest(t, url("/match/query?name=alex&age=123"), 200, "success")
	})
}

func mustStartServer(t *testing.T) (string, func()) {
	persistent.New(memory.New())
	address, done, err := server.New().StartFromFile(context.Background(), "mock.yml")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	return address, done
}

func assertHTTPPOSTRequest(t *testing.T, url string, requestBody string, expectedStatusCode int, expectedResponseBody string) {
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, expectedResponseBody, string(body))
}

func assertHTTPGETRequest(t *testing.T, url string, expectedStatusCode int, expectedResponseBody string) {
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
