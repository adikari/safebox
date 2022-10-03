package store

import (
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
		Name:        aws.String(input.Name),
		Type:        aws.String(configType),
		Value:       aws.String(input.Value),
		Description: aws.String(input.Description),
		Overwrite:   aws.Bool(true),
	}

	_, err := s.svc.PutParameter(putParameterInput)

	if err != nil {
		return err
	}

	return nil
}

func (s *SSMStore) Delete(config ConfigInput) error {
	_, err := s.Get(config)

	if err != nil {
		return err
	}

	deleteParameterInput := &ssm.DeleteParameterInput{
		Name: aws.String(config.Name),
	}

	_, err = s.svc.DeleteParameter(deleteParameterInput)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSMStore) GetMany(configs []ConfigInput) ([]Config, error) {
	getParametersInput := &ssm.GetParametersInput{
		Names:          getNames(configs),
		WithDecryption: aws.Bool(true),
	}

	resp, err := s.svc.GetParameters(getParametersInput)

	if err != nil {
		return []Config{}, err
	}

	var params []Config

	for _, param := range resp.Parameters {
		params = append(params, parameterToConfig(param))
	}

	return params, nil
}

func (s *SSMStore) Get(config ConfigInput) (Config, error) {
	configs, err := s.GetMany([]ConfigInput{config})

	if err != nil {
		return Config{}, err
	}

	return configs[0], nil
}

func basePath(key string) string {
	pathParts := strings.Split(key, "/")
	if len(pathParts) == 1 {
		return pathParts[0]
	}
	end := len(pathParts) - 1
	return strings.Join(pathParts[0:end], "/")
}

func parameterToConfig(param *ssm.Parameter) Config {
	return Config{
		Name:     param.Name,
		Value:    param.Value,
		Modified: *param.LastModifiedDate,
		Version:  int(*param.Version),
		Type:     *param.Type,
		DataType: *param.DataType,
	}
}

func getNames(configs []ConfigInput) []*string {
	var keys []string

	for _, value := range configs {
		keys = append(keys, value.Name)
	}

	var names []*string

	for _, key := range keys {
		names = append(names, aws.String(key))
	}

	return names
}