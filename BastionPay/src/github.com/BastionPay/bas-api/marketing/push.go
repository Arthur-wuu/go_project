package marketing

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PushMsg struct{
	MsgType  int      `json:"msg_type"`
	Data    interface{} `json:"data,omitempty"`
}

//1--9 裂变
const CONST_PUSHTYPE_FISSION_ACTIVITY = 1
const CONST_PUSHTYPE_FISSION_Red = 2
const CONST_PUSHTYPE_FISSION_Robber = 3

//10-19 是抽奖
const CONST_PUSHTYPE_LUCKDRAW_DRAWER = 10