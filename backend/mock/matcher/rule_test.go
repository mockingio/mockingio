package matcher_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cfg "github.com/smockyio/smocky/backend/config"
	"github.com/smockyio/smocky/backend/mock/matcher"
	"github.com/smockyio/smocky/backend/session"
)

func TestRuleMatcher_Match(t *testing.T) {
	sess := session.New()
	sess.Set("POST", "/api/person/detail", "request_number", 2)

	tests := []struct {
		name    string
		session *session.Session
		route   *cfg.Route
		request *http.Request
		rule    *cfg.Rule
		matched bool
		error   bool
	}{
		{
			"operator = Equal, found match",
			nil,
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 123", Operator: cfg.Equal},
			true,
			false,
		},
		{
			"operator = Regex, found match",
			nil,
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Body, Modifier: ".address.street", Value: "^[0-9]+.*", Operator: cfg.Regex},
			true,
			false,
		},
		{
			"operator = Regex, found not match",
			nil,
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Body, Modifier: ".address.street", Value: "^[a-z]+.*", Operator: cfg.Regex},
			false,
			false,
		},
		{
			"operator = Regex, invalid Regex",
			nil,
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Body, Modifier: ".address.street", Value: "^[a-z+.*", Operator: cfg.Regex},
			false,
			true,
		},
		{
			"test Equal, found not match",
			nil,
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 1231", Operator: cfg.Equal},
			false,
			false,
		},
		{
			"test number request, found match",
			sess,
			&cfg.Route{
				Request: "",
			},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.RequestNumber, Value: "2", Operator: cfg.Equal},
			true,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, err := matcher.NewRuleMatcher(tt.route, tt.rule, tt.request, tt.session).Match()
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.matched, matched)
		})
	}
}

func TestRuleMatcher_GetTargetValue(t *testing.T) {
	sess := session.New()
	sess.Set("POST", "/api/person/detail", "request_number", 2)

	var route = &cfg.Route{
		Request: "GET /api/:object/:action",
	}

	tests := []struct {
		request       *http.Request
		rule          *cfg.Rule
		expectedValue string
	}{
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Header, Modifier: "Authorization"}, "Bearer 123"},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Header, Modifier: "Random"}, ""},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Cookie, Modifier: "Token"}, "Token 123"},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Cookie, Modifier: "Random"}, ""},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.QueryString, Modifier: "name"}, "joe"},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.QueryString, Modifier: "Random"}, ""},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ".name"}, "joe"},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ""}, `{"name": "joe","address": { "street": "123 Road", "postcode": "2234" }}`},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ".address.postcode"}, "2234"},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ".address.random"}, ""},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.RouteParam, Modifier: "object"}, "person"},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.RouteParam, Modifier: "action"}, "detail"},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.RouteParam, Modifier: "random"}, ""},
		{newHTTPRequest(), &cfg.Rule{Target: cfg.RequestNumber}, "2"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Target: %v, Modifier: %v", tt.rule.Target, tt.rule.Modifier), func(t *testing.T) {
			actual, err := matcher.NewRuleMatcher(route, tt.rule, tt.request, sess).GetTargetValue()

			require.NoError(t, err)
			assert.Equal(t, tt.expectedValue, actual)
		})
	}
}

func newHTTPRequest() *http.Request {
	req, _ := http.NewRequest(
		"POST",
		"https://hi.com/api/person/detail?name=joe",
		strings.NewReader(`{"name": "joe","address": { "street": "123 Road", "postcode": "2234" }}`),
	)
	req.Header.Set("Authorization", "Bearer 123")
	req.AddCookie(&http.Cookie{Name: "Token", Value: "Token 123"})

	return req
}
