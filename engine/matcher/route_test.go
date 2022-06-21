package matcher_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smockyio/smocky/engine/matcher"
	cfg "github.com/smockyio/smocky/engine/mock"
	"github.com/smockyio/smocky/engine/persistent"
	"github.com/smockyio/smocky/engine/persistent/memory"
)

func TestRouteMatcher_Match(t *testing.T) {
	singleRuleResponse := cfg.Response{
		Status:          http.StatusOK,
		RuleAggregation: cfg.And,
		Rules: []cfg.Rule{{
			Target:   cfg.Target("header"),
			Modifier: "Authorization",
			Value:    "Bearer",
			Operator: cfg.Operator("equal"),
		}},
	}

	multiANDRulesResponse := cfg.Response{
		Status:          http.StatusOK,
		RuleAggregation: cfg.And,
		Rules: []cfg.Rule{{
			Target:   cfg.Target("header"),
			Modifier: "Authorization",
			Value:    "Bearer",
			Operator: cfg.Operator("equal"),
		}, {
			Target:   cfg.Target("body"),
			Modifier: ".gender",
			Value:    "female",
			Operator: cfg.Operator("equal"),
		}},
	}

	multiORRulesResponse := cfg.Response{
		Status:          http.StatusOK,
		RuleAggregation: cfg.Or,
		Rules: []cfg.Rule{
			{
				Target:   cfg.Target("header"),
				Modifier: "Authorization",
				Value:    "Bearer",
				Operator: cfg.Operator("equal"),
			},
			{
				Target:   cfg.Target("body"),
				Modifier: ".gender",
				Value:    "female",
				Operator: cfg.Operator("equal"),
			},
		},
	}

	multiORRulesNotMatchedResponse := cfg.Response{
		Status:          http.StatusOK,
		RuleAggregation: cfg.Or,
		Rules: []cfg.Rule{
			{
				Target:   cfg.Target("header"),
				Modifier: "Authorization",
				Value:    "Bearer 123",
				Operator: cfg.Operator("equal"),
			},
			{
				Target:   cfg.Target("body"),
				Modifier: ".gender",
				Value:    "n/a",
				Operator: cfg.Operator("equal"),
			},
		},
	}

	httpGetReq, _ := http.NewRequest("GET", "", nil)
	httpPostReq, _ := http.NewRequest("POST", "https://example.com/how/are/you", nil)

	httpPostReqWithHeaderBody, _ := http.NewRequest("POST", "https://example.com/how/are/you", strings.NewReader(`{"gender":"female"}`))
	httpPostReqWithHeaderBody.Header.Add("Authorization", "Bearer")

	httpPostReqWithHeader, _ := http.NewRequest("POST", "https://example.com/how/are/you", nil)
	httpPostReqWithHeader.Header.Add("Authorization", "Bearer")

	tests := []struct {
		name             string
		httpReq          *http.Request
		route            *cfg.Route
		expectedResponse *cfg.Response
		expectedError    bool
	}{
		{
			"single rule matched, response returned",
			httpPostReqWithHeaderBody,
			&cfg.Route{Request: "POST /how/are/you", Responses: []cfg.Response{singleRuleResponse}},
			&singleRuleResponse,
			false,
		},
		{
			"URL wildcard matched, response returned",
			httpPostReqWithHeaderBody,
			&cfg.Route{Request: "POST /how/are/*", Responses: []cfg.Response{singleRuleResponse}},
			&singleRuleResponse,
			false,
		},
		{
			"multiple rule matched, response returned",
			httpPostReqWithHeaderBody,
			&cfg.Route{Request: "POST /how/are/you", Responses: []cfg.Response{multiANDRulesResponse}},
			&multiANDRulesResponse,
			false,
		},
		{
			"one of rules matched, response returned",
			httpPostReqWithHeaderBody,
			&cfg.Route{Request: "POST /how/are/you", Responses: []cfg.Response{multiORRulesResponse}},
			&multiORRulesResponse,
			false,
		},
		{
			"none of rules matched, no response returned",
			httpPostReqWithHeaderBody,
			&cfg.Route{Request: "POST /how/are/you", Responses: []cfg.Response{multiORRulesNotMatchedResponse}},
			nil,
			false,
		},
		{
			"not all rules matched, no response returned",
			httpPostReqWithHeader,
			&cfg.Route{Request: "POST /how/are/you", Responses: []cfg.Response{multiANDRulesResponse}},
			nil,
			false,
		},
		{
			"method not matched, no response returned",
			httpGetReq,
			&cfg.Route{Request: "POST /", Responses: []cfg.Response{singleRuleResponse}},
			nil,
			false,
		},
		{
			"url not matched, no response returned",
			httpPostReq,
			&cfg.Route{Request: "POST /", Responses: []cfg.Response{singleRuleResponse}},
			nil,
			false,
		},
		{
			"no rule responses, no response returned",
			httpPostReq,
			&cfg.Route{Request: "POST /", Responses: []cfg.Response{}},
			nil,
			false,
		},
		{
			"no rules matched, no response returned",
			httpPostReq,
			&cfg.Route{Request: "POST /how/are/you", Responses: []cfg.Response{singleRuleResponse}},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			persistent.New(memory.New())
			result, err := matcher.NewRouteMatcher(tt.route, matcher.Context{
				HTTPRequest: tt.httpReq,
			}).Match()
			assert.Equal(t, tt.expectedResponse, result)
			assert.Equal(t, tt.expectedError, err != nil)
		})
	}
}

