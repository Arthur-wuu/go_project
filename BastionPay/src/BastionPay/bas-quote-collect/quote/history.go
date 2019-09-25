package quote

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/collect"
	"BastionPay/bas-quote-collect/config"
	"BastionPay/bas-quote-collect/db"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type HistoryMgr struct {
	mKXian1Days map[string]*KXian //不需要锁，没有并发
	mRedis      *db.DbRedis
	mFlag       bool
}

//读取历史
func (this *HistoryMgr) Init(redis *db.DbRedis) (err error) {
	if redis == nil {
		return errors.New("nil redis in historyMgr")
	}
	this.mRedis = redis
	this.mKXian1Days, err = this.loadLastedKxian()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("loadLastedKxian err")
		return err
	}
	this.mFlag = true
	return nil
}

func (this *HistoryMgr) loadLastedKxian() (kXian1Days map[string]*KXian, err error) {
	kXian1Days = make(map[string]*KXian)

	allNum := 0
	succNum := 0
	var keys []string
	con := this.mRedis.GetConn()
	defer con.Close() //释放连接到连接池
	for cur := 0; true; {
		cur, keys, err = this.mRedis.Scan(cur, "MATCH", CONST_KXIAN_1Day_Prefix+"*", "COUNT", 50)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Int("cur", cur)).Error("Scan err")
			return nil, err
		}
		allNum += len(keys)
		for i := 0; i < len(keys); i++ {
			err = con.Send("LRANGE", keys[i], 0, 0)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.String("key", keys[i])).Error("LRANGE err")
				return nil, err
			}
		}
		err = con.Flush()
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("redis flush err")
			return nil, err
		}
		for i := 0; i < len(keys); i++ {
			repyArr, err := redis.ByteSlices(con.Receive())
			if err != nil {
				ZapLog().With(zap.Error(err), zap.String("key", keys[i])).Error("redis Receive err")
				return nil, err
			}
			if len(repyArr) == 0 {
				ZapLog().With(zap.String("key", keys[i])).Error("redis Receive nil")
				continue
			}
			kxian := new(KXian)
			err = json.Unmarshal(repyArr[0], kxian)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.String("key", keys[i]), zap.String("content", string(repyArr[0]))).Error("Unmarshal err")
				return nil, err
			}
			kXian1Days[keys[i]] = kxian
			succNum++
		}
		if cur == 0 {
			break
		}
	}
	ZapLog().With(zap.Int("allKeys", allNum), zap.Int("succKeys", succNum), zap.Int("mKXian1Days", len(kXian1Days))).Info("loadLastedKxian ok ")
	return
}

func (this *HistoryMgr) Set(tickerInfos *collect.TickerInfo, coin string) (gerr error) {
	con := this.mRedis.GetConn()
	defer con.Close()
	count := 0
	for i := 0; i < len(tickerInfos.IdDetailInfos); i++ {
		IdDetailInfo := tickerInfos.IdDetailInfos[i]
		nowDayTime := GenDay(IdDetailInfo.GetLast_updated())

		for my, moneyInfo := range IdDetailInfo.Quotes {
			if coin != my {
				continue
			}
			key := this.genKXian1DayKey(IdDetailInfo.Id, my)

			kxian, updateWay, updateFlag := this.updateKXian1Day(nowDayTime, key, &moneyInfo)
			if !updateFlag {
				continue
			}
			ZapLog().With(zap.String("key", key), zap.Bool("listPush", updateWay)).Debug("kxian set")
			content, err := json.Marshal(kxian)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("Marshal err")
				continue
			}
			num, err := this.batchPushKXian(con, updateWay, key, content)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("batchPushKXian err")
				return err
			}
			count += num
		}

	}
	if err := con.Flush(); err != nil {
		ZapLog().With(zap.Error(err)).Error("Flush err")
		return err
	}

	for i := 0; i < count; i++ {
		_, err := con.Receive()
		if err != nil {
			gerr = err
			ZapLog().With(zap.Error(err)).Error("Receive err")
			continue
		}
	}
	return gerr
}

//往redis中存 btcexa的moneyinfo， 做kxian数据
func (this *HistoryMgr) BtcexaCoinmeritSetRedis(moneyInfo *collect.MoneyInfo, coin, id string) (gerr error) {
	con := this.mRedis.GetConn()
	defer con.Close()
	count := 0
	nowDayTime := GenDay(moneyInfo.GetLast_updated())
	newid, err := strconv.Atoi(id)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("string to int  err")
	}
	key := this.genKXian1DayKey(newid, coin)
	kxian, updateWay, updateFlag := this.updateKXian1Day(nowDayTime, key, moneyInfo)
	if updateFlag {
		ZapLog().With(zap.String("key", key), zap.Bool("listPush", updateWay)).Debug("kxian set")
		content, err := json.Marshal(kxian)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("Marshal err")
		}
		num, err := this.batchPushKXian(con, updateWay, key, content)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("batchPushKXian err")
			return err
		}
		count += num
	}

	if err := con.Flush(); err != nil {
		ZapLog().With(zap.Error(err)).Error("Flush err")
		return err
	}

	for i := 0; i < count; i++ {
		_, err := con.Receive()
		if err != nil {
			gerr = err
			ZapLog().With(zap.Error(err)).Error("Receive err")
			continue
		}
	}
	return gerr
}

