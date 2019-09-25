package sdk_aws_sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"strings"
)

type SnsData struct {
	Recipient string
	Body      string
}

type SnsSdk struct {
	mSnsClient *sns.SNS
	mSess      *session.Session
}

func NewSnsSdk(region, accessKeyID, accessKey, accessToken string) *SnsSdk {
	var sess *session.Session
	if (len(accessKeyID) == 0) && (len(accessKey) == 0) {
		sess = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(region),
		}))
	}else{
		sess = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKeyID, accessKey, accessToken),
		}))
	}

	client := sns.New(sess)
	sdk := &SnsSdk{
		mSess:      sess,
		mSnsClient: client,
	}
	return sdk
}

func (this *SnsSdk) Send(body, recipient string, senderId *string) error {
	params := &sns.PublishInput{
		Message:     aws.String(body),
		PhoneNumber: aws.String(recipient),
	}
	if senderId != nil && len(*senderId) > 0 {
		if params.MessageAttributes == nil {
			params.MessageAttributes = make(map[string]*sns.MessageAttributeValue)
		}
		senderIdStr := strings.Replace(*senderId, " ", "", -1)
		if len(senderIdStr) > 11 {
			senderIdStr = senderIdStr[:11]
		}
		if len(senderIdStr) > 0 {
			attV := new(sns.MessageAttributeValue).SetDataType("String").SetStringValue(senderIdStr)
			params.MessageAttributes["AWS.SNS.SMS.SenderID"] = attV
		}
	}
	_, err := this.mSnsClient.Publish(params)
	if err != nil {
		return err
	}
	return nil
}
