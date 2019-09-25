package quote

import "BastionPay/bas-quote/collect"

const (
	DB_Quote_prefix = "qt"
	//	DB_Code_prefix      = "cd"
	//	DB_Code_Time_prefix = "code_lasttime"
	CONST_KXIAN_1Day_Prefix = "kxian_1day"
)

type CodeInfoArr []collect.CodeInfo

func (this CodeInfoArr) Len() int {
	return len(this)
}

//如果index为i的元素小于index为j的元素，则返回true，否则返回false
func (this CodeInfoArr) Less(i, j int) bool {
	if this[i].Id == this[j].Id {
		return false
	}
	if this[i].Id == nil {
		return true
	}
	if this[j].Id == nil {
		return false
	}
	return *this[i].Id < *this[j].Id
}

// Swap 交换索引为 i 和 j 的元素
func (this CodeInfoArr) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type CoinDetailInfo struct {
	Symbol     *string              `json:"symbol,omitempty"`
	Id         *int                 `json:"id,omitempty"`
	MoneyInfos []*collect.MoneyInfo `json:"quotes"`
}

type SymbolInfo struct {
	Id      *int    `json:"id,omitempty"`
	Name    string  `json:"name"`
	Symbol  *string `json:"symbol,omitempty"`
	Website string  `json:"website_slug"`
}

func (this *SymbolInfo) SetId(id int) {
	if this.Id == nil {
		this.Id = new(int)
	}
	*this.Id = id
}

func (this *SymbolInfo) ClearId() {
	this.Id = nil
}

func (this *SymbolInfo) SetSymbol(s string) {
	if this.Symbol == nil {
		this.Symbol = new(string)
	}
	*this.Symbol = s
}

type HistoryDetailInfo struct {
	Symbol      *string        `json:"symbol,omitempty" doc:"币简称"`
	Id          *int           `json:"id,omitempty" doc:"币id"`
	KXianDetail []*KXianDetail `json:"kxian,omitempty" doc:"币行情数据"`
}

func (this *HistoryDetailInfo) AddKXianDetail() *KXianDetail {
	if this.KXianDetail == nil {
		this.KXianDetail = make([]*KXianDetail, 0)
	}
	k := new(KXianDetail)
	this.KXianDetail = append(this.KXianDetail, k)
	return k
}

func (this *HistoryDetailInfo) GetSymbol() string {
	if this.Symbol == nil {
		return ""
	}
	return *this.Symbol
}

func (this *HistoryDetailInfo) GetId() int {
	if this.Id == nil {
		return 0
	}
	return *this.Id
}

type KXianDetail struct {
	Symbol *string  `json:"symbol,omitempty" doc:"币简称"`
	KXians []*KXian `json:"detail,omitempty" doc:"币历史kxian数据"`
}

func (this *KXianDetail) SetSymbol(s string) {
	if this.Symbol == nil {
		this.Symbol = new(string)
	}
	*this.Symbol = s
}

type KXian struct {
	Timestamp  *int64   `json:"timestamp,omitempty"`   //日期
	OpenPrice  *float64 `json:"open_price,omitempty"`  //开盘价
	ClosePrice *float64 `json:"close_price,omitempty"` //收盘价
	LastPrice  *float64 `json:"last_price,omitempty"`  //最新价
	HighPrice  *float64 `json:"high_price,omitempty"`  //最高
	LowPrice   *float64 `json:"low_price,omitempty"`   //最低
}

func (this *KXian) GetTimestamp() int64 {
	if this.Timestamp == nil {
		return 0
	}
	return *this.Timestamp
}

func (this *KXian) SetTimestamp(i int64) {
	if this.Timestamp == nil {
		this.Timestamp = new(int64)
	}
	*this.Timestamp = i
}

func (this *KXian) GetOpenPrice() float64 {
	if this.OpenPrice == nil {
		return 0
	}
	return *this.OpenPrice
}

func (this *KXian) SetOpenPrice(p float64) {
	if this.OpenPrice == nil {
		this.OpenPrice = new(float64)
	}
	*this.OpenPrice = p
}

func (this *KXian) GetClosePrice() float64 {
	if this.ClosePrice == nil {
		return 0
	}
	return *this.ClosePrice
}

func (this *KXian) SetClosePrice(p float64) {
	if this.ClosePrice == nil {
		this.ClosePrice = new(float64)
	}
	*this.ClosePrice = p
}

func (this *KXian) GetLastPrice() float64 {
	if this.LastPrice == nil {
		return 0
	}
	return *this.LastPrice
}

func (this *KXian) SetLastPrice(p float64) {
	if this.LastPrice == nil {
		this.LastPrice = new(float64)
	}
	*this.LastPrice = p
}

func (this *KXian) GetHighPrice() float64 {
	if this.HighPrice == nil {
		return 0
	}
	return *this.HighPrice
}

func (this *KXian) SetHighPrice(p float64) {
	if this.HighPrice == nil {
		this.HighPrice = new(float64)
	}
	*this.HighPrice = p
}

func (this *KXian) GetLowPrice() float64 {
	if this.LowPrice == nil {
		return 0
	}
	return *this.LowPrice
}

func (this *KXian) SetLowPrice(p float64) {
	if this.LowPrice == nil {
		this.LowPrice = new(float64)
	}
	*this.LowPrice = p
}
