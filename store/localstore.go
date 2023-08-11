package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var _ Store = &LocalStore{}

type LocalStore struct {
	filename string
	dir      string
	path     string
}

type LocalStoreConfig struct {
	Directory string
	Filename  string
}

func NewLocalStore(config LocalStoreConfig) (*LocalStore, error) {
	if config.Directory == "" || config.Filename == "" {
		return nil, errors.New("path and filename is required for local store")
	}

	dir := filepath.Clean(config.Directory)

	store := &LocalStore{
		filename: config.Filename,
		dir:      dir,
		path:     filepath.Join(dir, config.Filename),
	}

	if _, err := os.Stat(dir); err == nil {
		return store, nil
	}

	err := os.MkdirAll(dir, 0755)

	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *LocalStore) PutMany(input []ConfigInput) error {
	updates := []Config{}
	for _, c := range input {
		t := "String"

		if c.Secret == true {
			t = "SecureString"
		}

		now := time.Now()

		name := c.Name
		value := c.Value

		updates = append(updates, Config{
			Name:     &name,
			Value:    &value,
			Version:  c.Name,
			Type:     t,
			Created:  now,
			Modified: now,
		})
	}

	existing, _ := read(s.path)
	for _, e := range existing {
		found := find(*e.Name, updates)
		if found == nil {
			updates = append(updates, e)
		}
	}

	return write(updates, s.path)
}

func find(id string, all []Config) *Config {
	for _, c := range all {
		if *c.Name == id {
			return &c
		}
	}

	return nil
}

func (s *LocalStore) Put(input ConfigInput) error {
	return s.PutMany([]ConfigInput{input})
}

func (s *LocalStore) DeleteMany(input []ConfigInput) error {
	return errors.New("DeleteMany is not implemented")
}

func (s *LocalStore) GetMany(input []ConfigInput) ([]Config, error) {
	configs, err := read(s.path)

	if err != nil {
		return nil, nil
	}

	return configs, nil
}

func (s *LocalStore) Get(input ConfigInput) (*Config, error) {
	if configs, _ := s.GetMany([]ConfigInput{input}); configs != nil && len(configs) > 0 {
		return &configs[0], nil
	}

	return nil, nil
}

func (s *LocalStore) GetByPath(path string) ([]Config, error) {
	return []Config{}, errors.New("Get by path not implemented")
}

// Read a record from json file
func read(path string) ([]Config, error) {
	if _, err := stat(path); err != nil {
		return []Config{}, err
	}

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return []Config{}, err
	}

	configs := []Config{}

	err = json.Unmarshal(b, &configs)

	if err != nil {
		return nil, errors.New("failed to parse data in database")
	}

	return configs, nil
}

func write(configs []Config, path string) error {
	b, err := json.MarshalIndent(configs, "", "\t")

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		return err
	}

	return nil
}

func stat(path string) (fi os.FileInfo, err error) {
	// check for dir, if path isn't a directory check to see if it's a file
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path)
	}

	return
}
