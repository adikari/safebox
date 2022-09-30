package store

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var _ Store = &SSMStore{}

var (
	numberOfRetries = 10
	throttleDelay   = client.DefaultRetryerMinRetryDelay
)

type SSMStore struct {
	svc ssmiface.SSMAPI
}

func NewSSMStore() (*SSMStore, error) {
	ssmSession := session.Must(session.NewSession())

	retryer := client.DefaultRetryer{
		NumMaxRetries:    numberOfRetries,
		MinThrottleDelay: throttleDelay,
	}

	svc := ssm.New(ssmSession, &aws.Config{
		Retryer: retryer,
	})

	return &SSMStore{
		svc: svc,
	}, nil
}

func (s *SSMStore) PutMany(input []ConfigInput) error {
	for _, config := range input {
		err := s.Put(config)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SSMStore) Put(input ConfigInput) error {
	configType := "String"

	if input.Secret == true {
		configType = "SecureString"
	}

	putParameterInput := &ssm.PutParameterInput{
		Name:      aws.String(input.Key),
		Type:      aws.String(configType),
		Value:     aws.String(input.Value),
		Overwrite: aws.Bool(true),
	}

	_, err := s.svc.PutParameter(putParameterInput)

	if err != nil {
		return err
	}

	return nil
}

func (s *SSMStore) Delete(key string) error {
	_, err := s.Get(key)

	if err != nil {
		return err
	}

	deleteParameterInput := &ssm.DeleteParameterInput{
		Name: aws.String(key),
	}

	_, err = s.svc.DeleteParameter(deleteParameterInput)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSMStore) GetMany(keys []string) ([]Config, error) {
	return nil, nil
}

func (s *SSMStore) GetAll() ([]Config, error) {
	return nil, nil
}

func (s *SSMStore) Get(key string) (Config, error) {
	getParametersInput := &ssm.GetParametersInput{
		Names:          []*string{aws.String(key)},
		WithDecryption: aws.Bool(true),
	}

	resp, err := s.svc.GetParameters(getParametersInput)

	if err != nil {
		return Config{}, err
	}

	if len(resp.Parameters) == 0 {
		return Config{}, ConfigNotFoundError
	}

	param := resp.Parameters[0]
	var parameter *ssm.ParameterMetadata
	var describeParametersInput *ssm.DescribeParametersInput

	// There is no way to use describe parameters to get a single key
	// if that key uses paths, so instead get all the keys for a path,
	// then find the one you are looking for :(
	describeParametersInput = &ssm.DescribeParametersInput{
		ParameterFilters: []*ssm.ParameterStringFilter{
			{
				Key:    aws.String("Path"),
				Option: aws.String("OneLevel"),
				Values: []*string{aws.String(basePath(key))},
			},
		},
	}

	if err := s.svc.DescribeParametersPages(describeParametersInput, func(o *ssm.DescribeParametersOutput, lastPage bool) bool {
		for _, param := range o.Parameters {
			if *param.Name == key {
				parameter = param
				return false
			}
		}
		return true
	}); err != nil {
		return Config{}, err
	}

	if parameter == nil {
		return Config{}, ConfigNotFoundError
	}

	return Config{
		Value:    param.Value,
		Metadata: mapMetadata(parameter),
	}, nil
}

func basePath(key string) string {
	pathParts := strings.Split(key, "/")
	if len(pathParts) == 1 {
		return pathParts[0]
	}
	end := len(pathParts) - 1
	return strings.Join(pathParts[0:end], "/")
}

func mapMetadata(p *ssm.ParameterMetadata) Metadata {
	version := 0
	if p.Description != nil {
		version, _ = strconv.Atoi(*p.Description)
	}
	return Metadata{
		Created:   *p.LastModifiedDate,
		CreatedBy: *p.LastModifiedUser,
		Version:   version,
		Key:       *p.Name,
	}
}
