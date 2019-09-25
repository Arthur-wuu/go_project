package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/config"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Amount struct {
		Controllers
	}

	AmountResponse struct{
		Code       int                 `json:"err,omitempty"`
		Quotes     []QuoteDetailInfo   `json:"quotes,omitempty" doc:"币简称"`
	}

	AmountDetailInfo struct {
		Symbol     string              `json:"symbol,omitempty" doc:"币简称"`
		Id         int                 `json:"id,omitempty" doc:"币id"`
		MoneyInfos []MoneyInfos         `json:"detail,omitempty" doc:"币行情数据"`
	}

	MoneyInfos struct {
		Symbol             *string  `json:"symbol,omitempty" doc:"币简称"`
		Price              float64 `json:"price,omitempty" doc:"最新价"`
		Volume_24h         *float64 `json:"volume_24h,omitempty" doc:"24小时成交量"`
		Market_cap         *float64 `json:"market_cap,omitempty" doc:"总市值"`
		Percent_change_1h  *float64 `json:"percent_change_1h,omitempty" doc:"1小时涨跌幅"`
		Percent_change_24h *float64 `json:"percent_change_24h,omitempty" doc:"24小时涨跌幅"`
		Percent_change_7d  *float64 `json:"percent_change_7d,omitempty" doc:"7天涨跌幅"`
		Last_updated       *int64   `json:"last_updated,omitempty" doc:"最近更新时间"`
	}
)




func (this *Amount) GetAmount (ctx iris.Context) {

	param := new(api.GetAmount)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	coinAmountResult, err := base.HttpSend(config.GConfig.BastionpayUrl.QuoteUrl+"/api/v1/coin/exchange?from="+*param.Legal+"&to="+ *param.Symbol+"&amount="+ *param.Amount, nil,"GET", nil)
	if err != nil {
		ZapLog().Error( "select quote info err", zap.Error(err))
		return
	}
	fmt.Println("coinAmountResult",string(coinAmountResult))

	amountResponse := new(AmountResponse)
	json.Unmarshal(coinAmountResult, amountResponse)
	amountResponse.Quotes[0].MoneyInfos[0].Fee = &config.GConfig.Fee.Coin2cash

	this.Response(ctx, amountResponse.Quotes[0].MoneyInfos)
}

//func  GetAmount (legal, legalNum, assets string)  *float64{
//
//	coinAmountResult, err := base.HttpSend("http://test-quote.rkuan.com/api/v1/coin/exchange?from="+legal+"&to="+ assets+"&amount="+legalNum, nil,"GET", nil)
//	if err != nil {
//		ZapLog().Error( "select quote info err", zap.Error(err))
//		return nil
//	}
//	fmt.Println("coinAmountResult",string(coinAmountResult))
//
//	amountResponse := new(AmountResponse)
//	json.Unmarshal(coinAmountResult, amountResponse)
//
//
//	return amountResponse.Quotes[0].MoneyInfos[0].Price
//}


