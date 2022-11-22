package store

import (
	a "github.com/adikari/safebox/v2/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

var _ Store = &SecretsManagerStore{}

type SecretsManagerStore struct {
	svc secretsmanageriface.SecretsManagerAPI
}

var secretsmanagerService *secretsmanager.SecretsManager

func NewSecretsManagerStore() (*SecretsManagerStore, error) {
	if secretsmanagerService == nil {
		retryer := client.DefaultRetryer{
			NumMaxRetries:    numberOfRetries,
			MinThrottleDelay: throttleDelay,
		}

		secretsmanagerService = secretsmanager.New(a.Session, &aws.Config{
			Retryer: retryer,
		})
	}

	return &SecretsManagerStore{
		svc: secretsmanagerService,
	}, nil
}

func (s *SecretsManagerStore) Put(input ConfigInput) error {
	return nil
}

func (s *SecretsManagerStore) PutMany(inputs []ConfigInput) error {
	return nil
}

func (s *SecretsManagerStore) Get(input ConfigInput) (Config, error) {
	return Config{}, nil
}

func (s *SecretsManagerStore) GetMany(inputs []ConfigInput) ([]Config, error) {
	return []Config{}, nil
}

func (s *SecretsManagerStore) GetByPath(path string) ([]Config, error) {
	return []Config{}, nil
}

func (s *SecretsManagerStore) Delete(input ConfigInput) error {
	return nil
}

func (s *SecretsManagerStore) DeleteMany(inputs []ConfigInput) error {
	return nil
}
