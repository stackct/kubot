package command

import (
	"errors"
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeploy(t *testing.T) {
	testCases := []struct {
		args    []string
		error   error
		command *Deploy
	}{
		{args: nil, error: errors.New("Deploy requires 2 arguments"), command: nil},
		{args: []string{}, error: errors.New("Deploy requires 2 arguments"), command: nil},
		{args: []string{"foo"}, error: errors.New("Deploy requires 2 arguments"), command: nil},
		{args: []string{"foo", "1.0.0"}, error: nil, command: &Deploy{product: "foo", version: "1.0.0"}},
	}

	for _, tc := range testCases {
		d, err := NewDeploy(tc.args)

		assert.Equal(t, err, tc.error)
		assert.Equal(t, d, tc.command)
	}
}

func TestDeploy_Execute(t *testing.T) {
	out := make(chan string)
	config.Conf = config.NewMockConfig()

	go Deploy{product: "Foo", version: "1.0.0"}.Execute(out, Context{Environment: config.Environment{Name: "local"}})

	assert.Equal(t, "Deploying *Foo* (local)...", <-out)
	assert.Equal(t, "*Foo* deployment complete", <-out)
}
