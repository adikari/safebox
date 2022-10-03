package config

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/adikari/safebox/v2/store"
	"gopkg.in/yaml.v2"
)

type rawConfig struct {
	Provider string
	Service  string

	Config map[string]map[string]string
	Secret map[string]map[string]string
}

type Config struct {
	Provider string
	Service  string
	All      []store.ConfigInput
	Configs  []store.ConfigInput
	Secrets  []store.ConfigInput
}

type LoadConfigInput struct {
	Path  string
	Stage string
}

func Load(param LoadConfigInput) (*Config, error) {
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

func parseConfig(rc rawConfig, c *Config, param LoadConfigInput) {
	variables := map[string]string{
		"stage":   param.Stage,
		"service": c.Service,
	}

	for key, value := range rc.Config["defaults"] {
		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatPath(param.Stage, c.Service, key),
			Value:  interpolate(value, variables),
			Secret: false,
		})
	}

	for key, value := range rc.Config["shared"] {
		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatSharedPath(param.Stage, key),
			Value:  interpolate(value, variables),
			Secret: false,
		})
	}

	for key, value := range rc.Config[param.Stage] {
		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatPath(param.Stage, c.Service, key),
			Value:  interpolate(value, variables),
			Secret: false,
		})
	}

	c.Configs = removeDuplicate(c.Configs)

	for key, value := range rc.Secret["defaults"] {
		c.Secrets = append(c.Secrets, store.ConfigInput{
			Name:        formatPath(param.Stage, c.Service, key),
			Description: value,
			Secret:      true,
		})
	}

	for key, value := range rc.Secret["shared"] {
		c.Secrets = append(c.Secrets, store.ConfigInput{
			Name:        formatSharedPath(param.Stage, key),
			Description: value,
			Secret:      true,
		})
	}

	c.All = append(c.Secrets, c.Configs...)
}

func formatSharedPath(stage string, key string) string {
	return fmt.Sprintf("/%s/shared/%s", stage, key)
}

func formatPath(stage string, service string, key string) string {
	return fmt.Sprintf("/%s/%s/%s", stage, service, key)
}

func interpolate(value string, variables map[string]string) string {
	var result bytes.Buffer
	tmpl, _ := template.New("interpolate").Parse(value)

	tmpl.Execute(&result, variables)
	return result.String()
}

func removeDuplicate(input []store.ConfigInput) []store.ConfigInput {
	var unique []store.ConfigInput

loop:
	for _, v := range input {
		for i, u := range unique {
			if v.Name == u.Name {
				unique[i] = v
				continue loop
			}
		}
		unique = append(unique, v)
	}

	return unique
}
