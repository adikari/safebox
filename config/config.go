package config

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/adikari/safebox/v2/aws"
	"github.com/adikari/safebox/v2/store"
	"github.com/adikari/safebox/v2/util"
	a "github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type rawConfig struct {
	Provider             string
	Service              string
	Prefix               string
	Generate             []Generate `yaml:"generate"`
	Config               map[string]map[string]string
	Secret               map[string]map[string]string
	CloudformationStacks []string `yaml:"cloudformation-stacks"`
	Region               string   `yaml:"region"`
	DBDir                string   `yaml:"db_dir"`
}

type Config struct {
	Provider string
	Service  string
	Stage    string
	Prefix   string
	Generate []Generate
	Region   string
	All      []store.ConfigInput
	Configs  []store.ConfigInput
	Secrets  []store.ConfigInput
	Stacks   []string
	Filepath string
}

type Generate struct {
	Type string
	Path string
}

type LoadConfigInput struct {
	Path  string
	Stage string
}

var defaultConfigPaths = []string{"safebox.yml", "safebox.yaml"}

func Load(param LoadConfigInput) (*Config, error) {
	yamlFile, err := readConfigFile(param.Path)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	rc := rawConfig{}

	err = yaml.Unmarshal(yamlFile, &rc)

	if err != nil {
		fmt.Printf("%v", err)
		return nil, fmt.Errorf("could not parse safebox config file %s", param.Path)
	}

	err = validateConfig(rc)

	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	c := Config{
		Service:  rc.Service,
		Stage:    param.Stage,
		Provider: rc.Provider,
	}

	if c.Provider == "" {
		c.Provider = util.SsmProvider
	}

	if c.Provider == util.GpgProvider {
		c.Filepath = getFilePath(c, rc)
	}

	variables, err := loadVariables(&c, rc)

	if c.Region == "" {
		c.Region = "local"
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to load variables for interpolation")
	}

	c.Prefix, err = Interpolate(getPrefix(param.Stage, c.Service, rc.Prefix), variables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to interpolate prefix")
	}

	for key, value := range rc.Config["defaults"] {
		val, err := Interpolate(value, variables)

		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to interpolate config.defaults.%s", key))
		}

		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatPath(c.Prefix, key),
			Value:  val,
			Secret: false,
		})
	}

	for _, value := range rc.Generate {
		path, err := Interpolate(value.Path, variables)

		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to interpolate generate: type: %s, path: %s", value.Type, value.Path))
		}

		c.Generate = append(c.Generate, Generate{
			Type: value.Type,
			Path: path,
		})
	}

	for key, value := range rc.Config["shared"] {
		val, err := Interpolate(value, variables)

		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to interpolate config.shared.%s", key))
		}

		c.Configs = append(c.Configs, store.ConfigInput{
			Name:   formatSharedPath(param.Stage, key),
			Value:  val,
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
	if stage != "" {
		return fmt.Sprintf("/%s/shared/%s", stage, key)
	}
	return fmt.Sprintf("/shared/%s", key)
}

func formatPath(prefix string, key string) string {
	return fmt.Sprintf("%s%s", prefix, key)
}

func getPrefix(stage string, service string, defaultPrefix string) string {
	if defaultPrefix != "" {
		// TODO: validate prefix starts and ends with /
		return defaultPrefix
	}

	if stage != "" {
		return fmt.Sprintf("/%s/%s/", stage, service)
	}

	return fmt.Sprintf("/%s/", service)
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

// loadVariables for interpolation
// TODO: in future as we support more stores, this many need to be refactored to handled each
func loadVariables(c *Config, rc rawConfig) (map[string]string, error) {
	if !util.IsAwsProvider(c.Provider) {
		return map[string]string{}, nil
	}

	session := aws.NewSession(a.Config{Region: &rc.Region})
	st := aws.NewSts(session)
	c.Region = *session.Config.Region

	id, err := st.GetCallerIdentity()

	if err != nil {
		return nil, errors.New("Failed to login to AWS")
	}

	variables := map[string]string{
		"stage":   c.Stage,
		"service": c.Service,
		"region":  c.Region,
		"account": *id.Account,
	}

	for _, name := range rc.CloudformationStacks {
		value, err := Interpolate(name, variables)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to interpolate cloudformation-stacks[%s]", name))
		}
		c.Stacks = append(c.Stacks, value)
	}

	// add cloudformation outputs to variables available for interpolation
	if len(c.Stacks) > 0 {
		cf := aws.NewCloudformation(session)
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

func readConfigFile(path string) ([]byte, error) {
	if path != "" {
		s, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("missing file %s", path)
		}
		return s, nil
	}

	for _, c := range defaultConfigPaths {
		if s, err := ioutil.ReadFile(c); err == nil {
			return s, nil
		}
	}

	return nil, fmt.Errorf("missing file %s", strings.Join(defaultConfigPaths, " or "))
}

func getFilePath(config Config, rc rawConfig) string {
	d := rc.DBDir
	if d == "" {
		ex, err := os.Executable()
		exPath := "."
		if err == nil {
			exPath = filepath.Dir(ex)
		}
		d = exPath
	}

	dir := filepath.Clean(d)

	usr, _ := user.Current()
	homedir := usr.HomeDir
	if dir == "~" {
		dir = homedir
	} else if strings.HasPrefix(dir, "~/") {
		dir = filepath.Join(homedir, dir[2:])
	}

	filename := fmt.Sprintf("%s-%s", config.Stage, config.Service)
	if config.Stage == "" {
		filename = fmt.Sprintf("%s", config.Service)
	}

	return filepath.Join(dir, filename)
}
