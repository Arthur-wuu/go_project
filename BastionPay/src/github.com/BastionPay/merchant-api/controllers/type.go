package controllers

type ResMsg struct{
	Err      int                 `json:"err" doc:"错误码。0成功；10001参数错误；其余错误"`
	ErrMsg   *string             `json:"errmsg,omitempty" doc:"错误相关信息"`
	Quotes   []*QuoteDetailInfo  `json:"quotes,omitempty" doc:"行情信息"`
	Codes    []*CodeInfo         `json:"codes,omitempty" doc:"码表信息"`
	Historys []*HistoryDetailInfo `json:"historys,omitempty" doc:"K线信息"`
}


type QuoteDetailInfo struct {
	Symbol     *string              `json:"symbol,omitempty" doc:"币简称"`
	Id         *int                 `json:"id,omitempty" doc:"币id"`
	MoneyInfos []*MoneyInfo         `json:"detail,omitempty" doc:"币行情数据"`
}

type MoneyInfo struct {
	Symbol             *string `json:"symbol,omitempty" doc:"币简称"`
	Price              *float64 `json:"price,omitempty" doc:"最新价"`
	Volume_24h         *float64 `json:"volume_24h,omitempty" doc:"24小时成交量"`
	Market_cap         *float64 `json:"market_cap,omitempty" doc:"总市值"`
	Percent_change_1h  *float64 `json:"percent_change_1h,omitempty" doc:"1小时涨跌幅"`
	Percent_change_24h *float64 `json:"percent_change_24h,omitempty" doc:"24小时涨跌幅"`
	Percent_change_7d  *float64 `json:"percent_change_7d,omitempty" doc:"7天涨跌幅"`
	Last_updated       *int64  `json:"last_updated,omitempty" doc:"最近更新时间"`
	Fee                *string  `json:"fee,omitempty" doc:"手续费"`
}


type CodeInfo struct {
	Id           *int    `json:"id,omitempty" doc:"币id"`
	Name         *string `json:"name,omitempty" doc:"币全称"`
	Symbol       *string `json:"symbol,omitempty" doc:"币简称，行情查询采用该字段"`
	Website_slug *string `json:"website_slug,omitempty" doc:"币相关网站信息"`
	Timestamp    *int64  `json:"timestamp,omitempty" doc:"最近更新时间"`
}

type HistoryDetailInfo struct {
	Symbol     *string              `json:"symbol,omitempty" doc:"币简称"`
	Id         *int                 `json:"id,omitempty" doc:"币id"`
	KXianDetail  []*KXianDetail        `json:"kxian,omitempty" doc:"币行情数据"`
}

type KXianDetail struct{
	Symbol     *string              `json:"symbol,omitempty" doc:"币简称"`
	KXians     []*KXian             `json:"detail,omitempty" doc:"币历史kxian数据"`
}

type KXian struct{
	Timestamp *int64   `json:"timestamp,omitempty"` //日期
	OpenPrice *float64  `json:"open_price,omitempty"` //开盘价
	ClosePrice *float64 `json:"close_price,omitempty"` //收盘价
	LastPrice *float64  `json:"last_price,omitempty"` //最新价
	HighPrice *float64  `json:"high_price,omitempty"` //最高
	LowPrice  *float64  `json:"low_price,omitempty"` //最低
}
