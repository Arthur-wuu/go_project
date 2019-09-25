package models


import (
"sync"
"BastionPay/bas-tv-proxy/common"
"time"
"BastionPay/bas-tv-proxy/models/kbspirit"
"BastionPay/bas-tv-proxy/api"

	"BastionPay/bas-tv-proxy/config"
	"BastionPay/bas-tv-proxy/type"
	"strings"
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
)

const(
	NAME_OBJ_SEARCHER = "obj_search"
)

var GKBSpirit KBSpirit

type KBSpirit struct {
	objSearcher *kbspirit.ObjSearcher
	sync.Mutex
}

func (this *KBSpirit) Init() {
	this.objSearcher = kbspirit.NewObjSearcher()
	this.objSearcher.Init()
}

func (this *KBSpirit) Start() {
	go this.run()
}

func (this *KBSpirit) run() {
	defer common.PanicPrint()
	time.Sleep(time.Second * 5)
	for true {
		this.loadMarkets()
		ZapLog().Info("loadMarkets ok")
		//启动时爬取交易对信息，全量爬取，24小时爬一次。存在问题是 启动时万一后端网站出问题了，那基本就没戏了，行情也没有了，要不要存下最新的一笔数据。
		time.Sleep(time.Hour * 6) //爬的频率要低，免得主接口被限制了。
	}
}

func (this * KBSpirit) GetObjs(env *kbspirit.Env) []*api.JPBShuChu {
	this.Lock()
	searcher := this.objSearcher
	this.Unlock()
	if searcher == nil {
		return nil
	}
	data := searcher.Search(env.Input, env)
	return data
}

// removeDuplicate 数据去重
func (self *KBSpirit) removeDuplicate(data []*api.JPBShuJu) []*api.JPBShuJu {
	found := make(map[string]bool)
	j := 0
	for i, d := range data {
		if !found[d.GetDaiMa()] {
			found[d.GetDaiMa()] = true
			data[j] = data[i]
			j++
		}
	}
	return data[:j]
}

func (this *KBSpirit) loadMarkets() {
	searcher := kbspirit.NewObjSearcher()
	searcher.Init()
	markets := config.GPreConfig.MarketMap
	for _, info := range markets {
		if  strings.ToUpper(info.Name) == "BTCEXA" {
			btcexaObjList ,err := GBtcExaModels.HttpObjList()
			if err != nil {
				ZapLog().Error("GBtcexaModels.HttpObjList err", zap.String("market", info.Name), zap.Error(err))
				continue
			}
			searcher.Update(info.Abb, BtcexaExaCurrencyPairs2ApiJPBShuJu(btcexaObjList))
			ZapLog().Info("loadMarket ok", zap.String("market", info.Name), zap.Int("num", len(btcexaObjList.Result)))
			continue
		}
		objList, err := GCoinMeritModels.HttpObjList(strings.ToLower(info.Name))
		if err != nil {
			ZapLog().Error("GCoinMeritModels.HttpObjList err", zap.String("market", info.Name), zap.Error(err))
			continue
		}
		searcher.Update(info.Abb, CoinMeritExaCurrencyPairs2ApiJPBShuJu(&objList.Data))
		//time.Sleep(time.Second*1)
		ZapLog().Info("loadMarket ok", zap.String("market", info.Name), zap.Int("num", len(objList.Data.CurrencyPairs)))
		continue
	}
	this.Lock()
	defer this.Unlock()
	this.objSearcher = searcher
}

func CoinMeritExaCurrencyPairs2ApiJPBShuJu(p *_type.CoinMeritExaCurrencyPairs) []*api.JPBShuJu {
	arr := make([]*api.JPBShuJu, len(p.CurrencyPairs), len(p.CurrencyPairs))
	for i:=0; i <len(p.CurrencyPairs); i++ {
		p.CurrencyPairs[i] = strings.ToUpper(p.CurrencyPairs[i])
		pairs := strings.Split(p.CurrencyPairs[i], "_")
		zhName := ""
		for j:=0; j < len(pairs); j++ {
			if len(pairs[j]) == 0 {
				continue
			}
			detail, ok := config.GPreConfig.CoinDetailMap[pairs[j]]
			if !ok {
				continue
			}
			if len(zhName) != 0 {
				zhName += "_"+detail.ZhName
			}else{
				zhName += detail.ZhName
			}
		}
		arr[i] = &api.JPBShuJu{
			DaiMa: &p.CurrencyPairs[i],
			MingCheng: &zhName,
		}
	}
	return arr
}


func BtcexaExaCurrencyPairs2ApiJPBShuJu(p *_type.ResBtcExaObjList) []*api.JPBShuJu {
	arr := make([]*api.JPBShuJu, len(p.Result), len(p.Result))
	result := p.Result
	CurrencyPairs := make([]string,0)

	for _ , detailInfo := range  result{
		CurrencyPairs = append(CurrencyPairs,detailInfo.Name)
	}
	for i:=0; i <len(CurrencyPairs); i++ {
		CurrencyPairs[i] = strings.ToUpper(CurrencyPairs[i])
		pairs := strings.Split(CurrencyPairs[i], "_")
		zhName := ""
		for j:=0; j < len(pairs); j++ {
			if len(pairs[j]) == 0 {
				continue
			}
			detail, ok := config.GPreConfig.CoinDetailMap[pairs[j]]
			if !ok {
				continue
			}
			if len(zhName) != 0 {
				zhName += "_"+detail.ZhName
			}else{
				zhName += detail.ZhName
			}
		}
		arr[i] = &api.JPBShuJu{
			DaiMa: &CurrencyPairs[i],
			MingCheng: &zhName,
		}
	}
	return arr




}