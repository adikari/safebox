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
	Put(input ConfigInput) error
	PutMany(input []ConfigInput) error
	Get(key string) (Config, error)
	GetMany(keys []string) ([]Config, error)
	GetAll() ([]Config, error)
	Delete(key string) error
}
