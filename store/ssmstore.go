package store

import (
	"fmt"

	a "github.com/adikari/safebox/v2/aws"
	"github.com/adikari/safebox/v2/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var _ Store = &SSMStore{}

var svc *ssm.SSM

type SSMStore struct {
	svc ssmiface.SSMAPI
}

func NewSSMStore() (*SSMStore, error) {
	if svc == nil {
		svc = ssm.New(a.Session, &aws.Config{
			Retryer: a.Retryer,
		})
	}

	return &SSMStore{
		svc: svc,
	}, nil
}

func (s *SSMStore) PutMany(input []ConfigInput) error {
	for _, config := range input {
		if err := s.Put(config); err != nil {
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
	if _, err := s.Get(config); err != nil {
		return err
	}

	deleteParameterInput := &ssm.DeleteParameterInput{
		Name: aws.String(config.Name),
	}

	if _, err := s.svc.DeleteParameter(deleteParameterInput); err != nil {
		return err
	}

	return nil
}

func (s *SSMStore) DeleteMany(configs []ConfigInput) error {
	if len(configs) <= 0 {
		return nil
	}

	for _, config := range configs {
		if err := s.Delete(config); err != nil {
			return err
		}
	}

	return nil
}

func (s *SSMStore) GetMany(configs []ConfigInput) ([]Config, error) {
	if len(configs) <= 0 {
		return []Config{}, nil
	}

	get := func(c []ConfigInput) ([]Config, error) {
		var result []Config

		getParametersInput := &ssm.GetParametersInput{
			Names:          getNames(c),
			WithDecryption: aws.Bool(true),
		}

		resp, err := s.svc.GetParameters(getParametersInput)

		if err != nil {
			return []Config{}, err
		}

		for _, param := range resp.Parameters {
			result = append(result, parameterToConfig(param))
		}

		return result, nil
	}

	var params []Config
	for _, chunk := range util.ChunkSlice(configs, 10) {
		p, err := get(chunk)

		if err != nil {
			return []Config{}, err
		}

		params = append(params, p...)
	}

	return params, nil
}

func (s *SSMStore) Get(config ConfigInput) (*Config, error) {
	configs, err := s.GetMany([]ConfigInput{config})

	if err != nil {
		return nil, err
	}

	return &configs[0], nil
}

func (s *SSMStore) GetByPath(path string) ([]Config, error) {
	var result []Config

	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		WithDecryption: aws.Bool(true),
	}

	var recursiveGet func()
	recursiveGet = func() {
		resp, err := s.svc.GetParametersByPath(input)

		if err != nil {
			return
		}

		for _, param := range resp.Parameters {
			result = append(result, parameterToConfig(param))
		}

		if resp.NextToken != nil {
			input.NextToken = resp.NextToken
			recursiveGet()
		}
	}

	recursiveGet()

	return result, nil
}

func parameterToConfig(param *ssm.Parameter) Config {
	return Config{
		Name:     param.Name,
		Value:    param.Value,
		Modified: *param.LastModifiedDate,
		Version:  fmt.Sprint(*param.Version),
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
