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

	expected := "Deploying *Foo*..."

	go func() {
		if msg := <-out; expected != msg {
			t.Fatalf("Assertion failed; wanted: %s, but got: %s", expected, msg)
		}
		<-out
		<-out
		<-out
	}()

	Deploy{product: "Foo"}.Execute(out)
}
