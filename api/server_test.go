package api_test

import (
	"context"
	"net"
	"testing"

	"github.com/mockingio/mockingio/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_Start(t *testing.T) {
	srv := api.NewServer(nil, nil)
	url, stop, err := srv.Start(context.Background(), "0")
	require.NoError(t, err)

	_, err = net.Dial("tcp", url)
	assert.NoError(t, err, "dial tcp should have no error")

	stop()
	_, err = net.Dial("tcp", url)
	assert.Error(t, err, "dial tcp should have error")

	assert.True(t, url != "", "url is empty")
}
