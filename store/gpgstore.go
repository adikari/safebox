package store

var _ Store = &GpgStore{}

type GpgStore struct{}

func NewGpgStore() (*GpgStore, error) {
	return &GpgStore{}, nil
}

func (s *GpgStore) PutMany(input []ConfigInput) error {
	return nil
}

func (s *GpgStore) Put(input ConfigInput) error {
	return nil
}

func (s *GpgStore) Delete(input ConfigInput) error {
	return nil
}

func (s *GpgStore) DeleteMany(input []ConfigInput) error {
	return nil
}

func (s *GpgStore) GetMany(input []ConfigInput) ([]Config, error) {
	return []Config{}, nil
}

func (s *GpgStore) Get(input ConfigInput) (*Config, error) {
	return &Config{}, nil
}

func (s *GpgStore) GetByPath(path string) ([]Config, error) {
	return []Config{}, nil
}
