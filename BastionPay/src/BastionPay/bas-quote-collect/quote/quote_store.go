package quote

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/collect"
	"BastionPay/bas-quote-collect/config"
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	//"os"
	//"sort"
	"strings"
)

func (this *QuoteMgr) putQuote(tickerInfos *collect.TickerInfo, coin string) error {
	if tickerInfos == nil {
		return nil
	}
	//pipeline过多问题
	//batch := new(leveldb.Batch)
	redisParam := make([]interface{}, 0)
	for i := 0; i < len(tickerInfos.IdDetailInfos); i++ {
		IdDetailInfo := tickerInfos.IdDetailInfos[i]
		lastupdated := IdDetailInfo.Last_updated

		for my, moneyInfo := range IdDetailInfo.Quotes {
			if coin != my {
				continue
			}
			if moneyInfo.Price == nil || *moneyInfo.Price <= 0.00000000000001 {
				continue
			}
			moneyInfo.Last_updated = lastupdated
			content, err := json.Marshal(moneyInfo)

			if err != nil {
				ZapLog().With(zap.Error(err), zap.Any("moneyInfo", moneyInfo)).Error("Marshal err")
				continue
			}
			key := this.genQuoteKey(IdDetailInfo.Id, my)
			ZapLog().With(zap.String("key", string(key))).Debug("redis put")

			redisParam = append(redisParam, key)
			redisParam = append(redisParam, content)
		}
	}
	_, err := this.mRedis.Do("MSET", redisParam...)
	if err != nil {
		return err
	}
	return nil
}

func (this *QuoteMgr) putHuilv(moneyInfo *collect.MoneyInfo, cCode string) error {
	if moneyInfo == nil {
		return nil
	}
	if moneyInfo.Price == nil || *moneyInfo.Price <= 0.00000000000001 {
		return nil
	}
	moneyInfo.Symbol = &cCode
	content, err := json.Marshal(moneyInfo)

	if err != nil {
		ZapLog().With(zap.Error(err), zap.Any("moneyInfo", moneyInfo)).Error("Marshal err")
		return err
	}
	key := this.genHuilvKey(cCode)
	ZapLog().With(zap.String("key", string(key))).Debug("redis put")
	_, err1 := this.mRedis.Do("SET", key, content)
	if err1 != nil {
		return err
	}
	return nil
}

//配置文件中的币的行情
func (this *QuoteMgr) putCoinQoute(moneyInfo *collect.MoneyInfo, cSymbol, cNum, cTo string) error {
	if moneyInfo == nil {
		return nil
	}
	if moneyInfo.Price == nil || *moneyInfo.Price <= 0.00000000000001 {
		return nil
	}
	moneyInfo.Symbol = &cSymbol
	content, err := json.Marshal(moneyInfo)

	if err != nil {
		ZapLog().With(zap.Error(err), zap.Any("moneyInfo", moneyInfo)).Error("Marshal err")
		return err
	}
	key := this.genCoinQouteKey(cNum, cTo)
	ZapLog().With(zap.String("key", string(key))).Debug("redis put")
	_, err1 := this.mRedis.Do("SET", key, content)
	if err1 != nil {
		return err
	}
	return nil
}

func (this *QuoteMgr) putCoinMarketCapQouteAll(res *collect.ResCoinMarketCapAll) error {
	if res == nil {
		return nil
	}
	redisParam := make([]interface{}, 0)
	//fmt.Println("len(res.Data)",len(res.Data))
	for j := 0; j < len(res.Data); j++ {

		IdDetailInfo := res.Data[j]
		lastupdated := res.Status.Timestamp
		timestampInt := collect.TimeToTimestamp(lastupdated)
		moneyInfo := new(collect.MoneyInfo)
		symbol := res.Data[j].Symbol
		price := res.Data[j].Quote["USD"].Price
		if price <= 0.00000000000001 {
			continue
		}

		moneyInfo.Symbol = &symbol
		moneyInfo.Price = &price
		moneyInfo.Last_updated = &timestampInt

		content, err := json.Marshal(moneyInfo)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("moneyInfo", moneyInfo)).Error("Marshal err")
			continue
		}
		//fmt.Println("id", IdDetailInfo.Id )
		key := this.genCoinMarketCapQouteKey(IdDetailInfo.Id)
		ZapLog().With(zap.String("key", string(key))).Debug("redis put")
		redisParam = append(redisParam, key)
		redisParam = append(redisParam, content)
	}

	_, err := this.mRedis.Do("MSET", redisParam...)
	if err != nil {
		return err
	}
	return nil
}

//码表中的币
func (this *QuoteMgr) putCoinMarketCapQoutePart(res *collect.ResCoinMarketCapPart, codes string) error {
	if res == nil {
		return nil
	}
	redisParam := make([]interface{}, 0)
	//fmt.Println("len(res.Data)",len(res.Data))
	codeArr := strings.Split(codes, ",")
	codeIntArr, _ := StringArrToIntArr(codeArr)
	codeIntArr = BubbleSort(codeIntArr)
	if len(codeIntArr) != len(res.Data) {
		ZapLog().Error("err", zap.Any("codeArr", codeArr), zap.Any("codeIntArr", codeIntArr), zap.Int("codeIntArrLen", len(codeIntArr)), zap.Int("codeIntArrlen", len(codeIntArr)))
		ZapLog().With(zap.Any("errmessage", "len(codeArr) != len(res.Data)")).Error("data err")
	}
	for j := 0; j < len(res.Data); j++ {
		IdDetailInfo := res.Data[strconv.Itoa(codeIntArr[j])]
		lastupdated := res.Status.Timestamp
		timestampInt := collect.TimeToTimestamp(lastupdated)
		moneyInfo := new(collect.MoneyInfo)
		symbol := res.Data[strconv.Itoa(codeIntArr[j])].Symbol
		price := res.Data[strconv.Itoa(codeIntArr[j])].Quote["USD"].Price
		if price <= 0.00000000000001 {
			continue
		}

		moneyInfo.Symbol = &symbol
		moneyInfo.Price = &price
		moneyInfo.Last_updated = &timestampInt

		content, err := json.Marshal(moneyInfo)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("moneyInfo", moneyInfo)).Error("Marshal err")
			continue
		}
		//fmt.Println("id", IdDetailInfo.Id )
		key := this.genCoinMarketCapQouteKey(IdDetailInfo.Id)
		ZapLog().With(zap.String("key", string(key))).Debug("redis put")
		redisParam = append(redisParam, key)
		redisParam = append(redisParam, content)
	}

	_, err := this.mRedis.Do("MSET", redisParam...)
	if err != nil {
		return err
	}
	return nil
}

//put huilv 行情  sina
func (this *QuoteMgr) putSinaHuilvQoute(monInfos []*collect.MoneyInfo) error {
	if monInfos == nil || len(monInfos) == 0 {
		return nil
	}
	redisParam := make([]interface{}, 0)

	for j := 0; j < len(monInfos); j++ {
		if monInfos[j].Price == nil || *monInfos[j].Price <= 0.00000000000001 {
			continue
		}
		content, err := json.Marshal(monInfos[j])
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("moneyInfo", monInfos[j])).Error("Marshal err")
			continue
		}
		key := this.genHuilvKey(config.GPreConfig.CountryCodeArr[j])
		ZapLog().With(zap.String("key", string(key))).Debug("redis put")
		redisParam = append(redisParam, key)
		redisParam = append(redisParam, content)
	}

	_, err := this.mRedis.Do("MSET", redisParam...)
	if err != nil {
		return err
	}
	return nil
}
