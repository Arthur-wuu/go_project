package controllers

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-tv-proxy/api"
	"BastionPay/bas-tv-proxy/base"
	"BastionPay/bas-tv-proxy/common"
	"BastionPay/bas-tv-proxy/config"
	"BastionPay/bas-tv-proxy/models"
	"BastionPay/bas-tv-proxy/type"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

var GCoinMerit CoinMerit

type CoinMerit struct {
	mConf *config.Config
}

func (this *CoinMerit) Init(c *config.Config) error {
	this.mConf = c
	if err := models.GCoinMeritModels.Init(&base.GWsServerMgr); err != nil {
		return err
	}
	if err := models.GCoinMeritModels.Start(); err != nil {
		ZapLog().Error("models.GCoinMeritModels.Start() err:" + err.Error())
	} //不需要任何错误处理
	ZapLog().Sugar().Infof("Ctrl CoinMerit Init ok")
	return nil
}

/**********ws*************/
func (this *CoinMerit) HandleWsKXian(request base.Requester) {
	ZapLog().Info("HandleWsKXian")

	exa := strings.ToLower(request.GetParamValue("market"))
	objs := strings.Split(strings.ToLower(request.GetParamValue("obj")), ",")
	period := strings.ToLower(request.GetParamValue("period"))
	start := request.GetParamValue("start")
	count := request.GetParamValue("count")

	info, ok := config.GPreConfig.MarketMap[strings.ToUpper(exa)]
	if !ok {
		request.OnResponseWithPack(api.ErrCode_Param, nil)
		return
	}
	exa = strings.ToLower(info.Name)

	quotekxians := make([]*api.QuoteKlineSingle, 0, len(objs))
	if !request.IsUnSub() {
		for i := 0; i < len(objs); i++ {
			cmKline, err := models.GCoinMeritModels.HttpKXian(exa, objs[i], period, count, start)
			if err != nil {
				ZapLog().Error("GCoinMeritModels HttpKXian err", zap.Error(err))
				request.OnResponseWithPack(api.ErrCode_InerServer, nil)
				return
			}
			quotekxians = append(quotekxians, api.NewQuoteKlineSingle(objs[i], CoinMeritKXianToApiKXian(cmKline)))
		}
	}
	if !request.IsSimpleReq() {
		for i := 0; i < len(objs); i++ {
			cmReq, err := _type.NewReqCoinMeritSub(strings.ToLower(exa), objs[i], strings.ToLower(request.GetParamValue("period")), request.GetSub())
			if err != nil {
				ZapLog().Error("NewReqCoinMeritSub err", zap.Error(err), zap.String("conid", request.GetConId()), zap.Any("cm_req", request))
				request.OnResponseWithPack(api.ErrCode_Param, nil)
				return
			}
			err = models.GCoinMeritModels.WsSend(request.GetFirstPath(), request.GetUuid(), request.GetSub(), cmReq)
			if err != nil {
				ZapLog().Error("GCoinMeritModels WsSend err", zap.Error(err), zap.String("conid", request.GetConId()), zap.Any("cm_req", request))
				request.OnResponseWithPack(api.ErrCode_InerServer, nil)
				return
			}
		}
	}

	if request.IsUnSub() {
		request.OnResponseWithPack(api.ErrCode_Success, nil)
		return
	}

	resmsg := new(api.MSG)
	resmsg.AddQuoteKlineSingle(quotekxians...)

	request.OnResponseWithPack(api.ErrCode_Success, resmsg)
}

/****************http*************/
func (this *CoinMerit) HandleHttpKXian(ctx iris.Context) {
	qid := ctx.URLParam("qid")
	exa := strings.ToLower(ctx.URLParam("market"))
	objs := strings.Split(strings.ToLower(ctx.URLParam("obj")), ",")
	period := strings.ToLower(ctx.URLParam("period"))
	start := ctx.URLParam("start")
	count := ctx.URLParam("count")
	if len(exa) == 0 || len(objs) == 0 || len(period) == 0 {
		ZapLog().Error(" HandleHttpKXian param err")
		ctx.JSON(api.NewErrResponse(qid, api.ErrCode_Param))
		return
	}
	info, ok := config.GPreConfig.MarketMap[strings.ToUpper(exa)]
	if !ok {
		ctx.JSON(api.NewErrResponse(qid, api.ErrCode_Param))
		return
	}
	exa = strings.ToLower(info.Name)

	quotekxians := make([]*api.QuoteKlineSingle, 0, len(objs))
	for i := 0; i < len(objs); i++ {
		cmKline, err := models.GCoinMeritModels.HttpKXian(exa, objs[i], period, count, start)
		if err != nil {
			ZapLog().Error("GCoinMeritModels HttpKXian err", zap.Error(err))
			ctx.JSON(api.NewErrResponse(qid, api.ErrCode_InerServer))
			return
		}
		quotekxians = append(quotekxians, api.NewQuoteKlineSingle(objs[i], CoinMeritKXianToApiKXian(cmKline)))
	}

	//	ZapLog().Info("cm_res", zap.Any("cmres", *cmKline), zap.Any("apires", apiKLine))
	resmsg := new(api.MSG)
	resmsg.AddQuoteKlineSingle(quotekxians...)
	common.CtxJson(ctx, qid, resmsg)
}

