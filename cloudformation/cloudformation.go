package cloudformation

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var (
	numberOfRetries = 10
	throttleDelay   = client.DefaultRetryerMinRetryDelay
)

type Cloudformation struct {
	client *cloudformation.CloudFormation
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

func NewCloudformation() Cloudformation {
	cfSession := session.Must(session.NewSession())

	retryer := client.DefaultRetryer{
		NumMaxRetries:    numberOfRetries,
		MinThrottleDelay: throttleDelay,
	}

	c := cloudformation.New(cfSession, &aws.Config{
		Retryer: retryer,
	})

	return Cloudformation{client: c}
}
