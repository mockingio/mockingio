package mock

import (
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
}

func (r Route) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Path, validation.Required),
		validation.Field(&r.ResponseMode, validation.In(DefaultResponse, ResponseRandomly, ResponseSequentially)),
		validation.Field(&r.Responses, validation.Required),
	)
}
