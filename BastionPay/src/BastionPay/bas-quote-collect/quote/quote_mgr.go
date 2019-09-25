package quote

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/collect"
	"BastionPay/bas-quote-collect/config"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
	//"os"
	//"sort"
	"BastionPay/bas-quote-collect/db"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type QuoteMgr struct {
	mHuilv      collect.HuiLv
	mCoinMarket collect.CoinMarket
	//mLvDb       *leveldb.DB
	mRedis   *db.DbRedis
	mSqlDb   db.DbMgr
	mHistory HistoryMgr
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
	if err := this.mCoinMarket.Init(this); err != nil {
		return err
	}
	if err := this.mHuilv.Init(); err != nil {
		return err
	}
	this.mCodeTable = make(map[string]*collect.CodeInfo, 0)

	this.mRedis = &db.GRedis

	err = this.mRedis.Init(config.GConfig.Redis.Host,
		config.GConfig.Redis.Port, config.GConfig.Redis.Password,
		config.GConfig.Redis.Database)
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
	ZapLog().Sugar().Infof("Start load codeTable from SqlDb")
	for i := 0; i < len(arr); i++ {
		v := CodeTableToCodeInfo(&arr[i])
		this.mCodeTable[arr[i].Symbol] = v
		ZapLog().Sugar().Infof("Symbol[%v] Info[%v]", v.GetSymbol(), v.ToPrintStr())
	}
	ZapLog().Sugar().Infof("End load codeTable from SqlDb")
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
	this.mRedis.Close()
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

func (this *QuoteMgr) GetKxian1Day(id int, to string, start, limit int) ([]*KXian, error) {
	return this.mHistory.GetKXian1Days(id, to, start, limit)
}

//
func (this *QuoteMgr) GetKxian1DayUseUSD(to string, start, limit int) ([]*KXian, error) {
	key := this.genKXian1DayKeyUSD(to)
	bytesArr, err := redis.ByteSlices(db.GRedis.Do("LRANGE", key, start, limit))
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

func (this *QuoteMgr) GetQuoteUseId(id int, to string) (*collect.MoneyInfo, error) {
	key := this.genQuoteKey(id, to)
	ZapLog().Debug("GetQuoteUseId", zap.String("leveldb_key", string(key)))

	data, err := this.mRedis.Do("GET", key)
	//data, err := this.mLvDb.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("nil data")
	}
	byteData, ok := data.([]byte)
	if !ok {
		//ZapLog().Error("wrong type")
		return nil, errors.New("wrong type")
	}
	moneyInfo := new(collect.MoneyInfo)
	err = json.Unmarshal(byteData, moneyInfo)
	if err != nil {
		return nil, err
	}
	moneyInfo.SetSymbol(to)
	//fmt.Println(data)
	return moneyInfo, nil
}

//GetQuoteHuilv   USD==>法币
func (this *QuoteMgr) GetQuoteHuilv(to string) (*collect.MoneyInfo, error) {
	to = fmt.Sprintf("qt_USD_%s", to)

	if to == "qt_USD_USD" {
		moneyInfo := new(collect.MoneyInfo)
		moneyInfo.SetPrice(1)
		moneyInfo.SetSymbol(to)
		return moneyInfo, nil
	}

	data, err := this.mRedis.Do("GET", to)
	//fmt.Println("data : " ,data)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("nil data")
	}
	byteData, ok := data.([]byte)
	if !ok {
		//ZapLog().Error("wrong type")
		return nil, errors.New("wrong type")
	}
	moneyInfo := new(collect.MoneyInfo)
	err = json.Unmarshal(byteData, moneyInfo)
	//fmt.Println("moneyInfo : " ,moneyInfo)
	if err != nil {
		return nil, err
	}
	moneyInfo.SetSymbol(to)
	//fmt.Println(data)
	return moneyInfo, nil
}

