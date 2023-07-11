package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
)

var ses *session.Session

func NewSession(cfg aws.Config) *session.Session {
	if ses == nil {
		if cfg.Retryer == nil {
			cfg.Retryer = Retryer
		}

		if cfg.Region == nil {
			region := os.Getenv("AWS_REGION")
			cfg.Region = &region
		}

		ses = session.Must(session.NewSession(&cfg))
	}

	return ses
}

var Retryer = client.DefaultRetryer{
	NumMaxRetries:    2,
	MinThrottleDelay: client.DefaultRetryerMaxRetryDelay,
}
