package config

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var mockConfig Config

func init() {
	mockConfig = Config{
		Environments: []Environment{
			{Name: "e1", Users: []string{}, Channel: "ch1", Variables: map[string]string{"foo1": "bar", "foo2": "baz"}},
			{Name: "e2", Users: []string{"john.doe", "mary.foo"}, Channel: "ch2"},
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
			{
				Name: "deploy",
				Commands: []Command{
					{
						Name: "echo",
						Args: []string{"deploy", "${productRepo}/${product}", "--version", "${version}", "--timeout", "${deployTimeoutSeconds}", "-f", "${chartFile}"},
					},
				},
			},
		},
		Init: []Command{
			{Name: "echo", Args: []string{"initialized"}},
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

func TestConfig_GetSlackTokenFromEnvironment(t *testing.T) {
	os.Setenv("SLACK_TOKEN", "abcde")
	assert.Equal(t, "abcde", Config{}.GetSlackToken())
}

func TestConfig_GetLogging(t *testing.T) {
	logging := Logging{}
	c := Config{Logging: logging}
	assert.Equal(t, logging, c.GetLogging())
}

func TestConfig_HasAccess(t *testing.T) {
	c := Config{
		Environments: []Environment{
			{Name: "e1", Users: []string{"u1"}},
			{Name: "e2", Users: []string{"u2"}},
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
			{Name: "cmd1"},
			{Name: "cmd2"},
			{Name: "cmd2"},
		},
	}

	assert.Equal(t, []string{"cmd1", "cmd2"}, c.GetCommands())
}

func TestConfig_GetCommandConfig_Is_Thread_Safe(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	config := Config{CommandConfig: data}

	commandConfig := config.GetCommandConfig()
	assert.Equal(t, data, commandConfig)

	commandConfig["foo"] = "baz"
	commandConfig = config.GetCommandConfig()
	assert.Equal(t, "bar", commandConfig["foo"])
}

func TestConfig_GetCommand_FromNameAndProduct(t *testing.T) {
	config := Config{
		Commands: []Command{
			{Name: "cmd"},
			{Name: "cmd", Product: "foo"},
			{Name: "cmd2"},
		},
	}
	cmd, err := config.GetCommand("cmd", "foo")
	assert.Equal(t, err, nil)
	assert.Equal(t, cmd.Name, "cmd")
	assert.Equal(t, cmd.Product, "foo")
}

func TestConfig_GetCommand_DefaultProduct(t *testing.T) {
	config := Config{
		Commands: []Command{
			{Name: "cmd", Product: "foo"},
			{Name: "cmd"},
			{Name: "cmd2"},
		},
	}
	cmd, err := config.GetCommand("cmd", "default")
	assert.Equal(t, err, nil)
	assert.Equal(t, cmd.Name, "cmd")
	assert.Equal(t, cmd.Product, "")
}

func TestConfig_GetCommand_NotFound(t *testing.T) {
	config := Config{
		Commands: []Command{
			{Name: "cmd"},
			{Name: "cmd2"},
		},
	}
	_, err := config.GetCommand("cmd3", "")
	assert.Equal(t, err, errors.New("command not found: cmd3"))
}

func TestConfig_GetCommandPrefix(t *testing.T) {
	assert.Equal(t, "!", Config{CommandPrefix: "!"}.GetCommandPrefix())
}

func TestConfig_GetInit(t *testing.T) {
	init := []Command{{Name: "cmd1"}}
	assert.Equal(t, init, Config{Init: init}.GetInit())
}
