package models

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/email"
	"BastionPay/bas-notify/models/table"
	"github.com/juju/errors"
	"go.uber.org/zap"
	"strings"
)

type EmailMsg struct {
	GroupName *string `valid:"optional" json:"group_name,omitempty"` //groupname+lang 联合使用
	GroupId   *int    `valid:"optional" json:"group_id,omitempty"`   //groupid+lang联合使用
	Lang      *string `valid:"required" json:"lang,omitempty"`
	//GroupAlias *string               `valid:"optional" json:"group_alias,omitempty"`//GroupAlias+lang
	//TempAlias *string                `valid:"optional" json:"temp_alias,omitempty"`     //可单独使用，重复则选其一
	//TempId    *int                  `valid:"optional" json:"temp_id,omitempty"`    //唯一，可单独使用
	Params    map[string]interface{} `valid:"-" json:"params,omitempty"`           //optional
	Recipient []string               `valid:"optional" json:"recipient,omitempty"` //require
	//Level     *int                  `valid:"optional" json:"level,omitempty"`     //optional
	AppName             *string `valid:"optional" json:"app_name,omitempty"`
	UseDefaultRecipient *bool   `valid:"optional" json:"use_default_recipient,omitempty"`
	SenderId            *string `valid:"optional" json:"sender_name,omitempty"`
	GroupSubName        *string `valid:"optional" json:"group_sub_name,omitempty"`
}

func (this *EmailMsg) Send(recordFlag bool) (int, error) {
	tmplate, errCode, err := this.GetValidTemplate()
	if err != nil {
		return errCode, err
	}

	body := ""
	if tmplate.Content != nil {
		if this.Params == nil {
			this.Params = make(map[string]interface{})
		}
		this.Params["title_key"] = tmplate.Title
		_, body, err = ParseHtmlTemplate(BodyToHtml(*tmplate.Content), this.Params)
		if err != nil {
			ZapLog().With(zap.Any("tempParam", this.Params), zap.Error(err), zap.Int("Id", *tmplate.Id)).Error("ParseHtmlTemplate err")
			go RecordHistory(*tmplate.GroupId, Notify_Type_Mail, 0, len(this.Recipient), recordFlag)
			return apibackend.BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL.Code(), errors.Annotate(err, "ParseTextTemplate")
		}
	}

	failCount := 0
	var gerr error
	for i := 0; i < len(this.Recipient); i++ {
		if len(this.Recipient[i]) == 0 {
			continue
		}
		err = email.GMailMgr.DirectSend(*tmplate.Title, body, this.Recipient[i], this.SenderId)
		if err != nil {
			failCount++
			gerr = err
			ZapLog().With(zap.Error(err), zap.String("recipient", this.Recipient[i])).Error("DirectSend err")
		}
	}

	if gerr != nil {
		go RecordHistory(*tmplate.GroupId, Notify_Type_Mail, len(this.Recipient)-failCount, failCount, recordFlag)
		return apibackend.BASERR_BASNOTIFY_AWS_ERR.Code(), gerr
	}

	go RecordHistory(*tmplate.GroupId, Notify_Type_Sms, len(this.Recipient), 0, recordFlag)
	ZapLog().With(zap.Any("param", *this)).Info("MailSend success")
	return 0, nil
}

func (this *EmailMsg) AppendRecipient(recipients string) {
	if len(recipients) < 4 { //手机和邮箱 何止4字符
		return
	}
	recipients = strings.Replace(recipients, " ", "", len(recipients))
	recipientsArr := strings.Split(recipients, ",")
	if this.Recipient == nil {
		this.Recipient = make([]string, 0)
	}
	for i := 0; i < len(recipientsArr); i++ {
		if len(recipientsArr[i]) == 0 {
			continue
		}
		this.Recipient = append(this.Recipient, recipientsArr[i])
	}
}

func (this *EmailMsg) GetValidTemplate() (*table.Template, int, error) {
	var group *table.TemplateGroup
	var err error
	if this.GroupId != nil {
		group, err = new(TemplateGroup).GetByid(*this.GroupId)
		if err != nil {
			return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetById")
		}
	} else if this.GroupName != nil {
		group, err = new(TemplateGroup).GetByNameAndType(*this.GroupName, Notify_Type_Mail, this.GroupSubName)
		if err != nil {
			return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetByNameAndType")
		}
	} else {
		return nil, apibackend.BASERR_INVALID_PARAMETER.Code(), errors.Annotate(err, "GetUnSupport")
	}

	if group == nil {
		return nil, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), errors.Errorf("nofind tempGroup")
	}

	tempalate, err := new(Template).GetByGIdAndLang(*group.Id, *this.Lang)
	if err != nil {
		return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetByGIdAndLang")
	}
	if tempalate == nil {
		return nil, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), errors.Errorf("nofind temp")
	}
	if group.Alive == nil || *group.Alive == Notify_AliveMode_Dead {
		return nil, apibackend.BASERR_BASNOTIFY_TEMPLATE_DEAD.Code(), errors.Errorf("dead temp")
	}
	if group.Type == nil || *group.Type != Notify_Type_Mail {
		return nil, apibackend.BASERR_UNKNOWN_BUG.Code(), errors.Errorf("type not sms")
	}
	if this.UseDefaultRecipient == nil { //兼容
		this.UseDefaultRecipient = new(bool)
		*this.UseDefaultRecipient = true
	}
	if (this.UseDefaultRecipient != nil) && *this.UseDefaultRecipient && (group.DefaultRecipient != nil) {
		this.AppendRecipient(*group.DefaultRecipient)
	}
	if this.Recipient == nil || len(this.Recipient) == 0 {
		return nil, apibackend.BASERR_BASNOTIFY_RECIPIENT_EMPTY.Code(), errors.Errorf("no Recipient")
	}
	if group.SmsPlatform == nil {
		group.SmsPlatform = new(int)
		*group.SmsPlatform = SMSPlatform_AWS
	}
	tempalate.SmsPlatform = group.SmsPlatform
	return tempalate, 0, nil
}
