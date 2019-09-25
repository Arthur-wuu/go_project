package email

import (
	"BastionPay/bas-notify/config"
	ses "BastionPay/bas-tools/sdk.aws.ses"
	"fmt"
)

var GMailMgr MailMgr

type MailMgr struct {
	mSdk     *ses.SesSdk
	mSrcMail string
}

func (this *MailMgr) Init() error {
	awsConfig := config.GConfig.Aws
	this.mSrcMail = awsConfig.SrcEMail
	this.mSdk = ses.NewSesSdk(awsConfig.SesRegion, awsConfig.Accesskeyid, awsConfig.Accesskey, awsConfig.Accesstoken)
	return nil
}

func (this *MailMgr) DirectSend(subject, body, toEmail string, senderId *string) error {
	srcMail := this.mSrcMail
	if senderId != nil && len(*senderId) > 1 {
		srcMail = fmt.Sprintf("%s <%s>", *senderId, this.mSrcMail)
	}
	err := this.mSdk.Send(toEmail, subject, srcMail, "UTF-8", body)
	if err != nil {
		return err
	}
	return nil
}
