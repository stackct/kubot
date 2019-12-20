package command

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	testCases := []struct {
		message  string
		expected interface{}
		err      error
	}{
		{
			message:  "nomatch",
			expected: nil,
			err:      errors.New("unknown command"),
		},
		{
			message:  "!help",
			expected: &Help{},
			err:      nil,
		},
		{
			message:  "!deploy",
			expected: &Deploy{},
			err:      errors.New("Deploy requires 1 argument"),
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

type TestCommand struct{}

func (t TestCommand) Name() string            { return "mock" }
func (t TestCommand) Execute(out chan string) {}

func TestExecuteSuccess(t *testing.T) {
	r, w := io.Pipe()
	out := make(chan string)
	buf := new(bytes.Buffer)

	go func() { buf.ReadFrom(r) }()
	Execute(&TestCommand{}, w, out, "echo", "foo")
	assert.Equal(t, "foo\n", buf.String())
}

func TestExecuteCommandFails(t *testing.T) {
	r, w := io.Pipe()
	out := make(chan string)
	buf := new(bytes.Buffer)

	go func() { buf.ReadFrom(r) }()
	go Execute(&TestCommand{}, w, out, "ls", "foo")

	assert.Equal(t, "*mock* command failed", <-out)
	assert.Equal(t, "ls: cannot access 'foo': No such file or directory\n", buf.String())
}
