package matcher

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/minio/pkg/wildcard"
	"github.com/pkg/errors"

	"github.com/mockingio/mockingio/engine/database"
	cfg "github.com/mockingio/mockingio/engine/mock"
)

type RouteMatcher struct {
	route *cfg.Route
	req   Context
	db    database.EngineDB
}

func NewRouteMatcher(route *cfg.Route, req Context, db database.EngineDB) *RouteMatcher {
	return &RouteMatcher{
		route: route,
		req:   req,
		db:    db,
	}
}

func (r *RouteMatcher) Match() (*cfg.Response, error) {
	if r.route.Disabled {
		return nil, nil
	}

	httpRequest := r.req.HTTPRequest
	method := r.route.Method
	if r.route.Method == "" {
		method = http.MethodGet
	}

	if !strings.EqualFold(method, httpRequest.Method) {
		return nil, nil
	}

	wildcardPath := toWildcardPath(r.route.Path)
	if !wildcard.Match(wildcardPath, httpRequest.URL.Path) {
		return nil, nil
	}

	_, err := r.db.Increment(
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
	sequenceID := r.req.SequenceID()
	ctx := r.req.HTTPRequest.Context()

	switch r.route.ResponseMode {
	case cfg.ResponseSequentially:
		idx, err := r.db.GetInt(ctx, sequenceID)
		if err != nil {
			return nil, err
		}

		if idx+1 == len(responses) {
			if err := r.db.Set(ctx, sequenceID, 0); err != nil {
				return nil, err
			}
		} else {
			if err := r.db.Set(ctx, sequenceID, idx+1); err != nil {
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
		matched, err := NewResponseMatcher(r.route, &response, r.req, r.db).Match()
		if err != nil {
			return nil, err
		}

		if matched {
			responses = append(responses, &response)
		}
	}

	return responses, nil
}

func toWildcardPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part != "" && string(part[0]) == ":" {
			parts[i] = "*"
		}
	}

	return strings.Join(parts, "/")
}
