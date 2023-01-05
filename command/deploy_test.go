package command

import (
	"kubot/config"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewDeploy(t *testing.T) {
	testCases := []struct {
		args    []string
		error   error
		command *Deploy
	}{
		{args: nil, error: &CommandArgumentError{"usage: deploy <product> <version> [release]"}, command: nil},
		{args: []string{}, error: &CommandArgumentError{"usage: deploy <product> <version> [release]"}, command: nil},
		{args: []string{"foo"}, error: &CommandArgumentError{"usage: deploy <product> <version> [release]"}, command: nil},
		{args: []string{"foo", "1.0.0"}, error: nil, command: &Deploy{product: "foo", version: "1.0.0", release: "foo"}},
		{args: []string{"foo", "1.0.0", "bar"}, error: nil, command: &Deploy{product: "foo", version: "1.0.0", release: "bar"}},
	}

	for _, tc := range testCases {
		d, err := NewDeploy(tc.args)

		assert.Equal(t, err, tc.error)
		assert.Equal(t, d, tc.command)
	}
}

func TestDeploy_Execute(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	out := make(chan string)
	config.Conf = config.NewMockConfig()

	go Deploy{product: "Foo", version: "1.0.0", release: "Foo"}.Execute(out, Context{Environment: config.Environment{Name: "local"}, User: "foo_bar"})

	assert.Equal(t, "Deploying *Foo-1.0.0* to *local*...", <-out)
	assert.Equal(t, "*Foo-1.0.0* was successfully deployed to *local*", <-out)

	assert.Equal(t, 1, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, "deployed successfully", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, log.Fields{"product": "Foo", "version": "1.0.0", "environment": "local", "username": "foo_bar"}, e.Fields)
}

func TestDeployWithRelease_Execute(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	out := make(chan string)
	config.Conf = config.NewMockConfig()

	go Deploy{product: "Foo", version: "1.0.0", release: "Bar"}.Execute(out, Context{Environment: config.Environment{Name: "local"}, User: "foo_bar"})

	assert.Equal(t, "Deploying *Bar* with *Foo-1.0.0* to *local*...", <-out)
	assert.Equal(t, "*Bar* with *Foo-1.0.0* was successfully deployed to *local*", <-out)

	assert.Equal(t, 1, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, "deployed successfully", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, log.Fields{"product": "Foo", "version": "1.0.0", "environment": "local", "username": "foo_bar"}, e.Fields)
}

func TestDeployWithProhibitedCommand_Execute(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	out := make(chan string)
	config.Conf = config.NewMockConfig()

	Deploy{product: "prohibited", version: "1.0.0", release: "prohibited"}.Execute(out, Context{Environment: config.Environment{Name: "local"}, User: "foo_bar"})

	assert.Equal(t, 1, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, "Skipped deploy command for product prohibited because it was found on the prohibited command list", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
}

func TestDeployWithCommandNotFound_Execute(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	out := make(chan string)
	config.Conf = config.NewMockConfig()

	go Deploy{product: "nil", version: "1.0.0", release: "nil"}.Execute(out, Context{Environment: config.Environment{Name: "local"}, User: "foo_bar"})

	assert.Equal(t, "deploy command for product nil was not found", <-out)

	assert.Equal(t, 1, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, "command not found", e.Message)
	assert.Equal(t, log.ErrorLevel, e.Level)
}
