package store

import (
	a "github.com/adikari/safebox/v2/aws"
	"github.com/aws/aws-sdk-go/aws"
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
		secretsmanagerService = secretsmanager.New(a.Session, &aws.Config{
			Retryer: a.Retryer,
		})
	}

	return &SecretsManagerStore{
		svc: secretsmanagerService,
	}, nil
}

func (s *SecretsManagerStore) Put(input ConfigInput) error {
	param := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(input.Name),
		SecretString: aws.String(input.Value),
	}

	_, err := s.svc.PutSecretValue(param)

	if err != nil {
		return err
	}

	return nil
}

func (s *SecretsManagerStore) PutMany(inputs []ConfigInput) error {
	for _, config := range inputs {
		err := s.Put(config)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SecretsManagerStore) Get(input ConfigInput) (Config, error) {
	param := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(input.Name),
	}

	result, err := s.svc.GetSecretValue(param)

	if err != nil {
		return Config{}, err
	}

	return Config{
		Name:     result.Name,
		Value:    result.SecretString,
		Version:  *result.VersionId,
		Type:     "SecureString",
		DataType: "SecureString",
		Modified: *result.CreatedDate,
	}, nil
}

func (s *SecretsManagerStore) GetMany(inputs []ConfigInput) ([]Config, error) {
	if len(inputs) <= 0 {
		return []Config{}, nil
	}

	result := []Config{}

	for _, input := range inputs {
		res, _ := s.Get(input)
		result = append(result, res)
	}

	return result, nil
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
