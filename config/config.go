package config

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/adikari/safebox/v2/aws"
	"github.com/adikari/safebox/v2/store"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type rawConfig struct {
	Provider string
	Service  string

	Config               map[string]map[string]string
	Secret               map[string]map[string]string
	CloudformationStacks []string `yaml:"cloudformation-stacks"`
}

type Config struct {
	Provider string
	Service  string
	Stage    string
	Prefix   string
	All      []store.ConfigInput
	Configs  []store.ConfigInput
	Secrets  []store.ConfigInput
	Stacks   []string
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

	err = validateConfig(rc)

	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	c := Config{}
	c.Service = rc.Service
	c.Stage = param.Stage
	c.Provider = rc.Provider
	c.Prefix = getPrefix(param.Stage, c.Service)

	if c.Provider == "" {
		c.Provider = store.SsmProvider
	}

	variables, err := loadVariables(c, rc)

	if err != nil {
		return nil, errors.Wrap(err, "failed to variables for interpolation")
	}

	for key, value := range rc.Config["defaults"] {
		val, err := Interpolate(value, variables)

		if err != nil {
			return nil, errors.Wrap(err, "failed to interpolate template variables")
		}

		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatPath(c.Prefix, key),
			Value:  val,
			Secret: false,
		})
	}

	for key, value := range rc.Config["shared"] {
		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatSharedPath(param.Stage, key),
			Value:  value,
			Secret: false,
		})
	}

	for key, value := range rc.Config[param.Stage] {
		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatPath(c.Prefix, key),
			Value:  value,
			Secret: false,
		})
	}

	c.Configs = removeDuplicate(c.Configs)

	for key, value := range rc.Secret["defaults"] {
		c.Secrets = append(c.Secrets, store.ConfigInput{
			Name:        formatPath(c.Prefix, key),
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

	return &c, nil
}

func formatSharedPath(stage string, key string) string {
	return fmt.Sprintf("/%s/shared/%s", stage, key)
}

func formatPath(prefix string, key string) string {
	return fmt.Sprintf("%s%s", prefix, key)
}

func getPrefix(stage string, service string) string {
	return fmt.Sprintf("/%s/%s/", stage, service)
}

func validateConfig(rc rawConfig) error {
	if rc.Service == "" {
		return fmt.Errorf("'service' is missing")
	}

	if rc.Provider == "" {
		return fmt.Errorf("'provider' is missing")
	}

	return nil
}

func loadVariables(c Config, rc rawConfig) (map[string]string, error) {
	st := aws.NewSts()

	id, err := st.GetCallerIdentity()

	if err != nil {
		return nil, err
	}

	variables := map[string]string{
		"stage":   c.Stage,
		"service": c.Service,
		"region":  *aws.Session.Config.Region,
		"account": *id.Account,
	}

	for _, name := range rc.CloudformationStacks {
		value, err := Interpolate(name, variables)
		if err != nil {
			return nil, err
		}
		c.Stacks = append(c.Stacks, value)
	}

	// add cloudformation outputs to variables available for interpolation
	if len(c.Stacks) > 0 {
		cf := aws.NewCloudformation()
		outputs, err := cf.GetOutputs(c.Stacks)

		if err != nil {
			return nil, err
		}

		for key, value := range outputs {
			variables[key] = value
		}
	}

	return variables, nil
}

func Interpolate(value string, variables map[string]string) (string, error) {
	var result bytes.Buffer
	tmpl, _ := template.New("interpolate").Option("missingkey=error").Parse(value)

	err := tmpl.Execute(&result, variables)

	if err != nil {
		return "", err
	}

	return result.String(), nil
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
