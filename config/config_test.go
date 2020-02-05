package config

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var mockConfig Config

func init() {
	mockConfig = Config{
		Environments: []Environment{
			Environment{Name: "e1", Release: "r1", Users: []string{}, Channel: "ch1"},
			Environment{Name: "e2", Release: "r2", Users: []string{"john.doe", "mary.foo"}, Channel: "ch2"},
		},
		SlackToken: "some-token",
		Logging: Logging{
			File:  "kubot.log",
			Level: "INFO",
		},
		CommandPrefix: ".",
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
		Init: []Command{
			Command{Name: "echo", Args: []string{"initialized"}},
		},
	}
}

func TestParseFile(t *testing.T) {
	config, err := ParseFile("./resources/kubot.yml")
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(config, mockConfig))
}

func TestParseFile_File_Error(t *testing.T) {
	testCases := []struct {
		file     string
		expected bool
	}{
		{"/dev/null/foo", true},
		{"/tmp", true},
		{"./resources/kubot.yml", false},
	}

	for _, tc := range testCases {
		config, e := ParseFile(tc.file)
		_, ok := e.(*os.PathError)

		assert.Equal(t, tc.expected, ok)
		assert.NotNil(t, config)
	}
}

func TestParseFile_Parse_Error(t *testing.T) {
	testCases := []struct {
		file     string
		expected bool
	}{
		{"./resources/bad.yml", true},
	}

	for _, tc := range testCases {
		config, e := ParseFile(tc.file)
		_, ok := e.(*yaml.TypeError)

		assert.Equal(t, tc.expected, ok)
		assert.NotNil(t, config)
	}
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

func TestConfig_GetSlackToken(t *testing.T) {
	assert.Equal(t, "abcde", Config{SlackToken: "abcde"}.GetSlackToken())
}

func TestConfig_GetLogging(t *testing.T) {
	logging := Logging{}
	c := Config{Logging: logging}
	assert.Equal(t, logging, c.GetLogging())
}

func TestConfig_HasAccess(t *testing.T) {
	c := Config{
		Environments: []Environment{
			Environment{Name: "e1", Users: []string{"u1"}},
			Environment{Name: "e2", Users: []string{"u2"}},
		},
	}

	assert.True(t, c.HasAccess("u1", "e1"))
	assert.True(t, c.HasAccess("u2", "e2"))
	assert.False(t, c.HasAccess("u1", "e2"))
	assert.False(t, c.HasAccess("u2", "e1"))
}

func TestConfig_GetCommands(t *testing.T) {
	c := Config{
		Commands: []Command{
			Command{Name: "cmd1"},
			Command{Name: "cmd2"},
		},
	}

	assert.Equal(t, []string{"cmd1", "cmd2"}, c.GetCommands())
}

func TestConfig_GetCommandConfig(t *testing.T) {
	cfg := map[string]string{"foo": "bar"}
	c := Config{CommandConfig: cfg}

	assert.Equal(t, cfg, c.GetCommandConfig())
}

func TestConfig_GetCommand(t *testing.T) {
	c := Config{
		Commands: []Command{
			Command{Name: "cmd1"},
			Command{Name: "cmd2"},
		},
	}

	testCases := []struct {
		name string
		cmd  *Command
		err  error
	}{
		{"cmd1", &Command{Name: "cmd1"}, nil},
		{"cmd2", &Command{Name: "cmd2"}, nil},
		{"nil", nil, errors.New("command not found: nil")},
	}

	for _, tc := range testCases {
		cmd, err := c.GetCommand(tc.name)

		assert.Equal(t, tc.cmd, cmd)
		assert.Equal(t, tc.err, err)
	}
}

func TestConfig_GetCommandPrefix(t *testing.T) {
	assert.Equal(t, "!", Config{CommandPrefix: "!"}.GetCommandPrefix())
}

func TestConfig_GetInit(t *testing.T) {
	init := []Command{Command{Name: "cmd1"}}
	assert.Equal(t, init, Config{Init: init}.GetInit())
}
