package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockConfig(t *testing.T) {
	c := NewMockConfig()
	assert.True(t, c.HasAccess("foo", "bar"))

	env, err := c.GetEnvironmentByChannel("ch")
	assert.Nil(t, err)
	assert.Equal(t, "ch", env.Name)
	assert.Equal(t, "ch", env.Channel)
}

func TestMockConfig_WithOptions(t *testing.T) {
	c := NewMockConfig(func(mc *MockConfig) {
		mc.MockHasAccess = false
		mc.MockHasEnvironment = false
		mc.MockSlackToken = "some-token"
	})

	assert.False(t, c.HasAccess("foo", "bar"))

	env, err := c.GetEnvironmentByChannel("ch")
	assert.NotNil(t, err)
	assert.Nil(t, env)
	assert.Equal(t, "some-token", c.GetSlackToken())
	assert.Equal(t, []string{"cmd"}, c.GetCommands())
}

func TestMockConfig_GetCommand(t *testing.T) {
	testCases := []struct {
		name    string
		product string
		error   error
		command *Command
	}{
		{"fail", "product", nil, &Command{Name: "fail", Commands: []Command{{Name: "ls", Args: []string{"foo"}}}}},
		{"echo", "product", nil, &Command{Name: "echo", Commands: []Command{{Name: "echo", Args: []string{"foo", "${foo1}", "${foo2}", "${foo3}"}}}, Config: map[string]string{"foo2": "bar2"}}},
		{"other", "product", nil, &Command{}},
	}

	c := NewMockConfig()

	for _, tc := range testCases {
		cmd, err := c.GetCommand(tc.name, tc.product)
		assert.Equal(t, tc.error, err)
		assert.Equal(t, tc.command, cmd)
	}
}

func TestMockConfig_GetLogging(t *testing.T) {
	assert.Equal(t, Logging{File: "/dev/null", Level: "DEBUG"}, NewMockConfig().GetLogging())
}

func TestMockConfig_GetCommandConfig(t *testing.T) {
	assert.Equal(t, map[string]string{"foo1": "bar1"}, NewMockConfig().GetCommandConfig())
}
