package models

import (
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	ses "github.com/BastionPay/bas-tools/sdk.aws.ses"
	sns "github.com/BastionPay/bas-tools/sdk.aws.sns"
	l4g "github.com/alecthomas/log4go"
	"github.com/caibirdme/yql"
)

var GlobalNotifyMgr NotifyMgr

const (
	ConstUserEmailElemName = "email"
	ConstSysEmailElemName  = "sysemail"
	ConstUserPhoneElemName = "phone"
	ConstSysPhoneElemName  = "sysphone"
)

//注意点：yql的in元素必须大于等于2

//level=int,behavior=string,textid=int,sub=bool
type NotifyMgr struct {
	mUserNotifyRaws []string //是否发短信匹配
	mSysNotifyRaws  []string
	mMailIdRaws     map[string]string //短信模板号匹配Mail
	mSmsIdRaws      map[string]string
	mConfig         *tools.Config
}

func (this *NotifyMgr) Init(config *tools.Config) {
	this.mConfig = config
	this.mUserNotifyRaws = make([]string, 0)
	this.mSysNotifyRaws = make([]string, 0)
	this.mMailIdRaws = make(map[string]string, 0)
	this.mSmsIdRaws = make(map[string]string, 0)
	for i := 0; i < len(config.Notify.UserNotify); i++ {
		this.mUserNotifyRaws = append(this.mUserNotifyRaws, config.Notify.UserNotify[i])
	}
	for i := 0; i < len(config.Notify.SysNotify); i++ {
		this.mSysNotifyRaws = append(this.mSysNotifyRaws, config.Notify.SysNotify[i])
	}
	for i := 0; i+1 < len(config.Notify.MailId); i = i + 2 {
		this.mMailIdRaws[config.Notify.MailId[i]] = config.Notify.MailId[i+1]
	}
	for i := 0; i+1 < len(config.Notify.SmsId); i = i + 2 {
		this.mSmsIdRaws[config.Notify.SmsId[i]] = config.Notify.SmsId[i+1]
	}
	l4g.Info("mUserNotifyRaws[%v]", this.mUserNotifyRaws)
	l4g.Info("mSysNotifyRaws[%v]", this.mSysNotifyRaws)
	l4g.Info("mMailIdRaws[%v]", this.mMailIdRaws)
	l4g.Info("mSmsIdRaws[%v]", this.mSmsIdRaws)
}

func (this *NotifyMgr) Start() error {
	return nil
}

func (this *NotifyMgr) Close() error {
	return nil
}

func (this *NotifyMgr) Match(input map[string]interface{}) {
	defer this.ruleRecover()
	f := func() {
		l4g.Info("start to Rule Match")
		userFlag := this.MatchUserNotify(input)
		//		sysFlag := this.MatchSysNotify(input)
		l4g.Info("Rule result[userFlag=%v ", userFlag)
		if !userFlag {
			l4g.Info("no Match rule")
			return
		}
		if userFlag {
			input["mailstatus"] = true
			input["smsstatus"] = true
			mailid := this.MatchTextId(input, this.mMailIdRaws)
			smsid := this.MatchTextId(input, this.mSmsIdRaws)
			l4g.Info("mailid[%s] smsid[%s]", mailid, smsid)
			if (len(mailid) == 0) && (len(smsid) == 0) {
				l4g.Info("empty id so NoSend")
				return
			}
			bothFlag := (len(mailid) != 0) && (len(smsid) != 0)

			email, mailStatus := input[ConstUserEmailElemName]
			if mailStatus {
				mailStatus = this.toMail(email.(string), mailid)
			}
			phone, smsStatus := input[ConstUserPhoneElemName]
			if smsStatus {
				smsStatus = this.toSMS(phone.(string), smsid)
			}

			if bothFlag {
				return
			}
			if mailStatus && smsStatus {
				return
			}
			if !mailStatus {
				l4g.Info("mail fail and resend with sms")
				input["mailstatus"] = mailStatus
				smsid = this.MatchTextId(input, this.mSmsIdRaws)
				this.toSMS(phone.(string), smsid)
			}

			if !smsStatus {
				l4g.Info("sms fail and resend with mail")
				input["smsstatus"] = smsStatus
				mailid = this.MatchTextId(input, this.mMailIdRaws)
				this.toMail(email.(string), mailid)
			}

		}
		//if sysFlag {
		//	email, ok := input[ConstSysEmailElemName]
		//	if ok {
		//		this.toMail(email.(string), id)
		//	}
		//	phone, ok := input[ConstSysPhoneElemName]
		//	if ok {
		//		this.toSMS(phone.(string), id)
		//	}
		//}
		l4g.Info("end to Rule Match")
	}
	go f()
}

func (this *NotifyMgr) MatchMail(input map[string]interface{}) {
	defer this.ruleRecover()
	go func() {
		mailid := this.MatchTextId(input, this.mMailIdRaws)
		email, mailStatus := input[ConstUserEmailElemName]
		if mailStatus {
			mailStatus = this.toMail(email.(string), mailid)
		}
	}()
}

func (this *NotifyMgr) MatchSms(input map[string]interface{}) {
	defer this.ruleRecover()
	go func() {
		smsid := this.MatchTextId(input, this.mSmsIdRaws)
		phone, smsStatus := input[ConstUserPhoneElemName]
		if smsStatus {
			smsStatus = this.toSMS(phone.(string), smsid)
		}
	}()
}

func (this *NotifyMgr) DirectSend(input map[string]interface{}) {

}

