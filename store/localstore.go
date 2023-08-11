package store

import (
	"encoding/json"
	"errors"
	"fmt"
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

// Read a record from json file
func (s *LocalStore) Read() ([]Config, error) {
	if _, err := stat(s.path); err != nil {
		return []Config{}, err
	}

	b, err := ioutil.ReadFile(s.path)

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

func (s *LocalStore) PutMany(input []ConfigInput) error {
	configs := []Config{}

	for _, c := range input {
		t := "String"

		if c.Secret == true {
			t = "SecureString"
		}

		now := time.Now()

		configs = append(configs, Config{
			Name:     &c.Name,
			Value:    &c.Value,
			Version:  "1",
			Type:     t,
			Created:  now,
			Modified: now,
		})
	}

	return write(configs, s.path)
}

func (s *LocalStore) Put(input ConfigInput) error {
	return s.PutMany([]ConfigInput{input})
}

func (s *LocalStore) Delete(input ConfigInput) error {
	return errors.New("Delete is not implemented")
}

func (s *LocalStore) DeleteMany(input []ConfigInput) error {
	return errors.New("DeleteMany is not implemented")
}

func (s *LocalStore) GetMany(input []ConfigInput) ([]Config, error) {
	data, err := s.Read()

	if err != nil {
		return nil, nil
	}

	for _, d := range data {
		fmt.Printf("%v", *d.Value)
	}

	// records, err := s.db.ReadAll("s.filename")
	//
	// if err != nil {
	// 	return []Config{}, nil
	// }
	//
	// configs := []Config{}
	// for _, f := range records {
	// 	found := Config{}
	// 	if err := json.Unmarshal([]byte(f), &found); err != nil {
	// 		return []Config{}, err
	// 	}
	// 	configs = append(configs, found)
	// }

	return []Config{}, nil
}

func (s *LocalStore) Get(input ConfigInput) (*Config, error) {
	return &Config{}, errors.New("Get is not implemented")
}

func (s *LocalStore) GetByPath(path string) ([]Config, error) {
	return []Config{}, errors.New("Get by path not implemented")
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
