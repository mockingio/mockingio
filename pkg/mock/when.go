package mock

import (
	"net/http/httptest"

	"github.com/mockingio/mockingio/engine/mock"
)

type When struct {
	builder *Builder
}

// Start starts the mock server
func (w *When) Start() (*httptest.Server, error) {
	return w.builder.Start()
}

// And is used to combine multiple rules with AND operator
func (w *When) And(target, modifier, operator, value string) *And {
	w.builder.response.RuleAggregation = mock.And
	w.builder.response.Rules = append(w.builder.response.Rules, mock.Rule{
		Target:   mock.Target(target),
		Modifier: modifier,
		Operator: mock.Operator(operator),
		Value:    value,
	})

	return &And{
		builder: w.builder,
	}
}

// Or is used to combine multiple rules with OR operator
func (w *When) Or(target, modifier, operator, value string) *Or {
	w.builder.response.RuleAggregation = mock.Or
	w.builder.response.Rules = append(w.builder.response.Rules, mock.Rule{
		Target:   mock.Target(target),
		Modifier: modifier,
		Operator: mock.Operator(operator),
		Value:    value,
	})

	return &Or{
		builder: w.builder,
	}
}
