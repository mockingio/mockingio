package matcher

import (
	"github.com/pkg/errors"

	cfg "github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
)

func NewResponseMatcher(
	route *cfg.Route,
	response *cfg.Response,
	req Context,
	db persistent.EngineDB,
) *ResponseMatcher {
	return &ResponseMatcher{
		route:    route,
		response: response,
		req:      req,
		db:       db,
	}
}

type ResponseMatcher struct {
	route    *cfg.Route
	response *cfg.Response
	req      Context
	db       persistent.EngineDB
}

func (r *ResponseMatcher) Match() (bool, error) {
	if len(r.response.Rules) == 0 {
		return true, nil
	}

	aggregation := r.response.RuleAggregation
	if aggregation == "" {
		aggregation = cfg.And
	}

	for _, rule := range r.response.Rules {
		matched, err := NewRuleMatcher(r.route, &rule, r.req, r.db).Match() // matcher. rule.Match(route, request)
		if err != nil {
			return false, errors.Wrap(err, "matching rule")
		}

		if !matched && aggregation == cfg.And {
			return false, nil
		}

		if matched && aggregation == cfg.Or {
			return true, nil
		}
	}

	// Match all rules
	if aggregation == cfg.And {
		return true, nil
	}

	return false, nil
}
