package models

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/config"
	"BastionPay/bas-notify/models/table"
	"BastionPay/bas-notify/sms"
	"github.com/juju/errors"
	"go.uber.org/zap"
	"strings"
)

type (
	//优先级TempId(唯一)>>TempAlias(同一个groupid，同一语言 可以少量重复)
	//     >>(GroupId+lang)(lang可以少量重复)>>(GroupName+lang)(GroupName唯一，lang可少量重复)
	//     >>(GroupAlias+lang)(GroupAlias唯一，lang可少量重复)
	SmsMsg struct {
		GroupName    *string `valid:"optional" json:"group_name,omitempty"`     //groupname+lang 联合使用
		GroupSubName *string `valid:"optional" json:"group_sub_name,omitempty"` //groupname+lang 联合使用
		GroupId      *int    `valid:"optional" json:"group_id,omitempty"`       //groupid+lang联合使用
		Lang         *string `valid:"required" json:"lang,omitempty"`
		//GroupAlias *string               `valid:"optional" json:"group_alias,omitempty"`//GroupAlias+lang
		//TempAlias *string                `valid:"optional" json:"temp_alias,omitempty"`     //可单独使用，重复则选其一
		//TempId    *int                  `valid:"optional" json:"temp_id,omitempty"`    //唯一，可单独使用
		Params    map[string]interface{} `valid:"-" json:"params,omitempty"`           //optional
		Recipient []string               `valid:"optional" json:"recipient,omitempty"` //require
		//Level     *int                  `valid:"optional" json:"level,omitempty"`     //optional
		AppName             *string `valid:"optional" json:"app_name,omitempty"`
		UseDefaultRecipient *bool   `valid:"optional" json:"use_default_recipient,omitempty"` //使用默认收件人
		SenderId            *string `valid:"optional" json:"sender_name,omitempty"`
	}
)

func (this *SmsMsg) Send(recordFlag bool) (int, error) {
	tmplate, errCode, err := this.GetValidTemplate()
	if err != nil {
		return errCode, err
	}

	smsBody := ""
	if tmplate.Content != nil {
		smsBody, err = ParseTextTemplate(*tmplate.Content, this.Params)
		if err != nil {
			ZapLog().With(zap.Any("tempParam", this.Params), zap.Any("tempid", tmplate.Id), zap.Error(err)).Error("ParseTextTemplate err")
			go RecordHistory(*tmplate.GroupId, Notify_Type_Sms, 0, len(this.Recipient), recordFlag)
			return apibackend.BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL.Code(), errors.Annotate(err, "ParseTextTemplate")
		}
	}

	failCount := 0
	var gerr error

	switch *tmplate.SmsPlatform {
	case SMSPlatform_TWL:
		failCount, gerr = this.SendTwl(smsBody)
		errCode = apibackend.BASERR_BASNOTIFY_TWL_ERR.Code()
	case SMSPlatform_CHUANGLAN:
		smsBodyChina := ""
		if *this.Lang != "en-US" && *this.Lang != "zh-CN" {
			tmplate, errCode, err := this.GetValidTemplateByEn()
			if err != nil {
				return errCode, err
			}

			if tmplate.Content != nil {
				smsBodyChina, err = ParseTextTemplate(*tmplate.Content, this.Params)
				if err != nil {
					ZapLog().With(zap.Any("tempParam", this.Params), zap.Any("tempid", tmplate.Id), zap.Error(err)).Error("ParseTextTemplate err")
					go RecordHistory(*tmplate.GroupId, Notify_Type_Sms, 0, len(this.Recipient), recordFlag)
					return apibackend.BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL.Code(), errors.Annotate(err, "ParseTextTemplate")
				}
			}
		} else {
			smsBodyChina = smsBody
		}

		failCount, gerr = this.SendLanChuang(smsBody, smsBodyChina)
		errCode = apibackend.BASERR_BASNOTIFY_LANCHUANG_ERR.Code()
	case SMSPlatform_AWS:
		failCount, gerr = this.SendAws(smsBody)
		errCode = apibackend.BASERR_BASNOTIFY_AWS_ERR.Code()
	case SMSPlatform_Nexmo:
		failCount, gerr = this.SendNexmo(smsBody)
		errCode = apibackend.BASERR_BASNOTIFY_Nexmo_ERR.Code()
	case SMSPlatform_YunTongXun:
		failCount, gerr = this.SendYunTongXin(smsBody, tmplate)
	}

	if gerr != nil {
		go RecordHistory(*tmplate.GroupId, Notify_Type_Sms, len(this.Recipient)-failCount, failCount, recordFlag)
		return errCode, gerr
	}

	go RecordHistory(*tmplate.GroupId, Notify_Type_Sms, len(this.Recipient), 0, recordFlag)
	ZapLog().With(zap.Any("param", *this)).Info("SmsSend success")
	return 0, nil
}

