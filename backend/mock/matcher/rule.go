package matcher

import (
	"net/http"
	"regexp"

	"github.com/pkg/errors"

	cfg "github.com/smockyio/smocky/backend/mock/config"
	"github.com/smockyio/smocky/backend/session"
)

func NewRuleMatcher(route *cfg.Route, rule *cfg.Rule, httpRequest *http.Request, session *session.Session) *RuleMatcher {
	return &RuleMatcher{
		route:       route,
		rule:        rule,
		httpRequest: httpRequest,
		session:     session,
	}
}

type RuleMatcher struct {
	route       *cfg.Route
	rule        *cfg.Rule
	session     *session.Session
	httpRequest *http.Request
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
		return targetFn(r.route, r.httpRequest, r.rule.Modifier, r.session)
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
