package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/common"
	"BastionPay/merchant-api/models"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type (
	Pos struct {
		Controllers
	}
)

//Pos机 创建订单
func (this *Pos) CreatePos(ctx iris.Context) {
	param := new(api.PosTrade)
	times := time.Now().Local().Format("2006-01-02 15:04:05")
	param.TimeStamp = &times

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	param.MerchantPosNo = common.GenerateUuid()
	//传的是各种法币的币种，法币的数量，还有数字货币的币种
	amount, err := baspay.GetAmount(*param.Legal, *param.Amount, *param.Assets)
	if err != nil {
		ZapLog().Error("amount err", zap.Error(err))
		this.ExceptionSerive(ctx, 100501, "amount err")
		return
	}
	amountString := strconv.FormatFloat(*amount, 'f', -1, 64)
	param.Amount = &amountString

	//存一张表，pos机的订单表 //测试环境不添加表了
	//err = new(models.PosTrade).Parse(param).Add()
	//if err != nil {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "add to pos table err")
	//	ZapLog().Error( "add to pos table err", zap.Error(err))
	//	return
	//}

	res, err := new(baspay.PosTrade).PosParse(param).PosSend()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		this.ExceptionSerive(ctx, 100571, "pos err")
		return
	}
	if res.Code != 0 {
		this.ExceptionSerive(ctx, 100581, "pos err")
		return
	}

	err = new(models.PosTrade).UpdateByPosTradeNo(param.MerchantPosNo, 1)
	if err != nil {
		ZapLog().Error("UpdateByPosTradeNo err", zap.Error(err))
		ctx.JSON(Response{Code: 100591, Message: err.Error()})
		return
	}

	res.Data.TimeStamp = *param.TimeStamp
	res.Data.Uid = *param.PayVoucher
	res.Data.Legal = *param.Legal
	res.Data.LegalAmount = *param.Amount
	res.Data.Amount = Decimal8(res.Data.Amount)

	this.PosResponse(ctx, *res)
}

//pos机的订单查询
func (this *Pos) PosOrders(ctx iris.Context) {
	param := new(api.PosOrders)
	times := time.Now().Local().Format("2006-01-02 15:04:05")
	param.TimeStamp = &times

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	//param.MerchantPosNo = common.GenerateUuid()
	//传的是各种法币的币种，法币的数量，还有数字货币的币种
	//amount, err := baspay.GetAmount( *param.Legal, *param.Amount ,strings.ToLower(*param.Assets))
	//if err != nil {
	//	ZapLog().Error( "amount err", zap.Error(err))
	//	this.ExceptionSerive(ctx, 100501, "amount err")
	//	return
	//}
	//amountString := strconv.FormatFloat(*amount, 'g', -1, 64)
	//param.Amount = &amountString

	res, err := new(baspay.PosOrders).PosOrderParse(param).PosOrdersSend()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		this.ExceptionSerive(ctx, 100551, "pos err")
		return
	}
	if res.Code != 0 {
		this.ExceptionSerive(ctx, 100561, "pos err")
		return
	}
	//res.Data.TimeStamp = *param.TimeStamp
	//res.Data.Uid = *param.PayVoucher

	this.PosOrdersRes(ctx, *res)
}

func (this *Pos) List(ctx iris.Context) {
	param := new(models.PosList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "get pos list err:"+err.Error())
		ZapLog().Error("get pos list err", zap.Error(err))
		return
	}

	data, err := param.ListWithConds(param.Page, param.Size, param.TimeStart, param.TimeEnd)
	if err != nil {
		//l4g.Error("AddRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, data)
}

func Decimal8(value string) string {
	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return ""
	}
	f64, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", float), 64)
	s2 := strconv.FormatFloat(f64, 'g', -1, 64) //float64
	return s2
}
