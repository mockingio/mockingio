package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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
	Proxy  *Proxy   `yaml:"proxy,omitempty" json:"proxy,omitempty"`
	// all OPTIONS calls are responded with success if AutoCORS is true
	AutoCORS bool `yaml:"auto_cors,omitempty" json:"auto_cors,omitempty"`
	TLS      *TLS `yaml:"tls,omitempty" json:"tls,omitempty"`
	options  mockOptions
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

func (r Mock) Clone() Mock {
	result := r
	result.ID = newID()
	result.Routes = make([]*Route, len(r.Routes))
	result.options = r.options
	result.Proxy = r.Proxy.Clone()
	result.TLS = r.TLS.Clone()

	for i, rule := range r.Routes {
		result.Routes[i] = rule.Clone()
	}

	return result
}

func FromFile(file string, opts ...Option) (*Mock, error) {
	// TODO: Detects file type, support JSON
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read mock file")
	}

	return FromYaml(string(data), opts...)
}

func FromYaml(text string, opts ...Option) (*Mock, error) {
	decoder := yaml.NewDecoder(strings.NewReader(text))
	m := New(opts...)
	if err := decoder.Decode(m); err != nil {
		return nil, errors.Wrap(err, "decode yaml to mock")
	}
	defaultValues(m)
	if m.options.idGeneration {
		addIDs(m)
	}

	return m, nil
}

func (m Mock) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Routes, validation.Required),
	)
}

func (m Mock) ProxyEnabled() bool {
	return m.Proxy != nil && m.Proxy.Enabled
}

func (m Mock) TLSEnabled() bool {
	return m.TLS != nil && m.TLS.Enabled
}

func defaultValues(m *Mock) {
	for _, r := range m.Routes {
		if r.Method == "" {
			r.Method = http.MethodGet
		}
		for i, res := range r.Responses {
			if res.Status == 0 {
				res.Status = 200
			}
			r.Responses[i] = res
		}
	}
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

func jsonFieldName(field reflect.StructField) string {
	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "" {
		name = field.Name
	}
	return name
}

func patchStruct(resource interface{}, patches map[string]*json.RawMessage) error {
	value := reflect.ValueOf(resource)
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("can't operate on non-struct: %s", value.Kind().String())
	}
	if !value.CanAddr() {
		return errors.New("unaddressable struct value")
	}
	valueT := value.Type()
	for i := 0; i < valueT.NumField(); i++ {
		field := value.Field(i)
		if !field.CanAddr() || !field.CanInterface() {
			continue
		}
		if patch, ok := patches[jsonFieldName(valueT.Field(i))]; ok {
			field.Set(reflect.Zero(field.Type()))
			if err := json.Unmarshal(*patch, field.Addr().Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}