func (this *SmsMsg) GetValidTemplate() (*table.Template, int, error) {
	var group *table.TemplateGroup
	var err error
	if this.GroupId != nil {
		group, err = new(TemplateGroup).GetByid(*this.GroupId)
		if err != nil {
			return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetById")
		}
	} else if this.GroupName != nil {
		group, err = new(TemplateGroup).GetByNameAndType(*this.GroupName, Notify_Type_Sms, this.GroupSubName)
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
	if tempalate == nil || tempalate.Content == nil || len(*tempalate.Content) == 0 {
		tempalate, _ = new(Template).GetByGIdAndLang(*group.Id, "en-US")
		if tempalate == nil {
			return nil, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), errors.Errorf("nofind temp")
		}
	}
	if tempalate.Content == nil || len(*tempalate.Content) == 0 {
		return nil, apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code(), errors.Errorf("temp content is nil")
	}
	if group.Alive == nil || *group.Alive == Notify_AliveMode_Dead {
		return nil, apibackend.BASERR_BASNOTIFY_TEMPLATE_DEAD.Code(), errors.Errorf("dead temp")
	}
	if group.Type == nil || *group.Type != Notify_Type_Sms {
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

func (this *SmsMsg) GetValidTemplateByEn() (*table.Template, int, error) {
	var group *table.TemplateGroup
	var err error
	if this.GroupId != nil {
		group, err = new(TemplateGroup).GetByid(*this.GroupId)
		if err != nil {
			return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetById")
		}
	} else if this.GroupName != nil {
		group, err = new(TemplateGroup).GetByNameAndType(*this.GroupName, Notify_Type_Sms, this.GroupSubName)
		if err != nil {
			return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetByNameAndType")
		}
	} else {
		return nil, apibackend.BASERR_INVALID_PARAMETER.Code(), errors.Annotate(err, "GetUnSupport")
	}

	if group == nil {
		return nil, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), errors.Errorf("nofind tempGroup")
	}

	tempalate, err := new(Template).GetByGIdAndLang(*group.Id, "en-US")
	if err != nil {
		return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetByGIdAndLang")
	}
	if tempalate == nil || tempalate.Content == nil || len(*tempalate.Content) == 0 {
		tempalate, _ = new(Template).GetByGIdAndLang(*group.Id, "en-US")
		if tempalate == nil {
			return nil, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), errors.Errorf("nofind temp")
		}
	}
	if tempalate.Content == nil || len(*tempalate.Content) == 0 {
		return nil, apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code(), errors.Errorf("temp content is nil")
	}
	if group.Alive == nil || *group.Alive == Notify_AliveMode_Dead {
		return nil, apibackend.BASERR_BASNOTIFY_TEMPLATE_DEAD.Code(), errors.Errorf("dead temp")
	}
	if group.Type == nil || *group.Type != Notify_Type_Sms {
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

func (this *SmsMsg) SendTwl(smsBody string) (int, error) {
	failCount := 0
	var gerr error
	for i := 0; i < len(this.Recipient); i++ {
		if len(this.Recipient[i]) == 0 {
			continue
		}
		if err := sms.GSmsMgr.DirectSendTwl(smsBody, this.Recipient[i]); err != nil {
			failCount++
			gerr = err
			ZapLog().With(zap.Error(err), zap.String("recipient", this.Recipient[i])).Error("DirectSendTWL err")
			continue
		}
	}
	return failCount, gerr
}

func (this *SmsMsg) SendLanChuang(smsBody, smsBodyChina string) (int, error) {
	failCount := 0
	var gerr error
	zhPhones, noZhPhones := ChuangLanSplitPhones(this.Recipient)
	//newClSmsBody,newClParam := ParseSmsBodyChuanglan(tmplate.GetContent(), req.Params)
	if num, err := sms.GSmsMgr.DirectSendChuanglan(smsBodyChina, zhPhones, nil); err != nil {
		failCount += len(zhPhones)
		gerr = err
		ZapLog().With(zap.Error(err), zap.Any("recipient", zhPhones)).Error("DirectSendChuanglan err")
		//这里不能退出 还得继续发送
	} else {
		failCount += len(zhPhones) - num
	}

	for i := 0; i < len(noZhPhones); i++ {
		if len(noZhPhones[i]) == 0 {
			continue
		}
		err := sms.GSmsMgr.DirectSendAws(smsBody, noZhPhones[i], this.SenderId)
		if err != nil {
			failCount++
			gerr = err
			ZapLog().With(zap.Error(err), zap.String("recipient", noZhPhones[i])).Error("DirectSendAws err")
		}
	}
	return failCount, gerr
}

func (this *SmsMsg) SendYunTongXin(smsBody string, tmplate *table.Template) (int, error) {
	failCount := 0
	var gerr error
	zhPhones, noZhPhones := SplitChinaAndUnChinaPhones(this.Recipient)
	if tmplate.YuntongxinTempId != nil {
		_, params := SortMap(this.Params)

		if err := sms.GSmsMgr.DirectSendYunTongXin(*tmplate.YuntongxinTempId, config.GConfig.YunTongXun.AppId, strings.Join(zhPhones, ","), params); err != nil {
			failCount += len(zhPhones)
			gerr = err
			ZapLog().With(zap.Error(err), zap.Any("recipient", zhPhones)).Error("DirectSendChuanglan err")
			//这里不能退出 还得继续发送
		}
	} else {
		for i := 0; i < len(zhPhones); i++ {
			if len(zhPhones[i]) == 0 {
				continue
			}
			err := sms.GSmsMgr.DirectSendAws(smsBody, zhPhones[i], this.SenderId)
			if err != nil {
				failCount++
				gerr = err
				ZapLog().With(zap.Error(err), zap.String("recipient", zhPhones[i])).Error("DirectSendAws err")
			}
		}
	}

	for i := 0; i < len(noZhPhones); i++ {
		if len(noZhPhones[i]) == 0 {
			continue
		}
		err := sms.GSmsMgr.DirectSendAws(smsBody, noZhPhones[i], this.SenderId)
		if err != nil {
			failCount++
			gerr = err
			ZapLog().With(zap.Error(err), zap.String("recipient", noZhPhones[i])).Error("DirectSendAws err")
		}
	}
	return failCount, gerr
}

func (this *SmsMsg) SendAws(smsBody string) (int, error) {
	failCount := 0
	var gerr error
	for i := 0; i < len(this.Recipient); i++ {
		if len(this.Recipient[i]) == 0 {
			continue
		}
		err := sms.GSmsMgr.DirectSendAws(smsBody, this.Recipient[i], this.SenderId)
		if err != nil {
			failCount++
			gerr = err
			ZapLog().With(zap.Error(err), zap.String("recipient", this.Recipient[i])).Error("DirectSendAWS err")
		}
	}
	return failCount, gerr
}

func (this *SmsMsg) SendNexmo(smsBody string) (int, error) {
	failCount := 0
	var gerr error
	for i := 0; i < len(this.Recipient); i++ {
		if len(this.Recipient[i]) == 0 {
			continue
		}
		err := sms.GSmsMgr.DirectSendNexmo(smsBody, this.Recipient[i], this.SenderId)
		if err != nil {
			failCount++
			gerr = err
			ZapLog().With(zap.Error(err), zap.String("recipient", this.Recipient[i])).Error("DirectSendNexmo err")
		}
	}
	return failCount, gerr
}

func (this *SmsMsg) AppendRecipient(recipients string) {
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
