package quote

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/config"
	"fmt"
	"go.uber.org/zap"
	"time"
)

//coinmarketcap新接口 2100种币==>>USD的信息
func (this *QuoteMgr) loadCoinTicker() {
	defer PanicPrint()
	coinMarketSuccessNum := 0
	putRedisSuccessNum := 0
	for {
		globalErrFlag := false
		startDayTime := GenDay(time.Now().Unix())

		for j := 1; true; j++ {
			mCoinMarketErrFlag := false
			redisErrFlag := false

			resTickerInfo, endFlag, err := this.mCoinMarket.GetTickerList(1+(j-1)*100, 100) //1+(j-1)*100, 100
			if err != nil {
				ZapLog().With(zap.Error(err), zap.String("get ticker err message", resTickerInfo.Status.Error_message)).Error("mCoinMarketCap GetPartTicker err")
				mCoinMarketErrFlag = true
				break
			}
			if endFlag || resTickerInfo.Data == nil || len(resTickerInfo.Data) == 0 {
				ZapLog().Sugar().Debugf("coinmarketcap coinparam[%s] ", resTickerInfo.Status.Error_code)
				break
			}

			if err = this.putCoinMarketCapQouteAll(resTickerInfo); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("err message", resTickerInfo.Status.Error_message)).Error("putCoinMarketCapQoute err")
			}

			// 存到 redis kxain
			if err = this.mHistory.CoinMarketCapSetRedisAll(resTickerInfo); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("redis put kxian", resTickerInfo.Status.Error_message)).Error("History.Set err")
			}

			time.Sleep(time.Second * time.Duration(20))
			if !redisErrFlag {
				putRedisSuccessNum++
			}
			if !mCoinMarketErrFlag {
				coinMarketSuccessNum++
			} else {
				globalErrFlag = true
				break
			}
			ZapLog().Sugar().Debugf("loadCoinMarketCap  coinMarketSuccessNum[%d]putRedisSuccessNum[%d]", coinMarketSuccessNum, putRedisSuccessNum)
			fmt.Printf("loadCoinMarketCap  coinMarketSuccessNum[%d]putRedisSuccessNum[%d]", coinMarketSuccessNum, putRedisSuccessNum)

		}
		sleepTime := 1
		if globalErrFlag {
			sleepTime = 3600
		} else if GenDay(time.Now().Unix()) > startDayTime {
			sleepTime = 1
		}
		ZapLog().Sugar().Debugf(fmt.Sprintf("complete once loop , sleep %d s", sleepTime))
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

//coinmarketcap新接口 数据库的码表==>>USD的信息
func (this *QuoteMgr) loadPartCoinTicker() {
	defer PanicPrint()

	for {
		resTickerInfo, codes, err := this.mCoinMarket.GetPartTicker() //1+(j-1)*100, 100
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("mCoinMarketCap GetPartTicker err")
			time.Sleep(time.Second * 60)
			goto COINMARKETCAP_LOADPART
		}
		if resTickerInfo == nil || resTickerInfo.Data == nil || len(resTickerInfo.Data) == 0 {
			ZapLog().Error("mCoinMarketCap GetPartTicker err:nil data")
			time.Sleep(time.Second * 60)
			goto COINMARKETCAP_LOADPART
		}
		ZapLog().Sugar().Debug("coinmarketcap part load codes[%s] ok", codes)

		if err = this.putCoinMarketCapQoutePart(resTickerInfo, codes); err != nil {
			ZapLog().With(zap.Error(err)).Error("putCoinMarketCapQoute err")
			time.Sleep(time.Second * 60)
		}

		// 存到 redis kxain
		if err = this.mHistory.CoinMarketCapSetRedisPart(resTickerInfo, codes); err != nil {
			ZapLog().With(zap.Error(err), zap.String("redis put kxian", resTickerInfo.Status.Error_message)).Error("History.Set err")
			time.Sleep(time.Second * 60)
		}

		//不同的环境 配不同的时间，dev一个key即可    pro 五分钟更新一次  五个即可    可都配300s
	COINMARKETCAP_LOADPART:
		time.Sleep(time.Second * time.Duration(config.GConfig.Coinmarketcap.Diff_env_interval))

		ZapLog().Sugar().Debugf("loadCoinMarketCap  codes[%v]", codes)
		//fmt.Printf("loadCoinMarketCap  coinMarketSuccessNum[%d]putRedisSuccessNum[%d]", coinMarketSuccessNum, putRedisSuccessNum)
	}

}

