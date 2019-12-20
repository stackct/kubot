package command

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeploy(t *testing.T) {
	testCases := []struct {
		args    []string
		error   error
		command *Deploy
	}{
		{args: nil, error: errors.New("Deploy requires 1 argument"), command: nil},
		{args: []string{}, error: errors.New("Deploy requires 1 argument"), command: nil},
		{args: []string{"foo"}, error: nil, command: &Deploy{product: "foo"}},
	}

	for _, tc := range testCases {
		d, err := NewDeploy(tc.args)

		assert.Equal(t, err, tc.error)
		assert.Equal(t, d, tc.command)
	}
}

func TestDeployExecute(t *testing.T) {
	out := make(chan string)

	go Deploy{product: "Foo"}.Execute(out)

	assert.Equal(t, "Deploying *Foo*...", <-out)
	assert.Equal(t, "*Foo* deployment complete", <-out)
}
