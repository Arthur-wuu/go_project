package api

//为了扩展性好，采用这种方式，以后再改成 proto吧；以后所有协议 参考它

const (
	ErrCode_Success    = 0
	ErrCode_Param      = 10001
	ErrCode_InerServer = 10002
	ErrCode_UrlPath = 10003
)






type Market struct {
	Name *string   `json:"Name,omitempty"`
	Objs []string  `json:"Obj,omitempty"`
}

func (this *Market) SetName(n string) {
	if this.Name == nil {
		this.Name = new(string)
	}
	*this.Name = n
}

/*************************************************/