func (this *QuoteMgr) GetQuoteUseSymbol(from, to string) (*collect.CodeInfo, *collect.MoneyInfo, error) {
	this.Lock()
	codeInfo, ok := this.mCodeTable[from]
	this.Unlock()
	if !ok {
		return nil, nil, errors.New("nofind codeTable")
	}
	key := this.genQuoteKey(codeInfo.GetId(), to)
	data, err := this.mRedis.Do("GET", key)
	if err != nil {
		return nil, nil, err
	}
	if data == nil {
		return nil, nil, errors.New("nil data")
	}
	byteData, ok := data.([]byte)
	if !ok {
		//ZapLog().Error("wrong type")
		return nil, nil, errors.New("wrong type")
	}

	moneyInfo := new(collect.MoneyInfo)
	err = json.Unmarshal(byteData, moneyInfo)
	if err != nil {
		return nil, nil, err
	}
	moneyInfo.SetSymbol(to)
	//fmt.Println(data)
	return codeInfo, moneyInfo, nil
}

/***********************内部接口分割线*************************/
func (this *QuoteMgr) run() {
	this.Add(1)
	defer this.Done()

	coinsArr := strings.Split(config.GConfig.Coinmarketcap.Coins, ",")
	ZapLog().Sugar().Infof("Config coins[%v]", coinsArr)

	cNameArr := strings.Split(config.GConfig.Parities.Country_name, ",")
	ZapLog().Sugar().Infof("Config cName[%v]", cNameArr)

	cCodeArr := strings.Split(config.GConfig.Parities.Country_code, ",")
	ZapLog().Sugar().Infof("Config cCode[%v]", cCodeArr)

	cNumArr := strings.Split(config.GConfig.Collect.Coin_num, ",")
	ZapLog().Sugar().Infof("Config cNum[%v]", cNumArr)

	cSymbolArr := strings.Split(config.GConfig.Collect.Coin_from, ",")
	ZapLog().Sugar().Infof("Config cSymbol[%v]", cSymbolArr)

	cToArr := strings.Split(config.GConfig.Collect.Coin_to, ",")
	ZapLog().Sugar().Infof("Config cTo[%v]", cToArr)

	cPairsArr := strings.Split(config.GConfig.Collect.Coin_pairs, ",")
	ZapLog().Sugar().Infof("Config cPairsArr[%v]", cPairsArr)

	cEntrance := strings.Split(config.GConfig.Collect.Coin_entrance, ",")
	ZapLog().Sugar().Infof("Config cEntrance[%v]", cEntrance)

	cExchangeArr := strings.Split(config.GConfig.Collect.Coin_exchange, ",")
	ZapLog().Sugar().Infof("Config cExchange[%v]", cExchangeArr)

	cOpenFlagArr := strings.Split(config.GConfig.Collect.Open_Flag, ",")
	ZapLog().Sugar().Infof("Config cOpenFlag[%v]", cOpenFlagArr)

	//需要根据配置文件里的交易所入口来区分走下面哪个go程
	cExchangeArrs, cSymbolArrs, cNumArrs, cPairsArrs, cToArrs := this.diffCoinmeritExchange(cEntrance, cExchangeArr, cSymbolArr, cNumArr, cPairsArr, cToArr, cOpenFlagArr)
	bSymbolArr, bNumArr, bPairsArr, bToArr := this.diffBtcexaExchange(cEntrance, cExchangeArr, cSymbolArr, cNumArr, cPairsArr, cToArr, cOpenFlagArr)

	//从配置文件中加载币，从coinmerit拿数据
	go this.loadConfigCoinAndStore(cExchangeArrs, cSymbolArrs, cNumArrs, cPairsArrs, cToArrs)

	//从配置文件中加载币，从btcexa拿数据, coinmerit支持btcexa交易所后，直接用上面的
	go this.loadConfigCoinAndStoreBtcexa(bSymbolArr, bNumArr, bPairsArr, bToArr)

	//coinmarketcap 新接口，也是全量拉取 ,需要付费api-key， 暂时不用
	//go this.loadCoinTicker()

	//coinmarketcap 新接口，拉取数据库有的币  可配置免费的apikey
	if config.GConfig.Coinmarketcap.NewApiFlag {
		ZapLog().Info("start loadPartCoinTicker, newApi and part coins")
		go this.loadPartCoinTicker()
	} else {
		//启动拉取全量
		ZapLog().Info("start loadQuotesAndStore, oldApi and all coins")
		go this.loadQuotesAndStore(coinsArr)
	}

	if config.GConfig.Switch.FxBaiduFlag {
		//获取USD=>其它法币的汇率  baidu
		ZapLog().Info("start loadHuilvBaidu")
		go this.loadHuilvBaidu(cNameArr, cCodeArr)
	}

	if config.GConfig.Switch.FxSinaFlag {
		//sina 汇率
		ZapLog().Info("start loadHuilvSina")
		go this.loadHuilvSina()
	}
}

