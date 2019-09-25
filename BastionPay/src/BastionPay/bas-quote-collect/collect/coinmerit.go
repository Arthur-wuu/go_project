package collect

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/base"
	. "BastionPay/bas-quote-collect/config"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//配置文件中的币的行情 period limit
func (this *CoinMarket) GetCoinQoute(exchange, symbol, tp, size string) (*MoneyInfo, bool, error) {
	fmt.Println("exchange, symbol,tp, size ", exchange, symbol, tp, size)
	value := make(url.Values)
	if len(exchange) != 0 {
		value.Set("exchange", exchange)
	}
	if len(symbol) != 0 {
		value.Set("symbol", symbol)
	}
	if len(tp) != 0 {
		value.Set("type", tp)
	}
	if len(size) != 0 {
		count, _ := strconv.Atoi(size)
		if count == 0 {
			size = "13"
		} else if count > 2000 {
			size = "2000"
		}
		value.Set("size", size)
	} else {
		value.Set("size", "13")
	}
	//if len(since) != 0 {
	//	value.Set("since", since)
	//}

	param := value.Encode()
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	sign := this.genSign("kline", timeStamp, param)
	//	ZapLog().Info("models coinmerit "+config.GConfig.CoinMerit.HttpUrl+"/kline?"+param)
	content, err := base.HttpSend(GConfig.CoinMerit.HttpUrl+"/kline?"+param, nil, "GET", map[string]string{"ApiKey": GConfig.CoinMerit.ApiKey, "Timestamp": timeStamp, "ApiSign": sign})
	if err != nil {
		return nil, true, err
	}
	//fmt.Println("res content :",content)
	res := new(ResCoinMeritKLine)
	if err = json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("data", string(content)))
		return nil, false, err
	}
	//fmt.Println("res :",res)
	ZapLog().Info("res=", zap.Int("num", len(res.Data)))
	if res.Status_code != 200 {
		return nil, false, fmt.Errorf("%d %s", res.Status_code, res.Message)
	}
	if res.Data == nil || len(res.Data) == 0 {
		return nil, false, fmt.Errorf(symbol + " has nil data")
	}
	moneyInfo := CoinMeritKXianToMoneyInfo(res)
	//fmt.Println("moneyInfo: " ,moneyInfo)
	return moneyInfo, false, nil
}

func (this *CoinMarket) genSign(path, timestamp string, params string) string {
	md := md5.New()
	str := GConfig.CoinMerit.Secret_key + path + timestamp + params + GConfig.CoinMerit.Secret_key
	//	fmt.Println(str)
	//newStr := md.Sum([]byte(str))
	md.Write([]byte(str))
	newStr := md.Sum(nil) //这个只是追加
	//	fmt.Println(len(newStr))
	return strings.ToUpper(string(fmt.Sprintf("%X", newStr)))
}

func TimeToTimestamp(t string) int64 {
	datetime := t //待转化为时间戳的字符串

	//日期转化为时间戳
	timeLayout := "2006-01-02T15:04:05.000Z" //转化所需模板
	loc, _ := time.LoadLocation("GMT")       //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()
	return timestamp
}

func CoinMeritKXianToMoneyInfo(k1 *ResCoinMeritKLine) *MoneyInfo {
	if k1.Data == nil || len(k1.Data) == 0 || len(k1.Data[0]) < 5 {
		return nil
	}
	moneyInfo := new(MoneyInfo)
	moneyInfo.Price = &k1.Data[0][4]
	str := strconv.FormatFloat(k1.Data[0][0], 'f', -1, 64)
	times, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		ZapLog().Error("strconv.Atoi err")
	}
	moneyInfo.Last_updated = &times
	return moneyInfo
}
