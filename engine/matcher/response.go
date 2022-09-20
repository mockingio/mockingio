package matcher

import (
	"github.com/pkg/errors"

	"github.com/mockingio/mockingio/engine/database"
	cfg "github.com/mockingio/mockingio/engine/mock"
)

func NewResponseMatcher(
	mok *cfg.Mock,
	route *cfg.Route,
	response *cfg.Response,
	req Context,
	db database.EngineDB,
) *ResponseMatcher {
	return &ResponseMatcher{
		route:    route,
		response: response,
		req:      req,
		db:       db,
		mock:     mok,
	}
}

type ResponseMatcher struct {
	route    *cfg.Route
	response *cfg.Response
	req      Context
	db       database.EngineDB
	mock     *cfg.Mock
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
		matched, err := NewRuleMatcher(r.mock, r.route, &rule, r.req, r.db).Match() // matcher. rule.Match(route, request)
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
