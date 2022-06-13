# smocky

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