//往redis中存 coinmarketcap 新的 行情数据， 做kxian数据
func (this *HistoryMgr) CoinMarketCapSetRedisAll(res *collect.ResCoinMarketCapAll) (gerr error) {

	con := this.mRedis.GetConn()
	defer con.Close()
	count := 0

	for j := 0; j < len(res.Data); j++ {
		IdDetailInfo := res.Data[j]
		timestampInt := TimeToTimestamp(res.Status.Timestamp)
		nowDayTime := GenDay(timestampInt)

		for _, pInfo := range IdDetailInfo.Quote {

			moneyInfo := new(collect.MoneyInfo)
			id := res.Data[j].Id
			price := pInfo.Price

			moneyInfo.Price = &price
			moneyInfo.Last_updated = &timestampInt

			key := this.genKXian1DayKeyCoinMCap(id)

			kxian, updateWay, updateFlag := this.updateKXian1Day(nowDayTime, key, moneyInfo)
			if !updateFlag {
				continue
			}
			ZapLog().With(zap.String("key", key), zap.Bool("listPush", updateWay)).Debug("kxian set")
			content, err := json.Marshal(kxian)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("Marshal err")
				continue
			}
			num, err := this.batchPushKXian(con, updateWay, key, content)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("batchPushKXian err")
				return err
			}
			count += num
		}
	}
	//}
	if err := con.Flush(); err != nil {
		ZapLog().With(zap.Error(err)).Error("Flush err")
		return err
	}

	for i := 0; i < count; i++ {
		_, err := con.Receive()
		if err != nil {
			gerr = err
			ZapLog().With(zap.Error(err)).Error("Receive err")
			continue
		}
	}
	return gerr
}

//往redis中存 coinmarketcap 新的 码表中的有的币
func (this *HistoryMgr) CoinMarketCapSetRedisPart(res *collect.ResCoinMarketCapPart, codes string) (gerr error) {

	con := this.mRedis.GetConn()
	defer con.Close()
	count := 0
	codeArr := strings.Split(codes, ",")
	codeIntArr, _ := StringArrToIntArr(codeArr)
	codeIntArr = BubbleSort(codeIntArr)
	if len(codeIntArr) != len(res.Data) {
		ZapLog().With(zap.Any("err message", "len(codeArr) != len(res.Data)")).Error("data err")
	}

	for j := 0; j < len(res.Data); j++ {
		IdDetailInfo := res.Data[strconv.Itoa(codeIntArr[j])]
		timestampInt := TimeToTimestamp(res.Status.Timestamp)
		nowDayTime := GenDay(timestampInt)

		for _, pInfo := range IdDetailInfo.Quote {

			moneyInfo := new(collect.MoneyInfo)
			id := res.Data[strconv.Itoa(codeIntArr[j])].Id
			price := pInfo.Price

			moneyInfo.Price = &price
			moneyInfo.Last_updated = &timestampInt

			key := this.genKXian1DayKeyCoinMCap(id)

			kxian, updateWay, updateFlag := this.updateKXian1Day(nowDayTime, key, moneyInfo)
			if !updateFlag {
				continue
			}
			ZapLog().With(zap.String("key", key), zap.Bool("listPush", updateWay)).Debug("kxian set")
			content, err := json.Marshal(kxian)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("Marshal err")
				continue
			}
			num, err := this.batchPushKXian(con, updateWay, key, content)
			if err != nil {
				ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("batchPushKXian err")
				return err
			}
			count += num
		}
	}
	//}
	if err := con.Flush(); err != nil {
		ZapLog().With(zap.Error(err)).Error("Flush err")
		return err
	}

	for i := 0; i < count; i++ {
		_, err := con.Receive()
		if err != nil {
			gerr = err
			ZapLog().With(zap.Error(err)).Error("Receive err")
			continue
		}
	}
	return gerr
}

//往redis中存  baidu  huilv的moneyinfo， 做kxian数据
func (this *HistoryMgr) HuilvSetRedis(moneyInfo *collect.MoneyInfo, coin string) (gerr error) {
	con := this.mRedis.GetConn()
	defer con.Close()
	count := 0
	nowDayTime := GenDay(moneyInfo.GetLast_updated())

	//fmt.Println("HuilvSetRedis 1")
	key := this.genKXian1DayKeyHuilv(coin)
	kxian, updateWay, updateFlag := this.updateKXian1Day(nowDayTime, key, moneyInfo)
	//fmt.Println("HuilvSetRedis 2",kxian, updateWay, updateFlag)
	if updateFlag {
		ZapLog().With(zap.String("key", key), zap.Bool("listPush", updateWay)).Debug("kxian set")
		content, err := json.Marshal(kxian)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("Marshal err")
		}
		num, err := this.batchPushKXian(con, updateWay, key, content)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("batchPushKXian err")
			return err
		}
		count += num
	}

	if err := con.Flush(); err != nil {
		ZapLog().With(zap.Error(err)).Error("Flush err")
		return err
	}

	for i := 0; i < count; i++ {
		_, err := con.Receive()
		if err != nil {
			gerr = err
			ZapLog().With(zap.Error(err)).Error("Receive err")
			continue
		}
	}
	return gerr
}

