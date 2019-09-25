package controllers

import (
	"BastionPay/bas-api/apibackend"
	apiquote "BastionPay/bas-api/quote"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/quote"
	. "BastionPay/bas-quote-collect/utils"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"sort"
	"strings"
)

func NewCodeCtl(q *quote.QuoteMgr) *CodeCtl {
	return &CodeCtl{
		mQuote: q,
	}
}

type CodeCtl struct {
	mQuote *quote.QuoteMgr
}

func (this *CodeCtl) SetSymbol(ctx iris.Context) {
	defer PanicPrint()
	cc := make([]*apiquote.CodeInfo, 0)
	if err := ctx.ReadJSON(&cc); err != nil {
		ZapLog().With(zap.Error(err)).Error("param err")
		ctx.JSON(apiquote.NewResMsg(apibackend.BASERR_INVALID_PARAMETER.Code(), "param err"))
		return
	}
	var err error
	for i := 0; i < len(cc); i++ {
		if cc[i].Symbol != nil {
			//*cc[i].Symbol = strings.ToUpper(*cc[i].Symbol)
			*cc[i].Symbol = *cc[i].Symbol
		}
		if cc[i].Valid == nil {
			cc[i].Valid = new(int)
			*cc[i].Valid = 1
		}
		lc := ToLocalCodeInfo(cc[i])
		lc.SetTimestamp(quote.NowTimestamp())
		ZapLog().Info(lc.ToPrintStr())
		if len(cc[i].GetSymbol()) == 0 {
			ZapLog().With(zap.String("symbol", cc[i].GetSymbol())).Error("param err")
			ctx.JSON(apiquote.NewResMsg(apibackend.BASERR_INVALID_PARAMETER.Code(), "param err"))
			return
		}
		if err = this.mQuote.SetCodeTable(lc); err != nil {
			ZapLog().With(zap.Error(err), zap.Any("CodeInfo", lc.ToPrintStr())).Error("SetCodeTable err")
			ctx.JSON(apiquote.NewResMsg(apibackend.BASERR_DATABASE_ERROR.Code(), "SetCodeTable err:"+err.Error()))
			return
		}

	}
	ctx.JSON(apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), ""))
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
		for i := 0; i < len(symbols); i++ {
			codes = append(codes, ToApiCodeInfo(&symbols[i]))
		}
	} else {
		symbolsArr := strings.Split(symbols, ",")
		for i := 0; i < len(symbolsArr); i++ {
			if len(symbolsArr[i]) == 0 {
				continue
			}
			info := this.mQuote.GetSymbol(symbolsArr[i])
			if info == nil {
				ZapLog().With(zap.String("symbol", symbolsArr[i])).Warn("GetSymbol nofind err")
				continue
			}
			codes = append(codes, ToApiCodeInfo(info))
		}
	}
	sort.Sort(CodeInfoList(codes))
	resMsg.Codes = codes
	ctx.JSON(resMsg)
	ZapLog().Debug("deal handleListSymbols")
	return
}

type CodeInfoList []*apiquote.CodeInfo

func (d CodeInfoList) Len() int { return len(d) }

//func (d CodeInfoList) Less(i, j int) bool { return uintptr(unsafe.Pointer(d[i].Id))< uintptr(unsafe.Pointer(d[j].Id)) }
func (d CodeInfoList) Less(i, j int) bool { return d[i].GetId() < d[j].GetId() }
func (d CodeInfoList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
