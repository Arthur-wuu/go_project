package api

import(
	"encoding/json"
)

func NewErrResponse(qid string, err int32) *Response {
	res := new(Response)
	res.Qid = new(string)
	*res.Qid = qid
	res.SetErr(err)
	return res
}

func NewResponse(qid string) *Response {
	res := new(Response)
	res.Qid = new(string)
	*res.Qid = qid
	return res
}

type Response struct {
	Qid              *string `json:"Qid,omitempty"`     //请求id，流水号, 单次请求唯一
	Err              *int32  `json:"Err,omitempty"`     //错误号
	Counter          *uint64 `json:"Counter,omitempty"` //推送序号, 0表示正常请求，>=1 表示推送
	Data             interface{}    `json:"Data,omitempty"`    //数据
}

func (this *Response) SetCounter(e uint64) {
	if this.Counter == nil {
		this.Counter = new(uint64)
	}
	*this.Counter = e
}

func (this *Response) SetErr(e int32) {
	if this.Err == nil {
		this.Err = new(int32)
	}
	*this.Err = e
}

func (this *Response) Marshal(msg interface{})([]byte, error){
	if msg != nil {
		this.Data = msg
	}
	return json.Marshal(this)
}
