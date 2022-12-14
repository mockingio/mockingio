package mock

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RuleAggregation string

const (
	Or  RuleAggregation = "or"
	And RuleAggregation = "and"
)

type Response struct {
	ID              string            `yaml:"id,omitempty" json:"id,omitempty"`
	Status          int               `yaml:"status" json:"status,omitempty"`
	Delay           Delay             `yaml:"delay,omitempty" json:"delay,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body            string            `yaml:"body,omitempty" json:"body,omitempty"`
	FilePath        string            `yaml:"file_path,omitempty" json:"file_path,omitempty"`
	RuleAggregation RuleAggregation   `yaml:"rule_aggregation,omitempty" json:"rule_aggregation,omitempty"`
	Rules           []Rule            `yaml:"rules,omitempty" json:"rules,omitempty"`
	IsDefault       bool              `yaml:"is_default,omitempty" json:"is_default,omitempty"`
}

func (r Response) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Status, validation.Min(100), validation.Max(999)),
		validation.Field(&r.RuleAggregation, validation.In(Or, And)),
	)
}