//法币汇率查询，这尿性怕是多个实例一块跑就要被IP限制了,因此得加上随机值。最佳方案是采集和查询分离，采集服务一个，查询服务多个，数据存redis中共享
func (this *QuoteMgr) loadQuotesAndStore(coins []string) {
	defer PanicPrint()
	allNum := len(coins)
	coinMarketSuccessNum := 0
	putRedisSuccessNum := 0
	for {
		globalErrFlag := false
		startDayTime := GenDay(time.Now().Unix())
		for i := 0; i < len(coins); i++ {
			coinTickerNum := 0
			tickErrFlag := false
			redisErrFlag := false
			for j := 1; true; j++ {
				tickerInfos, endFlag, _, err := this.mCoinMarket.GetAllTicker(coins[i], 1+(j-1)*100, 100)
				if err != nil {
					ZapLog().With(zap.Error(err), zap.String("coin", coins[i])).Error("mCoinMarket GetAllTicker err")
					tickErrFlag = true
					break
				}
				if endFlag || tickerInfos == nil || len(tickerInfos.IdDetailInfos) == 0 {
					ZapLog().Sugar().Debugf("coinmarketcap[%s]num[%d] getTicker end", coins[i], coinTickerNum)
					break
				}
				coinTickerNum += len(tickerInfos.IdDetailInfos)
				ZapLog().Debug("mCoinMarket GetAllTicker ", zap.Int("num", len(tickerInfos.IdDetailInfos)), zap.String("coin", coins[i]))

				//存储的地方要改成redis
				if err = this.putQuote(tickerInfos, coins[i]); err != nil {
					redisErrFlag = true
					ZapLog().With(zap.Error(err), zap.String("coin", coins[i])).Error("putQuote err")
				}
				if err = this.mHistory.Set(tickerInfos, coins[i]); err != nil {
					redisErrFlag = true
					ZapLog().With(zap.Error(err), zap.String("coin", coins[i])).Error("History.Set err")
				}
				time.Sleep(time.Second * time.Duration(13))
			}
			if !redisErrFlag {
				putRedisSuccessNum++
			}
			if !tickErrFlag {
				coinMarketSuccessNum++
			} else {
				globalErrFlag = true
				break
			}
			time.Sleep(time.Second * time.Duration(13))
		}
		ZapLog().Sugar().Debugf("loadQuotesAndStore all[%d]coinMarketSuccessNum[%d]putRedisSuccessNum[%d]", allNum, coinMarketSuccessNum, putRedisSuccessNum)
		sleepTime := 13
		if globalErrFlag {
			sleepTime = 3600
		} else if GenDay(time.Now().Unix()) > startDayTime {
			sleepTime = 15
		}
		ZapLog().Sugar().Debugf(fmt.Sprintf("complete once loop , sleep %d s", sleepTime))
		time.Sleep(time.Second * time.Duration(sleepTime))
	}

}

