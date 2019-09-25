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
	"strconv"
	"strings"
)

var GBtcExa BtcExa

type BtcExa struct {
	mConf *config.Config
}

func (this *BtcExa) Init(c *config.Config) error {
	this.mConf = c
	if err := models.GBtcExaModels.Init(&base.GWsServerMgr); err != nil{
		return err
	}
	if err := models.GBtcExaModels.Start(); err != nil {
		ZapLog().Error("models.GBtcExaModels.Start() err:"+ err.Error())
	} //不需要任何错误处理
	ZapLog().Sugar().Infof("Ctrl BtcExa Init ok")
	return nil
}

/**********wsBtcExa*************/
func (this *BtcExa) HandleWsKXian(request base.Requester)  {
	ZapLog().Info("HandleWsKXian")
	objs := strings.Split(strings.ToLower(request.GetParamValue("obj")),",")
	if request.IsSimpleReq() {
		quotekxians := make([]*api.QuoteKlineSingle, 0, len(objs))
		for i:=0; i< len(objs); i++ {
			beKline,err := models.GBtcExaModels.HttpbtcKXian(objs[i], ToBtcExaPeriod(request.GetParamValue("period")), ToBtcExaPeriod(request.GetParamValue("count")))
			if err != nil {
				ZapLog().Error("GBtcExaModels HttpKXian err", zap.Error(err))
				request.OnResponseWithPack(api.ErrCode_InerServer, nil)
				return
			}
			quotekxians = append(quotekxians, api.NewQuoteKlineSingle(objs[i],BtcExaKXianToApiKXian(beKline)))
		}
		resmsg := new(api.MSG)
		resmsg.AddQuoteKlineSingle(quotekxians...)
		request.OnResponseWithPack(api.ErrCode_Success, resmsg)
		return
	}
	for i:=0; i < len(objs); i++ {
		btcexaReq, err := _type.NewReqBtcExaSub(request.GetParamValue("sub"), objs[i], ToBtcExaPeriod(request.GetParamValue("period")),ToBtcExaPeriod(request.GetParamValue("count")))
		if err != nil {
			ZapLog().Error("NewReqBtcExaSub err", zap.Error(err), zap.String("conid", request.GetConId()), zap.Any("btcexa_req", request))
			request.OnResponseWithPack(api.ErrCode_Param, nil)
			return
		}
		err = models.GBtcExaModels.WsSend(request.GetFirstPath(),request.GetUuid(), request.GetSub() ,btcexaReq)
		if err != nil {
			ZapLog().Error("GBtcExaModels WsSend err", zap.Error(err), zap.String("conid", request.GetConId()), zap.Any("btcexa_req", request))
			request.OnResponseWithPack(api.ErrCode_InerServer, nil)
			return
		}
	}

	request.OnResponseWithPack(api.ErrCode_Success, nil)
}



/****************http*************/
func (this *BtcExa) HandleHttpKXian(ctx iris.Context) {
	qid := ctx.URLParam("qid")
	objs := strings.Split(strings.ToLower(ctx.URLParam("obj")), ",")
	period := ToBtcExaPeriod(ctx.URLParam("period"))
	limit := ctx.URLParam("count")
	if  len(objs) == 0 || len(period) == 0 {
		ZapLog().Error(" HandleHttpKXian param err")
		ctx.JSON(api.NewErrResponse(qid, api.ErrCode_Param))
		return
	}
	quotekxians := make([]*api.QuoteKlineSingle, 0, len(objs))
	for i:=0; i< len(objs); i++ {
		beKline,err := models.GBtcExaModels.HttpbtcKXian( objs[i], period, limit)
		if err != nil {
			ZapLog().Error("GBtcExaModels HttpKXian err", zap.Error(err))
			ctx.JSON(api.NewErrResponse(qid, api.ErrCode_InerServer))
			return
		}
		quotekxians = append(quotekxians, api.NewQuoteKlineSingle(objs[i],BtcExaKXianToApiKXian(beKline)))
	}

	resmsg := new(api.MSG)
	resmsg.AddQuoteKlineSingle(quotekxians...)
	common.CtxJson(ctx, qid, resmsg)
}

func (this *BtcExa) HandleHttpObjList(ctx iris.Context) {
	qid := ctx.URLParam("qid")

	beObjs,err := models.GBtcExaModels.HttpObjList()
	if err != nil {
		ZapLog().Error("GBtcExaModels HttpObjList err", zap.Error(err))
		ctx.JSON(api.NewErrResponse(qid, api.ErrCode_InerServer))
		return
	}
	apiMarket := BtcExaCurrencyPairsToApiMarket(beObjs)
	ZapLog().Info("BtcExa_res", zap.Any("BtcExaRes", *beObjs), zap.Any("apires", apiMarket))
	resmsg := new(api.MSG)
	resmsg.AddMarket(apiMarket)
	common.CtxJson(ctx, qid, resmsg)
}

func BtcExaKXianToApiKXian(k1 *_type.ResBtcExaKLine) []*api.KXian {
	k2 := make([]*api.KXian, 0)
	for j:=0; j < len(k1.Result); j++ {
		tmp := new(api.KXian)

		for i := 0; i < len(k1.Result[j]); i++ {
			switch i {
			case 0:
				in, err := strconv.ParseInt(k1.Result[j][0], 10, 64)
				if err != nil {
					ZapLog().Error("string to int err", zap.Error(err))
				}
				tmp.SetShiJian(int64(in))
				break
			case 1:
				tmp.SetKaiPanJia(k1.Result[j][1])
				break
			case 2:
				tmp.SetZuiGaoJia(k1.Result[j][2])
				break
			case 3:
				tmp.SetZuiDiJia(k1.Result[j][3])
				break
			case 4:
				tmp.SetShouPanJia(k1.Result[j][4])
				break
			case 5:
				tmp.SetChengJiaoLiang(k1.Result[j][5])
				break
			}
		}
		k2 = append(k2, tmp)
	}

	return k2
}



func BtcExaCurrencyPairsToApiMarket(c1 *_type.ResBtcExaObjList) *api.Market {
	c2 :=new(api.Market)
	name := "BtcExa"
	c2.Name = &name

	obj:=make([]string,0)
	for i:=0; i < len(c1.Result); i++ {
		obj=append(obj, c1.Result[i].Name)
	}

	c2.Objs = obj
	fmt.Println("obj: ",obj)
	return c2
}








