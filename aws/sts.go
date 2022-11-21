package aws

import (
	"github.com/aws/aws-sdk-go/service/sts"
)

type Sts struct {
	client *sts.STS
}

var stsClient *sts.STS

func NewSts() Sts {
	if stsClient == nil {
		stsClient = sts.New(Session)
	}

	return Sts{client: stsClient}
}

func (s *Sts) GetCallerIdentity() (*sts.GetCallerIdentityOutput, error) {
	return s.client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}
