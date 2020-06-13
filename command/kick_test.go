package command

import (
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKick(t *testing.T) {
	testCases := []struct {
		args    []string
		error   error
		command *Kick
	}{
		{args: nil, error: &CommandArgumentError{"usage: kick <product>"}, command: nil},
		{args: []string{}, error: &CommandArgumentError{"usage: kick <product>"}, command: nil},
		{args: []string{"foo"}, error: nil, command: &Kick{product: "foo"}},
	}

	for _, tc := range testCases {
		d, err := NewKick(tc.args)

		assert.Equal(t, err, tc.error)
		assert.Equal(t, d, tc.command)
	}
}

func TestKick_Execute(t *testing.T) {
	out := make(chan string)
	config.Conf = config.NewMockConfig()

	go Kick{product: "Foo"}.Execute(out, Context{Environment: config.Environment{Name: "local"}})

	assert.Equal(t, "Kicking *Foo* in *local* environment...", <-out)
	assert.Equal(t, "*Foo* was successfully kicked.", <-out)
}
