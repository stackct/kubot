package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name   string   `yaml:"name"`
	Groups []string `yaml:"groups"`
	Users  []string `yaml:"users"`
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
