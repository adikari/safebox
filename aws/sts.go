package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Sts struct {
	client *sts.STS
}

func NewSts(session *session.Session) Sts {
	return Sts{client: sts.New(session)}
}

func (s *Sts) GetCallerIdentity() (*sts.GetCallerIdentityOutput, error) {
	return s.client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}
