package mock

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type responseMode string

const (
	ResponseSequentially responseMode = "sequential"
	ResponseRandomly     responseMode = "random"
	DefaultResponse      responseMode = ""
)

type Route struct {
	ID           string       `yaml:"id,omitempty" json:"id,omitempty"`
	Method       string       `yaml:"method" json:"method"`
	Path         string       `yaml:"path" json:"path"`
	Description  string       `yaml:"description" json:"description"`
	ResponseMode responseMode `yaml:"response_mode,omitempty" json:"response_mode,omitempty"`
	Responses    []Response   `yaml:"responses" json:"responses"`
	Disabled     bool         `yaml:"disabled,omitempty" json:"disabled,omitempty"`
}

func (r Route) Clone() *Route {
	result := r
	result.ID = newID()
	result.ResponseMode = r.ResponseMode
	result.Responses = make([]Response, len(r.Responses))
	for i, response := range r.Responses {
		result.Responses[i] = response.Clone()
	}

	return &result
}

func (r Route) PatchString(data string) (*Route, error) {
	var values map[string]*json.RawMessage
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		return nil, err
	}

	var newRoute = &Route{}
	if err := patchStruct(newRoute, values); err != nil {
		return nil, err
	}

	return newRoute, nil
}

func (r Route) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Path, validation.Required),
		validation.Field(&r.ResponseMode, validation.In(DefaultResponse, ResponseRandomly, ResponseSequentially)),
		validation.Field(&r.Responses, validation.Required),
	)
}
