package mock

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Mock struct {
	ID     string   `yaml:"id,omitempty" json:"id,omitempty"`
	Name   string   `yaml:"name,omitempty" json:"name,omitempty"`
	Port   string   `yaml:"port,omitempty" json:"port,omitempty"`
	Routes []*Route `yaml:"routes,omitempty" json:"routes,omitempty"`
	Proxy  *Proxy   `yaml:"proxy,omitempty" json:"proxy,omitempty"`
	// all OPTIONS calls are responded with success if AutoCORS is true
	AutoCORS bool `yaml:"auto_cors,omitempty" json:"auto_cors,omitempty"`
	TLS      *TLS `yaml:"tls,omitempty" json:"tls,omitempty"`
	options  mockOptions
	FilePath string `yaml:"-" json:"-"`
}

func New(opts ...Option) *Mock {
	m := &Mock{
		options: mockOptions{},
	}

	for _, opt := range opts {
		opt(&m.options)
	}

	return m
}

func FromFile(file string, opts ...Option) (*Mock, error) {
	// TODO: Detects file type, support JSON
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read mock file")
	}

	mok, err := FromYaml(string(data), opts...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read the mock file")
	}

	mok.FilePath = file

	return mok, err
}

func FromYaml(text string, opts ...Option) (*Mock, error) {
	decoder := yaml.NewDecoder(strings.NewReader(text))
	m := New(opts...)
	if err := decoder.Decode(m); err != nil {
		return nil, errors.Wrap(err, "decode yaml to mock")
	}
	m.ApplyDefault()
	if m.options.idGeneration {
		addIDs(m)
	}

	if err := m.Validate(); err != nil {
		return nil, errors.Wrap(err, "mock validation")
	}

	return m, nil
}

func (m Mock) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.ID, validation.Length(0, 100)),
		validation.Field(&m.Name, validation.Length(0, 255)),
		validation.Field(&m.Port, is.Port),
		validation.Field(&m.Routes, validation.Required),
	)
}

func (m Mock) ProxyEnabled() bool {
	return m.Proxy != nil && m.Proxy.Enabled
}

func (m Mock) TLSEnabled() bool {
	return m.TLS != nil && m.TLS.Enabled
}

func (m Mock) JSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", errors.Wrap(err, "marshal mock to json")
	}

	return string(data), nil
}

func (m Mock) ApplyDefault() Mock {
	for _, r := range m.Routes {
		if r.Method == "" {
			r.Method = http.MethodGet
		}

		r.Path = "/" + strings.TrimPrefix(r.Path, "/")

		for i, res := range r.Responses {
			if res.Status == 0 {
				res.Status = 200
			}
			r.Responses[i] = res
		}
	}

	return m
}

// addIDs Add ids for mock and routes, responses and rules
func addIDs(m *Mock) {
	if m.ID == "" {
		m.ID = newID()
	}
	for _, r := range m.Routes {
		if r.ID == "" {
			r.ID = newID()
		}

		for i, res := range r.Responses {
			if res.ID == "" {
				res.ID = newID()
				r.Responses[i] = res
			}

			for j, rule := range res.Rules {
				if rule.ID == "" {
					rule.ID = newID()
					res.Rules[j] = rule
				}
			}
		}
	}
}

func newID() string {
	return uuid.NewString()
}
