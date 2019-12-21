package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	})

	assert.False(t, c.HasAccess("foo", "bar"))

	env, err := c.GetEnvironmentByChannel("ch")
	assert.NotNil(t, err)
	assert.Nil(t, env)
}