func (this *NotifyMgr) DirectMail(fileDir, fileNameSuffix, toEmail, srcEmail string) error {
	path := fileDir + "/email_" + fileNameSuffix + ".html"
	subject, body, err := common.ParseHtmlTemplate(path, nil)
	if err != nil {
		l4g.Error("DirectMail ParseTextTemplate[%s] err[%s]", path, err.Error())
		return err
	}
	awsConfig := this.mConfig.Aws
	sdk := ses.NewSesSdk(awsConfig.SesRegion, awsConfig.AccessKeyId, awsConfig.AccessKey, awsConfig.AccessToken)
	err = sdk.Send(toEmail, subject, srcEmail, "UTF-8", body)
	if err != nil {
		l4g.Error("SendMail err:%v", err.Error())
		return err
	}
	l4g.Info("DirectMail[%s][%s] ok", subject, toEmail)
	return nil
}

func (this *NotifyMgr) DirectSms(fileDir, fileNameSuffix string, phone string) error {
	path := fileDir + "/sms_" + fileNameSuffix + ".txt"
	body, err := common.ParseTextTemplate(path, nil)
	if err != nil {
		l4g.Error("DirectSms ParseTextTemplate[%s] err[%s]", path, err.Error())
		return err
	}
	awsConfig := this.mConfig.Aws
	sdk := sns.NewSnsSdk(awsConfig.SnsRegion, awsConfig.AccessKeyId, awsConfig.AccessKey, awsConfig.AccessToken)
	if err := sdk.Send(body, phone); err != nil {
		l4g.Error("DirectSms Send[%s][%s] err[%s]", body, "", err.Error())
		return err
	}
	l4g.Info("DirectSms[%s][%s] ok", body, phone)
	return nil
}

func (this *NotifyMgr) MatchUserNotify(input map[string]interface{}) bool {
	l4g.Debug("rule input[%v]", input)
	for i := 0; i < len(this.mUserNotifyRaws); i++ {
		l4g.Debug("userRule[%v]", this.mUserNotifyRaws[i])
		if len(this.mUserNotifyRaws[i]) == 0 {
			continue
		}
		ruler, err := yql.Rule(this.mUserNotifyRaws[i])
		if err != nil {
			l4g.Error("yql Rule[%s] err[%s]", this.mUserNotifyRaws[i], err.Error())
			return false
		}
		result, err := ruler.Match(input)
		if err != nil {
			l4g.Error("yql Match[%s]input[%v] err[%s]", this.mUserNotifyRaws[i], input, err.Error())
			return false
		}
		if result {
			return true
		}
	}
	return false
}

func (this *NotifyMgr) MatchSysNotify(input map[string]interface{}) bool {
	for i := 0; i < len(this.mSysNotifyRaws); i++ {
		ruler, err := yql.Rule(this.mSysNotifyRaws[i])
		if err != nil {
			l4g.Error("yql Rule[%s] err[%s]", this.mSysNotifyRaws[i], err.Error())
			return false
		}
		result, err := ruler.Match(input)
		if err != nil {
			l4g.Error("yql Match[%s]input[%v] err[%s]", this.mSysNotifyRaws[i], input, err.Error())
			return false
		}
		if result {
			return true
		}
	}
	return false
}

func (this *NotifyMgr) MatchTextId(input map[string]interface{}, raws map[string]string) string {
	for k, v := range raws {
		ruler, err := yql.Rule(v)
		if err != nil {
			l4g.Error("yql Rule[%s:%s] err[%v]", k, v, err)
			return ""
		}
		result, err := ruler.Match(input)
		if err != nil {
			l4g.Error("yql Match[%s:%s] input[%v] err[%v]", k, v, input, err)
			return ""
		}
		if result {
			return k
		}
	}
	return ""
}

//id是邮件模板编号
func (this *NotifyMgr) toMail(toEmail, id string) bool {
	if len(id) == 0 {
		return true
	}
	awsConfig := this.mConfig.Aws
	sdk := ses.NewSesSdk(awsConfig.SesRegion, awsConfig.AccessKeyId, awsConfig.AccessKey, awsConfig.AccessToken)
	err := sdk.SendTemplate(toEmail, "", this.mConfig.Notify.SrcEmail, id, "{ \"name\":\"Users\"}", awsConfig.SesTimeout)
	if err != nil {
		l4g.Error("SendMail err:%v", err.Error())
		return false
	}
	l4g.Info("SendMail[%s][%s] ok", toEmail, id)
	return true
}

//https://docs.aws.amazon.com/zh_cn/sns/latest/dg/sms_preferences.html
func (this *NotifyMgr) toSMS(phone, id string) bool {
	if len(id) == 0 {
		return true
	}
	path := this.mConfig.Notify.TmplateDir + "/" + id + ".txt"
	body, err := common.ParseTextTemplate(path, nil)
	if err != nil {
		l4g.Error("toSMS ParseTextTemplate[%s] err[%s]", path, err.Error())
		return false
	}
	awsConfig := this.mConfig.Aws
	sdk := sns.NewSnsSdk(awsConfig.SesRegion, awsConfig.AccessKeyId, awsConfig.AccessKey, awsConfig.AccessToken)
	if err := sdk.Send(body, phone); err != nil {
		l4g.Error("toSMS Send[%s][%s] err[%s]", body, "", err.Error())
		return false
	}
	l4g.Info("SendSms[%s][%s] ok", phone, id)
	return true
}

func (this *NotifyMgr) ruleRecover() {
	if err := recover(); err != nil {
		l4g.Error("panic err[%v]", err)
	}
}
