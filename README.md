# smocky

[![CI](https://github.com/tuongaz/smocky/actions/workflows/main.yml/badge.svg)](https://github.com/tuongaz/smocky/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/tuongaz/smocky/branch/main/graph/badge.svg?token=0AXGI7UR85)](https://codecov.io/gh/tuongaz/smocky)
[![Docker Repository](https://img.shields.io/docker/pulls/tuongaz/smocky)](https://hub.docker.com/r/tuongaz/smocky)
[![Github Release](https://img.shields.io/github/v/release/tuongaz/smocky)](https://github.com/tuongaz/smocky/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tuongaz/smocky)](https://goreportcard.com/report/github.com/tuongaz/smocky)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


## Installation

### Go install

```go
go install github.com/smockyio/smocky@latest
```

### Homebrew

```shell
brew tap tuongaz/smocky-tap
brew install tuongaz/smocky-tap/smocky
```

### Docker hub

```shell
docker pull tuongaz/smocky

docker run -ti tuongaz/smocky --version
```

## Usage

### CLI

`smocky start --filename example/mock.yml`

### Go package

```go
package example_test

import (
	"net/http"
	"testing"

	"github.com/smockyio/smocky/backend/pkg/smocky"
)

func Test_Example(t *testing.T) {
	srv := smocky.
		New().
		Get("/hello").
		Response(http.StatusOK, "hello world").
		Start(t)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	client := &http.Client{}

	resp, err := client.Do(req)
}

func Test_Example_WithRules(t *testing.T) {
	srv := smocky.
		New().
		Get("/hello").
		When("cookie", "name", "equal", "Chocolate").
		Response(http.StatusOK, "hello world").
		Start(t)
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	client := &http.Client{}

	resp, err := client.Do(req)
}
```