package config

import (
	"errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name    string   `yaml:"name"`
	Users   []string `yaml:"users"`
	Channel string   `yaml:"channel"`
}

var AppConfig Config

func init() {
	AppConfig, _ = ParseFile(os.Getenv("KUBOT_CONFIG"))
}

func ParseFile(f string) (Config, error) {
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

	return nil, errors.New("Environment not found")
}

func (c Config) GetEnvironmentByChannel(ch string) (*Environment, error) {
	for _, e := range c.Environments {
		if e.Channel == ch {
			return &e, nil
		}
	}

	return nil, errors.New("Environment not found")
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
