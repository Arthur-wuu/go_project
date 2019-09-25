package collect

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/base"
	"go.uber.org/zap"

	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

//配置文件中的币的行情 ， 在btcexa拿数据
func (this *CoinMarket) GetCoinQouteBtcexa(symbol string) (*MoneyInfo, bool, error) {
	//fmt.Println(" symbol is trade_pair ", symbol)
	value := make(url.Values)
	if len(symbol) != 0 {
		value.Set("trading_pair", symbol)
	}
	value.Set("period", "1m")
	value.Set("limit", "1")

	param := value.Encode()

	content, err := base.HttpSend("https://www.btcexa.com/api/market/kline?"+param, nil, "GET", nil)
	if err != nil {
		//		ZapLog().Error("http://www.btcexa.com/api/market/kline?"+param)
		return nil, true, err
	}
	//fmt.Println("res content :",content)
	res := new(ResBtcExaKLine)
	if err = json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("data", string(content)))
		return nil, false, err
	}
	//fmt.Println("res :",res)
	ZapLog().Debug("res=", zap.Int("num", len(res.Result)))
	if res.Status.Code != 0 {
		return nil, false, fmt.Errorf("%d %s", res.Status.Code, res.Status.Msg)
	}
	if res.Result == nil || len(res.Result) == 0 {
		return nil, false, fmt.Errorf(symbol + " has nil data")
	}
	moneyInfo := BtcexaKXianToMoneyInfo(res)
	//fmt.Println("moneyInfo: " ,moneyInfo)
	return moneyInfo, false, nil
}

func BtcexaKXianToMoneyInfo(k1 *ResBtcExaKLine) *MoneyInfo {
	if k1.Result == nil || len(k1.Result) == 0 || len(k1.Result[0]) < 5 {
		return nil
	}
	moneyInfo := new(MoneyInfo)
	str := k1.Result[0][4]
	//fmt.Println("str :" ,str)
	float, err := strconv.ParseFloat(str, 64)
	//fmt.Println("float :" ,float)
	times, err := strconv.ParseInt(k1.Result[0][0], 10, 64)
	times = times / 1000
	//	t,err :=strconv.Atoi(k1.Result[0][0])
	if err != nil {
		ZapLog().Error("strconv.Atoi err")
	}
	moneyInfo.Price = &float
	moneyInfo.Last_updated = &times
	return moneyInfo
}
