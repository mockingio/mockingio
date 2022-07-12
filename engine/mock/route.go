package mock

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type responseMode string

const (
	ResponseSequentially responseMode = "sequential"
	ResponseRandomly     responseMode = "random"
	DefaultResponse      responseMode = ""
)

type Route struct {
	ID           string       `yaml:"id,omitempty" json:"id,omitempty"`
	Request      string       `yaml:"request" json:"request"`
	Description  string       `yaml:"description" json:"description"`
	ResponseMode responseMode `yaml:"response_mode,omitempty" json:"response_mode,omitempty"`
	Responses    []Response   `yaml:"responses" json:"responses"`
}

func (r Route) RequestParts() (string, string) {
	parts := strings.Split(r.Request, " ")
	return strings.ToUpper(strings.Trim(parts[0], " ")), strings.Trim(parts[1], " ")
}

func (r Route) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Request, validation.Required),
		validation.Field(&r.ResponseMode, validation.In(DefaultResponse, ResponseRandomly, ResponseSequentially)),
		validation.Field(&r.Responses, validation.Required),
	)
}
