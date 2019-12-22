package config

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockConfig Config

func init() {
	mockConfig = Config{
		Environments: []Environment{
			Environment{Name: "e1", Users: []string{}, Channel: "ch1"},
			Environment{Name: "e2", Users: []string{"john.doe", "mary.foo"}, Channel: "ch2"},
		},
		SlackToken: "some-token",
		Logging: Logging{
			File:  "kubot.log",
			Level: "INFO",
		},
		CommandConfig: map[string]string{
			"productRepo":          "repo",
			"deployTimeoutSeconds": "900",
			"chartFile":            "chart",
		},
		Commands: []Command{
			Command{
				Name: "deploy",
				Commands: []Command{
					Command{
						Name: "echo",
						Args: []string{"deploy", "${productRepo}/${product}", "--version", "${version}", "--timeout", "${deployTimeoutSeconds}", "-f", "${chartFile}"},
					},
				},
			},
		},
	}
}

func TestParseFile(t *testing.T) {
	config, err := ParseFile("./resources/kubot.yml")
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(config, mockConfig))
}

func TestConfig_GetEnvironment(t *testing.T) {
	testCases := []struct {
		name        string
		environment *Environment
		error       error
	}{
		{"e1", &mockConfig.Environments[0], nil},
		{"e2", &mockConfig.Environments[1], nil},
		{"e3", nil, errors.New("Environment 'e3' not found")},
	}

	for _, tc := range testCases {
		env, err := mockConfig.GetEnvironment(tc.name)
		assert.Equal(t, tc.error, err)
		assert.IsType(t, tc.environment, env)
	}
}

func TestConfig_GetEnvironmentByChannel(t *testing.T) {
	testCases := []struct {
		channel     string
		environment *Environment
		error       error
	}{
		{"ch1", &mockConfig.Environments[0], nil},
		{"ch2", &mockConfig.Environments[1], nil},
		{"ch3", nil, errors.New("Environment not found for channel: 'ch3'")},
	}

	for _, tc := range testCases {
		env, err := mockConfig.GetEnvironmentByChannel(tc.channel)
		assert.Equal(t, tc.error, err)
		assert.IsType(t, tc.environment, env)
	}
}

func TestEnvironment_ContainsUser(t *testing.T) {
	e := Environment{Users: []string{"u1"}}

	assert.True(t, e.ContainsUser("u1"))
	assert.False(t, e.ContainsUser("u2"))
}
