package matcher

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/pkg/errors"

	"github.com/mockingio/mockingio/engine/database"
	cfg "github.com/mockingio/mockingio/engine/mock"
)

var targets = map[cfg.Target]getTargetValueFn{
	cfg.Header:        getValueFromHeader,
	cfg.Cookie:        getValueFromCookie,
	cfg.QueryString:   getValueFromQueryString,
	cfg.RequestNumber: getRequestNumber,
	cfg.RouteParam:    getValueFromRouteParam,
	cfg.Body:          getValueFromBody,
}

type getTargetValueFn func(mok *cfg.Mock, route *cfg.Route, modifier string, req Context, db database.EngineDB) (string, error)

func getValueFromHeader(_ *cfg.Mock, _ *cfg.Route, modifier string, req Context, _ database.EngineDB) (string, error) {
	return req.HTTPRequest.Header.Get(modifier), nil
}

func getValueFromCookie(_ *cfg.Mock, _ *cfg.Route, modifier string, req Context, _ database.EngineDB) (string, error) {
	cookies := req.HTTPRequest.Cookies()
	for _, c := range cookies {
		if c.Name == modifier {
			return c.Value, nil
		}
	}
	return "", nil
}

func getValueFromQueryString(_ *cfg.Mock, _ *cfg.Route, modifier string, req Context, _ database.EngineDB) (string, error) {
	return req.HTTPRequest.URL.Query().Get(modifier), nil
}

func getRequestNumber(mok *cfg.Mock, _ *cfg.Route, _ string, req Context, db database.EngineDB) (string, error) {
	value, err := db.GetInt(req.HTTPRequest.Context(), mok.ID, req.CountID())
	if err != nil {
		return "", err
	}
	return strconv.Itoa(value), nil
}

func getValueFromRouteParam(_ *cfg.Mock, route *cfg.Route, modifier string, req Context, _ database.EngineDB) (string, error) {
	templateParts := strings.Split(route.Path, "/")
	actualParts := strings.Split(req.HTTPRequest.URL.Path, "/")
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

func getValueFromBody(_ *cfg.Mock, _ *cfg.Route, modifier string, req Context, _ database.EngineDB) (string, error) {
	httpRequest := req.HTTPRequest
	if httpRequest.Body == nil {
		return "", nil
	}

	value, err := ioutil.ReadAll(httpRequest.Body)
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
