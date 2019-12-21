package config

import (
	"errors"
)

type MockConfig struct {
	MockHasAccess      bool
	MockHasEnvironment bool
}

type MockConfigOption func(*MockConfig)

func NewMockConfig(options ...MockConfigOption) *MockConfig {
	c := &MockConfig{MockHasAccess: true, MockHasEnvironment: true}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *MockConfig) HasAccess(s string, e string) bool {
	return c.MockHasAccess
}

func (c *MockConfig) GetEnvironmentByChannel(ch string) (*Environment, error) {
	if c.MockHasEnvironment {
		return &Environment{Name: ch, Channel: ch}, nil
	}

	return nil, errors.New("Environment not found")
}
