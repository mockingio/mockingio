package mock

import (
	"io/ioutil"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Mock struct {
	ID     string   `yaml:"id,omitempty" json:"id,omitempty"`
	Name   string   `yaml:"name,omitempty" json:"name,omitempty"`
	Port   string   `yaml:"port,omitempty" json:"port,omitempty"`
	Routes []*Route `yaml:"routes,omitempty" json:"routes,omitempty"`
}

func New() *Mock {
	return &Mock{
		ID: newID(),
	}
}

func FromFile(file string) (*Mock, error) {
	// TODO: check the file extension, support loading mock from JSON
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read mock file")
	}

	return FromYaml(string(data))
}

func FromYaml(text string) (*Mock, error) {
	decoder := yaml.NewDecoder(strings.NewReader(text))
	m := &Mock{}
	if err := decoder.Decode(m); err != nil {
		return nil, errors.Wrap(err, "decode yaml to mock")
	}
	addIds(m)
	defaultVals(m)

	return m, nil
}

func defaultVals(m *Mock) {
	for _, r := range m.Routes {
		for i, res := range r.Responses {
			if res.Status == 0 {
				res.Status = 200
			}
			r.Responses[i] = res
		}
	}
}

// Add ids for mock and routes, responses and rules
func addIds(m *Mock) {
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

func (c Mock) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Routes, validation.Required),
	)
}

func newID() string {
	return uuid.NewString()
}
