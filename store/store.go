package store

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Name     *string
	Value    *string
	Modified time.Time
	Version  int
	Type     string
	DataType string
}

const (
	SsmProvider            = "ssm"
	SecretsManagerProvider = "secrets-manager"
)

type ConfigInput struct {
	Name        string
	Value       string
	Secret      bool
	Description string
}

var (
	ConfigNotFoundError = errors.New("config not found")
)

type Store interface {
	Put(input ConfigInput) error
	PutMany(input []ConfigInput) error
	Get(input ConfigInput) (Config, error)
	GetMany(inputs []ConfigInput) ([]Config, error)
	GetByPath(path string) ([]Config, error)
	Delete(input ConfigInput) error
	DeleteMany(inputs []ConfigInput) error
}

func GetStore(provider string) (Store, error) {
	switch provider {
	case SsmProvider:
		return NewSSMStore()
	case SecretsManagerProvider:
		return NewSecretsManagerStore()
	default:
		return nil, fmt.Errorf("invalid provider `%s`", provider)
	}
}

func (c *Config) Key() string {
	parts := strings.Split(*c.Name, "/")
	return parts[len(parts)-1]
}

func (c *ConfigInput) Key() string {
	parts := strings.Split(c.Name, "/")
	return parts[len(parts)-1]
}

func (c *Config) Path() string {
	parts := strings.Split(*c.Name, "/")
	return strings.Join(parts[0:len(parts)-1], "/")
}