func (this *QuoteMgr) GenCodeTable() (map[string]*collect.CodeInfo, error) {
	arr, err := this.mSqlDb.GetAllCode()
	if err != nil {
		ZapLog().Sugar().With(zap.Error(err)).Error("load codeTable from SqlDb err")
		return nil, err
	}
	codeTable := make(map[string]*collect.CodeInfo, 0)

	ZapLog().Sugar().Infof("Start load codeTable from SqlDb")
	for i := 0; i < len(arr); i++ {
		v := CodeTableToCodeInfo(&arr[i])
		codeTable[arr[i].Symbol] = v
		ZapLog().Sugar().Infof("Symbol[%v] Info[%v]", v.GetSymbol(), v.ToPrintStr())
	}
	ZapLog().Sugar().Infof("End load codeTable from SqlDb")
	fmt.Println("codetabkle**********", &codeTable)
	return codeTable, nil
}

func (this *QuoteMgr) genQuoteKey(id int, money string) []byte {
	return []byte(fmt.Sprintf("%s_%d_%s", DB_Quote_prefix, id, money))
}

func (this *QuoteMgr) genHuilvKey(money string) []byte {
	return []byte(fmt.Sprintf("%s_USD_%s", DB_Quote_prefix, money))
}

func (this *QuoteMgr) genCoinQouteKey(num, to string) []byte {
	return []byte(fmt.Sprintf("%s_%s_%s", DB_Quote_prefix, num, to))
}

func (this *QuoteMgr) genCoinMarketCapQouteKey(num int64) []byte {
	return []byte(fmt.Sprintf("%s_%d_USD", DB_Quote_prefix, num))
}

//func (this *QuoteMgr) getConfigCoins() []string {
//	moneys := strings.TrimSpace(config.GConfig.Coinmarketcap.Coins)
//	moneys = strings.ToUpper(moneys)
//	moneys = strings.TrimRight(moneys, ",")
//	return strings.Split(moneys, ",")
//}
//
//func (this *QuoteMgr) getConfigCountryName() []string {
//	cName := strings.TrimSpace(config.GConfig.Parities.Country_name)
//	cName = strings.ToUpper(cName)
//	cName = strings.TrimRight(cName, ",")
//	return strings.Split(cName, ",")
//}
//
//func (this *QuoteMgr) getConfigCountryCode() []string {
//	cCode := strings.TrimSpace(config.GConfig.Parities.Country_code)
//	cCode = strings.ToUpper(cCode)
//	cCode = strings.TrimRight(cCode, ",")
//	return strings.Split(cCode, ",")
//}

