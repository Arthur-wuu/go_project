package kbspirit

import (
	//"regexp"
	"strings"
	//"time"
	//"os"
	"BastionPay/bas-base/log/zap"
	"BastionPay/bas-tv-proxy/api"
	qs "BastionPay/bas-tv-proxy/base/quicksearch"
	"go.uber.org/zap"
	"time"
)

// 证券搜索器
type ObjSearcher struct {
	prefixTree map[string]qs.IQuickSearch //全代码、去掉市场的代码、去掉市场的拼音 建前缀树 (无法区分拼音还是obj代码)
	suffixTree map[string]qs.IQuickSearch //去掉市场代码, 取反， 后缀树
	//pingyinTree    map[string]qs.IQuickSearch //拼音,无市场，前缀树
	hanziTree map[string]qs.IQuickSearch //汉子 无市场，前缀树

	stocks      map[string][]*api.JPBShuJu
	updateTimes map[string]int64
}

// NewObjSearcher 新建证券搜索器
func NewObjSearcher() *ObjSearcher {
	// 按市场区分
	prefixTree := make(map[string]qs.IQuickSearch)
	suffixTree := make(map[string]qs.IQuickSearch)
	hanziTree := make(map[string]qs.IQuickSearch)
	hanziTree["*"] = qs.New(qs.QSType_Trie)

	stocks := make(map[string][]*api.JPBShuJu)
	updateTimes := make(map[string]int64)

	return &ObjSearcher{
		prefixTree:  prefixTree,
		suffixTree:  suffixTree,
		hanziTree:   hanziTree,
		stocks:      stocks,
		updateTimes: updateTimes,
	}
}

// Init 初始化
func (this *ObjSearcher) Init() error {
	return nil
}

//这是云平台的方式，存在的问题是 无序，并且需要各种去重。
//优化 不区分市场建树，前缀树，后缀树（_USDT）,汉字树 三颗即可。并且采用数组形式。 树的节点包含儿子及子孙节点并且无重复。
func (this *ObjSearcher) Update(market string, arrs []*api.JPBShuJu) {
	prefixTree := qs.New(qs.QSType_Trie)
	suffixTree := qs.New(qs.QSType_Trie)
	hanziTree := qs.New(qs.QSType_Trie)

	market = strings.ToUpper(market)

	for i := 0; i < len(arrs); i++ {
		objInfo := arrs[i]
		daima := strings.Replace(objInfo.GetDaiMa(), " ", "", -1)
		daima = strings.ToUpper(strings.Replace(daima, "_", "", -1))
		mingcheng := strings.Replace(objInfo.GetMingCheng(), " ", "", -1)
		mingcheng = strings.ToUpper(strings.Replace(mingcheng, "_", "", -1))
		if len(daima) == 0 {
			continue
		}
		prefixTree.Insert(daima, i)
		prefixTree.Insert(market+daima, i)
		suffixTree.Insert(Reverse(daima), i)

		if len(mingcheng) == 0 {
			continue
		}
		hanziTree.Insert(mingcheng, i)
		hanziTree.Insert(market+mingcheng, i)
		//halfname := Half(mingcheng)   //拼音在币中文名称不全时最好不要用
		//for _, name := range FirstLetter(halfname) {
		//	name = strings.ToUpper(name)
		//	prefixTree.Insert(name, i)
		//	num := -1
		//	name, num = DeleteNotLetter(name)
		//	if num > 0{
		//		prefixTree.Insert(name, i)
		//	}
		//}
	}

	this.prefixTree[market] = prefixTree
	this.suffixTree[market] = suffixTree
	this.hanziTree[market] = hanziTree
	this.stocks[market] = arrs
	this.updateTimes[market] = time.Now().Unix()
}

//
func (this *ObjSearcher) Search(key string, env *Env) []*api.JPBShuChu {
	prefixTree := this.prefixTree
	suffixTree := this.suffixTree
	hanziTree := this.hanziTree
	updateTimes := this.updateTimes

	result := make([]*api.JPBShuChu, 0, int(env.Count)+20)
	needCount := env.Count
	for _, market := range env.Markets {
		stockArr, ok := this.stocks[market]
		if !ok {
			continue
		}
		if env.ChineseFlag {
			tempArr := make([]*api.JPBShuJu, 0)
			trieIndex, _ := hanziTree[market]
			for _, object := range trieIndex.ValueForPrefix(key, []int64{env.Count}...) {
				index := object.(int)
				stock := stockArr[index]
				tempArr = append(tempArr, stock)
			}
			if len(tempArr) == 0 {
				continue
			}
			tempArr = RemoveDuplicate(tempArr)
			shuchu := api.NewJPBShuChu(0, tempArr, market, "", updateTimes[market])
			result = append(result, shuchu)
			needCount -= int64(len(tempArr))
			if needCount <= 0 {
				break
			}
			continue
		}
		tempArr := make([]*api.JPBShuJu, 0)
		if env.SuffixFlag {
			log.ZapLog().Info("use suffix Tree")
			trieIndex, _ := suffixTree[market]
			revsKey := Reverse(key)
			for _, object := range trieIndex.ValueForPrefix(revsKey, []int64{env.Count}...) {
				index := object.(int)
				stock := stockArr[index]
				tempArr = append(tempArr, stock)
			}
			log.ZapLog().Info("use suffix Tree", zap.Int("num", len(tempArr)))
			tempArr = RemoveDuplicate(tempArr)
			needCount -= int64(len(tempArr))
			if needCount <= 0 {
				shuchu := api.NewJPBShuChu(0, tempArr, market, "", updateTimes[market])
				result = append(result, shuchu)
				break
			}
		}
		trieIndex, _ := prefixTree[market]
		for _, object := range trieIndex.ValueForPrefix(key, []int64{env.Count}...) {
			index := object.(int)
			stock := stockArr[index]
			tempArr = append(tempArr, stock)
		}
		if len(tempArr) == 0 {
			continue
		}
		tempArr = RemoveDuplicate(tempArr)
		shuchu := api.NewJPBShuChu(0, tempArr, market, "", updateTimes[market])
		result = append(result, shuchu)
		needCount -= int64(len(tempArr))
		if needCount <= 0 {
			break
		}
	}

	return result
}

//
//// UnInit 反初始化
func (self *ObjSearcher) UnInit() {
}
