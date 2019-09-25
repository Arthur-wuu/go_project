package quote

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote/collect"
	"BastionPay/bas-quote/config"
	"encoding/json"
	"fmt"
	//"github.com/bluele/gcache"
	"go.uber.org/zap"
	//"sort"
	"BastionPay/bas-quote/db"
	"errors"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

type QuoteMgr struct {
	mHistory    HistoryMgr
	mRedis      *db.Redis
	mSqlDb      db.DbMgr
	//	mRptSymbols []collect.CodeInfo
	//mUniSymbols []collect.CodeInfo
	mCodeTable map[string]*collect.CodeInfo
	mExitCh    chan bool
	mRunFlag   bool
	sync.WaitGroup
	sync.Mutex
}

func (this *QuoteMgr) Init() (err error) {
	if this.mRunFlag {
		return nil
	}
	this.mExitCh = make(chan bool)
	this.mCodeTable = make(map[string]*collect.CodeInfo, 0)
	//	this.mRptSymbols = make([]collect.CodeInfo, 0)
	//makedir
	err = db.GRedis.Init(config.GConfig.Redis.Host, config.GConfig.Redis.Port, config.GConfig.Redis.Password, config.GConfig.Redis.Database)
	if err != nil {
		return err
	}
	err = this.mSqlDb.Init(&db.DbOptions{
		Host:        config.GConfig.Db.Host,
		Port:        config.GConfig.Db.Port,
		User:        config.GConfig.Db.User,
		Pass:        config.GConfig.Db.Pwd,
		DbName:      config.GConfig.Db.Quote_db,
		MaxIdleConn: config.GConfig.Db.Max_idle_conn,
		MaxOpenConn: config.GConfig.Db.Max_open_conn,
	})
	if err != nil {
		return err
	}
    //init
	db.GCache.SetQuoteCacheFunc(GetCacheQuote)
	db.GCache.SetHuilvCacheFunc(GetCacheHuilv)

	this.mRedis = &db.GRedis
	if err = this.mHistory.Init(this.mRedis); err != nil {
		return err
	}
	this.mRunFlag = true
	return nil
}

func (this *QuoteMgr) Start() error {
	arr, err := this.mSqlDb.GetAllCode()
	if err != nil {
		return err
	}
	if this.mCodeTable == nil {
		this.mCodeTable = make(map[string]*collect.CodeInfo, 0)
	}
	ZapLog().Sugar().Infof("Start load codeTable from SqlDb count")
	for i := 0; i < len(arr); i++ {
		v := CodeTableToCodeInfo(&arr[i])
		this.mCodeTable[arr[i].Symbol] = v
		ZapLog().Sugar().Infof("Symbol[%v] Info[%v]", v.GetSymbol(), v.ToPrintStr())
	}
	ZapLog().Sugar().Infof("End load codeTable from SqlDb count[%d]",  len(this.mCodeTable))
	go this.run()
	return nil
}

func (this *QuoteMgr) Stop() error {
	select {
	case <-this.mExitCh:
		return nil
	default:
	}

	close(this.mExitCh)
	this.Wait()
	db.GRedis.Close()
	this.mSqlDb.Close()
	return nil
}

func (this *QuoteMgr) ListSymbols() []collect.CodeInfo {
	this.Lock()
	mm := this.mCodeTable
	this.Unlock()
	arr := make([]collect.CodeInfo, 0)
	for _, v := range mm {
		arr = append(arr, *v)
	}
	return arr
}

func (this *QuoteMgr) GetSymbol(symbol string) *collect.CodeInfo {
	this.Lock()
	defer this.Unlock()
	return this.mCodeTable[symbol]
}

func (this *QuoteMgr) SetCodeTable(cc *collect.CodeInfo) error {
	this.Lock()
	this.mCodeTable[cc.GetSymbol()] = cc
	this.Unlock()
	tt := CodeInfoToCodeTable(cc)
	return this.mSqlDb.AddCode(tt)
}

//行情的回调
func GetCacheQuote(key interface{}) (interface{}, *time.Duration, error) {
	expire := config.GConfig.Cache.Timeout * time.Second
	data ,err:= redis.Bytes(db.GRedis.Do("GET",key))
	if err != nil {
		return  nil, nil, err
	}
	if data == nil {
		return nil, nil, errors.New("nil data")
	}
	return data, &expire, nil
}

func (this *QuoteMgr) GetQuoteUseId(id int, to string) ( *collect.MoneyInfo, error) {

	key := this.genQuoteKey(id, to)
	ZapLog().Debug("GetQuoteUseId",zap.String("key", string(key)))

	value, err := db.GCache.QuoteCache.Get(string(key))
	if err != nil {
		return nil, err
	}
	if value != nil {
		moneyInfo := new(collect.MoneyInfo)
		err = json.Unmarshal(value.([]byte), moneyInfo)
		if err != nil {
			return  nil, err
		}
		moneyInfo.SetSymbol(to)
		return  moneyInfo, nil
	}
	return  nil, nil
}

//汇率的回调
func GetCacheHuilv (key interface{}) (interface{}, *time.Duration, error) {
	expire := config.GConfig.Cache.Timeout * time.Second
	data ,err:= redis.Bytes(db.GRedis.Do("GET",key))
	if err != nil {
		return  nil, nil, err
	}
	if data == nil {
		return nil, nil, errors.New("nil data")
	}
	return data, &expire, nil
}

