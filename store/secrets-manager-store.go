package store

import (
	a "github.com/adikari/safebox/v2/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/pkg/errors"
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

func (s *SecretsManagerStore) Create(input ConfigInput) error {
	param := &secretsmanager.CreateSecretInput{
		Name:         aws.String(input.Name),
		SecretString: aws.String(input.Value),
	}

	if _, err := s.svc.CreateSecret(param); err != nil {
		return errors.Wrap(err, input.Name)
	}

	return nil
}

func (s *SecretsManagerStore) Update(input ConfigInput) error {
	param := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(input.Name),
		SecretString: aws.String(input.Value),
	}

	if _, err := s.svc.UpdateSecret(param); err != nil {
		return errors.Wrap(err, input.Name)
	}

	return nil
}

func (s *SecretsManagerStore) Put(input ConfigInput) error {
	found, _ := s.Get(input)

	var err error
	if found != nil {
		err = s.Update(input)
	} else {
		err = s.Create(input)
	}

	if err != nil {
		return errors.Wrap(err, input.Name)
	}

	return nil
}

func (s *SecretsManagerStore) PutMany(inputs []ConfigInput) error {
	for _, config := range inputs {
		if err := s.Put(config); err != nil {
			return err
		}
	}

	return nil
}

func (s *SecretsManagerStore) Get(input ConfigInput) (*Config, error) {
	param := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(input.Name),
	}

	result, err := s.svc.GetSecretValue(param)

	if err != nil {
		return nil, err
	}

	return &Config{
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
		if res != nil {
			result = append(result, *res)
		}
	}

	return result, nil
}

func (s *SecretsManagerStore) GetByPath(path string) ([]Config, error) {
	var result []Config

	input := &secretsmanager.ListSecretsInput{
		Filters: []*secretsmanager.Filter{
			{
				Key:    aws.String("name"),
				Values: []*string{aws.String(path)},
			},
		},
	}

	var recursiveGet func()
	recursiveGet = func() {
		resp, err := s.svc.ListSecrets(input)

		if err != nil {
			return
		}

		for _, secret := range resp.SecretList {
			result = append(result, Config{Name: secret.Name})
		}

		if resp.NextToken != nil {
			input.NextToken = resp.NextToken
			recursiveGet()
		}
	}

	recursiveGet()

	return result, nil
}

func (s *SecretsManagerStore) Delete(input ConfigInput) error {
	param := &secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(7),
		SecretId:             aws.String(input.Name),
	}

	if _, err := s.svc.DeleteSecret(param); err != nil {
		return err
	}

	return nil
}

func (s *SecretsManagerStore) DeleteMany(inputs []ConfigInput) error {
	if len(inputs) <= 0 {
		return nil
	}

	for _, input := range inputs {
		if err := s.Delete(input); err != nil {
			return err
		}
	}

	return nil
}
