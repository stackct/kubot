package config

import (
	"errors"
)

type MockConfig struct {
	Config
	MockHasAccess      bool
	MockHasEnvironment bool
	MockSlackToken     string
}

type MockConfigOption func(*MockConfig)

func NewMockConfig(options ...MockConfigOption) *MockConfig {
	c := &MockConfig{MockHasAccess: true, MockHasEnvironment: true}
	c.CommandPrefix = "!"

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *MockConfig) HasAccess(s string, e string) bool {
	return c.MockHasAccess
}

func (c *MockConfig) GetSlackToken() string {
	return c.MockSlackToken
}

func (c *MockConfig) GetEnvironmentByChannel(ch string) (*Environment, error) {
	if c.MockHasEnvironment {
		return &Environment{Name: ch, Channel: ch, Variables: map[string]string{"CLUSTER": "mock-cluster"}}, nil
	}

	return nil, errors.New("Environment not found")
}

func (c *MockConfig) GetLogging() Logging {
	return Logging{File: "/dev/null", Level: "DEBUG"}
}

func (c *MockConfig) GetCommands() []string {
	return []string{"cmd"}
}

func (c *MockConfig) GetCommand(name string, product string) (*Command, error) {
	if "fail" == name {
		return &Command{
			Name:     name,
			Commands: []Command{{Name: "ls", Args: []string{"foo"}}},
		}, nil
	}
	if "echo" == name && "stdout" == product {
		return &Command{
			Name:          name,
			Commands:      []Command{{Name: name, Args: []string{"foo"}}},
			Config:        map[string]string{},
			ChannelStdout: true,
		}, nil
	}
	if "echo" == name {
		return &Command{
			Name:     name,
			Commands: []Command{{Name: name, Args: []string{"foo", "${foo1}", "${foo2}", "${foo3}"}}},
			Config:   map[string]string{"foo2": "bar2"},
		}, nil
	}
	if "deploy" == name && "prohibited" == product {
		return nil, &ProhibitedCmdError{}
	}
	if "deploy" == name && "nil" == product {
		return nil, errors.New("command not found")
	}
	return &Command{}, nil
}

func (c *MockConfig) GetCommandConfig() map[string]string {
	return map[string]string{"foo1": "bar1"}
}
