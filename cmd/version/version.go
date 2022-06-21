package version

import (
	"fmt"
	"runtime"
)

const (
	DevelopmentVersion = "devel"
)

var (
	Version  = DevelopmentVersion
	Revision = ""
)

func Short() string {
	return fmt.Sprintf("smocky %s (%s)", Version, Revision)
}

func Long() string {
	return fmt.Sprintf("%s %s/%s", Short(), runtime.GOOS, runtime.GOARCH)
}
