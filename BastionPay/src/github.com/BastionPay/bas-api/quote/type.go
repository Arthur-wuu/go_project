package quote

type MoneyInfo struct {
	Id         *int             `json:"id,omitempty" doc:"币id"`
	Symbol             *string `json:"symbol,omitempty" doc:"币简称"`
	Price              *float64 `json:"price,omitempty" doc:"最新价"`
	Volume_24h         *float64 `json:"volume_24h,omitempty" doc:"24小时成交量"`
	Market_cap         *float64 `json:"market_cap,omitempty" doc:"总市值"`
	Percent_change_1h  *float64 `json:"percent_change_1h,omitempty" doc:"1小时涨跌幅"`
	Percent_change_24h *float64 `json:"percent_change_24h,omitempty" doc:"24小时涨跌幅"`
	Percent_change_7d  *float64 `json:"percent_change_7d,omitempty" doc:"7天涨跌幅"`
	Last_updated       *int64  `json:"last_updated,omitempty" doc:"最近更新时间"`
}

func (this * MoneyInfo)GetSymbol()string {
	if this.Symbol == nil {
		return ""
	}
	return *this.Symbol
}

func (this *MoneyInfo) SetSymbol(coin string) {
	if this.Symbol == nil {
		this.Symbol = new(string)
	}
	*this.Symbol = coin
}

func (this *MoneyInfo) SetPrice(p float64) {
	if this.Price == nil {
		this.Price = new(float64)
	}
	*this.Price = p
}

type CodeInfo struct {
	Id           *int    `json:"id,omitempty" doc:"币id"`
	Name         *string `json:"name,omitempty" doc:"币全称"`
	Symbol       *string `json:"symbol,omitempty" doc:"币简称，行情查询采用该字段"`
	Website_slug *string `json:"website_slug,omitempty" doc:"币相关网站信息"`
	Timestamp    *int64  `json:"timestamp,omitempty" doc:"最近更新时间"`
	Valid        *int  `json:"valid,omitempty"  doc:"是否有效"`
}

func (this * CodeInfo)GetSymbol() string{
	if this.Symbol == nil {
		return ""
	}
	return *this.Symbol
}

func (this * CodeInfo)GetId() int{
	if this.Id == nil {
		return 0
	}
	return *this.Id
}

type QuoteDetailInfo struct {
	Symbol     *string              `json:"symbol,omitempty" doc:"币简称"`
	Id         *int                 `json:"id,omitempty" doc:"币id"`
	MoneyInfos []*MoneyInfo         `json:"detail,omitempty" doc:"币行情数据"`
}

func (this *QuoteDetailInfo) SetSymbol(s string) {
	if this.Symbol == nil {
		this.Symbol = new(string)
	}
	*this.Symbol = s
}

func (this *QuoteDetailInfo) SetId(i int) {
	if this.Id == nil {
		this.Id = new(int)
	}
	*this.Id = i
}

func (this *QuoteDetailInfo) AddMoneyInfo(info *MoneyInfo) {
	if this.MoneyInfos == nil {
		this.MoneyInfos = make([]*MoneyInfo, 0)
	}
	this.MoneyInfos = append(this.MoneyInfos, info)
}

type ResMsg struct{
	Err      int                 `json:"err" doc:"错误码。0成功；10001参数错误；其余错误"`
	ErrMsg   *string             `json:"errmsg,omitempty" doc:"错误相关信息"`
	Quotes   []*QuoteDetailInfo  `json:"quotes,omitempty" doc:"行情信息"`
	Codes    []*CodeInfo         `json:"codes,omitempty" doc:"码表信息"`
	Historys []*HistoryDetailInfo `json:"historys,omitempty" doc:"K线信息"`
	Coins    *CoinInfo         `json:"coins,omitempty" doc:"码表信息"`
}

type CoinInfo struct{
	Digitals  []*CodeInfo         `json:"digital_coin,omitempty" doc:"码表信息"`
	Legals    []*CodeInfo           `json:"legal_coin,omitempty" doc:"码表信息"`
}

func NewResMsg(e int, msg string) *ResMsg {
	res := new(ResMsg)
	res.Err = e
	res.SetErrMsg(msg)
	return res
}

func (this *ResMsg) SetErrMsg(info string) {
	if len(info) == 0 {
		return
	}
	if this.ErrMsg == nil {
		this.ErrMsg = new(string)
	}
	*this.ErrMsg = info
}

func (this *ResMsg) GenQuoteDetailInfo() *QuoteDetailInfo {
	if this.Quotes == nil {
		this.Quotes = make([]*QuoteDetailInfo, 0)
	}
	info := new(QuoteDetailInfo)
	this.Quotes = append(this.Quotes, info)
	return info
}

func (this *ResMsg) GenHistoryDetailInfo() *HistoryDetailInfo {
	if this.Historys == nil {
		this.Historys = make([]*HistoryDetailInfo, 0)
	}
	info := new(HistoryDetailInfo)
	this.Historys = append(this.Historys, info)
	return info
}

type HistoryDetailInfo struct {
	Symbol     *string              `json:"symbol,omitempty" doc:"币简称"`
	Id         *int                 `json:"id,omitempty" doc:"币id"`
	KXianDetail  []*KXianDetail        `json:"kxian,omitempty" doc:"币行情数据"`
}

func (this *HistoryDetailInfo) AddKXianDetail() *KXianDetail {
	if this.KXianDetail == nil {
		this.KXianDetail = make([]*KXianDetail, 0)
	}
	k := new(KXianDetail)
	this.KXianDetail = append(this.KXianDetail, k)
	return k
}

func (this * HistoryDetailInfo) GetSymbol() string {
	if this.Symbol == nil {
		return ""
	}
	return *this.Symbol
}

func (this * HistoryDetailInfo) GetId() int {
	if this.Id == nil {
		return 0
	}
	return *this.Id
}

type KXianDetail struct{
	Symbol     *string              `json:"symbol,omitempty" doc:"币简称"`
	KXians     []*KXian             `json:"detail,omitempty" doc:"币历史kxian数据"`
}

func (this *KXianDetail)SetSymbol(s string){
	if this.Symbol == nil {
		this.Symbol = new(string)
	}
	*this.Symbol = s
}


type KXian struct{
	Timestamp *int64   `json:"timestamp,omitempty"` //日期
	OpenPrice *float64  `json:"open_price,omitempty"` //开盘价
	ClosePrice *float64 `json:"close_price,omitempty"` //收盘价
	LastPrice *float64  `json:"last_price,omitempty"` //最新价
	HighPrice *float64  `json:"high_price,omitempty"` //最高
	LowPrice  *float64  `json:"low_price,omitempty"` //最低
}

func (this * KXian) GetTimestamp() int64 {
	if this.Timestamp == nil {
		return 0
	}
	return *this.Timestamp
}

func (this * KXian) GetOpenPrice() float64 {
	if this.OpenPrice == nil {
		return 0
	}
	return *this.OpenPrice
}

func (this * KXian) GetClosePrice() float64 {
	if this.ClosePrice == nil {
		return 0
	}
	return *this.ClosePrice
}

func (this * KXian) GetLastPrice() float64 {
	if this.LastPrice == nil {
		return 0
	}
	return *this.LastPrice
}

func (this * KXian) GetHighPrice() float64 {
	if this.HighPrice == nil {
		return 0
	}
	return *this.HighPrice
}

func (this * KXian) GetLowPrice() float64 {
	if this.Timestamp == nil {
		return 0
	}
	return *this.LowPrice
}