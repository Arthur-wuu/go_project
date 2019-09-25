package controllers

import (
	"BastionPay/bas-api/apibackend"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-api/quote"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type CoinMarketController struct {
	config    *config.Config
	symbolids map[string]int
	sync.Mutex
}

func NewCoinMarketController(config *config.Config) *CoinMarketController {
	cm := &CoinMarketController{
		config:    config,
		symbolids: make(map[string]int),
	}

	// TODO: read from bastion database
	// this is read from json file
	//b, err := ioutil.ReadFile(cm.config.CoinMarket.IdPath)
	//if err == nil {
	//	param := struct {
	//		Data []struct {
	//			Id     int    `json:"id"`
	//			Symbol string `json:"symbol"`
	//		} `json:"data"`
	//	}{}
	//
	//	if err := json.Unmarshal(b, &param); err == nil {
	//		for _, v := range param.Data {
	//			cm.symbolids[strings.ToLower(v.Symbol)] = v.Id
	//		}
	//	}
	//}

	return cm
}

type QuoteDetailInfoData struct {
	QuoteDetailInfo QuoteDetailInfo `json:"data"`
}

type QuoteDetailInfo struct {
	Symbol     *string                     `json:"symbol,omitempty"`
	Id         *int                        `json:"id,omitempty"`
	MoneyInfos map[string]*quote.MoneyInfo `json:"quotes"`
}

//协议转换，兼容
func (cm *CoinMarketController) Ticker(ctx iris.Context) {
	params := struct {
		Symbols  []string `json:"symbols"`
		Converts []string `json:"converts"`
	}{}
	results := struct {
		Data []interface{} `json:"data"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "COINMARKET_GET_PARAMS_ERRORS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	if len(params.Symbols) == 0 {
		ZapLog().Error("ReadJSON err:symbols is nil")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "COINMARKET_GET_PARAMS_ERRORS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	_, err = common.JwtParse(cm.config.Token.Secret, cm.config.Token.Expiration, ctx.GetHeader("Authorization"))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtParse err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "AUTHENTICATION_FAILED", apibackend.BASERR_TOKEN_INVALID.Code()))
		return
	}

	from := strings.Join(params.Symbols, ",")
	to := strings.Join(params.Converts, ",")
	quoteParam := fmt.Sprintf("from=%s&to=%s", from, to)
	err, res := cm.redirectCoinQuote(ctx, "/v1/coin/quote", quoteParam, nil, "GET")
	if err != nil {
		ZapLog().With(zap.Error(err), zap.String("symbol", from), zap.String("convert", to)).Error("getTicker  err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BAS_QUOTE_ERRORS", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}
	if res == nil {
		ZapLog().With(zap.Error(err)).Error("redirectCoinQuote res is nil")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BAS_QUOTE_ERRORS", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	for i := 0; i < len(res.Quotes); i++ {
		quoteDetailInfoData := new(QuoteDetailInfoData)
		quoteDetailInfoData.QuoteDetailInfo.MoneyInfos = make(map[string]*quote.MoneyInfo)
		q := &quoteDetailInfoData.QuoteDetailInfo
		q.Id = res.Quotes[i].Id
		q.Symbol = res.Quotes[i].Symbol
		for j := 0; j < len(res.Quotes[i].MoneyInfos); j++ {
			q.MoneyInfos[res.Quotes[i].MoneyInfos[j].GetSymbol()] = res.Quotes[i].MoneyInfos[j]
		}
		results.Data = append(results.Data, quoteDetailInfoData)
	}

	ctx.JSON(common.NewSuccessResponse(ctx, results))
	return
}

func (this *CoinMarketController) redirectCoinQuote(ctx iris.Context, quotePath string, param string, body io.Reader, method string) (error, *quote.ResMsg) {
	if len(this.config.Bas_quote.Addr) < 5 {
		ZapLog().Info("config not set BasQuote")
		return errors.New("config not set BasQuote"), nil
	}
	//GetBody在服务端无效
	ZapLog().Debug("start redirectCoinQuote")

	// 不用ctx.Redirect()，因为Response结构不一样，还得有业务处理
	if len(method) == 0 {
		method = "GET"
	}
	adminUrl := this.config.Bas_quote.Addr + quotePath
	if len(param) != 0 {
		adminUrl = adminUrl + "?" + param
	}
	ZapLog().Sugar().Debug("url: %s", adminUrl)
	req, err := http.NewRequest(method, adminUrl, body)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("http NewRequest  err")
		return err, nil
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("http Do  err")
		return err, nil
	}

	if resp.StatusCode != 200 {
		ZapLog().With(zap.String("status", resp.Status)).Error("http response err")
		return errors.New(resp.Status), nil
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Body readAll  err")
		return err, nil
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		ZapLog().With(zap.String("url", adminUrl)).Error(" admin.api response is null")
		return err, nil
	}
	if string(content) == "Not Found" {
		ZapLog().With(zap.String("usrl", adminUrl)).Error("username[%s] admin.api[%s] response is Not Found")
		return errors.New("Not Found"), nil
	}
	ZapLog().With(zap.String("content", string(content))).Debug("BasQuote response content")

	adminRes := new(quote.ResMsg)
	if err := json.Unmarshal(content, adminRes); err != nil {
		ZapLog().With(zap.Error(err), zap.String("content", string(content))).Error("Unmarshal username[%s] content[%s] err[%s]")
		return err, nil
	}

	if adminRes.Err != 0 {
		errMsg := ""
		if adminRes.ErrMsg != nil {
			errMsg = *adminRes.ErrMsg
		}
		ZapLog().With(zap.Int("err", adminRes.Err), zap.String("error", errMsg)).Error("CoinQuote Response.Status.Code err")
		return errors.New("from CoinQuote response code != 0, " + errMsg), nil
	}
	return nil, adminRes
}
