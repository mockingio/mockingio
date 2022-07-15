package smocky

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tuongaz/smocky-engine/engine"
	"github.com/tuongaz/smocky-engine/engine/mock"
	"github.com/tuongaz/smocky-engine/engine/persistent"
	"github.com/tuongaz/smocky-engine/engine/persistent/memory"
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
		config: &mock.Mock{},
	}
}

type Builder struct {
	response *mock.Response
	route    *mock.Route
	config   *mock.Mock
}

func (b *Builder) Start(t *testing.T) *httptest.Server {
	b.clear()
	if err := b.config.Validate(); err != nil {
		t.Errorf("invalid config: %v", err)
	}
	id := uuid.NewString()
	b.config.ID = id

	mem := memory.New()
	persistent.New(mem)
	_ = mem.SetMock(context.Background(), b.config)
	_ = mem.SetActiveSession(context.Background(), id, "session-id")

	m := engine.New(id)

	return httptest.NewServer(http.HandlerFunc(m.Handler))
}

func (b *Builder) Port(port string) *Builder {
	b.config.Port = port
	return b
}

func (b *Builder) Post(url string) *Builder {
	b.clear()
	b.route = &mock.Route{
		Request: "POST " + url,
	}
	return b
}

func (b *Builder) Get(url string) *Builder {
	b.clear()
	b.route = &mock.Route{
		Request: "GET " + url,
	}
	return b
}

func (b *Builder) Put(url string) *Builder {
	b.clear()
	b.route = &mock.Route{
		Request: "PUT " + url,
	}
	return b
}

func (b *Builder) Delete(url string) *Builder {
	b.clear()
	b.route = &mock.Route{
		Request: "DELETE " + url,
	}
	return b
}

func (b *Builder) Option(url string) *Builder {
	b.clear()
	b.route = &mock.Route{
		Request: "OPTION " + url,
	}
	return b
}

func (b *Builder) Response(status int, body string, headers ...Headers) *Response {
	if len(headers) == 0 {
		headers = append(headers, map[string]string{})
	}

	b.response = &mock.Response{
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
