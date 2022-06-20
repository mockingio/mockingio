package smocky

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smockyio/smocky/backend/mock"
	"github.com/smockyio/smocky/backend/mock/config"
	"github.com/smockyio/smocky/backend/persistent/memory"
)

const (
	Header        = "header"
	Body          = "body"
	QueryString   = "query_string"
	Cookie        = "cookie"
	RouteParam    = "route_param"
	RequestNumber = "request_number"
)

const (
	Equal = "equal"
	Regex = "regex"
)

type Headers map[string]string

func New() *Builder {
	return &Builder{
		config: &config.Config{},
	}
}

type Builder struct {
	response *config.Response
	route    *config.Route
	config   *config.Config
}

func (b *Builder) Start(t *testing.T) *httptest.Server {
	b.clear()
	if err := b.config.Validate(); err != nil {
		t.Errorf("invalid config: %v", err)
	}
	id := uuid.NewString()

	mem := memory.New()
	_ = mem.SetConfig(context.Background(), id, b.config)

	m, err := mock.New(id, "id", mem)
	if err != nil {
		t.Errorf("fail to create mock: %v", err)
	}

	return httptest.NewServer(http.HandlerFunc(m.Handler))
}

func (b *Builder) Port(port string) *Builder {
	b.config.Port = port
	return b
}

func (b *Builder) Post(url string) *Builder {
	b.clear()
	b.route = &config.Route{
		Request: "POST " + url,
	}
	return b
}

func (b *Builder) Get(url string) *Builder {
	b.clear()
	b.route = &config.Route{
		Request: "GET " + url,
	}
	return b
}

func (b *Builder) Put(url string) *Builder {
	b.clear()
	b.route = &config.Route{
		Request: "PUT " + url,
	}
	return b
}

func (b *Builder) Delete(url string) *Builder {
	b.clear()
	b.route = &config.Route{
		Request: "DELETE " + url,
	}
	return b
}

func (b *Builder) Option(url string) *Builder {
	b.clear()
	b.route = &config.Route{
		Request: "OPTION " + url,
	}
	return b
}

func (b *Builder) Response(status int, body string, headers ...Headers) *Response {
	if len(headers) == 0 {
		headers = append(headers, map[string]string{})
	}

	b.response = &config.Response{
		Body:    body,
		Status:  status,
		Headers: headers[0],
	}

	resp := &Response{
		builder: b,
	}

	return resp
}

func (b *Builder) clear() {
	if b.response != nil {
		b.route.Responses = append(b.route.Responses, *b.response)
		b.response = nil
	}

	if b.route != nil {
		b.config.Routes = append(b.config.Routes, b.route)
		b.route = nil
	}
}