////配置文件内的参数
//func (this *QuoteMgr) getConfigCoinNum() []string {
//	cNum := strings.TrimSpace(config.GConfig.Collect.Coin_num)
//	cNum = strings.TrimRight(cNum, ",")
//	return strings.Split(cNum, ",")
//}
//
//func (this *QuoteMgr) getConfigCoinFrom() []string {
//	cFrom := strings.TrimSpace(config.GConfig.Collect.Coin_from)
//	cFrom = strings.ToUpper(cFrom)
//	cFrom = strings.TrimRight(cFrom, ",")
//	return strings.Split(cFrom, ",")
//}
//
//func (this *QuoteMgr) getConfigCoinExchange() []string {
//	cExchange := strings.TrimSpace(config.GConfig.Collect.Coin_exchange)
//	cExchange = strings.TrimRight(cExchange, ",")
//	return strings.Split(cExchange, ",")
//}
//
//func (this *QuoteMgr) getConfigCoinPairs() []string {
//	cPairs := strings.TrimSpace(config.GConfig.Collect.Coin_pairs)
//	cPairs = strings.TrimRight(cPairs, ",")
//	return strings.Split(cPairs, ",")
//}
//
//func (this *QuoteMgr) getConfigCoinTo() []string {
//	cTo := strings.TrimSpace(config.GConfig.Collect.Coin_to)
//	cTo = strings.ToUpper(cTo)
//	cTo = strings.TrimRight(cTo, ",")
//	return strings.Split(cTo, ",")
//}
//
//func (this *QuoteMgr) getConfigCoinEntrance() []string {
//	cEntrance := strings.TrimSpace(config.GConfig.Collect.Coin_entrance)
//	cEntrance = strings.TrimRight(cEntrance, ",")
//	return strings.Split(cEntrance, ",")
//}

func (this *QuoteMgr) diffCoinmeritExchange(cEntranceArr, cExchangeArr, cSymbolArr, cNumArr, cPairsArr, cToArr, cOpenFlagArr []string) ([]string, []string, []string, []string, []string) {
	numCoinmerit := make([]int, 0)
	for index, entrance := range cEntranceArr {
		if cOpenFlagArr[index] != "1" {
			continue
		}
		if entrance == "coinmerit" {
			numCoinmerit = append(numCoinmerit, index)
		}
	}
	tmpExchangeArr := make([]string, 0)
	tmpSymbolArr := make([]string, 0)
	tmpNumArr := make([]string, 0)
	tmpPairsArr := make([]string, 0)
	tmpToArr := make([]string, 0)
	for i := 0; i < len(numCoinmerit); i++ {
		tmpExchangeArr = append(tmpExchangeArr, cExchangeArr[numCoinmerit[i]])
		tmpSymbolArr = append(tmpSymbolArr, cSymbolArr[numCoinmerit[i]])
		tmpNumArr = append(tmpNumArr, cNumArr[numCoinmerit[i]])
		tmpPairsArr = append(tmpPairsArr, cPairsArr[numCoinmerit[i]])
		tmpToArr = append(tmpToArr, cToArr[numCoinmerit[i]])
	}

	return tmpExchangeArr, tmpSymbolArr, tmpNumArr, tmpPairsArr, tmpToArr
}

func (this *QuoteMgr) diffBtcexaExchange(cEntranceArr, cExchangeArr, cSymbolArr, cNumArr, cPairsArr, cToArr, cOpenFlagArr []string) ([]string, []string, []string, []string) {
	numCoinmerit := make([]int, 0)
	for index, entrance := range cEntranceArr {
		if cOpenFlagArr[index] != "1" {
			continue
		}
		if entrance == "btcexa" {
			numCoinmerit = append(numCoinmerit, index)
		}
	}
	//fmt.Println("len of num",numCoinmerit)
	tmpSymbolArr := make([]string, 0)
	tmpNumArr := make([]string, 0)
	tmpPairsArr := make([]string, 0)
	tmpToArr := make([]string, 0)
	for i := 0; i < len(numCoinmerit); i++ {
		tmpNumArr = append(tmpNumArr, cNumArr[numCoinmerit[i]])
		tmpSymbolArr = append(tmpSymbolArr, cSymbolArr[numCoinmerit[i]])
		tmpPairsArr = append(tmpPairsArr, cPairsArr[numCoinmerit[i]])
		tmpToArr = append(tmpToArr, cToArr[numCoinmerit[i]])
	}
	return tmpSymbolArr, tmpNumArr, tmpPairsArr, tmpToArr
}

func (this *QuoteMgr) genKXian1DayKeyUSD(coin string) string {
	return fmt.Sprintf("%s_%s_%s", CONST_KXIAN_1Day_Prefix, "USD", coin)
}
