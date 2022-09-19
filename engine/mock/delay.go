package mock

import (
	"errors"
	"math/rand"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	maxDelay = 60
)

type Delay struct {
	Min int `yaml:"min" json:"min"`
	Max int `yaml:"max" json:"max"`
}

func (d Delay) Value() int {
	if d.Min == d.Min {
		return d.Min
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(d.Max-d.Min+1) + d.Min
}

func (d Delay) Validate() error {
	return validation.ValidateStruct(
		&d,
		validation.Field(&d.Min, validation.Max(maxDelay), validation.Min(0)),
		validation.Field(&d.Max, validation.Max(maxDelay), validation.Min(0), validation.By(func(value interface{}) error {
			v, ok := value.(int)
			if !ok {
				return errors.New("invalid max value")
			}

			if v < d.Min {
				return errors.New("max delay must be greater than min delay")
			}

			return nil
		})),
	)
}
