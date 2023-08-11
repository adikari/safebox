package store

import (
	"errors"

	"github.com/sonyarouje/simdb"
)

var _ Store = &LocalStore{}

type LocalStore struct {
	db *simdb.Driver
}

type LocalStoreConfig struct {
	Path     string
	Filename string
}

func NewLocalStore(config LocalStoreConfig) (*LocalStore, error) {
	if config.Path == "" || config.Filename == "" {
		return nil, errors.New("path and filename is required for local store")
	}

	if db, err := simdb.New(config.Filename); err == nil {
		return &LocalStore{db: db}, nil
	}

	return nil, errors.New("Failed to initialize local db")
}

func (s *LocalStore) PutMany(input []ConfigInput) error {
	return nil
}

func (s *LocalStore) Put(input ConfigInput) error {
	return nil
}

func (s *LocalStore) Delete(input ConfigInput) error {
	return nil
}

func (s *LocalStore) DeleteMany(input []ConfigInput) error {
	return nil
}

func (s *LocalStore) GetMany(input []ConfigInput) ([]Config, error) {
	return []Config{}, nil
}

func (s *LocalStore) Get(input ConfigInput) (*Config, error) {
	return &Config{}, nil
}

func (s *LocalStore) GetByPath(path string) ([]Config, error) {
	return []Config{}, nil
}
