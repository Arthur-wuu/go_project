package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/models"
	"encoding/json"
	"fmt"
	"strconv"
	//"BastionPay/merchant-api/db"
	//l4g "github.com/alecthomas/log4go"
	//"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	PrePareParam struct {
		Controllers
	}
)

type PrepareParam struct {
	DeviceNo *string `valid:"required" json:"device_no"`
	Symbol   *string `valid:"required" json:"symbol,omitempty"`
	Count    *int64  `valid:"required" json:"count,omitempty"`
	Discount *string `valid:"optional" json:"discount,omitempty"`
}

func (this *PrePareParam) GetPrice(ctx iris.Context) {
	param := new(PrepareParam)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "create trade param err")
		ZapLog().Error("create trade param err", zap.Error(err))
		return
	}

	url := "http://quote.rkuan.com/api/v1/coin/quote?from=" + *param.Symbol + "&to=cny,usd"

	result, err := base.HttpSend(url, nil, "GET", nil)

	res := new(ResMsg)

	json.Unmarshal(result, res)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "unmarshal err")
		ZapLog().Error("unmarshal err", zap.Error(err))
		return
	}

	if len(res.Quotes) == 0 {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "get price err")
		ZapLog().Error("get price err", zap.Error(err))
		return
	}

	price := res.Quotes[0].MoneyInfos[0].Price

	coinPrice := 1 / *price //一块钱等同的币的价值，再拿到后台设置的价格价格，再乘以数量

	configPrice, err := models.GetPrice(*param.DeviceNo)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "get config err")
		ZapLog().Error("get config err", zap.Error(err))
		return
	}

	FinalPrice := coinPrice * *configPrice * float64(*param.Count)

	if param.Discount != nil {
		discount, err := strconv.ParseFloat(*param.Discount, 64)
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "string to float err")
			ZapLog().Error("string to float err", zap.Error(err))
			return
		}
		FinalPrice = FinalPrice * discount * 0.1
	}
	//// 等于多少usd
	usdPrice := res.Quotes[0].MoneyInfos[1].Price
	usdShows := Decimal(FinalPrice) * *usdPrice
	usdShows = Decimal3(usdShows)
	////

	this.PriceResponse(ctx, Decimal(FinalPrice), usdShows)
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", value), 64)
	return value
}

func Decimal3(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	return value
}
