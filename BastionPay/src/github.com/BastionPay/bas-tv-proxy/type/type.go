package _type

import (
	"BastionPay/bas-tv-proxy/api"
	"errors"
	"fmt"
)

const (
	EVENTCoinMerit_sub   = "subscribe"
	EVENTCoinMerit_unsub = "unsubscribe"
	EVENTBtcExa_sub      = "sub"
	EVENTBtcExa_unsub    = "uns"
)

func NewReqCoinMeritSub(exa, obj, period, sub string) (*ReqCoinMeritSub, error) {
	if len(exa) == 0 || len(obj) == 0 || len(period) == 0 {
		return nil, errors.New("param err")
	}
	req := new(ReqCoinMeritSub)
	req.Topic = "cm_" + exa + "_" + obj + "_kline_" + period
	switch sub {
	case "":
		req.Event = EVENTCoinMerit_sub
		break
	case api.EVENT_sub:
		req.Event = EVENTCoinMerit_sub
		break
	case api.EVENT_unsub:
		req.Event = EVENTCoinMerit_unsub
		break
	default:
		req.Event = EVENTCoinMerit_sub
	}
	return req, nil
}

func NewReqBtcExaSub(sub, obj, period, count string) (*ReqBtcExaSub, error) {
	if len(sub) == 0 || len(obj) == 0 || len(period) == 0 {
		return nil, errors.New("param err")
	}
	req := new(ReqBtcExaSub)
	switch sub {
	case "":
		sub = EVENTBtcExa_sub
		break
	case api.EVENT_sub:
		sub = EVENTBtcExa_sub
		break
	case api.EVENT_unsub:
		sub = EVENTBtcExa_unsub
		break
	default:
		sub = EVENTBtcExa_sub
	}
	fmt.Println("count:", count)
	count = "1"
	if len(count) == 0 {
		count = "1"
		req.Event = sub
		req.Topic = "market." + obj + ".kline|{\"period\":\"" + period + "\",\"limit\":" + count + "}"
		fmt.Println("req.Topic:" + req.Topic)
	} else {
		req.Event = sub
		req.Topic = "market." + obj + ".kline|{\"period\":\"" + period + "\",\"limit\":" + count + "}"
		fmt.Println("req.Topic:" + req.Topic)
	}
	return req, nil
}

type ReqCoinMeritSub struct {
	Event string `json:"event,omitempty"`
	Topic string `json:"topic,omitempty"`
}

type ResCoinMeritSub struct {
	Result     *bool               `json:"result, omitempty"`
	Error_code *int32              `json:"error_code, omitempty"`
	Topic      string              `json:"topic,omitempty"`
	Data       []CoinMeritSubKXian `json:"data,omitempty"`
}

type CoinMeritSubKXian struct {
	Pair string    `json:"pair,omitempty"`
	Data []float64 `json:"data,omitempty"`
}

type ResCoinMeritExchanges struct {
	Status_code int32    `json:"status_code, omitempty"`
	Message     string   `json:"message, omitempty"`
	Data        []string `json:"data, omitempty"`
}

type ResCoinMeritKLine struct {
	Status_code int32       `json:"status_code, omitempty"`
	Message     string      `json:"message, omitempty"`
	Data        [][]float64 `json:"data, omitempty"`
}

type ResCoinMeritCurrencyPairs struct {
	Status_code int32                     `json:"status_code, omitempty"`
	Message     string                    `json:"message, omitempty"`
	Data        CoinMeritExaCurrencyPairs `json:"data, omitempty"`
}

type CoinMeritExaCurrencyPairs struct {
	Exchange      string   `json:"exchange,omitempty"`
	CurrencyPairs []string `json:"currency_pairs,omitempty"`
}

//btcexa
type ReqBtcExaSub struct {
	Event string
	Topic string //`json:"topic,omitempty"`
}

func (this *ReqBtcExaSub) GetTopic() string {
	return this.Topic
}

func (this *ReqBtcExaSub) GetReqStr() string {
	return this.Event + "." + this.Topic
}

type ResBtcExaSubU struct {
	Topic  string   `json:"topic,omitempty"`
	Status string   `json:"status, omitempty"`
	Type   string   `json:"type, omitempty"`
	Data   []string `json:"data,omitempty"`
}

type ResBtcExaSubI struct {
	Topic  string     `json:"topic,omitempty"`
	Status string     `json:"status, omitempty"`
	Type   string     `json:"type, omitempty"`
	Data   [][]string `json:"data,omitempty"`
}

type ResBtcExaKLine struct {
	Status BtcExaStatus `json:"status, omitempty"`
	Result [][]string   `json:"result, omitempty"`
}

type BtcExaStatus struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type ResBtcExaObjList struct {
	Status BtcExaObjStatus `json:"status, omitempty"`
	Result []BtcExaObjInfo `json:"result, omitempty"`
}

type BtcExaObjStatus struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type BtcExaObjInfo struct {
	Name          string `json:"name, omitempty"`
	Base          string `json:"base, omitempty"`
	Base_decimal  int32  `json:"base_decimal, omitempty"`
	Quote         string `json:"quote, omitempty"`
	Quote_decimal int32  `json:"quote_decimal, omitempty"`
	Quantity_min  string `json:"quantity_min, omitempty"`
	Amount_min    string `json:"amount_min, omitempty"`
	Can_trading   bool   `json:"can_trading, omitempty"`
}
