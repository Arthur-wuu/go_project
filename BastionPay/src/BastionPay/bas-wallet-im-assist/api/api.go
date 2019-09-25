package api

type (
	//用户状态 变更回调
	UserStatus struct {
		CallbackCommand string         `json:"CallbackCommand,omitempty"`
		UserStatusInfo  UserStatusInfo `json:"Info,omitempty"`
	}

	UserStatusInfo struct {
		Action    string `json:"Action,omitempty"`
		ToAccount string `json:"To_Account,omitempty"`
		Reason    string `json:"Reason,omitempty"`
	}

	//单聊 回调
	SingleChat struct {
		CallbackCommand string    `json:"CallbackCommand,omitempty"`
		FromAccount     string    `json:"From_Account,omitempty"`
		ToAccount       string    `json:"To_Account,omitempty"`
		MsgBody         []MsgBody `json:"MsgBody,omitempty"`
	}

	MsgBody struct {
		MsgType    string     `json:"MsgType,omitempty"`
		MsgContent MsgContent `json:"MsgContent,omitempty"`
	}

	MsgContent struct {
		Text string `json:"Text,omitempty"`
	}
)

const (
	CONST_User_Status_On  = 0
	CONST_User_Status_Off = 1
)
