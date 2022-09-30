package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type rawConfig struct {
	Provider string
	Service  string

	Config map[string]map[string]string
}

type Config struct {
	Provider string
	Service  string
	Config   map[string]string
}

type LoadParam struct {
	Path  string
	Stage string
}

func Load(param LoadParam) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(param.Path)

	if err != nil {
		return nil, errors.New("Could not find config file: " + param.Path)
	}

	rc := rawConfig{}

	err = yaml.Unmarshal(yamlFile, &rc)

	if err != nil {
		return nil, errors.New("Could not parse config file")
	}

	c := Config{}
	c.Service = rc.Service
	c.Provider = rc.Provider
	parseConfig(rc, &c, param)

	return &c, nil
}

func parseConfig(rc rawConfig, c *Config, param LoadParam) {
	c.Config = map[string]string{}

	defaultConfig := rc.Config["defaults"]
	sharedConfig := rc.Config["shared"]
	envConfig := rc.Config[param.Stage]

	for key, value := range defaultConfig {
		c.Config[formatKey(param.Stage, c.Service, key)] = value
	}

	for key, value := range sharedConfig {
		c.Config[formatSharedKey(param.Stage, key)] = value
	}

	for key, value := range envConfig {
		c.Config[formatKey(param.Stage, c.Service, key)] = value
	}
}

func formatSharedKey(stage string, key string) string {
	return fmt.Sprintf("/%s/shared/%s", stage, key)
}

func formatKey(stage string, service string, key string) string {
	return fmt.Sprintf("/%s/%s/%s", stage, service, key)
}
