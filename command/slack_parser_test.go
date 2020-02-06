package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	SlackCommandPrefix = "!"
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
