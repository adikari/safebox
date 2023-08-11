package store

var _ Store = &GpgStore{}

type GpgStore struct {
	svc *LocalStore
}

func NewGpgStore(config LocalStoreConfig) (*GpgStore, error) {
	svc, err := NewLocalStore(config)

	if err != nil {
		return nil, err
	}

	return &GpgStore{svc: svc}, nil
}

func (s *GpgStore) PutMany(input []ConfigInput) error {
	return s.svc.PutMany(input)
}

func (s *GpgStore) Put(input ConfigInput) error {
	return s.svc.Put(input)
}

func (s *GpgStore) Delete(input ConfigInput) error {
	return s.svc.Delete(input)
}

func (s *GpgStore) DeleteMany(input []ConfigInput) error {
	return s.svc.DeleteMany(input)
}

func (s *GpgStore) GetMany(input []ConfigInput) ([]Config, error) {
	return s.svc.GetMany(input)
}

func (s *GpgStore) Get(input ConfigInput) (*Config, error) {
	return s.svc.Get(input)
}

func (s *GpgStore) GetByPath(path string) ([]Config, error) {
	return s.svc.GetByPath(path)
}
