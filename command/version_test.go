package command

import (
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersion(t *testing.T) {
	testCases := []struct {
		args    []string
		error   error
		command *Version
	}{
		{args: nil, error: &CommandArgumentError{"usage: version <product>"}, command: nil},
		{args: []string{}, error: &CommandArgumentError{"usage: version <product>"}, command: nil},
		{args: []string{"foo"}, error: nil, command: &Version{product: "foo"}},
	}

	for _, tc := range testCases {
		d, err := NewVersion(tc.args)

		assert.Equal(t, err, tc.error)
		assert.Equal(t, d, tc.command)
	}
}

func TestVersion_Execute(t *testing.T) {
	out := make(chan string)
	config.Conf = config.NewMockConfig()

	go Version{}.Execute(out, Context{Environment: config.Environment{Name: "local"}})

	assert.Equal(t, "", <-out)
}
