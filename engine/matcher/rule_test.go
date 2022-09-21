package matcher_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mockingio/mockingio/engine/database/memory"
	"github.com/mockingio/mockingio/engine/matcher"
	cfg "github.com/mockingio/mockingio/engine/mock"
)

func TestRuleMatcher_Match(t *testing.T) {
	sessionID := "123456"

	req := matcher.Context{
		HTTPRequest: newHTTPRequest(),
		SessionID:   sessionID,
	}

	tests := []struct {
		name    string
		route   *cfg.Route
		request *http.Request
		rule    *cfg.Rule
		matched bool
		error   bool
	}{
		{
			"operator = Equal, found match",
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 123", Operator: cfg.Equal},
			true,
			false,
		},
		{
			"operator = Regex, found match",
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Body, Modifier: ".address.street", Value: "^[0-9]+.*", Operator: cfg.Regex},
			true,
			false,
		},
		{
			"operator = Regex, found not match",
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Body, Modifier: ".address.street", Value: "^[a-z]+.*", Operator: cfg.Regex},
			false,
			false,
		},
		{
			"operator = Regex, invalid Regex",
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Body, Modifier: ".address.street", Value: "^[a-z+.*", Operator: cfg.Regex},
			false,
			true,
		},
		{
			"test Equal, found not match",
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 1231", Operator: cfg.Equal},
			false,
			false,
		},
		{
			"test number request, found match",
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.RequestNumber, Value: "2", Operator: cfg.Equal},
			true,
			false,
		},
		{
			"non exist operator, no match",
			&cfg.Route{},
			newHTTPRequest(),
			&cfg.Rule{Target: cfg.RequestNumber, Value: "2", Operator: cfg.Operator("random")},
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.New()
			_ = mem.Set(context.Background(), "", req.CountID(), "2")

			matched, err := matcher.NewRuleMatcher(&cfg.Mock{}, tt.route, tt.rule, matcher.Context{
				HTTPRequest: tt.request,
				SessionID:   sessionID,
			}, mem).Match()
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
	sessionID := "123456"

	req := matcher.Context{
		HTTPRequest: newHTTPRequest(),
		SessionID:   sessionID,
	}

	mem := memory.New()
	_ = mem.Set(context.Background(), "", req.CountID(), "2")

	var route = &cfg.Route{
		Path: "/api/:object/:action",
	}

	tests := []struct {
		route         *cfg.Route
		request       *http.Request
		rule          *cfg.Rule
		expectedValue string
	}{
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Header, Modifier: "Authorization"}, "Bearer 123"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Header, Modifier: "Random"}, ""},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Cookie, Modifier: "Token"}, "Token 123"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Cookie, Modifier: "Random"}, ""},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.QueryString, Modifier: "name"}, "joe"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.QueryString, Modifier: "Random"}, ""},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ".name"}, "joe"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ""}, `{"name": "joe","address": { "street": "123 Road", "postcode": "2234" }}`},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ".address.postcode"}, "2234"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Body, Modifier: ".address.random"}, ""},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.RouteParam, Modifier: "object"}, "person"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.RouteParam, Modifier: "action"}, "detail"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.RouteParam, Modifier: "random"}, ""},
		{&cfg.Route{Path: "/api/:object/:action/:something"}, newHTTPRequest(), &cfg.Rule{Target: cfg.RouteParam, Modifier: "random"}, ""},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.RequestNumber}, "2"},
		{route, newHTTPRequest(), &cfg.Rule{Target: cfg.Target("random target")}, ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Target: %v, Modifier: %v", tt.rule.Target, tt.rule.Modifier), func(t *testing.T) {
			actual, err := matcher.NewRuleMatcher(&cfg.Mock{}, tt.route, tt.rule, matcher.Context{
				HTTPRequest: tt.request,
				SessionID:   sessionID,
			}, mem).GetTargetValue()

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