//往redis中存 sina huilv， 做kxian数据
func (this *HistoryMgr) SinaHuilvSetRedis(monInfos []*collect.MoneyInfo) (gerr error) {
	if monInfos == nil || len(monInfos) == 0 {
		return
	}
	con := this.mRedis.GetConn()
	defer con.Close()
	count := 0

	for i := 0; i < len(monInfos); i++ {
		nowDayTime := GenDay(monInfos[i].GetLast_updated())
		key := this.genKXian1DayKeyHuilv(config.GPreConfig.CountryCodeArr[i])

		kxian, updateWay, updateFlag := this.updateKXian1Day(nowDayTime, key, monInfos[i])
		if !updateFlag {
			continue
		}
		ZapLog().With(zap.String("key", key), zap.Bool("listPush", updateWay)).Debug("kxian set")
		content, err := json.Marshal(kxian)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("Marshal err")
			continue
		}
		num, err := this.batchPushKXian(con, updateWay, key, content)
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("key", key)).Error("batchPushKXian err")
			return err
		}
		count += num
	}

	if err := con.Flush(); err != nil {
		ZapLog().With(zap.Error(err)).Error("Flush err")
		return err
	}

	for i := 0; i < count; i++ {
		_, err := con.Receive()
		if err != nil {
			gerr = err
			ZapLog().With(zap.Error(err)).Error("Receive err")
			continue
		}
	}
	return gerr
}

func (this *HistoryMgr) batchPushKXian(con redis.Conn, updateWay bool, key string, content []byte) (int, error) {
	if updateWay {
		err := con.Send("LPUSH", key, content)
		if err != nil {
			return 0, err
		}
		err = con.Send("LTRIM", key, 0, 8)
		if err != nil {
			return 1, err
		}

		return 2, nil
	} else {
		err := con.Send("LSET", key, 0, content)
		if err != nil {
			return 0, err
		}
		return 1, err
	}
	return 0, nil
}

//从左往右
func (this *HistoryMgr) GetKXian1Days(id int, coin string, start, count int) ([]*KXian, error) {
	//LRANGE key start stop
	key := this.genKXian1DayKey(id, coin)
	bytesArr, err := redis.ByteSlices(this.mRedis.Do("LRANGE", key, start, count))
	if err != nil {
		return nil, err
	}
	arr := make([]*KXian, 0)
	for i := 0; i < len(bytesArr); i++ {
		kxian := new(KXian)
		err := json.Unmarshal(bytesArr[i], kxian)
		if err != nil {
			return nil, err
		}
		arr = append(arr, kxian)
	}
	return arr, nil
}

//返回值，最新的kxian，是否是新日期，是否需要更新
func (this *HistoryMgr) updateKXian1Day(nowDayTime int64, key string, moneyInfo *collect.MoneyInfo) (*KXian, bool, bool) {
	price := moneyInfo.GetPrice()
	newFlag := false
	kxian, ok := this.mKXian1Days[key]
	if !ok {
		newFlag = true
	} else {
		if nowDayTime < kxian.GetTimestamp() {
			return nil, false, false
		}
		if nowDayTime > kxian.GetTimestamp() {
			newFlag = true
		}
	}

	if newFlag {
		kxian = new(KXian)
		this.mKXian1Days[key] = kxian
		kxian.SetClosePrice(price)
		kxian.SetLastPrice(price)
		kxian.SetHighPrice(price)
		kxian.SetLowPrice(price)
		kxian.SetOpenPrice(price)
		kxian.SetTimestamp(nowDayTime)
	} else {
		kxian.SetClosePrice(price)
		kxian.SetLastPrice(price)
		if price > kxian.GetHighPrice() {
			kxian.SetHighPrice(price)
		}
		if price < kxian.GetLowPrice() {
			kxian.SetLowPrice(price)
		}
	}
	return kxian, newFlag, true
}

func (this *HistoryMgr) genKXian1DayKey(id int, coin string) string {
	return fmt.Sprintf("%s_%d_%s", CONST_KXIAN_1Day_Prefix, id, coin)
}

//huilv  在redis中的key
func (this *HistoryMgr) genKXian1DayKeyHuilv(coin string) string {
	return fmt.Sprintf("%s_%s_%s", CONST_KXIAN_1Day_Prefix, "USD", coin)
}

func (this *HistoryMgr) genKXian1DayKeyCoinMCap(id int64) string {
	return fmt.Sprintf("%s_%d_%s", CONST_KXIAN_1Day_Prefix, id, "USD")
}

func GenDay(t int64) int64 {
	t1 := time.Unix(t, 0).Year()  //年
	t2 := time.Unix(t, 0).Month() //月
	t3 := time.Unix(t, 0).Day()   //日
	lc, _ := time.LoadLocation("UTC")
	currentTimeData := time.Date(t1, t2, t3, 0, 0, 0, 0, lc)
	return currentTimeData.Unix()
}
