package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/base"
	"encoding/json"
	//"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Quote struct {
		Controllers
	}

	QuoteResponse struct{
		Code       int                 `json:"err,omitempty"`
		Quotes     []QuotDetailInfo   `json:"quotes,omitempty" doc:"币简称"`
	}

	QuotDetailInfo struct {
		Symbol     string              `json:"symbol,omitempty" doc:"币简称"`
		Id         int                 `json:"id,omitempty" doc:"币id"`
		MoneyInfos []QuoteInfos         `json:"detail,omitempty" doc:"币行情数据"`
	}

	QuoteInfos struct {
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

const QuoteUrl  = "http://quote.rkuan.com/api/v1/coin/quote?from="




func (this *Quote) GetQuote (ctx iris.Context) {

	param := new(api.Quote)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	quoteResult, err := base.HttpSend(QuoteUrl+*param.Symbol+"&to="+ *param.Legal, nil,"GET", nil)
	if err != nil {
		ZapLog().Error( "select quote info err", zap.Error(err))
		return
	}
	ZapLog().Info("quoteResult", zap.Any("quoteResult", quoteResult))

	quoteResponse := new(QuoteResponse)
	json.Unmarshal(quoteResult, quoteResponse)

	this.Response(ctx, quoteResponse.Quotes[0].MoneyInfos)
}



