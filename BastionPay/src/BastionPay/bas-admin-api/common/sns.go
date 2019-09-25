package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsData struct {
	Recipient string
	Body      string
}

type Sns struct {
	*AwsConfig
	*SnsData
}

func NewSns(ac *AwsConfig) *Sns {
	return &Sns{
		AwsConfig: ac,
	}
}

func (s *Sns) Send(sd *SnsData) error {
	var err error

	s.SnsData = sd

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.Region),
		Credentials: credentials.NewStaticCredentials(s.AccessKeyId, s.SecretKey, "")})
	if err != nil {
		return err
	}

	svc := sns.New(sess)

	params := &sns.PublishInput{
		Message:     aws.String(s.Body),
		PhoneNumber: aws.String(s.Recipient),
	}
	_, err = svc.Publish(params)
	if err != nil {
		return err
	}

	return nil
}
