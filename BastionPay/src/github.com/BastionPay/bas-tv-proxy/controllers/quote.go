package controllers


import (
	"BastionPay/bas-tv-proxy/config"
	"BastionPay/bas-tv-proxy/base"
	"BastionPay/bas-tv-proxy/api"
	"github.com/kataras/iris"
	"strings"
)


var GQuoteCtrl QuoteCtrl

type QuoteCtrl struct {
	mConf *config.Config
}

func (this *QuoteCtrl) Init(c *config.Config) error {
	this.mConf = c
	//ZapLog().Sugar().Infof("Ctrl CoinMerit Init ok")
	return nil
}

func (this *QuoteCtrl) HandleWsKXian(request base.Requester) {
	exa := request.GetParamValue("market")
	_,ok := config.GPreConfig.MarketMap[strings.ToUpper(exa)]
	if !ok {
		request.OnResponseWithPack(api.ErrCode_Param, nil)
		return
	}
	switch strings.ToUpper(exa) {
	case "BTE":
		GBtcExa.HandleWsKXian(request)
		break
	default:
		GCoinMerit.HandleWsKXian(request)
		return
	}
}

func (this *QuoteCtrl) HandleHttpKXian(ctx iris.Context) {
	exa := ctx.URLParam("market")
	_,ok := config.GPreConfig.MarketMap[strings.ToUpper(exa)]
	if !ok {
		ctx.JSON(api.NewErrResponse(ctx.URLParam("qid"), api.ErrCode_Param))
		return
	}
	switch strings.ToUpper(exa) {
	case "BTE":
		GBtcExa.HandleHttpKXian(ctx)
		break
	default:
		GCoinMerit.HandleHttpKXian(ctx)
		break
	}
}