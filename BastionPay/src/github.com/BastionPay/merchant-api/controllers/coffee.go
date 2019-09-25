package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/common"
	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/config"
	"encoding/json"
	"strings"

	//"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strconv"
)
//获取行情价格，咖啡机设置的人民币价格转到对应币种的价格
const quoteUrl  = "http://quote.rkuan.com/api/v1/coin/quote?from="

type (
	CoffeeTrade struct {
		Controllers
	}
)


//咖啡二维码订单
func (this *CoffeeTrade) CoffeeQr(ctx iris.Context) {
	coffeeParams := new(CoffeeParams)
	err := Tools.ShouldBindQuery(ctx, coffeeParams)
	if err != nil {
		ZapLog().Error( " form param err", zap.Error(err))
		ctx.JSON(Response{Code: 1001, Message: err.Error()})
		return
	}
	name      := coffeeParams.Name
	price     := coffeeParams.Price
	channelid := coffeeParams.Channelid
	orderId   := coffeeParams.Orderid

	fmt.Println("name,price,id",name,price,channelid,orderId)
	//验签


	//需要的参数
	var Assets string
	var Amount string
	var PayId string
	var MerchantId string
	var Remark string
	var ProductDtail string
	var ExpireTime int64

	//通道是36的  usdt   可以在这加，也可以配置文件的形式

	priceFloat, _ := strconv.ParseFloat(price, 64)
	//Amount = strconv.FormatFloat(priceFloat / 100, 'g', -1, 64)
	MerchantId = config.GConfig.PayeeId.MerchantId    //咖啡机商户
	PayId = config.GConfig.PayeeId.PayId              //咖啡机收款人
	ProductDtail = "coffee"
	Remark = "remark coffee"
	if channelid == "36"{
		Assets = "USDT"
		Amount = GetPrice(Assets, priceFloat / 100)

	}
	if channelid == "37"{
		Assets = "BTC"
		Amount = GetPrice(Assets, priceFloat / 100)

	}
	if channelid == "38"{
		Assets = "OKG"
		Amount = GetPrice(Assets, priceFloat / 100)

	}
	if channelid == "39"{
		Assets = "SHINE"
		Amount = GetPrice(Assets, priceFloat / 100)

	}
	ZapLog().Info("amount:", zap.Any("amount:", Amount))

	//请求open-api的参数
	uuid := common.GenerateUuid()
	ExpireTime = 900
	param := new(api.CoffeeQrTrade)

	param.MerchantTradeNo = uuid
	param.ExpireTime = ExpireTime
	param.Assets = &Assets
	param.Amount = &Amount
	param.ProductName = &name
	param.MerchantId = &MerchantId
	param.PayeeId = &PayId
	param.ProductDetail = &ProductDtail
	param.Remark = &Remark

	res, err := new(baspay.CoffeeQrTrade).CoffeeQrParse(param).SendCoffeeQr()
	if err != nil {
		ZapLog().Error( "baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}
	if res.Code != 0 {
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "create qr order err"})
		return
	}

	coffeeStruct := new(CoffeeStruct)
	coffeeStruct.Code = "1"
	coffeeStruct.Msg = "msg content"
	coffeeStruct.Orderid = orderId
	coffeeStruct.Torderid = uuid
	coffeeStruct.Twocode = res.Data.Qr_code

	fmt.Println("here coffee response", *coffeeStruct)

	this.CoffeeResponse(ctx, *coffeeStruct)
}



//查询coffee订单状态
func (this *CoffeeTrade) OrderStatus (ctx iris.Context) {
	params := new(CoffeeStatusParams)
	Tools.ShouldBindQuery(ctx, params)
	torderid := params.Torderid
	orderId := params.Orderid
	MerchantId := "2"

	param := new(api.TradeInfo)
	param.MerchantId = &MerchantId
	param.MerchantTradeNo = &torderid

	res,err := new(baspay.TradeInfo).Parse(param).SendCoffee()
	if err != nil {
		ZapLog().Error( "select trade info err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}
	if res != "3"{
		return
	}

	coffeeStatus := new(CoffeeStatus)
	coffeeStatus.Code = "1"
	coffeeStatus.Msg = "msg content"
	coffeeStatus.Orderid = orderId
	coffeeStatus.Torderid = torderid

	this.CoffeeStatusRes(ctx, *coffeeStatus)
}


type CoffeeStruct struct {
	Orderid  string  `json:"orderid,omitempty"`
	Torderid string  `json:"torderid,omitempty"`
	Code     string  `json:"code,omitempty"`
	Msg      string  `json:"msg,omitempty"`
	Twocode  string  `json:"twocode,omitempty"`
}

//查询结果
type CoffeeStatus struct {
	Orderid  string  `json:"orderid,omitempty"`
	Torderid string  `json:"torderid,omitempty"`
	Code     string  `json:"code"`
	Msg      string  `json:"msg,omitempty"`
}

//创建订单
type CoffeeParams struct {
	Ver             string     `json:"ver" form:"ver"`
	Orderid         string     `json:"orderid" form:"orderid"`
	Machid          string     `json:"machid" form:"machid"`
	Trackno         string     `json:"trackno" form:"trackno"`
	Name            string     `json:"name" form:"name"`
	Price           string     `json:"price" form:"price"`
	Channelid       string     `json:"channelid" form:"channelid"`
	Randstr         string     `json:"randstr" form:"randstr"`
	Timestamp       string     `json:"timestamp" form:"timestamp"`
	Sign            string     `json:"sign" form:"sign"`
}

//查询订单状态的
type CoffeeStatusParams struct {
	Ver             string     `json:"ver" form:"ver"`
	Orderid         string     `json:"orderid" form:"orderid"`
	Torderid        string     `json:"torderid" form:"torderid"`
	Machid          string     `json:"machid" form:"machid"`
	Channelid       string     `json:"channelid" form:"channelid"`
	Randstr         string     `json:"randstr" form:"randstr"`
	Timestamp       string     `json:"timestamp" form:"timestamp"`
	Sign            string     `json:"sign" form:"sign"`
}

type BaspayRes struct {
	Merchant_trade_no  string  `json:"merchant_trade_no,omitempty"`
	Qr_code            string  `json:"qr_code,omitempty"`
}




func  GetPrice(symbol string, amount float64) string {

	url := quoteUrl+strings.ToUpper(symbol)+"&to=cny"

	result, err := base.HttpSend(url, nil, "GET", nil)

	res := new(ResMsg)

	json.Unmarshal(result, res)
	if err != nil {
		ZapLog().Error( "unmarshal err", zap.Error(err))
		return ""
	}

	if len(res.Quotes) == 0  {
		ZapLog().Error( "get price err", zap.Error(err))
		return ""
	}

	price := res.Quotes[0].MoneyInfos[0].Price

	coinPrice := 1 / *price   //一块钱等同的币的价值，再拿到后台设置的价格价格，再乘以数量

	FinalPrice := coinPrice * amount

	priceString := strconv.FormatFloat(FinalPrice, 'f', -1, 64)


	return  priceString
}

