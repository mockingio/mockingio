package mock

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Target string
type Operator string

const (
	Header        Target = "header"
	Body          Target = "body"
	QueryString   Target = "query_string"
	Cookie        Target = "cookie"
	RouteParam    Target = "route_param"
	RequestNumber Target = "request_number"
)

const (
	Equal Operator = "equal"
	Regex Operator = "regex"
)

type Rule struct {
	ID       string   `yaml:"id,omitempty" json:"id,omitempty"`
	Target   Target   `yaml:"target" json:"target"`
	Modifier string   `yaml:"modifier" json:"modifier,omitempty"`
	Value    string   `yaml:"value" json:"value"`
	Operator Operator `yaml:"operator" json:"operator"`
}

func (r Rule) Clone() Rule {
	result := r
	result.ID = newID()

	return result
}

func (r Rule) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Target, validation.Required, validation.In(Body, QueryString, Header, Cookie, RouteParam, RequestNumber)),
		validation.Field(&r.Value, validation.Required),
		validation.Field(&r.Operator, validation.Required, validation.In(Equal, Regex)),
	)
}
