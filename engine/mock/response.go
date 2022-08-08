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
	Status          int               `yaml:"status" json:"status"`
	Delay           int64             `yaml:"delay,omitempty" json:"delay,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body            string            `yaml:"body,omitempty" json:"body,omitempty"`
	RuleAggregation RuleAggregation   `yaml:"rule_aggregation,omitempty" json:"rule_aggregation,omitempty"`
	Rules           []Rule            `yaml:"rules,omitempty" json:"rules,omitempty"`
	IsDefault       bool              `yaml:"is_default,omitempty" json:"is_default,omitempty"`
}

func (r Response) Clone() Response {
	result := r
	result.ID = newID()
	result.RuleAggregation = r.RuleAggregation
	result.Rules = make([]Rule, len(r.Rules))

	for i, rule := range r.Rules {
		result.Rules[i] = rule.Clone()
	}

	return result
}

func (r Response) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Status, validation.Required),
		validation.Field(&r.RuleAggregation, validation.In(Or, And)),
	)
}
