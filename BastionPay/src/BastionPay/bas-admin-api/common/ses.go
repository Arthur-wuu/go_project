package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SesData struct {
	Recipient string
	Body      string
	Subject   string
	Sender    string
	CharSet   string
}

type Ses struct {
	*AwsConfig
	*SesData
}

func NewSes(ac *AwsConfig) *Ses {
	return &Ses{
		AwsConfig: ac,
	}
}

func (s *Ses) Send(sd *SesData) error {
	s.SesData = sd

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.Region),
		Credentials: credentials.NewStaticCredentials(s.AccessKeyId, s.SecretKey, "")})
	if err != nil {
		return err
	}

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(s.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(s.CharSet),
					Data:    aws.String(s.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(s.CharSet),
				Data:    aws.String(s.Subject),
			},
		},
		Source: aws.String(s.Sender),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)
	if err != nil {
		return err
	}

	return nil
}
