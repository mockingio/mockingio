package mock

import (
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/mockingio/mockingio/engine/mock"
)

type Response struct {
	builder *Builder
}

func (r *Response) Delay(delay int64) *Response {
	r.builder.response.Delay = delay
	return r
}

func (r *Response) Start() (*httptest.Server, error) {
	return r.builder.Start()
}

func (r *Response) When(target, modifier, operator, value string) *When {
	r.builder.response.RuleAggregation = mock.And
	r.builder.response.Rules = append(r.builder.response.Rules, mock.Rule{
		ID:       uuid.NewString(),
		Target:   mock.Target(target),
		Modifier: modifier,
		Operator: mock.Operator(operator),
		Value:    value,
	})

	return &When{
		builder: r.builder,
	}
}

type And struct {
	builder *Builder
}

func (a *And) Start() (*httptest.Server, error) {
	return a.builder.Start()
}

func (a *And) And(target, modifier, operator, value string) *And {
	a.builder.response.Rules = append(a.builder.response.Rules, mock.Rule{
		ID:       uuid.NewString(),
		Target:   mock.Target(target),
		Modifier: modifier,
		Operator: mock.Operator(operator),
		Value:    value,
	})

	return a
}

type Or struct {
	builder *Builder
}

func (o *Or) Or(target, modifier, operator, value string) *Or {
	o.builder.response.Rules = append(o.builder.response.Rules, mock.Rule{
		ID:       uuid.NewString(),
		Target:   mock.Target(target),
		Modifier: modifier,
		Operator: mock.Operator(operator),
		Value:    value,
	})

	return o
}

func (o *Or) Start() (*httptest.Server, error) {
	return o.builder.Start()
}
