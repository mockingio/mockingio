package matcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/pkg/errors"

	cfg "github.com/smockyio/smocky/backend/config"
	"github.com/smockyio/smocky/backend/session"
)

var targets = map[cfg.Target]getTargetValueFn{
	cfg.Header:        getValueFromHeader,
	cfg.Cookie:        getValueFromCookie,
	cfg.QueryString:   getValueFromQueryString,
	cfg.RequestNumber: getRequestNumber,
	cfg.RouteParam:    getValueFromRouteParam,
	cfg.Body:          getValueFromBody,
}

type getTargetValueFn func(route *cfg.Route, request *http.Request, modifier string, session *session.Session) (string, error)

func getValueFromHeader(_ *cfg.Route, request *http.Request, modifier string, _ *session.Session) (string, error) {

	return request.Header.Get(modifier), nil
}

func getValueFromCookie(_ *cfg.Route, request *http.Request, modifier string, _ *session.Session) (string, error) {
	cookies := request.Cookies()
	for _, c := range cookies {
		if c.Name == modifier {
			return c.Value, nil
		}
	}
	return "", nil
}

func getValueFromQueryString(_ *cfg.Route, request *http.Request, modifier string, _ *session.Session) (string, error) {
	return request.URL.Query().Get(modifier), nil
}

func getRequestNumber(_ *cfg.Route, request *http.Request, _ string, session *session.Session) (string, error) {
	return strconv.Itoa(session.GetRequestNumber(request)), nil
}

func getValueFromRouteParam(route *cfg.Route, request *http.Request, modifier string, _ *session.Session) (string, error) {
	_, path := route.RequestParts()
	templateParts := strings.Split(path, "/")
	actualParts := strings.Split(request.URL.Path, "/")
	if len(templateParts) != len(actualParts) {
		return "", nil
	}

	for i, templatePart := range templateParts {
		if p, ok := param(templatePart); ok {
			if p == modifier {
				return actualParts[i], nil
			}
		}
	}

	return "", nil
}

func getValueFromBody(_ *cfg.Route, request *http.Request, modifier string, _ *session.Session) (string, error) {
	if request.Body == nil {
		return "", nil
	}

	value, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", errors.Wrap(err, "read request body")
	}

	if string(value) == "" {
		return "", nil
	}

	if modifier == "" {
		return string(value), nil
	}

	input := map[string]interface{}{}
	if err := json.Unmarshal(value, &input); err != nil {
		return "", errors.Wrap(err, "unmarshal body")
	}

	query, err := gojq.Parse(modifier)
	if err != nil {
		return "", nil
	}

	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if v == nil {
			return "", nil
		}

		if err, ok := v.(error); ok {
			return "", errors.Wrapf(err, "unable to parse json query: %v", modifier)
		}

		return v.(string), nil // nolint
	}
	return "", nil
}
