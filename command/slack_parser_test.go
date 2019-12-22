package command

import (
	"errors"
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
			err:      errors.New("input does not match command syntax"),
		},
		{
			message:  "!help",
			expected: &Help{},
			err:      nil,
		},
		{
			message:  "!deploy",
			expected: &Deploy{},
			err:      errors.New("Deploy requires 2 arguments"),
		},
		{
			message:  "!unknown",
			expected: nil,
			err:      errors.New("unknown command"),
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
