package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var cfClient *cloudformation.CloudFormation

type Cloudformation struct {
	client *cloudformation.CloudFormation
}

func NewCloudformation(session *session.Session) Cloudformation {
	if cfClient == nil {
		cfClient = cloudformation.New(session)
	}

	return Cloudformation{client: cfClient}
}

func (c *Cloudformation) GetOutput(stackname string) (map[string]string, error) {
	result := map[string]string{}

	resp, _ := c.client.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackname),
	})

	if len(resp.Stacks) <= 0 {
		return nil, fmt.Errorf("%s stack does not exist", stackname)
	}

	stack := resp.Stacks[0]

	for _, output := range stack.Outputs {
		result[*output.OutputKey] = *output.OutputValue
	}

	return result, nil
}

func (c *Cloudformation) GetOutputs(stacknames []string) (map[string]string, error) {
	result := map[string]string{}

	for _, stackname := range stacknames {
		outputs, err := c.GetOutput(stackname)

		if err != nil {
			continue
		}

		for key, value := range outputs {
			result[key] = value
		}
	}

	return result, nil
}
