package models

import (
	"BastionPay/bas-notify/base"
	"BastionPay/bas-notify/config"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	//"BastionPay/bas-notify/sms"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/models/table"
	"github.com/juju/errors"
)

const(
	robotUrl  = "https://oapi.dingtalk.com/robot/send?access_token="

)
type Msg struct {
	Msgtype  *string `json:"msgtype,omitempty"`
	Text      Content `json:"text,omitempty"`
}

type Content struct {
	Content  *string `json:"content,omitempty"`
}

type MsgRes struct {
	Errcode  *int `json:"errcode,omitempty"`
	Errmsg   *string `json:"errmsg,omitempty"`
}



type(
	//优先级TempId(唯一)>>TempAlias(同一个groupid，同一语言 可以少量重复)
	//     >>(GroupId+lang)(lang可以少量重复)>>(GroupName+lang)(GroupName唯一，lang可少量重复)
	//     >>(GroupAlias+lang)(GroupAlias唯一，lang可少量重复)
	DingDingMsg struct {
		GroupName *string                 `valid:"optional" json:"group_name,omitempty"` //groupname+lang 联合使用
		GroupId   *int                  `valid:"optional" json:"group_id,omitempty"` //groupid+lang联合使用
		Lang      *string                 `valid:"required" json:"lang,omitempty"`
		//GroupAlias *string               `valid:"optional" json:"group_alias,omitempty"`//GroupAlias+lang
		//TempAlias *string                `valid:"optional" json:"temp_alias,omitempty"`     //可单独使用，重复则选其一
		//TempId    *int                  `valid:"optional" json:"temp_id,omitempty"`    //唯一，可单独使用
		Params    map[string]interface{}  `valid:"-" json:"params,omitempty"`    //optional
		//Recipient []string               `valid:"optional" json:"recipient,omitempty"` //require
		Level     *int                     `valid:"optional" json:"level,omitempty"`     //optional
		AppName   *string                  `valid:"optional" json:"app_name,omitempty"`
		//UseDefaultRecipient *bool        `valid:"optional" json:"use_default_recipient,omitempty"` //使用默认收件人
		//SenderId   *string                `valid:"optional" json:"sender_name,omitempty"`
	}
)

func (this *DingDingMsg) Send(recordFlag bool) (int, error){
	tmplate, errCode, err := this.GetValidTemplate()
	if err != nil {
		return errCode,err
	}

	dingDingMsgBody := ""
	if tmplate.Content != nil {
		dingDingMsgBody, err = ParseTextTemplate(*tmplate.Content, this.Params)
		if err != nil {
			ZapLog().With(zap.Any("tempParam", this.Params),zap.Any("tempid", tmplate.Id),zap.Error(err)).Error("ParseTextTemplate err")
			//go RecordHistory(*tmplate.GroupId, Notify_Type_DDing, 0,len(this.Recipient), recordFlag)
			return apibackend.BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL.Code(), errors.Annotate(err, "ParseTextTemplate")
		}
	}

	code ,gerr := this.SendDingDing(dingDingMsgBody, tmplate.DingQunName)
	errCode = apibackend.BASERR_BASNOTIFY_TWL_ERR.Code()

	if gerr != nil {
	//	go RecordHistory(*tmplate.GroupId, Notify_Type_Sms, len(this.Recipient)-failCount ,failCount, recordFlag)
		return errCode, gerr
	}

	//go RecordHistory(*tmplate.GroupId,Notify_Type_Sms, len(this.Recipient),0, recordFlag)
	ZapLog().With(zap.Any("param", *this)).Info("SmsSend success")
	return code,nil
}

