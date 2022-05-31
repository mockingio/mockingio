package config

import (
	"io/ioutil"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Name   string   `yaml:"name,omitempty" json:"name,omitempty"`
	Port   string   `yaml:"port,omitempty" json:"port,omitempty"`
	Routes []*Route `yaml:"routes,omitempty" json:"routes,omitempty"`
}

func FromYamlFile(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read mock file")
	}

	return FromYaml(string(data))
}

func FromYaml(text string) (*Config, error) {
	decoder := yaml.NewDecoder(strings.NewReader(text))
	cfg := &Config{}
	if err := decoder.Decode(cfg); err != nil {
		return nil, errors.Wrap(err, "decode yaml to mock")
	}

	return cfg, nil
}

func (c Config) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Routes, validation.Required),
	)
}
