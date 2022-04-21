package matcher

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/minio/pkg/wildcard"
	"github.com/pkg/errors"

	cfg "github.com/smockyio/smocky/backend/config"
	sess "github.com/smockyio/smocky/backend/session"
)

type RouteMatcher struct {
	route       *cfg.Route
	session     *sess.Session
	httpRequest *http.Request
}

func NewRouteMatcher(route *cfg.Route, session *sess.Session, req *http.Request) *RouteMatcher {
	return &RouteMatcher{
		route:       route,
		session:     session,
		httpRequest: req,
	}
}

func (r *RouteMatcher) Match() (*cfg.Response, error) {
	request := r.httpRequest
	method, path := r.route.RequestParts()

	if !strings.EqualFold(method, request.Method) {
		return nil, nil
	}

	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part != "" && string(part[0]) == ":" {
			parts[i] = "*"
		}
	}

	if !wildcard.Match(strings.Join(parts, "/"), request.URL.Path) {
		return nil, nil
	}

	r.session.IncreaseRequestNumber(r.httpRequest)

	responses, err := r.findMatches(request)
	if err != nil {
		return nil, errors.Wrap(err, "matching route")
	}

	return r.pickResponse(responses), nil
}

func (r *RouteMatcher) pickResponse(responses []*cfg.Response) *cfg.Response {
	if len(responses) == 0 {
		return nil
	}

	switch r.route.ResponseMode {
	case cfg.ResponseSequentially:
		idx := r.session.NextResponseIndex(r.httpRequest)
		if idx+1 == len(responses) {
			r.session.SetNextResponseIndex(r.httpRequest, 0)
		} else {
			r.session.SetNextResponseIndex(r.httpRequest, idx+1)
		}
		return responses[idx]
	case cfg.ResponseRandomly:
		rand.Seed(time.Now().UnixNano())
		return responses[rand.Intn(len(responses))]
	case cfg.DefaultResponse:
		fallthrough
	default:
		for _, response := range responses {
			if response.IsDefault {
				return response
			}
		}
		return responses[0] // No default setup, pick first one
	}
}

func (r *RouteMatcher) findMatches(request *http.Request) ([]*cfg.Response, error) {
	var responses []*cfg.Response

	for _, response := range r.route.Responses {
		response := response
		matched, err := NewResponseMatcher(r.route, &response, request, r.session).Match()
		if err != nil {
			return nil, err
		}

		if matched {
			responses = append(responses, &response)
		}
	}

	return responses, nil
}
