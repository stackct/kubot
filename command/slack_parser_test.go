package command

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	SlackCommandPrefix = "!"
}

func TestRemoveDirectMessages(t *testing.T) {
	testCases := []struct {
		message     string
		contextName string
		expected    string
		err         error
	}{
		{
			message:     "!deploy foo x.y.z",
			contextName: "bar",
			expected:    "!deploy foo x.y.z",
			err:         nil,
		},
		{
			message:     "!deploy foo x.y.z @baz",
			contextName: "bar",
			expected:    "",
			err:         errors.New("command is not intended for this instance"),
		},
		{
			message:     "@baz !deploy foo x.y.z",
			contextName: "bar",
			expected:    "",
			err:         errors.New("command is not intended for this instance"),
		},
		{
			message:     "@bar !deploy foo x.y.z",
			contextName: "bar",
			expected:    "!deploy foo x.y.z",
			err:         nil,
		},
		{
			message:     "!deploy foo x.y.z @bar",
			contextName: "bar",
			expected:    "!deploy foo x.y.z",
			err:         nil,
		},
		{
			message:     "!deploy foo x.y.z @baz @bar",
			contextName: "bar",
			expected:    "!deploy foo x.y.z",
			err:         nil,
		},
		{
			message:     "!deploy foo x.y.z @baz @foobar",
			contextName: "bar",
			expected:    "",
			err:         errors.New("command is not intended for this instance"),
		},
	}

	for _, tt := range testCases {
		command, err := RemoveDirectMessages(tt.message, tt.contextName)
		assert.Equal(t, tt.expected, command)
		assert.Equal(t, tt.err, err)
	}
}

func TestParseCommand(t *testing.T) {
	testCases := []struct {
		message  string
		expected interface{}
		err      error
	}{
		{
			message:  "nomatch",
			expected: nil,
			err:      &UnknownCommandError{"input does not match command syntax"},
		},
		{
			message:  "!help",
			expected: &Help{},
			err:      nil,
		},
		{
			message:  "!deploy",
			expected: &Deploy{},
			err:      &CommandArgumentError{"usage: deploy <product> <version> [release]"},
		},
		{
			message:  "!bounce",
			expected: &Bounce{},
			err:      nil,
		},
		{
			message:  "!kick",
			expected: &Kick{},
			err:      &CommandArgumentError{"usage: kick <product>"},
		},
		{
			message:  "!version",
			expected: &Version{},
			err:      &CommandArgumentError{"usage: version <product>"},
		},
		{
			message:  "!unknown",
			expected: nil,
			err:      &UnknownCommandError{"unknown command"},
		},
	}

	for _, tt := range testCases {
		command, err := SlackCommandParser{}.Parse(tt.message)
		assert.Equal(t, tt.err, err)

		if tt.err == nil {
			assert.Equal(t, tt.expected, command)
		}
	}
}
