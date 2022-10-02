package config

import (
	"fmt"
	"io/ioutil"

	"github.com/adikari/safebox/v2/store"
	"gopkg.in/yaml.v2"
)

type rawConfig struct {
	Provider string
	Service  string

	Config map[string]map[string]string
	Secret map[string]string
}

type Config struct {
	Provider string
	Service  string
	Configs  []store.ConfigInput
}

type LoadParam struct {
	Path  string
	Stage string
}

func Load(param LoadParam) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(param.Path)

	if err != nil {
		return nil, fmt.Errorf("missing safebox config file %s", param.Path)
	}

	rc := rawConfig{}

	err = yaml.Unmarshal(yamlFile, &rc)

	if err != nil {
		return nil, fmt.Errorf("could not parse safebox config file %s", param.Path)
	}

	if rc.Service == "" {
		return nil, fmt.Errorf("'service' missing in config file%s", param.Path)
	}

	c := Config{}
	c.Service = rc.Service
	c.Provider = rc.Provider

	if c.Provider == "" {
		c.Provider = store.SsmProvider
	}

	parseConfig(rc, &c, param)

	return &c, nil
}

func parseConfig(rc rawConfig, c *Config, param LoadParam) {
	temp := map[string]string{}
	c.Configs = []store.ConfigInput{}

	defaultConfig := rc.Config["defaults"]
	sharedConfig := rc.Config["shared"]
	envConfig := rc.Config[param.Stage]

	for key, value := range defaultConfig {
		temp[formatPath(param.Stage, c.Service, key)] = value
	}

	for key, value := range sharedConfig {
		temp[formatSharedPath(param.Stage, key)] = value
	}

	for key, value := range envConfig {
		temp[formatPath(param.Stage, c.Service, key)] = value
	}

	for key, value := range temp {
		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   key,
			Value:  value,
			Secret: false,
		})
	}

	for key, value := range rc.Secret {
		c.Configs = append(c.Configs, store.ConfigInput{
			Name:        key,
			Description: value,
			Secret:      true,
		})
	}
}

func formatSharedPath(stage string, key string) string {
	return fmt.Sprintf("/%s/shared/%s", stage, key)
}

func formatPath(stage string, service string, key string) string {
	return fmt.Sprintf("/%s/%s/%s", stage, service, key)
}
