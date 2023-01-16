package mock

import (
	"net/http/httptest"

	"github.com/mockingio/mockingio/engine/mock"
)

type Response struct {
	builder *Builder
}

// Delay is a response delay in milliseconds
func (r *Response) Delay(min, max int) *Response {
	r.builder.response.Delay = mock.Delay{
		Min: min,
		Max: max,
	}
	return r
}

// Start starts the mock server
func (r *Response) Start() (*httptest.Server, error) {
	return r.builder.Start()
}

// When is a response rule.
// It can be used to match a request.
func (r *Response) When(target, modifier, operator, value string) *When {
	r.builder.response.RuleAggregation = mock.And
	r.builder.response.Rules = append(r.builder.response.Rules, mock.Rule{
		Target:   mock.Target(target),
		Modifier: modifier,
		Operator: mock.Operator(operator),
		Value:    value,
	})

	return &When{
		builder: r.builder,
	}
}

// WhenBodyEq is a response rule. It can be used to match a request body with the given value.
func (r *Response) WhenBodyEq(value string) *When {
	return r.When(Body, "", Equal, value)
}

// WhenPathInBodyEq is a response rule. It can be used to match a child body with the given value.
func (r *Response) WhenPathInBodyEq(field string, value string) *When {
	return r.When(Body, field, Equal, value)
}

// WhenHeaderEq is a response rule. It can be used to match a request header with the given value.
func (r *Response) WhenHeaderEq(headerName, headerValue string) *When {
	return r.When(Header, headerName, Equal, headerValue)
}

// WhenQueryStringEq is a response rule. It can be used to match a request query string with the given value.
func (r *Response) WhenQueryStringEq(queryStringName, queryStringValue string) *When {
	return r.When(QueryString, queryStringName, Equal, queryStringValue)
}

// WhenRouteParamEq is a response rule. It can be used to match a request route parameter with the given value.
func (r *Response) WhenRouteParamEq(routeParamName, routeParamValue string) *When {
	return r.When(Header, routeParamName, Equal, routeParamValue)
}

type And struct {
	builder *Builder
}

// Start starts the mock server
func (a *And) Start() (*httptest.Server, error) {
	return a.builder.Start()
}

// And is used to combine multiple rules with AND operator
func (a *And) And(target, modifier, operator, value string) *And {
	a.builder.response.Rules = append(a.builder.response.Rules, mock.Rule{
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

// Or is used to combine multiple rules with OR operator
func (o *Or) Or(target, modifier, operator, value string) *Or {
	o.builder.response.Rules = append(o.builder.response.Rules, mock.Rule{
		Target:   mock.Target(target),
		Modifier: modifier,
		Operator: mock.Operator(operator),
		Value:    value,
	})

	return o
}

// Start starts the mock server
func (o *Or) Start() (*httptest.Server, error) {
	return o.builder.Start()
}