//拉取存储汇率  -- baidu
func (this *QuoteMgr) loadHuilvBaidu(cName, cCode []string) {
	defer PanicPrint()
	allNum := len(cName)
	HuilvSuccessNum := 0
	putRedisSuccessNum := 0
	for {
		globalErrFlag := false
		startDayTime := GenDay(time.Now().Unix())
		for i := 0; i < len(cName); i++ {
			huilvNum := 0
			huilvErrFlag := false
			redisErrFlag := false
			moneyInfo, endFlag, err := this.mHuilv.GetHuiLv(cName[i])
			if err != nil {
				ZapLog().With(zap.Error(err), zap.String("cName", cName[i])).Error("mHuilv GetHuilvInfo err")
				huilvErrFlag = true
				break
			}
			if endFlag {
				ZapLog().Sugar().Debugf("huilvInfo[%s]num[%d] getHuilvTicker end", cName[i], huilvNum)
				break
			}
			//存储在redis，
			if err = this.putHuilv(moneyInfo, cCode[i]); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("cName", cName[i])).Error("putHuilv err")
			}
			//fmt.Println("putHuilv succ。。。",cName[i])
			if err = this.mHistory.HuilvSetRedis(moneyInfo, cCode[i]); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("coin", cCode[i])).Error("History.Set err")
			}
			//fmt.Println("putHuilv redis kxian succ。。。",cCode[i])
			time.Sleep(time.Second * time.Duration(10))
			if !redisErrFlag {
				putRedisSuccessNum++
			}
			if !huilvErrFlag {
				HuilvSuccessNum++
			} else {
				globalErrFlag = true
				break
			}
			ZapLog().Sugar().Debugf("loadBaiduHuilv all[%d]huilvSuccessNum[%d]putRedisSuccessNum[%d]", allNum, HuilvSuccessNum, putRedisSuccessNum)
		}
		sleepTime := 1200
		if globalErrFlag {
			sleepTime = 1800
		} else if GenDay(time.Now().Unix()) > startDayTime {
			sleepTime = 60
		}
		ZapLog().Sugar().Debugf(fmt.Sprintf("complete once loop , sleep %d s", sleepTime))
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

//拉取存储汇率  -- sina
func (this *QuoteMgr) loadHuilvSina() {
	defer PanicPrint()
	for {
		moneyInfos, err := this.mHuilv.GetHuiLvSina()
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Error(err)).Error("GetHuiLvSina err")
			time.Sleep(time.Second * 300)
			continue
		}
		if moneyInfos == nil || len(moneyInfos) == 0 {
			ZapLog().Sugar().Debugf("GetHuiLvSina return nil data")
			goto SINAHUILV
		}

		if err = this.putSinaHuilvQoute(moneyInfos); err != nil {
			ZapLog().With(zap.Error(err), zap.String("err message", "put HuiLvSina err")).Error("put HuiLvSina err")
		}
		// 存到 redis kxain
		if err = this.mHistory.SinaHuilvSetRedis(moneyInfos); err != nil {
			ZapLog().With(zap.Error(err), zap.String("redis put kxian", "put HuiLvSina kxian err")).Error("History.Set err")
		}
	SINAHUILV:
		time.Sleep(time.Second * time.Duration(40))
		//fmt.Println("6")
		ZapLog().Sugar().Debugf("loadSinaHuilv huilvSuccessNum[%d]putRedisSuccessNum[%d]", len(moneyInfos), len(moneyInfos))
	}
}

