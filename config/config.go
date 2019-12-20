package config

import yaml "gopkg.in/yaml.v2"

type Config struct {
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name   string   `yaml:"name"`
	Groups []string `yaml:"groups"`
	Users  []string `yaml:"users"`
}

func Parse(bs []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(bs, &config)

	return config, err
}
