package aws

import "github.com/aws/aws-sdk-go/aws/session"

var Session = session.Must(session.NewSession())
