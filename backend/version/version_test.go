package version

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShort(t *testing.T) {
	require.Equal(t, Short(), "smocky devel ()")
}

func TestLong(t *testing.T) {
	require.Equal(t, Long(), fmt.Sprintf("smocky devel () %s/%s", runtime.GOOS, runtime.GOARCH))
}
