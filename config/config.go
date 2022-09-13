package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigParam struct {
	Name        string
	Description string
	Value       string
	Override    map[string]string
	Secret      bool
	Required    bool
}

type Config struct {
	Provider string

	Defaults struct {
		Path string
	}

	Params []ConfigParam
}

func loadConfig() (*Config, error) {
	path := "example/safebox.yml"
	yamlFile, err := ioutil.ReadFile(path)

	config := Config{}

	if err != nil {
		return nil, errors.New("Could not find config file: " + path)
	}

	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		return nil, errors.New("Could not parse config file")
	}

	return &config, nil
}

func GetConfig() []ConfigParam {
	return nil
}