func TestRouteMatcher_ResponseStrategy(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://example.com/how/are/you", nil)
	request.Header.Add("Authorization", "Bearer")

	response1 := cfg.Response{
		Status:          http.StatusOK,
		RuleAggregation: cfg.And,
		Rules: []cfg.Rule{{
			Target:   cfg.Target("header"),
			Modifier: "Authorization",
			Value:    "Bearer",
			Operator: cfg.Operator("equal"),
		}},
	}

	response2 := cfg.Response{
		Status:          http.StatusNotFound,
		RuleAggregation: cfg.And,
		Rules: []cfg.Rule{{
			Target:   cfg.Target("header"),
			Modifier: "Authorization",
			Value:    "Bearer",
			Operator: cfg.Operator("equal"),
		}},
	}

	response3 := cfg.Response{
		Status:          http.StatusInternalServerError,
		RuleAggregation: cfg.And,
		Rules: []cfg.Rule{{
			Target:   cfg.Target("header"),
			Modifier: "Authorization",
			Value:    "Bearer",
			Operator: cfg.Operator("equal"),
		}},
	}

	t.Run("multi requests until matched", func(t *testing.T) {
		route := &cfg.Route{
			Request: "GET /how/are/you",
			Responses: []cfg.Response{{
				Status: http.StatusOK,
				Rules: []cfg.Rule{{
					Target:   cfg.Target("request_number"),
					Value:    "3",
					Operator: cfg.Operator("equal"),
				}},
			}},
		}

		m := matcher.NewRouteMatcher(route, matcher.Context{
			HTTPRequest: request,
		})
		res, _ := m.Match()
		assert.Nil(t, res)
	})

	t.Run("no response strategy setup", func(t *testing.T) {
		route := &cfg.Route{
			Request:   "GET /how/are/you",
			Responses: []cfg.Response{response1, response2, response3},
		}

		result, err := matcher.NewRouteMatcher(route, matcher.Context{
			HTTPRequest: request,
		}).Match()

		require.NoError(t, err)
		assert.Equal(t, &response1, result)
	})

	t.Run("default strategy setup", func(t *testing.T) {
		defaultResponse := cfg.Response{
			Status:          http.StatusAccepted,
			RuleAggregation: cfg.And,
			IsDefault:       true,
			Rules: []cfg.Rule{{
				Target:   cfg.Target("header"),
				Modifier: "Authorization",
				Value:    "Bearer",
				Operator: cfg.Operator("equal"),
			}},
		}

		route := &cfg.Route{
			Request:   "GET /how/are/you",
			Responses: []cfg.Response{response1, response2, defaultResponse, response3},
		}

		result, err := matcher.NewRouteMatcher(route, matcher.Context{
			HTTPRequest: request,
		}).Match()
		require.NoError(t, err)
		assert.Equal(t, &defaultResponse, result)
	})

	t.Run("sequential strategy setup", func(t *testing.T) {
		route := &cfg.Route{
			Request:      "GET /how/are/you",
			ResponseMode: cfg.ResponseSequentially,
			Responses:    []cfg.Response{response1, response2, response3},
		}

		result1, err := matcher.NewRouteMatcher(route, matcher.Context{
			HTTPRequest: request,
		}).Match()
		require.NoError(t, err)
		assert.Equal(t, &response1, result1)

		result2, err := matcher.NewRouteMatcher(route, matcher.Context{
			HTTPRequest: request,
		}).Match()
		require.NoError(t, err)
		assert.Equal(t, &response2, result2)

		result3, err := matcher.NewRouteMatcher(route, matcher.Context{
			HTTPRequest: request,
		}).Match()
		require.NoError(t, err)
		assert.Equal(t, &response3, result3)
	})
}
