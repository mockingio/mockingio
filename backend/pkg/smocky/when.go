package smocky

import (
	"net/http/httptest"
	"testing"

	"github.com/smockyio/smocky/backend/mock/config"
)

type When struct {
	builder *Builder
}

func (w *When) Start(t *testing.T) *httptest.Server {
	return w.builder.Start(t)
}

func (w *When) And(target, modifier, operator, value string) *And {
	w.builder.response.RuleAggregation = config.And
	w.builder.response.Rules = append(w.builder.response.Rules, config.Rule{
		Target:   config.Target(target),
		Modifier: modifier,
		Operator: config.Operator(operator),
		Value:    value,
	})

	return &And{
		builder: w.builder,
	}
}

func (w *When) Or(target, modifier, operator, value string) *Or {
	w.builder.response.RuleAggregation = config.Or
	w.builder.response.Rules = append(w.builder.response.Rules, config.Rule{
		Target:   config.Target(target),
		Modifier: modifier,
		Operator: config.Operator(operator),
		Value:    value,
	})

	return &Or{
		builder: w.builder,
	}
}
