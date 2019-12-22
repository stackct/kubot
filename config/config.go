package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

var Conf Configurator
var Log *zap.Logger

func init() {
	Log, _ = zap.NewProduction()
	Conf, _ = ParseFile(os.Getenv("KUBOT_CONFIG"))
}

type Config struct {
	Environments  []Environment     `yaml:"environments"`
	SlackToken    string            `yaml:"slackToken"`
	CommandConfig map[string]string `yaml:"commandConfig"`
	Commands      []Command         `yaml:"commands"`
}

type Configurator interface {
	HasAccess(id string, env string) bool
	GetEnvironmentByChannel(ch string) (*Environment, error)
	GetSlackToken() string
	GetCommands() []string
	GetCommand(name string) (*Command, error)
	GetCommandConfig() map[string]string
}

type Environment struct {
	Name    string   `yaml:"name"`
	Users   []string `yaml:"users"`
	Channel string   `yaml:"channel"`
}

type Command struct {
	Name     string            `yaml:"name"`
	Commands []Command         `yaml:"commands"`
	Args     []string          `yaml:"args"`
	Config   map[string]string `yaml:"config"`
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

func (c Config) GetSlackToken() string {
	return c.SlackToken
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

func (c Config) GetCommands() []string {
	commands := []string{}
	for _, cmd := range c.Commands {
		commands = append(commands, cmd.Name)
	}
	return commands
}

func (c Config) GetCommandConfig() map[string]string {
	return c.CommandConfig
}

func (c Config) GetCommand(name string) (*Command, error) {
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return &cmd, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("command not found: %s", name))
}
