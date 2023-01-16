package matcher

import (
	"encoding/json"
	"regexp"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"github.com/mockingio/mockingio/engine/database"
	cfg "github.com/mockingio/mockingio/engine/mock"
)

func NewRuleMatcher(mok *cfg.Mock, route *cfg.Route, rule *cfg.Rule, req Context, db database.EngineDB) *RuleMatcher {
	return &RuleMatcher{
		route: route,
		rule:  rule,
		req:   req,
		db:    db,
		mock:  mok,
	}
}

type RuleMatcher struct {
	route *cfg.Route
	rule  *cfg.Rule
	req   Context
	db    database.EngineDB
	mock  *cfg.Mock
}

func (r *RuleMatcher) Match() (bool, error) {
	value, err := r.GetTargetValue()
	if err != nil {
		return false, errors.Wrap(err, "get target value")
	}

	rule := r.rule

	switch rule.Operator {
	case cfg.Regex:
		matched, err := regexp.MatchString(rule.Value, value)
		if err != nil {
			return false, errors.Wrap(err, "regex match string")
		}
		return matched, nil
	case cfg.Equal:
		// special treatment for target is Body, with JSON. We'll need to compare json
		if r.rule.Target == cfg.Body && r.req.HTTPRequest.Header.Get("Content-Type") == "application/json" {
			return matchJSON(value, rule.Value), nil
		}
		return value == rule.Value, nil
	default:
		return false, nil
	}
}

func (r *RuleMatcher) GetTargetValue() (string, error) {
	if targetFn, ok := targets[r.rule.Target]; ok {
		return targetFn(r.mock, r.route, r.rule.Modifier, r.req, r.db)
	}

	return "", nil
}

func param(p string) (string, bool) {
	if p == "" {
		return "", false
	}

	if string(p[0]) == ":" {
		return p[1:], true
	}

	return "", false
}

func matchJSON(actual, expected string) bool {
	var actualJSON, expectedJSON interface{}

	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		return false
	}

	return cmp.Equal(actualJSON, expectedJSON)
}
