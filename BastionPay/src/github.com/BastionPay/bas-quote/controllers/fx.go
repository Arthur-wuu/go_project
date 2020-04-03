package controllers

import (
	"BastionPay/bas-api/apibackend"
	apiquote "BastionPay/bas-api/quote"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote/collect"
	"BastionPay/bas-quote/quote"
	"BastionPay/bas-quote/utils"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

func NewFxCtl(q *quote.QuoteMgr) *FxCtl {
	return &FxCtl{
		mQuote: q,
	}
}

type FxCtl struct {
	mQuote *quote.QuoteMgr
}

func (this *FxCtl) Ticker(ctx iris.Context) {
	defer utils.PanicPrint()
	ZapLog().Debug("start handleHuilv ")

	//判断coin是数字还是字符
	from := strings.ToUpper(ctx.URLParam("from"))
	to := strings.ToUpper(ctx.URLParam("to"))
	if len(from) == 0 {
		from = "USD"
	}
	if len(to) == 0 {
		ZapLog().With(zap.String("to", to)).Error("param err")
		ctx.JSON(*apiquote.NewResMsg(apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc()))
		return
	}

	from = strings.Replace(from, " ", "", -1)
	to = strings.Replace(to, " ", "", -1)
	to = strings.TrimRight(to, ",")
	from = strings.TrimRight(from, ",")
	toArr := strings.Split(to, ",")
	fromArr := strings.Split(from, ",")

	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")
	for i := 0; i < len(fromArr); i++ {
		if len(fromArr[i]) == 0 {
			continue
		}
		fromInfo, err := this.mQuote.GetQuoteHuilv(fromArr[i])
		if err != nil {
			ZapLog().Error("get qt_USD_to[i] err", zap.Error(err), zap.String("huilv", fromArr[i]))
			continue
		}

		QuoteDetailInfo := resMsg.GenQuoteDetailInfo()
		QuoteDetailInfo.Symbol = &fromArr[i]
		for j := 0; j < len(toArr); j++ {
			if len(toArr[j]) == 0 {
				continue
			}
			moneyInfo := new(collect.MoneyInfo)
			toInfo, err := this.mQuote.GetQuoteHuilv(toArr[j])
			if err != nil {
				ZapLog().Error("get qt_USD_to[i] err", zap.Error(err), zap.String("huilv", toArr[j]))
				continue
			}
			moneyInfo.SetPrice(toInfo.GetPrice() / fromInfo.GetPrice())
			moneyInfo.SetSymbol(toArr[j])
			if toInfo.GetLast_updated() != 0 {
				moneyInfo.SetLast_updated(toInfo.GetLast_updated())
			} else {
				moneyInfo.SetLast_updated(fromInfo.GetLast_updated())
			}

			QuoteDetailInfo.AddMoneyInfo(utils.ToApiMoneyInfo(moneyInfo))
		}
	}
	ctx.JSON(resMsg)
	ZapLog().Debug("deal handleTicker ok")
}
