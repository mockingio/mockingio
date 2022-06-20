package matcher

import (
	"regexp"

	"github.com/pkg/errors"

	cfg "github.com/smockyio/smocky/backend/mock/config"
)

func NewRuleMatcher(route *cfg.Route, rule *cfg.Rule, req Request) *RuleMatcher {
	return &RuleMatcher{
		route: route,
		rule:  rule,
		req:   req,
	}
}

type RuleMatcher struct {
	route *cfg.Route
	rule  *cfg.Rule
	req   Request
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
		return value == rule.Value, nil
	default:
		return false, nil
	}
}

func (r *RuleMatcher) GetTargetValue() (string, error) {
	if targetFn, ok := targets[r.rule.Target]; ok {
		return targetFn(r.route, r.rule.Modifier, r.req)
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
