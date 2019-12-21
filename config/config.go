package config

import (
	"fmt"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var Log *zap.Logger

func init() {
	Log, _ = zap.NewProduction()
}

type Config struct {
	Environments []Environment `yaml:"environments"`
}

type Configurator interface {
	HasAccess(id string, env string) bool
	GetEnvironmentByChannel(ch string) (*Environment, error)
}

type Environment struct {
	Name    string   `yaml:"name"`
	Users   []string `yaml:"users"`
	Channel string   `yaml:"channel"`
}

func ParseFile(f string) (Configurator, error) {
	file, err := os.Open(f)
	if err != nil {
		return Config{}, err
	}

	input, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	bytes := []byte(input)
	config, err := Parse(bytes)
	if err != nil {
		return Config{}, err
	}

	Log.Info("configuration file loaded", zap.String("path", f), zap.Int("environments", len(config.Environments)))

	return config, nil
}

func Parse(bs []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(bs, &config)

	return config, err
}

func (c Config) GetEnvironment(env string) (*Environment, error) {
	for _, e := range c.Environments {
		if e.Name == env {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("Environment '%v' not found", env)
}

func (c Config) GetEnvironmentByChannel(ch string) (*Environment, error) {
	for _, e := range c.Environments {
		if e.Channel == ch {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("Environment not found for channel: '%v'", ch)
}

func (c Config) HasAccess(user string, env string) bool {
	e, err := c.GetEnvironment(env)
	if err != nil {
		return false
	}

	return e.ContainsUser(user)
}

func (e Environment) ContainsUser(u string) bool {
	for _, user := range e.Users {
		if user == u {
			return true
		}
	}

	return false
}
