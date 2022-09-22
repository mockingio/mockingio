package mock

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Proxy struct {
	Enabled            bool              `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Host               string            `yaml:"host,omitempty" json:"host"`
	RequestHeaders     map[string]string `yaml:"request_headers,omitempty" json:"request_headers,omitempty"`
	ResponseHeaders    map[string]string `yaml:"response_headers,omitempty" json:"response_headers,omitempty"`
	InsecureSkipVerify bool              `yaml:"insecure_skip_verify,omitempty" json:"insecure_skip_verify,omitempty"`
}

func (p Proxy) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Host, validation.Required),
	)
}
