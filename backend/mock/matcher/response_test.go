package matcher_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cfg "github.com/smockyio/smocky/backend/mock/config"
	"github.com/smockyio/smocky/backend/mock/matcher"
)

func TestResponseMatcher_Match(t *testing.T) {
	newRequest := func() *http.Request {
		request, _ := http.NewRequest("POST", "https://example.com/how/are/you", strings.NewReader(`{"name": "Joe"}`))
		request.Header.Add("Authorization", "Bearer 123")

		return request
	}

	tests := []struct {
		name      string
		response  *cfg.Response
		isMatched bool
	}{
		{
			name:      "rule is empty, return matched",
			response:  &cfg.Response{},
			isMatched: true,
		},
		{
			name: "aggregation is And, found response, all rule matched",
			response: &cfg.Response{
				RuleAggregation: cfg.And,
				Rules: []cfg.Rule{
					{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 123", Operator: cfg.Equal},
					{Target: cfg.Body, Modifier: ".name", Value: "Joe", Operator: cfg.Equal},
				},
			},
			isMatched: true,
		},
		{
			name: "aggregation is And, found no response, not all rule matched",
			response: &cfg.Response{
				RuleAggregation: cfg.And,
				Rules: []cfg.Rule{
					{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 123", Operator: cfg.Equal},
					{Target: cfg.Body, Modifier: ".name", Value: "random name", Operator: cfg.Equal},
				},
			},
			isMatched: false,
		},
		{
			name: "aggregation is Or, found response, one of rule matched",
			response: &cfg.Response{
				RuleAggregation: cfg.Or,
				Rules: []cfg.Rule{
					{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 123", Operator: cfg.Equal},
					{Target: cfg.Body, Modifier: ".name", Value: "random name", Operator: cfg.Equal},
				},
			},
			isMatched: true,
		},
		{
			name: "aggregation is Or, found no response, no rule matched",
			response: &cfg.Response{
				RuleAggregation: cfg.And,
				Rules: []cfg.Rule{
					{Target: cfg.Header, Modifier: "Authorization", Value: "random name", Operator: cfg.Equal},
					{Target: cfg.Body, Modifier: ".name", Value: "random name", Operator: cfg.Equal},
				},
			},
			isMatched: false,
		},
		{
			name: "aggregation is empty, found response, all rule matched",
			response: &cfg.Response{
				RuleAggregation: "",
				Rules: []cfg.Rule{
					{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 123", Operator: cfg.Equal},
					{Target: cfg.Body, Modifier: ".name", Value: "Joe", Operator: cfg.Equal},
				},
			},
			isMatched: true,
		},
		{
			name: "aggregation is empty, found no response, not all rule matched",
			response: &cfg.Response{
				RuleAggregation: "",
				Rules: []cfg.Rule{
					{Target: cfg.Header, Modifier: "Authorization", Value: "Bearer 123", Operator: cfg.Equal},
					{Target: cfg.Body, Modifier: ".name", Value: "random name", Operator: cfg.Equal},
				},
			},
			isMatched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isMatched, err := matcher.NewResponseMatcher(nil, tt.response, matcher.Request{
				HTTPRequest: newRequest(),
			}).Match()
			require.NoError(t, err)
			assert.Equal(t, tt.isMatched, isMatched)
		})
	}
}
