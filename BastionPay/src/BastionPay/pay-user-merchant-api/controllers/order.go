package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/pay-user-merchant-api/api"
	"BastionPay/pay-user-merchant-api/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"runtime/debug"
)

type (
	Trade struct {
		Controllers
	}
)

//交易订单 列表
func (this *Trade) ListTrade(ctx iris.Context) {
	param := new(api.TradeList)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	res, err := new(models.Trade).List(param.Page, param.Size, param.MerchantTradeNo, param.PayeeId, param.TradeNo)
	if err != nil {
		ZapLog().Error("List err", zap.Error(err))
		ctx.JSON(OrderResponse{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}
	this.Response(ctx, res)
}

//退单 列表
func (this *Trade) ReFundList(ctx iris.Context) {
	defer PanicPrint()
	param := new(api.RefundTradeList)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	if param.MerchantId == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "merchantid nil or original num nil err")
		ZapLog().Error("merchantid nil or original num nil err")
		return
	}

	if param.Size == nil || param.Size == nil {
		var page int64 = 1
		var size int64 = 10

		param.Page = &page
		param.Size = &size
	}

	res, originOrderList, err := new(models.ReFund).ParseList(param).List(*param.MerchantId, *param.Page, *param.Size)
	if err != nil {
		ZapLog().Error("UpdateTransferStatusByTradeNo err", zap.Error(err))
		ctx.JSON(OrderResponse{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "get refund list err"})
		return
	}

	if len(res) == 0 || len(originOrderList) == 0 {
		ZapLog().Info("this merchant no refund order ")
		ctx.JSON(OrderResponse{Code: 0, Message: "no refund"})
		return
	}

	//根据退单里的原始订单查询 原始订单里的金额，币种
	tradeInfos, err := new(models.Trade).GetByOriginNoList(originOrderList)
	if err != nil {
		ZapLog().Error("use refund order origin no get trade info err", zap.Error(err))
		ctx.JSON(OrderResponse{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "use refund order origin no get trade info err"})
		return
	}

	reFundWithInfo := make([]*models.ReFundWithInfo, 0)

	for i := 0; i < len(res); i++ {
		reInfo := new(models.ReFundWithInfo).Parse(res[i])
		reInfo.PayeeId = tradeInfos[i].PayeeId
		reInfo.Assets = tradeInfos[i].Assets
		reInfo.Amount = tradeInfos[i].Amount
		reFundWithInfo = append(reFundWithInfo, reInfo)
	}
	this.Response(ctx, &reFundWithInfo)
}

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}
