package config

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var mockConfig Config

func init() {
	mockConfig = Config{
		Environments: []Environment{
			Environment{Name: "e1", Users: []string{}, Channel: "ch1"},
			Environment{Name: "e2", Users: []string{"john.doe", "mary.foo"}, Channel: "ch2"},
		},
	}
}

func TestParseFile(t *testing.T) {
	config, _ := ParseFile("./resources/config.yml")

	assert.True(t, reflect.DeepEqual(config, mockConfig))
}

func TestConfig_GetEnvironment(t *testing.T) {
	testCases := []struct {
		name        string
		environment *Environment
		error       bool
	}{
		{"e1", &mockConfig.Environments[0], false},
		{"e2", &mockConfig.Environments[1], false},
		{"e3", nil, true},
	}

	for _, tc := range testCases {
		env, err := mockConfig.GetEnvironment(tc.name)
		if !tc.error {
			assert.Nil(t, err)
		}

		assert.IsType(t, tc.environment, env)
	}
}

func TestConfig_GetEnvironmentByChannel(t *testing.T) {
	testCases := []struct {
		channel     string
		environment *Environment
		error       bool
	}{
		{"ch1", &mockConfig.Environments[0], false},
		{"ch2", &mockConfig.Environments[1], false},
		{"ch3", nil, true},
	}

	for _, tc := range testCases {
		env, err := mockConfig.GetEnvironmentByChannel(tc.channel)
		if !tc.error {
			assert.Nil(t, err)
		}

		assert.IsType(t, tc.environment, env)
	}
}

func TestEnvironment_ContainsUser(t *testing.T) {
	e := Environment{Users: []string{"u1"}}

	assert.True(t, e.ContainsUser("u1"))
	assert.False(t, e.ContainsUser("u2"))
}
