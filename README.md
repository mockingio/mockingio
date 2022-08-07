# mockingio

[![CI](https://github.com/mockingio/mockingio/actions/workflows/auto-release.yml/badge.svg)](https://github.com/mockingio/mockingio/actions/workflows/auto_release.yml)
[![codecov](https://codecov.io/gh/mockingio/mockingio/branch/main/graph/badge.svg?token=0AXGI7UR85)](https://codecov.io/gh/mockingio/mockingio)
[![Docker Repository](https://img.shields.io/docker/pulls/mockingio/mockingio)](https://hub.docker.com/r/mockingio/mockingio)
[![Github Release](https://img.shields.io/github/v/release/mockingio/mockingio)](https://github.com/mockingio/mockingio/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/mockingio/mockingio)](https://goreportcard.com/report/github.com/mockingio/mockingio)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


## Installation

### Go install

```go
go install github.com/mockingio/mockingio@latest
```

### Homebrew

```shell
brew tap mockingio/mockingio-tap
brew install mockingio/mockingio-tap/mockingio
```

### Docker hub

```shell
docker pull mockingio/mockingio

docker run -ti mockingio/mockingio --version
```

## Usage

### CLI
```yaml
# mock.yml
name: Example mock 1
routes:
  - method: GET
    path: /products
    responses:
      - body: |
          [
            {
              "id": "1",
              "name": "Product 1",
              "price": "10.00"
            },
            {
              "id": "2",
              "name": "Product 2",
              "price": "20.00"
            }
          ]

```
`mockingio start --filename mock.yml`

### Go package usage

```go
import (
	"net/http"
	"testing"

	mock "github.com/mockingio/mock"
)

func main() {
	srv, _ := mock.
		New().
		Get("/hello").
		Response(http.StatusOK, "hello world").
		Start()
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	
    // With rules
    srv, _ := mock.
        New().
        Get("/hello").
        When("cookie", "name", "equal", "Chocolate").
        And("header", "Authorization", "equal", "Bearer 123").
        Response(http.StatusOK, "hello world").
        Start()
    defer srv.Close()
    
    req, _ := http.NewRequest("GET", srv.URL, nil)
    client := &http.Client{}
    
    resp, err := client.Do(req)
}
```