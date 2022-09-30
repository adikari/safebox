package store

import (
	"errors"
	"time"
)

type Metadata struct {
	Created   time.Time
	CreatedBy string
	Version   int
	Key       string
}

type Config struct {
	Value *string
	Metadata
}

type ConfigInput struct {
	Key    string
	Value  string
	Secret bool
}

var (
	ConfigNotFoundError = errors.New("config not found")
)

type Store interface {
	Write(input ConfigInput) error
	WriteMany(input []ConfigInput) error
	Read(key string) (Config, error)
	ReadMany(keys []string) ([]Config, error)
	ReadAll() ([]Config, error)
	Delete(key string) error
}
