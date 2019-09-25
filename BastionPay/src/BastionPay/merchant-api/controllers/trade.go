package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/common"
	"BastionPay/merchant-api/db"
	"BastionPay/merchant-api/models"
	"encoding/json"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"runtime/debug"
	"strconv"
)

type (
	Trade struct {
		Controllers
	}
)

func (this *Trade) CreateWeb(ctx iris.Context) {
	param := new(api.Trade)

	err := ctx.ReadJSON(param)
	//err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "create trade param err")
		ZapLog().Error("create trade param err", zap.Error(err))
		return
	}

	//添加一个价格判断，如果币种和币价存在差异较大，直接不让创建订单,
	intGameCoin, err := strconv.Atoi(param.GameCoin)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "string to int err")
		ZapLog().Error("string to int err", zap.Error(err))
		return
	}

	prePrice := this.CheckPricePre(*param.DeviceId, *param.Assets, intGameCoin)
	//比较价格，自己算的币的数量 上下波动 避免出现极小的数字货币支付成功

	amount, err := strconv.ParseFloat(*param.Amount, 64)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "string to float err")
		ZapLog().Error("string to float err", zap.Error(err))
		return
	}

	if amount > prePrice*0.2 && amount < amount*1.8 {
		goto Here
	} else {
		return
	}
	//

Here:
	param.MerchantTradeNo = common.GenerateUuid()
	param.ExpireTime = 900

	if err := new(models.Trade).Parse(param).Add(); err != nil { //唯一性错误判断
		ZapLog().Error("Add trade param err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	res, err := new(baspay.Trade).Parse(param).Send()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}
	if res.Code != 0 {
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "create h5 order err"})
		return
	}

	this.Response(ctx, res)
}

//创建二维码
func (this *Trade) CreateQr(ctx iris.Context) {
	param := new(api.QrTrade)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	param.MerchantTradeNo = common.GenerateUuid()
	param.ExpireTime = 900

	//if err := new(models.Trade).ParseQr(param).Add(); err != nil {//唯一性错误判断
	//	ZapLog().Error( "Add err", zap.Error(err))
	//	ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
	//	return
	//}

	res, err := new(baspay.QrTrade).QrParse(param).SendQr()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}
	if res.Code != 0 {
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "create qr order err"})
		return
	}

	//err = new(models.Trade).UpdateByTradeNo(*res.MerchantTradeNo, *res.QrCode, models.EUM_OPEN_STATUS_OPEN)
	//if err != nil {
	//	ZapLog().Error( "UpdateByTradeNo err", zap.Error(err))
	//	ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
	//	return
	//}

	this.Response(ctx, res.Data)
}

