package faker

import (
	"testing"

	"github.com/mockingio/mockingio/engine/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input            string
		expectedCommands []commandTuple
	}{
		{
			"something ${faker.person.name} here ${faker.lorem.text(100)} example",
			[]commandTuple{
				{
					placeholder: "${faker.person.name}",
					command: command{
						"person": command{
							"name": nil,
						},
					},
				},
				{
					placeholder: "${faker.lorem.text(100)}",
					command: command{
						"lorem": command{
							"text": []any{100},
						},
					},
				},
			},
		},
		{
			"${faker.person.name}",
			[]commandTuple{
				{
					placeholder: "${faker.person.name}",
					command: command{
						"person": command{
							"name": nil,
						},
					},
				},
			},
		},
		{
			"${faker.person.name('alex')}",
			[]commandTuple{
				{
					placeholder: "${faker.person.name('alex')}",
					command: command{
						"person": command{
							"name": []any{"alex"},
						},
					},
				},
			},
		},
		{
			"no commands",
			nil,
		},
	}

	for _, tt := range tests {
		commands, err := parseCommand(tt.input)
		require.NoError(t, err)
		assert.Equal(t, tt.expectedCommands, commands)
	}
}

func TestFaker_Response(t *testing.T) {
	plug := New()

	resp := &mock.Response{
		Body: "hi ${faker.person.name}",
		Headers: map[string]string{
			"X-Name": "token: ${faker.lorem.text(5)}",
			"X-Test": "no commands",
		},
	}
	plug.Response(resp)

	assert.NotEqual(t, "hi ${faker.person.name}", resp.Body)
	assert.NotEqual(t, map[string]string{
		"X-Name": "token: ${faker.lorem.text(5)}",
		"X-Test": "no commands",
	}, resp.Headers)
}
