package aws

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
)

var Session = session.Must(session.NewSession())

var (
	numberOfRetries = 10
	throttleDelay   = client.DefaultRetryerMinRetryDelay
)

var Retryer = client.DefaultRetryer{
	NumMaxRetries:    numberOfRetries,
	MinThrottleDelay: throttleDelay,
}
