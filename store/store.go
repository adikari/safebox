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

type ConfigId struct {
	Path string
	Key  string
}

type WriteConfigInput struct {
	ConfigId
	Value  string
	Secret bool
}

var (
	ConfigNotFoundError = errors.New("config not found")
)

type Store interface {
	Write(input WriteConfigInput) error
	// WriteMany(input []WriteConfigInput) error
	Read(id ConfigId) (Config, error)
	// ReadMany(id []ConfigId) ([]Config, error)
	// ReadAll() ([]Config, error)
	Delete(id ConfigId) error
}
