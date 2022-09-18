package faker

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type command map[string]any
type commandTuple struct {
	placeholder string
	command     command
}

// parseCommand parse the string and get commands
func parseCommand(input string) ([]commandTuple, error) {
	var commands []commandTuple
	comp, err := regexp.Compile(`\${faker\.([^{}]+)}`)
	if err != nil {
		return nil, err
	}

	groups := comp.FindAllStringSubmatch(input, -1)
	if len(groups) == 0 {
		return nil, nil
	}

	for _, group := range groups {
		if len(group) != 2 {
			continue
		}

		cmd, err := toCommand(strings.Split(group[1], "."))
		if err != nil {
			return nil, err
		}
		commands = append(commands, commandTuple{
			placeholder: group[0],
			command:     cmd,
		})
	}

	return commands, nil
}

func toArguments(args []string) ([]any, error) {
	var arguments []any
	for _, arg := range args {
		matched, err := regexp.Match(`^['"]+.*['"]+$`, []byte(arg))
		if err != nil {
			return nil, err
		}

		if matched {
			arguments = append(arguments, strings.Trim(arg, `"'`))
			continue
		}

		// if not matched, assume it's a number
		intVar, err := strconv.Atoi(arg)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, intVar)
	}

	return arguments, nil
}

func toCommand(tokens []string) (command, error) {
	if len(tokens) == 1 {
		token := tokens[0]
		if strings.Contains(token, "(") {
			comp, err := regexp.Compile(`([^(]+)\(([^)]+)\)`)
			if err != nil {
				return command{}, err
			}

			groups := comp.FindAllStringSubmatch(token, -1)
			if len(groups) != 1 {
				return command{}, errors.New("invalid command")
			}

			args, err := toArguments(groups[0][2:])
			if err != nil {
				return command{}, err
			}

			return command{
				groups[0][1]: args,
			}, nil
		}

		return command{
			token: nil,
		}, nil
	}

	newCmd, err := toCommand(tokens[1:])
	if err != nil {
		return command{}, err
	}

	return command{
		tokens[0]: newCmd,
	}, nil
}
