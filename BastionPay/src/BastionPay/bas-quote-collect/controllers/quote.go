package controllers

import (
	"BastionPay/bas-api/apibackend"
	apiquote "BastionPay/bas-api/quote"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/collect"
	"BastionPay/bas-quote-collect/config"
	"BastionPay/bas-quote-collect/quote"
	"BastionPay/bas-quote-collect/utils"
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

func NewQuoteCtl(q *quote.QuoteMgr) *QuoteCtl {
	return &QuoteCtl{
		mQuote: q,
	}
}

type QuoteCtl struct {
	mQuote *quote.QuoteMgr
}

func (this *QuoteCtl) Ticker(ctx iris.Context) {
	defer utils.PanicPrint()
	//from := strings.ToUpper(ctx.URLParam("from"))
	from := ctx.URLParam("from")
	if len(from) == 0 {
		ZapLog().With(zap.String("from", from)).Error("param err")
		ctx.JSON(*apiquote.NewResMsg(apibackend.BASERR_INVALID_PARAMETER.Code(), "param fail"))
		return
	}
	to := strings.ToUpper(ctx.URLParam("to"))
	if len(to) == 0 {
		to = "USD"
	}

	fromArr := strings.Split(strings.Replace(from, " ", "", -1), ",")
	toArr := strings.Split(strings.Replace(to, " ", "", -1), ",")

	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")
	var moneyInfo *apiquote.MoneyInfo
	var err error
	for i := 0; i < len(fromArr); i++ {
		if len(fromArr[i]) == 0 {
			continue
		}
		codeInfo := this.mQuote.GetSymbol(fromArr[i])
		if codeInfo == nil {
			ZapLog().With(zap.String("from", fromArr[i])).Info("GetSymbol nofind err")
			continue
		}
		QuoteDetailInfo := resMsg.GenQuoteDetailInfo()
		QuoteDetailInfo.Id = codeInfo.Id
		QuoteDetailInfo.Symbol = codeInfo.Symbol
		for j := 0; j < len(toArr); j++ {
			if len(toArr[j]) == 0 {
				continue
			}
			if codeInfo.GetId() >= 100000 {
				moneyInfo, err = this.getHighTicker(codeInfo.GetId(), fromArr[i], toArr[j])
			} else {
				moneyInfo, err = this.getLowTicker(codeInfo.GetId(), toArr[j])
			}
			if err != nil || moneyInfo == nil {
				continue
			}
			QuoteDetailInfo.AddMoneyInfo(moneyInfo)
		}
	}
	ctx.JSON(resMsg)
	ZapLog().Debug("deal handleTicker ok")
}

func (this *QuoteCtl) getLowTicker(from int, to string) (*apiquote.MoneyInfo, error) {
	moneyInfo := new(collect.MoneyInfo)
	moneyInfo2, err := this.mQuote.GetQuoteUseId(from, "USD")
	if err != nil {
		ZapLog().Info("get qt_X_USD err", zap.Error(err), zap.Int("from", from))
		return nil, err
	}
	if moneyInfo2 == nil {
		return nil, nil
	}

	moneyInfo3, err := this.mQuote.GetQuoteHuilv(to)
	if err != nil {
		ZapLog().Error("get qt_USD_to_X err", zap.Error(err), zap.String("huilv", to))
		return nil, err
	}
	if moneyInfo3 == nil {
		return nil, nil
	}

	moneyInfo.SetPrice((moneyInfo2.GetPrice()) * (moneyInfo3.GetPrice()))
	moneyInfo.SetSymbol(to)
	moneyInfo.SetLast_updated(moneyInfo2.GetLast_updated())
	return utils.ToApiMoneyInfo(moneyInfo), nil
}

func (this *QuoteCtl) getHighTicker(from int, fromStr, to string) (*apiquote.MoneyInfo, error) {
	SrcInfo, ok := config.GPreConfig.FromCollects[fromStr]
	if !ok {
		SrcInfo, ok = config.GPreConfig.IdsCollects[fmt.Sprintf("%d", from)]
		if !ok {
			ZapLog().Error(fromStr + " config not set")
			return nil, errors.New("config not set")
		}
	}
	moneyInfo := new(collect.MoneyInfo)
	moneyInfo1, err := this.mQuote.GetQuoteUseId(from, SrcInfo.Coin_to)
	if err != nil {
		ZapLog().Error("get qt_from_X_BTC err", zap.Error(err), zap.Int("id", from))
		return nil, err
	}
	if moneyInfo1 == nil {
		return nil, nil
	}
	codeInfo := this.mQuote.GetSymbol(SrcInfo.Coin_to)
	if codeInfo == nil {
		ZapLog().With(zap.String("from", SrcInfo.Coin_to)).Error("GetSymbol nofind err")
		return nil, errors.New("nofind in codeTable")
	}

	moneyInfo2, err := this.mQuote.GetQuoteUseId(codeInfo.GetId(), "USD")
	if err != nil {
		ZapLog().Error("get qt_X_USD err", zap.Error(err), zap.Int("id", codeInfo.GetId()))
		return nil, err
	}
	if moneyInfo2 == nil {
		return nil, nil
	}

	moneyInfo3, err := this.mQuote.GetQuoteHuilv(to)
	if err != nil {
		ZapLog().Error("get qt_USD_to_X err", zap.Error(err), zap.String("huilv", to))
		return nil, err
	}
	if moneyInfo3 == nil {
		return nil, nil
	}

	moneyInfo.SetPrice((moneyInfo1.GetPrice()) * (moneyInfo2.GetPrice()) * (moneyInfo3.GetPrice()))
	moneyInfo.SetSymbol(to)
	moneyInfo.SetLast_updated(moneyInfo1.GetLast_updated())

	return utils.ToApiMoneyInfo(moneyInfo), nil
}