//GetQuoteHuilv   USD==>法币
func (this *QuoteMgr) GetQuoteHuilv( to string) ( *collect.MoneyInfo, error) {
	to = fmt.Sprintf("qt_USD_%s", to)
	if to == "qt_USD_USD" {
		moneyInfo := new(collect.MoneyInfo)
		moneyInfo.SetPrice(1)
		moneyInfo.SetSymbol(to)
		return moneyInfo, nil
	}
	value, err := db.GCache.HuilvCache.Get(string(to))
	if err != nil {
		return nil, err
	}
	if value != nil {
		moneyInfo := new(collect.MoneyInfo)
		err = json.Unmarshal(value.([]byte), moneyInfo)
		if err != nil {
			return nil, err
		}
		moneyInfo.SetSymbol(to)
		return moneyInfo, nil
	}

	return nil, nil
}
func (this *QuoteMgr) GetKxian1Day(id int, to string, start, limit int) ([]*KXian, error) {
	key := this.genKXian1DayKey(id, to)
	bytesArr, err := redis.ByteSlices(db.GRedis.Do("LRANGE", key, start, limit))
	if err != nil {
		return nil, err
	}
	arr := make([]*KXian, 0)
	for i:=0; i < len(bytesArr); i++ {
		kxian:= new(KXian)
		err := json.Unmarshal(bytesArr[i], kxian)
		if err != nil {
			return nil, err
		}
		arr = append(arr, kxian)
	}
	return arr,nil

}
//得到一天的kxian
func (this *QuoteMgr) GetKxian1DayUseUSD( to string, start, limit int) ([]*KXian, error) {
	key := this.genKXian1DayKeyUSD(to)
	bytesArr, err := redis.ByteSlices(db.GRedis.Do("LRANGE", key, start, limit))
	if err != nil {
		return nil, err
	}
	arr := make([]*KXian, 0)
	for i:=0; i < len(bytesArr); i++ {
		kxian:= new(KXian)
		err := json.Unmarshal(bytesArr[i], kxian)
		if err != nil {
			return nil, err
		}
		arr = append(arr, kxian)
	}
	return arr,nil

}

func (this *QuoteMgr) GetQuoteUseSymbol(from, to string) (*collect.CodeInfo, *collect.MoneyInfo, error) {
	this.Lock()
	codeInfo, ok := this.mCodeTable[from]
	this.Unlock()
	if !ok {
		return nil, nil, errors.New("nofind codeTable")
	}
	key := this.genQuoteKey(codeInfo.GetId(), to)
	data ,err:= redis.Bytes(db.GRedis.Do("GET",key))
	if err != nil {
		return nil, nil, err
	}
	if data == nil {
		return nil,nil, errors.New("nil data")
	}
	moneyInfo := new(collect.MoneyInfo)
	err = json.Unmarshal(data, moneyInfo)
	if err != nil {
		return nil, nil, err
	}
	moneyInfo.SetSymbol(to)
	return codeInfo, moneyInfo, nil
}

/***********************内部接口分割线*************************/
func (this *QuoteMgr) run() {
	this.Add(1)
	defer this.Done()
	codeTicker := time.NewTicker(time.Second * time.Duration(3600+5))

	for true {
		select {
		case <-this.mExitCh:
			ZapLog().Sugar().Infof("QuoteMgr exist ok")
			break
		case <-codeTicker.C:
			arr, err := this.genCodeTable()
			if err != nil {
				break
			}
			this.Lock()
			this.mCodeTable = arr
			this.Unlock()
			break
		}

	}
}

func (this*QuoteMgr) genCodeTable() (map[string]*collect.CodeInfo, error){
	arr, err := this.mSqlDb.GetAllCode()
	if err != nil {
		ZapLog().Sugar().With(zap.Error(err)).Error("load codeTable from SqlDb err")
		return nil,err
	}
	codeTable := make(map[string]*collect.CodeInfo, 0)

	ZapLog().Sugar().Infof("Start load codeTable from SqlDb")
	for i := 0; i < len(arr); i++ {
		v := CodeTableToCodeInfo(&arr[i])
		codeTable[arr[i].Symbol] = v
		ZapLog().Sugar().Debugf("Symbol[%v] Info[%v]", v.GetSymbol(), v.ToPrintStr())
	}
	ZapLog().Sugar().Infof("End load codeTable from SqlDb", zap.Int("count", len(this.mCodeTable)))
	return codeTable,nil
}

func (this *QuoteMgr) genQuoteKey(id int, money string) []byte {
	return []byte(fmt.Sprintf("%s_%d_%s", DB_Quote_prefix, id, money))

}

func (this *QuoteMgr) genKXian1DayKey(id int, coin string) string {
	return fmt.Sprintf("%s_%d_%s", CONST_KXIAN_1Day_Prefix, id, coin)
}

func (this *QuoteMgr) genKXian1DayKeyUSD( coin string) string {
	return fmt.Sprintf("%s_%s_%s", CONST_KXIAN_1Day_Prefix, "USD", coin)
}