package mock

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	"github.com/mockingio/mockingio/engine"
	"github.com/mockingio/mockingio/engine/database/memory"
	"github.com/mockingio/mockingio/engine/mock"
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
		config: mock.New(),
	}
}

type Builder struct {
	response *mock.Response
	route    *mock.Route
	config   *mock.Mock
}

func (b *Builder) Start() (*httptest.Server, error) {
	b.clear()
	if err := b.config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}
	id := uuid.NewString()
	b.config.ID = id

	mem := memory.New()
	_ = mem.SetMock(context.Background(), b.config)
	_ = mem.SetActiveSession(context.Background(), id, "session-id")

	m := engine.New(id, mem)

	return httptest.NewServer(http.HandlerFunc(m.Handler)), nil
}

func (b *Builder) Post(url string) *Method {
	b.clear()
	b.route = &mock.Route{
		Method: "POST",
		Path:   url,
	}
	return &Method{
		builder: b,
	}
}

func (b *Builder) Get(url string) *Method {
	b.clear()
	b.route = &mock.Route{
		Method: "GET",
		Path:   url,
	}
	return &Method{
		builder: b,
	}
}

func (b *Builder) Put(url string) *Method {
	b.clear()
	b.route = &mock.Route{
		Method: "PUT",
		Path:   url,
	}
	return &Method{
		builder: b,
	}
}

func (b *Builder) Delete(url string) *Method {
	b.clear()
	b.route = &mock.Route{
		Method: "DELETE",
		Path:   url,
	}
	return &Method{
		builder: b,
	}
}

func (b *Builder) Option(url string) *Method {
	b.clear()
	b.route = &mock.Route{
		Method: "OPTIONS",
		Path:   url,
	}
	return &Method{
		builder: b,
	}
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