//配置的数字货币，coinmerit K线数据     param:   交易所 ，Symbol ，对应的编码，交易对，Symbol对应的币
func (this *QuoteMgr) loadConfigCoinAndStore(cExchange, cSymbol, cNum, cPairsArr, cToArr []string) {
	defer PanicPrint()
	allNum := len(cNum)
	CoinSuccessNum := 0
	putRedisSuccessNum := 0
	for {
		startDayTime := GenDay(time.Now().Unix())
		for i := 0; i < len(cNum); i++ {
			mCoinmerit := 0
			mCoinmeritErrFlag := false
			redisErrFlag := false
			moneyInfo, endFlag, err := this.mCoinMarket.GetCoinQoute(cExchange[i], cPairsArr[i], "1min", "1")
			if err != nil {
				ZapLog().With(zap.Error(err), zap.String("cNum", cNum[i])).Error("mCoinmerit GetCoinQouteInfo err")
				mCoinmeritErrFlag = true
				continue
			}
			if endFlag {
				ZapLog().Sugar().Debugf("ConfigCoinInfo[%s]num[%d] getConfigCoinTicker end", cNum[i], mCoinmerit)
				break
			}
			//存储在redis
			if err = this.putCoinQoute(moneyInfo, cSymbol[i], cNum[i], cToArr[i]); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("cNum", cNum[i])).Error("putCoinQoute err")
			}
			//fmt.Println("putCoinQoute succ -- from coinmerit ",cSymbol[i])
			if err = this.mHistory.BtcexaCoinmeritSetRedis(moneyInfo, cToArr[i], cNum[i]); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("coin", cToArr[i])).Error("History.Set err")
			}
			time.Sleep(time.Second * time.Duration(15))
			if !redisErrFlag {
				putRedisSuccessNum++
			}
			if !mCoinmeritErrFlag {
				CoinSuccessNum++
			}
			ZapLog().Sugar().Debugf("loadConfigCoinAndStore all[%d]ConfigCoinSuccessNum[%d]putRedisSuccessNum[%d]", allNum, CoinSuccessNum, putRedisSuccessNum)
		}
		ZapLog().Sugar().Debugf("complete coinQoute once loop , sleep one hour")
		if GenDay(time.Now().Unix()) > startDayTime {
			time.Sleep(time.Second * time.Duration(60))
		} else {
			time.Sleep(time.Second * time.Duration(60))
		}
	}

}

//配置的数字货币，btcexa K线数据     param:   交易所 ，Symbol ，对应的编码，交易对，Symbol对应的币
func (this *QuoteMgr) loadConfigCoinAndStoreBtcexa(cSymbol, cNum, cPairsArr, cToArr []string) {
	defer PanicPrint()
	allNum := len(cNum)
	CoinSuccessNum := 0
	putRedisSuccessNum := 0
	for {
		startDayTime := GenDay(time.Now().Unix())
		for i := 0; i < len(cNum); i++ {
			mCoinmerit := 0
			mCoinmeritErrFlag := false
			redisErrFlag := false
			moneyInfo, endFlag, err := this.mCoinMarket.GetCoinQouteBtcexa(cPairsArr[i])
			if err != nil {
				ZapLog().With(zap.Error(err), zap.String("cNum", cNum[i])).Error("mBtcexa GetCoinQouteInfo err")
				mCoinmeritErrFlag = true
				continue
			}
			if endFlag {
				ZapLog().Sugar().Debugf("ConfigCoinInfo[%s]num[%d] getConfigCoinTicker end", cNum[i], mCoinmerit)
				break
			}
			//存储在redis
			if err = this.putCoinQoute(moneyInfo, cSymbol[i], cNum[i], cToArr[i]); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("cNum", cNum[i])).Error("putCoinQoute err")
			}
			//fmt.Println("putCoinQoute succ -- from btcexa ",cSymbol[i])
			if err = this.mHistory.BtcexaCoinmeritSetRedis(moneyInfo, cToArr[i], cNum[i]); err != nil {
				redisErrFlag = true
				ZapLog().With(zap.Error(err), zap.String("coin", cToArr[i])).Error("History.Set err")
			}
			time.Sleep(time.Second * time.Duration(15))
			if !redisErrFlag {
				putRedisSuccessNum++
			}
			if !mCoinmeritErrFlag {
				CoinSuccessNum++
			}
			ZapLog().Sugar().Debugf("loadConfigCoinAndStoreBtcexa all[%d]ConfigCoinSuccessNum[%d]putRedisSuccessNum[%d]", allNum, CoinSuccessNum, putRedisSuccessNum)
		}
		ZapLog().Sugar().Debugf("complete coinQoute once loop , sleep one hour")
		if GenDay(time.Now().Unix()) > startDayTime {
			time.Sleep(time.Second * time.Duration(60))
		} else {
			time.Sleep(time.Second * time.Duration(60))
		}
	}

}
