package matcher

import (
	"github.com/smockyio/smocky/backend/persistent"
	"math/rand"
	"strings"
	"time"

	"github.com/minio/pkg/wildcard"
	"github.com/pkg/errors"

	cfg "github.com/smockyio/smocky/backend/mock/config"
)

type RouteMatcher struct {
	route *cfg.Route
	req   Context
}

func NewRouteMatcher(route *cfg.Route, req Context) *RouteMatcher {
	return &RouteMatcher{
		route: route,
		req:   req,
	}
}

func (r *RouteMatcher) Match() (*cfg.Response, error) {
	method, path := r.route.RequestParts()
	httpRequest := r.req.HTTPRequest
	db := persistent.GetDefault()

	if !strings.EqualFold(method, httpRequest.Method) {
		return nil, nil
	}

	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part != "" && string(part[0]) == ":" {
			parts[i] = "*"
		}
	}

	if !wildcard.Match(strings.Join(parts, "/"), httpRequest.URL.Path) {
		return nil, nil
	}

	_, err := db.Increment(
		httpRequest.Context(),
		r.req.CountID(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "increase request times")
	}

	responses, err := r.findMatches()
	if err != nil {
		return nil, errors.Wrap(err, "matching route")
	}

	return r.pickResponse(responses)
}

func (r *RouteMatcher) pickResponse(responses []*cfg.Response) (*cfg.Response, error) {
	if len(responses) == 0 {
		return nil, nil
	}
	db := persistent.GetDefault()
	sequenceID := r.req.SequenceID()
	ctx := r.req.HTTPRequest.Context()

	switch r.route.ResponseMode {
	case cfg.ResponseSequentially:
		idx, err := db.GetInt(ctx, sequenceID)
		if err != nil {
			return nil, err
		}

		if idx+1 == len(responses) {
			if err := db.Set(ctx, sequenceID, 0); err != nil {
				return nil, err
			}
		} else {
			if err := db.Set(ctx, sequenceID, idx+1); err != nil {
				return nil, err
			}
		}

		return responses[idx], nil
	case cfg.ResponseRandomly:
		rand.Seed(time.Now().UnixNano())
		return responses[rand.Intn(len(responses))], nil
	case cfg.DefaultResponse:
		fallthrough
	default:
		for _, response := range responses {
			if response.IsDefault {
				return response, nil
			}
		}
		return responses[0], nil // No default setup, pick first one
	}
}

func (r *RouteMatcher) findMatches() ([]*cfg.Response, error) {
	var responses []*cfg.Response

	for _, response := range r.route.Responses {
		response := response
		matched, err := NewResponseMatcher(r.route, &response, r.req).Match()
		if err != nil {
			return nil, err
		}

		if matched {
			responses = append(responses, &response)
		}
	}

	return responses, nil
}
