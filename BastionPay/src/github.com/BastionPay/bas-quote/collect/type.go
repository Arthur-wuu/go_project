package collect

import (
	"fmt"
)

type CodeListInfo struct {
	Data     []CodeInfo   `json:"data"`
	MetaData MetaDataInfo `json:"metadata"`
}

type MetaDataInfo struct {
	Timestamp            uint64 `json:"timestamp"`
	Num_cryptocurrencies int    `json:"num_cryptocurrencies"`
	Error                string `json:"error"`
}

type CodeInfo struct {
	Id           *int    `json:"id"`
	Name         *string `json:"name"`
	Symbol       *string `json:"symbol"`
	Website_slug *string `json:"website_slug"`
	Timestamp    *int64  `json:"timestamp"`
	Valid        *int    `json:"valid"`
}

func (this *CodeInfo) ToPrintStr() string {
	return fmt.Sprintf("id[%d]Name[%s]Symbol[%s]Website_slug[%s]Timestamp[%d]valid[%d]", this.GetId(), this.GetName(), this.GetSymbol(), this.GetWebsite_slug(), this.GetTimestamp(), this.GetValid())
}

func (this *CodeInfo) GetName() string {
	if this.Name == nil {
		return ""
	}
	return *this.Name
}

func (this *CodeInfo) GetWebsite_slug() string {
	if this.Website_slug == nil {
		return ""
	}
	return *this.Website_slug
}

func (this *CodeInfo) GetTimestamp() int64 {
	if this.Timestamp == nil {
		return 0
	}
	return *this.Timestamp
}

func (this *CodeInfo) SetTimestamp(t int64) {
	if this.Timestamp == nil {
		this.Timestamp = new(int64)
	}
	*this.Timestamp = t
}

func (this *CodeInfo) GetValid() int {
	if this.Valid == nil {
		return 0
	}
	return *this.Valid
}

func (this *CodeInfo) GetSymbol() string {
	if this.Symbol == nil {
		return ""
	}
	return *this.Symbol
}

func (this *CodeInfo) GetId() int {
	if this.Id == nil {
		return 0
	}
	return *this.Id
}

type TickerInfo struct {
	IdDetailInfos []IdDetailInfo `json:"data"`
	MetaDataInfo  MetaDataInfo   `json:"metadata"`
}

type IdDetailInfo struct {
	Id                 int                  `json:"id"`
	Name               string               `json:"name"`
	Symbol             string               `json:"symbol"`
	Website_slug       string               `json:"website_slug"`
	Rank               int                  `json:"rank"`
	Circulating_supply float64              `json:"circulating_supply"`
	Total_supply       float64              `json:"total_supply"`
	Max_supply         float64              `json:"max_supply"`
	Quotes             map[string]MoneyInfo `json:"quotes"`
	Last_updated       *int64               `json:"last_updated,omitempty"`
}

func (this *IdDetailInfo) GetLast_updated() int64 {
	if this.Last_updated == nil {
		return 0
	}
	return *this.Last_updated
}

//type QuotesInfo struct {
//	MoneyInfo map[string]MoneyInfo `json:"quotes"`
//}

type MoneyInfo struct {
	Symbol             *string  `json:"symbol"`
	Price              *float64 `json:"price"`
	Volume_24h         *float64 `json:"volume_24h"`
	Market_cap         *float64 `json:"market_cap"`
	Percent_change_1h  *float64 `json:"percent_change_1h"`
	Percent_change_24h *float64 `json:"percent_change_24h"`
	Percent_change_7d  *float64 `json:"percent_change_7d"`
	Last_updated       *int64   `json:"last_updated"`
}

func (this *MoneyInfo) GetPrice() float64 {
	if this.Price == nil {
		return 0
	}
	return *this.Price
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

func (this *MoneyInfo) GetLast_updated() int64 {
	if this.Last_updated == nil {
		return 0
	}
	return *this.Last_updated
}

func (this *MoneyInfo) SetLast_updated(p int64) {
	if this.Last_updated == nil {
		this.Last_updated = new(int64)
	}
	*this.Last_updated = p
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
