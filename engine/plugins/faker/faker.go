package faker

import (
	"reflect"
	"strings"

	"github.com/jaswdr/faker"

	"github.com/mockingio/mockingio/engine/mock"
)

func New() *Faker {
	return &Faker{}
}

type Faker struct{}

func (f *Faker) Response(response *mock.Response) {
	apply := applier(faker.New())

	response.Body = apply(response.Body)
	for k, v := range response.Headers {
		response.Headers[k] = apply(v)
	}
}

func applier(fakerInc faker.Faker) func(string) string {
	return func(input string) string {
		tuples, err := parseCommand(input)
		if err != nil {
			return ""
		}

		for _, tuple := range tuples {
			value := runCmd(fakerInc, tuple.command)
			input = strings.Replace(input, tuple.placeholder, value, 1)
		}

		return input
	}
}

func runCmd(obj any, com command) string {
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if next, ok := com[strings.ToLower(m.Name)]; ok {
			if subCmd, ok := next.(command); ok {
				return runCmd(m.Func.Call([]reflect.Value{reflect.ValueOf(obj)})[0].Interface(), subCmd)
			}

			if next == nil {
				return m.Func.Call(append([]reflect.Value{reflect.ValueOf(obj)}))[0].String()
			}

			if args, ok := next.([]any); ok {
				return m.Func.Call(append([]reflect.Value{reflect.ValueOf(obj)}, convertArgs(args)...))[0].String()
			}

			return ""
		}
	}
	return ""
}

func convertArgs(args []any) []reflect.Value {
	var converted []reflect.Value
	for _, arg := range args {
		converted = append(converted, reflect.ValueOf(arg))
	}
	return converted
}
