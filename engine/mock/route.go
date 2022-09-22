package mock

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
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
	Method       string       `yaml:"method" json:"method"`
	Path         string       `yaml:"path" json:"path"`
	Description  string       `yaml:"description" json:"description"`
	ResponseMode responseMode `yaml:"response_mode,omitempty" json:"response_mode,omitempty"`
	Responses    []Response   `yaml:"responses" json:"responses"`
	Disabled     bool         `yaml:"disabled,omitempty" json:"disabled,omitempty"`
}

func (r Route) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Method, validation.By(func(value interface{}) error {
			for _, method := range validMethods {
				if strings.EqualFold(value.(string), method) {
					return nil
				}
			}
			return errors.New("invalid request method")
		})),
		validation.Field(&r.Path, validation.Required),
		validation.Field(&r.ResponseMode, validation.In(DefaultResponse, ResponseRandomly, ResponseSequentially)),
		validation.Field(&r.Responses, validation.Required),
	)
}

var validMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodOptions,
}
