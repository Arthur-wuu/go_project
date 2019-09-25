package controllers

import (

	"sort"
	"go.uber.org/zap"
	"github.com/kataras/iris"
	"strings"
	. "BastionPay/bas-base/log/zap"
	apiquote "BastionPay/bas-api/quote"
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-quote/quote"
	. "BastionPay/bas-quote/utils"
	"BastionPay/bas-quote/config"
)

func NewCodeCtl(q *quote.QuoteMgr) *CodeCtl {
	cc:= &CodeCtl{
		mQuote:q,
		mLegals :   make([]*apiquote.CodeInfo,0),
	}
	for i:=0; i< len(config.GPreConfig.CountryCodeArr); i++ {
		info := &apiquote.CodeInfo{
			Symbol: &config.GPreConfig.CountryCodeArr[i],
			Name: &config.GPreConfig.CountryNameArr[i],
		}
		cc.mLegals = append(cc.mLegals, info)
	}
	return cc
}

type CodeCtl struct{
	mQuote *quote.QuoteMgr
	mLegals    []*apiquote.CodeInfo
}

func (this *CodeCtl) ListSymbols(ctx iris.Context) {
	defer PanicPrint()
	ZapLog().Debug("start handleListSymbols")
	//symbols := strings.ToUpper(ctx.URLParam("symbols"))
	symbols := ctx.URLParam("symbols")
	symbols = strings.TrimSpace(symbols)
	symbols = strings.TrimRight(symbols, ",")

	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")
	codes := make([]*apiquote.CodeInfo, 0)
	if len(symbols) == 0 {
		symbols := this.mQuote.ListSymbols()
		for i:=0; i < len(symbols); i++{
			codes = append(codes, ToApiCodeInfo(&symbols[i]))
		}
	}else{
		symbolsArr := strings.Split(symbols, ",")
		for i:=0; i< len(symbolsArr); i++ {
			if len(symbolsArr[i]) == 0 {
				continue
			}
			info := this.mQuote.GetSymbol(symbolsArr[i])
			if info == nil {
				ZapLog().With(zap.String("symbol", symbolsArr[i])).Warn("GetSymbol nofind err")
				continue
			}
			codes = append(codes,  ToApiCodeInfo(info))
		}
	}
	sort.Sort(CodeInfoList(codes))
	resMsg.Codes = codes
	ctx.JSON(resMsg)
	ZapLog().Debug("deal handleListSymbols")
	return

}

type CodeInfoList []*apiquote.CodeInfo
func (d CodeInfoList) Len() int           { return len(d) }
//func (d CodeInfoList) Less(i, j int) bool { return uintptr(unsafe.Pointer(d[i].Id))< uintptr(unsafe.Pointer(d[j].Id)) }
func (d CodeInfoList) Less(i, j int) bool { return d[i].GetId()< d[j].GetId()}
func (d CodeInfoList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func (this *CodeCtl) ListCoinAndFx(ctx iris.Context) {
	defer PanicPrint()
	ZapLog().Debug("start handleListSymbols")

	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")
	coins := new(apiquote.CoinInfo)
	Digitals  := make([]*apiquote.CodeInfo, 0)
	symbols := this.mQuote.ListSymbols()
	for i:=0; i < len(symbols); i++{
		Digitals = append(Digitals, ToApiCodeInfo(&symbols[i]))
	}
	sort.Sort(CodeInfoList(Digitals))
	coins.Digitals = Digitals
	coins.Legals = this.mLegals
	resMsg.Coins = coins
	ctx.JSON(resMsg)
	ZapLog().Debug("deal handleListSymbols")
	return
}