func (this *DingDingMsg) SendDingDing(dingBody string, dingQunName *string) (int, error) {

	msg := new(Msg)
	text := "text"
	msg.Msgtype = &text
	msg.Text.Content = &dingBody

	reqBody, err := json.Marshal(msg)
	if err != nil {
		ZapLog().With(zap.Any("marshal", msg),zap.Any("msg", msg),zap.Error(err)).Error("marshal error")
		return apibackend.BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL.Code(), errors.Annotate(err, "marshal error")
	}

	var res []byte
	switch *dingQunName {
	case config.GConfig.DingDing[0].QunName:
		res, err = base.HttpSend(robotUrl+config.GConfig.DingDing[0].RobToken, bytes.NewBuffer(reqBody), "POST",nil)
	case config.GConfig.DingDing[1].QunName:
		res, err = base.HttpSend(robotUrl+config.GConfig.DingDing[1].RobToken, bytes.NewBuffer(reqBody), "POST",nil)
	case config.GConfig.DingDing[2].QunName:
		res, err = base.HttpSend(robotUrl+config.GConfig.DingDing[2].RobToken, bytes.NewBuffer(reqBody), "POST",nil)
	}

	robotRes := new(MsgRes)
	err = json.Unmarshal(res, robotRes)
	if err != nil {
		ZapLog().With(zap.Any("unmarshal", msg),zap.Any("robotRes", robotRes),zap.Error(err)).Error("unmarshal error")
		return apibackend.BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL.Code(), errors.Annotate(err, "unmarshal error")
	}

	return *robotRes.Errcode, err
}

func (this *DingDingMsg) GetValidTemplate() (*table.Template, int, error) {
	var group *table.TemplateGroup
	var err error
	if this.GroupId != nil {
		group,err = new(TemplateGroup).GetByid(*this.GroupId)
		if err != nil {
			return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetById")
		}
	}else if this.GroupName != nil {                                   //需要改动 2
		group,err = new(TemplateGroup).GetByNameAndType(*this.GroupName, Notify_Type_Ding, nil)
		if err != nil {
			return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetByNameAndType")
		}
	}else{
		return nil, apibackend.BASERR_INVALID_PARAMETER.Code(), errors.Annotate(err, "GetUnSupport")
	}

	if group == nil {
		return nil, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), errors.Errorf("nofind tempGroup")
	}
	if this.Lang == nil {
		this.Lang = new(string)
		*this.Lang = "zh-CN"
	}

	tempalate,err := new(Template).GetByGIdAndLang(*group.Id, *this.Lang)
	if err != nil {
		return nil, apibackend.BASERR_DATABASE_ERROR.Code(), errors.Annotate(err, "GetByGIdAndLang")
	}
	if (*this.Lang != "zh-CN") &&(tempalate == nil || tempalate.Content == nil || len(*tempalate.Content) == 0){
		tempalate,_ = new(Template).GetByGIdAndLang(*group.Id, "zh-CN")
		if tempalate == nil {
			return nil, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), errors.Errorf("nofind temp")
		}
	}
	if tempalate.Content == nil || len(*tempalate.Content) == 0 {
		return nil, apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code(), errors.Errorf("temp content is nil")
	}
	//if tempalate.Alive == nil || *tempalate.Alive == Notify_AliveMode_Dead {
	//	return nil, apibackend.BASERR_BASNOTIFY_TEMPLATE_DEAD.Code(), errors.Errorf("dead temp")
	//}
	//if tempalate.Type == nil || *tempalate.Type != Notify_Type_Sms {
	//	return nil, apibackend.BASERR_UNKNOWN_BUG.Code(), errors.Errorf("type not sms")
	//}
	//if this.UseDefaultRecipient == nil { //兼容
	//	this.UseDefaultRecipient = new(bool)
	//	*this.UseDefaultRecipient = true
	//}
	//if (this.UseDefaultRecipient!=nil) && *this.UseDefaultRecipient  && (tempalate.DefaultRecipient != nil) {
	//	this.AppendRecipient( *tempalate.DefaultRecipient)
	//}
	//if this.Recipient == nil || len(this.Recipient) == 0 {
	//	return nil, apibackend.BASERR_BASNOTIFY_RECIPIENT_EMPTY.Code(),  errors.Errorf("no Recipient")
	//}
	//if tempalate.SmsPlatform == nil {
	//	tempalate.SmsPlatform = new(int)
	//	*tempalate.SmsPlatform = SMSPlatform_AWS
	//}
	return tempalate, 0,nil
}

//
//func (this *DingDingMsg)AppendRecipient(recipients string) {
//	if len(recipients) < 4 {//手机和邮箱 何止4字符
//		return
//	}
//	recipients = strings.Replace(recipients," ","", len(recipients))
//	recipientsArr := strings.Split(recipients, ",")
//	if this.Recipient == nil {
//		this.Recipient = make([]string, 0)
//	}
//	for i:=0; i< len(recipientsArr); i++ {
//		if len(recipientsArr[i]) == 0 {
//			continue
//		}
//		this.Recipient = append(this.Recipient, recipientsArr[i])
//	}
//}