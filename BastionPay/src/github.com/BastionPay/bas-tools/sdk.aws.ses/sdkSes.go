package sdk_aws_ses

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	pkgSess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"time"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

/*
参考网址：https://docs.aws.amazon.com/zh_cn/ses/latest/DeveloperGuide/send-personalized-email-api.html
https://docs.aws.amazon.com/zh_cn/ses/latest/DeveloperGuide/send-personalized-email-api.html
*/

func NewSesSdk(region, accessKeyID, accessKey, accessToken string) *SesSdk {
	var sess *pkgSess.Session
	if (len(accessKeyID) == 0) && (len(accessKey) == 0) {
		sess = pkgSess.Must(pkgSess.NewSession(&aws.Config{
			Region:      aws.String(region),
		}))
	}else{
		sess = pkgSess.Must(pkgSess.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKeyID, accessKey, accessToken),
		}))
	}
	client := ses.New(sess)
	sdk := &SesSdk{
		mSess:      sess,
		mSesClient: client,
	}
	return sdk
}

type SesSdk struct {
	mSess      *pkgSess.Session
	mSesClient *ses.SES
}

func (this *SesSdk) Close() {

}

//阻塞是一定的，考虑用单独go程 10m大小限制, templateDate是替换模板部分内容的
func (this *SesSdk) SendTemplate(toAddr, replyAddr, srcAddr, templateName, templateData string, timeout int) error {
	if timeout <= 0 {
		timeout = 60 //default 60s as http timeout
	}
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancelFn()

	des := &ses.Destination{
		ToAddresses: []*string{
			aws.String(toAddr),
		},
	}
	var reply []*string
	if len(replyAddr) != 0 {
		reply = []*string{aws.String(replyAddr)}
	}

	input := &ses.SendTemplatedEmailInput{
		Destination:      des,
		ReplyToAddresses: reply,
		Source:           aws.String(srcAddr),
		Template:         aws.String(templateName),
		TemplateData:     aws.String(templateData),
	}
	_, err := this.mSesClient.SendTemplatedEmailWithContext(ctx, input)
	if err != nil { //CanceledErrorCode说明超时了
		return err
	}
	return nil
}

func (this *SesSdk) Send(toMail, subject, sender, charSet, body string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(toMail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the email.
	_, err := this.mSesClient.SendEmail(input)
	if err != nil {
		return err
	}
	return nil
}