func (this *CoinMerit) HandleHttpExa(ctx iris.Context) {
	qid := ctx.URLParam("qid")
	cmExa, err := models.GCoinMeritModels.HttpExa()
	if err != nil {
		ZapLog().Error("GCoinMeritModels HttpExa err", zap.Error(err))
		ctx.JSON(api.NewErrResponse(qid, api.ErrCode_InerServer))
		return
	}
	apiMarket := CoinMeritExchangesToApiMarket(cmExa)
	ZapLog().Info("cm_res", zap.Any("cmres", *cmExa), zap.Any("apires", apiMarket))
	resmsg := new(api.MSG)
	resmsg.AddMarket(apiMarket...)
	common.CtxJson(ctx, qid, resmsg)
}

func (this *CoinMerit) HandleHttpObjList(ctx iris.Context) {
	qid := ctx.URLParam("qid")
	exa := strings.ToLower(ctx.URLParam("market"))
	if len(exa) == 0 {
		ZapLog().Error("HttpObjList param err:market is nil")
		ctx.JSON(api.NewErrResponse(qid, api.ErrCode_Param))
		return
	}
	cmObjs, err := models.GCoinMeritModels.HttpObjList(exa)
	if err != nil {
		ZapLog().Error("GCoinMeritModels HttpObjList err", zap.Error(err))
		ctx.JSON(api.NewErrResponse(qid, api.ErrCode_InerServer))
		return
	}
	apiMarket := CoinMeritCurrencyPairsToApiMarket(cmObjs)
	ZapLog().Info("cm_res", zap.Any("cmres", *cmObjs), zap.Any("apires", apiMarket))
	resmsg := new(api.MSG)
	resmsg.AddMarket(apiMarket)
	common.CtxJson(ctx, qid, resmsg)
}

func CoinMeritKXianToApiKXian(k1 *_type.ResCoinMeritKLine) []*api.KXian {
	k2 := make([]*api.KXian, 0)
	for j := 0; j < len(k1.Data); j++ {
		tmp := new(api.KXian)

		for i := 0; i < len(k1.Data[j]); i++ {
			switch i {
			case 0:
				tmp.SetShiJian(int64(k1.Data[j][0]))
				break
			case 1:
				tmp.SetKaiPanJia(fmt.Sprintf("%f", k1.Data[j][1]))
				break
			case 2:
				tmp.SetZuiGaoJia(fmt.Sprintf("%f", k1.Data[j][2]))
				break
			case 3:
				tmp.SetZuiDiJia(fmt.Sprintf("%f", k1.Data[j][3]))
				break
			case 4:
				tmp.SetShouPanJia(fmt.Sprintf("%f", k1.Data[j][4]))
				break
			case 5:
				tmp.SetChengJiaoLiang(fmt.Sprintf("%f", k1.Data[j][5]))
				break
			}

		}
		k2 = append(k2, tmp)
	}
	return k2
}

func CoinMeritExchangesToApiMarket(c1 *_type.ResCoinMeritExchanges) []*api.Market {
	c2 := make([]*api.Market, 0)
	for i := 0; i < len(c1.Data); i++ {
		tmp := new(api.Market)
		tmp.SetName(c1.Data[i])
		c2 = append(c2, tmp)
	}
	return c2
}

func CoinMeritCurrencyPairsToApiMarket(c1 *_type.ResCoinMeritCurrencyPairs) *api.Market {
	c2 := new(api.Market)
	c2.SetName(c1.Data.Exchange)
	c2.Objs = c1.Data.CurrencyPairs
	return c2
}

// 1min 5min         1hour            1Day 7Day 1Month
// 1m, 5m, 15m, 30m, 1h, 2h, 4h, 12h, 1D, 7D, 1M,
func ToBtcExaPeriod(s string) string {
	p := strings.ToLower(s)
	if strings.Contains(p, "min") {
		return strings.Replace(p, "min", "m", len(p))
	}
	if strings.Contains(p, "hour") {
		return strings.Replace(p, "hour", "h", len(p))
	}
	if strings.Contains(p, "day") {
		return strings.Replace(p, "day", "D", len(p))
	}
	if strings.Contains(p, "month") {
		return strings.Replace(p, "month", "M", len(p))
	}
	return s
}
