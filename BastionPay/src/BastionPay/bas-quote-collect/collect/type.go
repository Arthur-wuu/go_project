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
	Valid        *int    `json:"vaild"`
}

func (this *CodeInfo) GetValid() int {
	if this.Valid == nil {
		return 0
	}
	return *this.Valid
}

func (this *CodeInfo) SetValid(v int) {
	if this.Valid == nil {
		this.Valid = new(int)
	}
	*this.Valid = v
}

//汇率

//type HuiLvInfo struct {
//	Result     []HuiLvInfo1   `json:"Result"`
//}
//
//type HuiLvInfo1 struct {
//	DisplayData     HuiLvInfo2   `json:"DisplayData"`
//}
//
//type HuiLvInfo2 struct {
//	ResultData     HuiLvInfo3   `json:"resultData"`
//}
//
//type HuiLvInfo3 struct {
//	TplData     HuiLvInfo4   `json:"tplData"`
//}
//
//type HuiLvInfo4 struct {
//	Money2_num 	 	 string    `json:"money2_num"`     //与美元对应的价格
//	Pcdata  	 	 Time   	`json:"pcdata"`         //更新的时间
//}
//
//type Time struct {
//	UpdateTime string `json:"_update_time"`
//}

type HuiLvInfo struct {
	Data []HuiLvInfo2 `json:"data"`
}

type HuiLvInfo2 struct {
	Money2_num string `json:"number2"`      //与美元对应的价格
	UpdateTime string `json:"_update_time"` //更新的时间
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

func (this *MoneyInfo) SetSymbol(coin string) {
	if this.Symbol == nil {
		this.Symbol = new(string)
	}
	*this.Symbol = coin
}

type ResCoinMeritKLine struct {
	Status_code int32       `json:"status_code, omitempty"`
	Message     string      `json:"message, omitempty"`
	Data        [][]float64 `json:"data, omitempty"`
}

type ResBtcExaKLine struct {
	Status BtcExaStatus `json:"status, omitempty"`
	Result [][]string   `json:"result, omitempty"`
}
type BtcExaStatus struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

//new coinmarketcap response
type ResCoinMarketCapAll struct {
	Status Status       `json:"status, omitempty"`
	Data   []DetailInfo `json:"data, omitempty"`
}

type Status struct {
	Timestamp     string `json:"timestamp, omitempty"`
	Error_code    int32  `json:"error_code, omitempty"`
	Error_message string `json:"error_message, omitempty"`
	Elapsed       int32  `json:"elapsed, omitempty"`
	Credit_count  int32  `json:"credit_count, omitempty"`
}

type DetailInfo struct {
	Id                 int64                `json:"id, omitempty"`
	Name               string               `json:"name, omitempty"`
	Symbol             string               `json:"symbol, omitempty"`
	Slug               string               `json:"slug, omitempty"`
	Circulating_supply float64              `json:"circulating_supply, omitempty"`
	Total_supply       float64              `json:"total_supply, omitempty"`
	Max_supply         float64              `json:"max_supply, omitempty"`
	Date_added         string               `json:"date_added, omitempty"`
	Num_market_pairs   int64                `json:"num_market_pairs, omitempty"`
	Cmc_rank           int64                `json:"cmc_rank, omitempty"`
	Tags               []string             `json:"tags, omitempty"`
	Last_updated       string               `json:"last_updated, omitempty"`
	Platform           interface{}          `json:"platform, omitempty"`
	Quote              map[string]PriceInfo `json:"quote, omitempty"`
}

type PriceInfo struct {
	Price              float64 `json:"price, omitempty"`
	Volume_24h         float64 `json:"volume_24h, omitempty"`
	Percent_change_1h  float64 `json:"percent_change_1h, omitempty"`
	Percent_change_24h float64 `json:"percent_change_24h, omitempty"`
	Percent_change_7d  float64 `json:"percent_change_7d, omitempty"`
	Market_cap         float64 `json:"market_cap, omitempty"`
	Last_updated       string  `json:"last_updated, omitempty"`
}

type ResCoinMarketCapPart struct {
	Status Status                `json:"status, omitempty"`
	Data   map[string]DetailInfo `json:"data, omitempty"`
}

type CodeTableIntf interface {
	ListSymbols() []CodeInfo
	SetCodeTable(cc *CodeInfo) error
}