//SDK订单
func (this *Trade) CreateSdk(ctx iris.Context) {
	param := new(api.Trade)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "create trade param err")
		ZapLog().Error("create trade param err", zap.Error(err))
		return
	}
	param.MerchantTradeNo = common.GenerateUuid()
	param.ExpireTime = 900

	if err := new(models.Trade).Parse(param).Add(); err != nil {
		ZapLog().Error("Add trade param err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	res, err := new(baspay.Trade).Parse(param).Send()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, res)
}

//查询订单状态
func (this *Trade) TradeInfo(ctx iris.Context) {

	param := new(api.TradeInfo)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	if param.MerchantId == nil || param.MerchantTradeNo == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "merchantid nil or tradeno nil err")
		ZapLog().Error("merchantid nil or tradeno nil err")
		return
	}

	res, err := new(baspay.TradeInfo).Parse(param).Send()
	if err != nil {
		ZapLog().Error("select trade info err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, res)
}

//回调
func (this *Trade) Pay(ctx iris.Context) {
	//baspay 请求
	//更新订单状态，此时具体状态未知等回调。单独线程必须定时扫描已经支付但是无回调的订单
	param := new(api.Pay)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	//要加个锁，或者好方法是 弄个数据库的锁
	tx := db.GDbMgr.Get().Begin()
	status, err := new(models.Trade).GetStatus(tx, *param.MerchantTradeNo)
	if err != nil {
		tx.Rollback()
		ZapLog().Error("UpdateTransferStatusByTradeNo err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}
	if (status.OpenStatus != nil && *status.OpenStatus != models.EUM_OPEN_STATUS_OPEN) || (status.TransferStatus != nil && (*status.TransferStatus != models.EUM_TRANSFER_STATUS_INIT || *status.TransferStatus != models.EUM_TRANSFER_STATUS_FAIL)) {
		tx.Rollback()
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		ZapLog().Error("BASERR_STATE_NOT_ALLOW err")
		return
	}

	err = new(models.Trade).UpdateTransferStatusByTradeNo(tx, *param.MerchantTradeNo, models.EUM_TRANSFER_STATUS_APPLY)
	if err != nil {
		tx.Rollback()
		ZapLog().Error("UpdateTransferStatusByTradeNo err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}
	res, err := new(baspay.Transfer).Parse(param).Send()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		//恢复状态
		//tx ??
		return
	}

	this.Response(ctx, &api.ResPay{MerchantTransferNo: res.MerchantTransferNo, TransferNo: res.TransferNo})
}

//退款
func (this *Trade) ReFund(ctx iris.Context) {
	param := new(api.RefundTrade)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	if param.MerchantId == nil || param.OriginalMerchantTradeNo == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "merchantid nil or original num nil err")
		ZapLog().Error("merchantid nil or original num nil err")
		return
	}

	if err = new(models.ReFund).Parse(param).Add(); err != nil {
		ZapLog().Error("UpdateTransferStatusByTradeNo err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	res, err := new(baspay.RefundTrade).Parse(param).Send()
	if err != nil {
		ZapLog().Error("baspay refund Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}
	this.RefundResponse(ctx, *res)
}

func (this *Trade) FundOut(ctx iris.Context) {
	//提现

	param := new(api.FundOut)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	param.MerchantTradeNo = common.GenerateUuid()

	if err := new(models.FundOut).Parse(param).Add(); err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	res, err := new(baspay.FundOut).Parse(param).Send()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}

	err = new(models.FundOut).UpdateByTradeNo(*res.MerchantFundoutNo, *res.FundoutNo, models.EUM_FUNDOUT_STATUS_APPLY)
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}
	this.Response(ctx, &api.ResFundOut{MerchantFundoutNo: res.MerchantFundoutNo, FundoutNo: res.FundoutNo})
}

func (this *Trade) ListTrade(ctx iris.Context) {
	//查询交易 列表
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
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, res)
}

func (this *Trade) CheckPricePre(deviceNo, asserts string, count int) float64 {

	url := "http://quote.rkuan.com/api/v1/coin/quote?from=" + asserts + "&to=cny"

	result, err := base.HttpSend(url, nil, "GET", nil)

	res := new(ResMsg)

	json.Unmarshal(result, res)
	if err != nil {
		this.ExceptionSerive(nil, apibackend.BASERR_INVALID_PARAMETER.Code(), "unmarshal err")
		ZapLog().Error("unmarshal err", zap.Error(err))
		return 0
	}

	if len(res.Quotes) == 0 {
		this.ExceptionSerive(nil, apibackend.BASERR_INVALID_PARAMETER.Code(), "get price err")
		ZapLog().Error("get price err", zap.Error(err))
		return 0
	}

	price := res.Quotes[0].MoneyInfos[0].Price

	coinPrice := 1 / *price //一块钱等同的币的价值，再拿到后台设置的价格价格，再乘以数量

	configPrice, err := models.GetPrice(deviceNo)
	if err != nil {
		this.ExceptionSerive(nil, apibackend.BASERR_INVALID_PARAMETER.Code(), "get config err")
		ZapLog().Error("get config err", zap.Error(err))
		return 0
	}

	FinalPrice := coinPrice * *configPrice * float64(count)

	//if param.Discount != nil {
	//	discount, err := strconv.ParseFloat(*param.Discount, 64)
	//	if err != nil {
	//		this.ExceptionSerive(nil, apibackend.BASERR_INVALID_PARAMETER.Code(), "string to float err")
	//		ZapLog().Error( "string to float err", zap.Error(err))
	//		return
	//	}
	//	FinalPrice = FinalPrice * discount * 0.1
	//}
	//// 等于多少usd
	//usdPrice := res.Quotes[0].MoneyInfos[1].Price
	//usdShows := Decimal(FinalPrice) * *usdPrice
	//usdShows = Decimal3(usdShows)
	////
	return Decimal(FinalPrice)
}

//退款订单列表
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
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "get refund list err"})
		return
	}

	if len(res) == 0 || len(originOrderList) == 0 {
		ZapLog().Info("this merchant no refund order ")
		ctx.JSON(Response{Code: 0, Message: "no refund"})
		return
	}

	//根据退单里的原始订单查询 原始订单里的金额，币种
	tradeInfos, err := new(models.Trade).GetByOriginNoList(originOrderList)
	if err != nil {
		ZapLog().Error("use refund order origin no get trade info err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "use refund order origin no get trade info err"})
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
