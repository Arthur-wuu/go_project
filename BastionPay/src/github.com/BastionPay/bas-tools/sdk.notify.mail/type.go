package sdk_notify_mail

const (
	Notify_Level_Normal = 0
	Notify_Level_Low    = 0
	Notify_Level_High   = 1

	ErrCode_Success             = 0
	ErrCode_Param               = 10001
	ErrCode_InerServer          = 10002
	ErrCode_Exist               = 10003
	ErrCode_NoAlive             = 10004
	ErrCode_NoFindAliveTemplate = 10005
	ErrCode_SubRecipient        = 10006
)

type ReqNotifyMsg struct {
	GroupName           *string                `json:"group_name,omitempty"` //groupname+lang 联合使用
	GroupId             *uint                  `json:"group_id,omitempty"`   //groupid+lang联合使用
	Lang                *string                `json:"lang,omitempty"`
	GroupAlias          *string                `json:"group_alias,omitempty"` //GroupAlias+lang
	TempAlias           *string                `json:"temp_alias,omitempty"`  //可单独使用，重复则选其一
	TempId              *uint                  `json:"temp_id,omitempty"`     //唯一，可单独使用
	Params              map[string]interface{} `json:"params,omitempty"`      //optional
	Recipient           []string               `json:"recipient,omitempty"`   //require
	Level               *uint                  `json:"level,omitempty"`       //optional
	AppName             *string                `json:"app_name,omitempty"`
	UseDefaultRecipient *bool                  `json:"use_default_recipient,omitempty"` //使用默认收件人，默认是使用的
}

func (this *ReqNotifyMsg) SetUseDefaultRecipient(n bool) {
	if this.UseDefaultRecipient == nil {
		this.UseDefaultRecipient = new(bool)
	}
	*this.UseDefaultRecipient = n
}

func (this *ReqNotifyMsg) SetGroupName(n string) {
	if this.GroupName == nil {
		this.GroupName = new(string)
	}
	*this.GroupName = n
}

func (this *ReqNotifyMsg) SetTempId(id uint) {
	if this.TempId == nil {
		this.TempId = new(uint)
	}
	*this.TempId = id
}

func (this *ReqNotifyMsg) SetLevel(l uint) {
	if this.Level == nil {
		this.Level = new(uint)
	}
	*this.Level = l
}

func (this *ReqNotifyMsg) SetAppName(a string) {
	if this.AppName == nil {
		this.AppName = new(string)
	}
	*this.AppName = a
}

func (this *ReqNotifyMsg) AddRecipient(a string) {
	if this.Recipient == nil {
		this.Recipient = make([]string, 0)
	}
	this.Recipient = append(this.Recipient, a)
}

func (this *ReqNotifyMsg) SetTempAlias(a string) {
	if this.TempAlias == nil {
		this.TempAlias = new(string)
	}
	*this.TempAlias = a
}

func (this *ReqNotifyMsg) SetLang(l string) {
	if this.Lang == nil {
		this.Lang = new(string)
	}
	*this.Lang = l
}

func (this *ReqNotifyMsg) SetGroupId(id uint) {
	if this.GroupId == nil {
		this.GroupId = new(uint)
	}
	*this.GroupId = id
}

type ResNotifyMsg struct {
	Err    *int    `json:"err,omitempty"`
	ErrMsg *string `json:"errmsg,omitempty"`
}

func (this *ResNotifyMsg) IsOk() bool {
	if this.Err == nil {
		return true
	}
	if *this.Err != 0 {
		return false
	}
	return true
}

func (this *ResNotifyMsg) HaveSubErr() bool {
	if this.Err == nil {
		return false
	}
	if *this.Err == ErrCode_SubRecipient {
		return true
	}
	return false
}

func (this *ResNotifyMsg) GetErr() int {
	if this.Err == nil {
		return 0
	}
	return *this.Err
}

func (this *ResNotifyMsg) GetErrMsg() string {
	if this.ErrMsg == nil {
		return ""
	}
	return *this.ErrMsg
}
