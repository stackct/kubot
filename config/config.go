package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	yaml "gopkg.in/yaml.v2"
)

var Conf Configurator

func init() {
	Conf, _ = ParseFile(os.Getenv("KUBOT_CONFIG"))
}

func InitLogging(logFilename string, logLevel string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		log.SetHandler(json.New(logFile))
	} else {
		log.WithError(err).WithField("logfile", logFilename).Error("Failed to create log file, using console instead")
		log.SetHandler(text.New(os.Stdout))
	}
	log.SetLevelFromString(logLevel)
	return logFile, err
}

type Config struct {
	Environments  []Environment     `yaml:"environments"`
	SlackToken    string            `yaml:"slackToken"`
	Logging       Logging           `yaml:"logging"`
	Commands      []Command         `yaml:"commands"`
	CommandConfig map[string]string `yaml:"commandConfig"`
	CommandPrefix string            `yaml:"commandPrefix"`
	Init          []Command         `yaml:"init"`
}

type Configurator interface {
	HasAccess(id string, env string) bool
	GetEnvironmentByChannel(ch string) (*Environment, error)
	GetSlackToken() string
	GetLogging() Logging
	GetCommand(name string, product string) (*Command, error)
	GetCommands() []string
	GetCommandConfig() map[string]string
	GetCommandPrefix() string
	GetInit() []Command
}

type Command struct {
	Name          string            `yaml:"name"`
	Product       string            `yaml:"product"`
	ChannelStdout bool              `yaml:"channelStdout"`
	Commands      []Command         `yaml:"commands"`
	Args          []string          `yaml:"args"`
	Config        map[string]string `yaml:"config"`
}

type Logging struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
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

	log.WithField("path", f).WithField("environments", len(config.Environments)).Info("configuration file loaded")

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

func (c Config) GetLogging() Logging {
	return c.Logging
}

func (c Config) HasAccess(user string, env string) bool {
	e, err := c.GetEnvironment(env)
	if err != nil {
		return false
	}

	return e.ContainsUser(user)
}

func (c Config) GetCommands() []string {
	commands := []string{}
	for _, cmd := range c.Commands {
		commands = append(commands, cmd.Name)
	}
	return commands
}

func (c Config) GetCommandConfig() map[string]string {
	config := map[string]string{}
	for k, v := range c.CommandConfig {
		config[k] = v
	}
	return config
}

func (c Config) GetCommand(name string, product string) (*Command, error) {
	for _, cmd := range c.Commands {
		if cmd.Name == name && cmd.Product == product {
			return &cmd, nil
		}
	}
	for _, cmd := range c.Commands {
		if cmd.Name == name && cmd.Product == "" {
			return &cmd, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("command not found: %s", name))
}

func (c Config) GetCommandPrefix() string {
	return c.CommandPrefix
}

func (c Config) GetInit() []Command {
	return c.Init
}
