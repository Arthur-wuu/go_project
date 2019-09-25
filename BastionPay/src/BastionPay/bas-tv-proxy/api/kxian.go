package api

/*******************************************************************************/

/*
urlpath协议可带上版本号
协议处理，转发省去版本号。版本号对应返回数据结构，对业务处理无影响。
*/

/* 协议规则 定义如下：
   /版本号/公司(网站名)/业务功能?sub=1&qid=10&output=json,
   /v1/coinmerit/quote/kxian?exa=huobi&obj=eth_usdt&sub=1&qid=&period=
   sub=0 单独请求，sub=1 订阅， sub=2取消订阅。如果有 qid的话，单个命令 cancel?qid=就能搞定

   /v1/coinmerit/exa

   /v1/coinmerit/objlist?exa=zb

btcexa交易所同样遵循上面定义
例：http：  /v1/btcexa/quote/kxian?qid=1&obj=eth_btc&period=1m&limit=3
   ws： /api/v1/btcexa/quote/kxian?qid=1&obj=eth_btc&period=1m&sub=1
*/

func NewQuoteKlineSingle(obj string, kxians []*KXian) *QuoteKlineSingle {
	return &QuoteKlineSingle{
		Obj:  &obj,
		Data: kxians,
	}
}

type QuoteKlineSingle struct {
	Obj              *string  `json:"Obj,omitempty"`
	Data             []*KXian `json:"Data,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (this *QuoteKlineSingle) Append() {

}

func (m *QuoteKlineSingle) Reset() { *m = QuoteKlineSingle{} }

func (m *QuoteKlineSingle) GetObj() string {
	if m != nil && m.Obj != nil {
		return *m.Obj
	}
	return ""
}

func (m *QuoteKlineSingle) GetData() []*KXian {
	if m != nil {
		return m.Data
	}
	return nil
}

type KXian struct {
	ShiJian        *int64  `json:"ShiJian,omitempty"`
	KaiPanJia      *string `json:"KaiPanJia,omitempty"`
	ZuiGaoJia      *string `json:"ZuiGaoJia,omitempty"`
	ZuiDiJia       *string `json:"ZuiDiJia,omitempty"`
	ShouPanJia     *string `json:"ShouPanJia,omitempty"`
	ChengJiaoLiang *string `json:"ChengJiaoLiang,omitempty"`
	ChengJiaoE     *string `json:"ChengJiaoE,omitempty"`
	ChengJiaoBiShu *string `json:"ChengJiaoBiShu,omitempty"`
}

func (this *KXian) SetShiJian(s int64) {
	if this.ShiJian == nil {
		this.ShiJian = new(int64)
	}
	*this.ShiJian = s
}

func (this *KXian) SetKaiPanJia(s string) {
	if this.KaiPanJia == nil {
		this.KaiPanJia = new(string)
	}
	*this.KaiPanJia = s
}

func (this *KXian) SetZuiGaoJia(s string) {
	if this.ZuiGaoJia == nil {
		this.ZuiGaoJia = new(string)
	}
	*this.ZuiGaoJia = s
}

func (this *KXian) SetZuiDiJia(s string) {
	if this.ZuiDiJia == nil {
		this.ZuiDiJia = new(string)
	}
	*this.ZuiDiJia = s
}

func (this *KXian) SetShouPanJia(s string) {
	if this.ShouPanJia == nil {
		this.ShouPanJia = new(string)
	}
	*this.ShouPanJia = s
}

func (this *KXian) SetChengJiaoLiang(s string) {
	if this.ChengJiaoLiang == nil {
		this.ChengJiaoLiang = new(string)
	}
	*this.ChengJiaoLiang = s
}

func (this *KXian) SetChengJiaoE(s string) {
	if this.ChengJiaoE == nil {
		this.ChengJiaoE = new(string)
	}
	*this.ChengJiaoE = s
}

func (this *KXian) SetChengJiaoBiShu(s string) {
	if this.ChengJiaoBiShu == nil {
		this.ChengJiaoBiShu = new(string)
	}
	*this.ChengJiaoBiShu = s
}
