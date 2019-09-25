package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-api/quote"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	l4g "github.com/alecthomas/log4go"

	"github.com/BastionPay/bas-bkadmin-api/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"net/http"

	"BastionPay/bas-api/apibackend"
	"io"
	"strings"
	"sync"
)

func NewCoinMarket(config *tools.Config) *CoinMarket {
	cm := &CoinMarket{
		config:    config,
		symbolids: make(map[string]int, 0),
	}

	//b, err := ioutil.ReadFile(cm.config.CoinMarket.IdPath)
	//if err != nil {
	//	l4g.Error("ReadFile IdPath:", cm.config.CoinMarket.IdPath, " err:", err.Error())
	//}
	//if (err == nil) && (len(b) != 0) {
	//	param := struct {
	//		Data []struct {
	//			Id     int    `json:"id"`
	//			Symbol string `json:"symbol"`
	//		} `json:"data"`
	//	}{}
	//
	//	if err := json.Unmarshal(b, &param); err == nil {
	//		l4g.Info("Init Id and Symbols")
	//		for _, v := range param.Data {
	//			cm.symbolids[strings.ToLower(v.Symbol)] = v.Id
	//			l4g.Info("symbol[%s]id[%d]", strings.ToLower(v.Symbol), v.Id)
	//		}
	//	}
	//}

	return cm
}

type AdminResponse struct {
	ctx iris.Context
	// response status
	Status struct {
		// response code
		Code int `json:"code"`
		// response msg
		Msg string `json:"msg"`
	} `json:"status"`
	// response result
	Result interface{} `json:"result"`
}

type CoinMarket struct {
	config    *tools.Config
	symbolids map[string]int
	sync.Mutex
}

func (this *CoinMarket) Handler(ctx iris.Context) {
	l4g.Debug("start deal CoinMarket Handler username[%s]", utils.GetValueUserName(ctx))
	param := ctx.Params().Get("param")
	if len(param) == 0 {
		l4g.Error("Context Params get[param] username[%s] err[is nil]", utils.GetValueUserName(ctx))
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: "COINMARKET_GET_PATH_ERRORS: Wrong Path"})
		return
	}
	if len(this.config.BasQuote.Addr) < 3 {
		l4g.Error("redirectCoinQuote username[%s] err[config.BasQuote.Addr not set]", utils.GetValueUserName(ctx))
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code(), Message: "redirectCoinQuote_ERROR:config.BasQuote.Addr not set"})
		return
	}

	var res *Response
	switch param {
	case "ticker":
		this.Ticker(ctx)
		break
	case "getids":
		this.GetIds(ctx)
		break
	case "addids":
		this.UpdateIds(ctx)
		break
	case "delids":
		//		this.DelIds(ctx)
		break
	case "updateids":
		this.UpdateIds(ctx)
		break
	case "kxian":
		this.GetKxian(ctx)
		break
	default:
		l4g.Error("unregister path[%s] redirect BasQuote", param)
		//err, res := this.redirectAdmin(ctx)
		//if err != nil {
		//	ctx.JSON(*res)
		//	return
		//}
		//ctx.JSON(*res)
		ctx.JSON(&Response{Code: apibackend.BASERR_UNSUPPORTED_METHOD.Code(), Message: "CoinQuote_PATH_NOFIND"})
		return
	}
	l4g.Debug("deal CoinMarket Handler username[%s] ok, result[%v]", utils.GetValueUserName(ctx), res)
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
func (cm *CoinMarket) Ticker(ctx iris.Context) {
	l4g.Debug("start deal Ticker username[%s]", utils.GetValueUserName(ctx))
	params := struct {
		Symbols  []string `json:"symbols"`
		Converts []string `json:"converts"`
	}{}
	results := struct {
		Data []interface{} `json:"data"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err:%s", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: "COINMARKET_GET_PARAMS_ERRORS"})
		return
	}

	if len(params.Symbols) == 0 {
		l4g.Error("ReadJSON err:symbols is nil")
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: "COINMARKET_GET_PARAMS_ERRORS"})
		return
	}

	from := strings.Join(params.Symbols, ",")
	to := strings.Join(params.Converts, ",")
	quoteParam := fmt.Sprintf("from=%s&to=%s", from, to)
	err, res := cm.redirectCoinQuote(ctx, cm.config.BasQuote.Addr, "/v1/coin/quote", quoteParam, nil, "GET")
	if err != nil {
		l4g.Error("redirectCoinQuote username[%s] symbol[%s] convert[%s] err[%s]", utils.GetValueUserName(ctx), from, to, err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "redirectCoinQuote_ERROR:" + err.Error()})
		return
	}
	if res == nil {
		l4g.Error("redirectCoinQuote username[%s] err[res is nil]", utils.GetValueUserName(ctx))
		ctx.JSON(&Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: "redirectCoinQuote_ERROR:Res is nil"})
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

	ctx.JSON(Response{Code: 0, Data: results})
	l4g.Debug("deal Ticker username[%s] ok, result[%v]", utils.GetValueUserName(ctx), results)
}

//包含业务处理，不能完全透传
func (this *CoinMarket) redirectAdmin(ctx iris.Context) (error, *Response) {
	//GetBody在服务端无效
	l4g.Debug("start redirectAdmin")
	bodyBytes, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		l4g.Error("Body read username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		return err, &Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: "COINMARKET_READ_BODY_ERROR:" + err.Error()}
	}
	newBody := bytes.NewReader(bodyBytes)

	l4g.Debug("redirectAdmin username[%s] body[%s]", utils.GetValueUserName(ctx), string(bodyBytes))

	// 不用ctx.Redirect()，因为Response结构不一样，还得有业务处理

	adminUrl := this.config.BasAmin.Url + ctx.Path()
	l4g.Debug("url: %s", adminUrl)
	req, err := http.NewRequest(ctx.Method(), adminUrl, newBody)
	if err != nil {
		l4g.Error("http NewRequest username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		return err, &Response{Code: apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code(), Message: "COINMARKET_NewRequest_ERROR:" + err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l4g.Error("http Do username[%s] url[%s] err[%s]", utils.GetValueUserName(ctx), adminUrl, err.Error())
		return err, &Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "COINMARKET_HTTP_DO_ERROR:" + err.Error()}
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l4g.Error("Body readAll username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		return err, &Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "COINMARKET_READ_BODY_ERROR:" + "From admin.api, " + err.Error()}
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		l4g.Error("username[%s] admin.api[%s] response is null", utils.GetValueUserName(ctx), adminUrl)
		return err, &Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "COINMARKET_REDIRECT_ERROR: From admin.api response body content null"}
	}
	if string(content) == "Not Found" {
		l4g.Error("username[%s] admin.api[%s] response is Not Found", utils.GetValueUserName(ctx), adminUrl)
		return errors.New("Not Found"), &Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "COINMARKET_REDIRECT_ERROR:From admin.api response is Not Found"}
	}
	l4g.Debug("admin.api response content[%s]", string(content))
	adminRes := new(AdminResponse)
	if err := json.Unmarshal(content, adminRes); err != nil {
		l4g.Error("Unmarshal username[%s] content[%s] err[%s]", utils.GetValueUserName(ctx), string(content), err.Error())
		return err, &Response{Code: apibackend.BASERR_DATA_UNPACK_ERROR.Code(), Message: "COINMARKET_REDIRECT_ERROR:From admin.api response cannot Unmarshal, " + err.Error()}
	}

	if adminRes.Status.Code != 0 {
		l4g.Error("admin.api username[%s] Response.Status.Code[%s] err[%s]", utils.GetValueUserName(ctx), adminRes.Status.Code, adminRes.Status.Msg)
		return errors.New("from admin.api response code != 0, " + adminRes.Status.Msg), &Response{Code: adminRes.Status.Code, Message: "COINMARKET_REDIRECT_ERROR:From admin.api, " + adminRes.Status.Msg, Data: adminRes.Result}
	}
	ctx.Request().Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

	return nil, &Response{Code: adminRes.Status.Code, Message: adminRes.Status.Msg, Data: adminRes.Result}
}

//包含业务的处理，不能完全透传
func (this *CoinMarket) redirectCoinQuote(ctx iris.Context, addr string, quotePath string, param string, body io.Reader, method string) (error, *quote.ResMsg) {
	//GetBody在服务端无效
	l4g.Debug("start redirectCoinQuote")

	// 不用ctx.Redirect()，因为Response结构不一样，还得有业务处理
	if len(method) == 0 {
		method = "GET"
	}
	adminUrl := addr + quotePath
	if len(param) != 0 {
		adminUrl = adminUrl + "?" + param
	}
	l4g.Debug("url: %s", adminUrl)
	req, err := http.NewRequest(method, adminUrl, body)
	if err != nil {
		l4g.Error("http NewRequest username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		return err, nil
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l4g.Error("http Do username[%s] url[%s] err[%s]", utils.GetValueUserName(ctx), adminUrl, err.Error())
		return err, nil
	}

	if resp.StatusCode != 200 {
		l4g.Error("http Do response username[%s] url[%s] err[%s]", utils.GetValueUserName(ctx), adminUrl, resp.Status)
		return errors.New(resp.Status), nil
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l4g.Error("Body readAll username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		return err, nil
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		l4g.Error("username[%s] admin.api[%s] response is null", utils.GetValueUserName(ctx), adminUrl)
		return errors.New("nil content and must bug"), nil
	}
	if string(content) == "Not Found" {
		l4g.Error("username[%s] admin.api[%s] response is Not Found", utils.GetValueUserName(ctx), adminUrl)
		return errors.New("Not Found"), nil
	}
	l4g.Debug("BasQuote response content[%s]", string(content))

	adminRes := new(quote.ResMsg)
	if err := json.Unmarshal(content, adminRes); err != nil {
		l4g.Error("Unmarshal username[%s] content[%s] err[%s]", utils.GetValueUserName(ctx), string(content), err.Error())
		return err, nil
	}

	if adminRes.Err != 0 {
		errMsg := ""
		if adminRes.ErrMsg != nil {
			errMsg = *adminRes.ErrMsg
		}
		l4g.Error("CoinQuote username[%s] Response.Status.Code[%s] err[%s]", utils.GetValueUserName(ctx), adminRes.Err, errMsg)
		return errors.New("from CoinQuote response code != 0, " + errMsg), nil
	}
	return nil, adminRes
}

type CoinMarketIdSymbol struct {
	Id     int    `json:"id"`
	Symbol string `json:"symbol"`
}

func (this *CoinMarket) AddIds(ctx iris.Context) error {
	l4g.Debug("start deal AddIds username[%s]", utils.GetValueUserName(ctx))
	//	if strings.Count(this.config.BasQuote.Addr, ".")==3 {//ip地址
	err, _ := this.redirectCoinQuote(ctx, this.config.BasQuote.Addr, "/v1/coin/set", "", ctx.Request().Body, "POST")
	if err != nil {
		l4g.Error("redirectCoinQuote username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "redirectCoinQuote_ERROR:" + err.Error()})
		return err
	}
	ctx.Next()
	ctx.JSON(Response{Code: 0, Message: "ok"})
	l4g.Debug("deal AddIds username[%s] ok", utils.GetValueUserName(ctx))
	return nil
	//	}
	//	ipAddr, err := net.LookupHost(this.config.BasQuote.Addr)
	//	if err != nil {
	//		l4g.Error("Domain LookupHost[%s] err[%s]", this.config.BasQuote.Addr, err.Error())
	//		ctx.JSON(&Response{Code: -1, Message: "Domain_LookupHost_ERROR:" + err.Error()})
	//		return err
	//	}
	//	ff :=func(addr string) error {
	//		err, _ := this.redirectCoinQuote(ctx, addr, "/v1/set", "", ctx.Request().Body, "POST")
	//		if err != nil {
	//			l4g.Error("redirectCoinQuote username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
	//			ctx.JSON(&Response{Code: -1, Message: "redirectCoinQuote_ERROR:" + err.Error()})
	//			return err
	//		}
	//		return nil
	//	}
	//	for i:=0; i < len(ipAddr); i++ {//端口默认80
	//		if i < 1 {
	//			if err := ff(ipAddr[i]); err != nil {
	//				return err
	//			}
	//		}else{
	//			go ff(ipAddr[i])
	//		}
	//	}
	ctx.JSON(Response{Code: 0, Message: "ok"})
	l4g.Debug("deal AddIds username[%s] ok", utils.GetValueUserName(ctx))
	return nil
}

//
func (this *CoinMarket) GetIds(ctx iris.Context) error {
	l4g.Debug("start deal GetIds username[%s]", utils.GetValueUserName(ctx))
	symbolStr := ctx.FormValue("symbols")
	err, res := this.redirectCoinQuote(ctx, this.config.BasQuote.Addr, "/v1/coin/code", "symbols="+symbolStr, nil, "GET")
	if err != nil {
		l4g.Error("redirectCoinQuote username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "RedirectCoinQuote_Err:" + err.Error()})
		return err
	}

	if res.Codes == nil {
		ctx.JSON(&Response{Code: apibackend.BASERR_SUCCESS.Code()}) //[]byte切片会进行base64编码
		l4g.Debug("deal GetIds username[%s] ok, result[%s]", utils.GetValueUserName(ctx))
		return nil
	}

	content, err := json.Marshal(res.Codes)
	if err != nil {
		l4g.Error("Marshal username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: "COINMARKET_Marshal_IDS_ERROR:" + err.Error()})
		return err
	}
	ctx.JSON(&Response{Code: 0, Data: string(content)}) //[]byte切片会进行base64编码
	l4g.Debug("deal GetIds username[%s] ok, result[%s]", utils.GetValueUserName(ctx), string(content))
	return nil
}

//history
func (this *CoinMarket) GetKxian(ctx iris.Context) error {
	l4g.Debug("start deal getkxian username[%s]", utils.GetValueUserName(ctx))
	periodStr := ctx.FormValue("period")
	limitStr := ctx.FormValue("limit")
	startStr := ctx.FormValue("start")
	fromStr := ctx.FormValue("from")
	toStr := ctx.FormValue("to")

	urlParam := "period=" + periodStr + "&limit=" + limitStr + "&start=" + startStr + "&from=" + fromStr + "&to=" + toStr
	err, res := this.redirectCoinQuote(ctx, this.config.BasQuote.Addr, "/v1/coin/kxian", urlParam, nil, "GET")
	if err != nil {
		l4g.Error("redirectCoinQuote username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "RedirectCoinQuote_Err:" + err.Error()})
		return err
	}

	if res.Historys == nil {
		ctx.JSON(&Response{Code: 0}) //[]byte切片会进行base64编码
		l4g.Debug("deal getkxian username[%s] ok, result[%s]", utils.GetValueUserName(ctx))
		return nil
	}

	content, err := json.Marshal(res.Historys)
	if err != nil {
		l4g.Error("Marshal username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: "COINMARKET_Marshal_IDS_ERROR:" + err.Error()})
		return err
	}
	ctx.JSON(&Response{Code: 0, Data: string(content)}) //[]byte切片会进行base64编码
	l4g.Debug("deal getkxian username[%s] ok, result[%s]", utils.GetValueUserName(ctx), string(content))
	return nil
}

//没有异步消息通知系统所以只能广播一遍，第一次成功即可，如果其中某次失败，一致性让bas-quote定时更新码表即可
func (this *CoinMarket) UpdateIds(ctx iris.Context) error {
	l4g.Debug("start deal UpdateIds username[%s]", utils.GetValueUserName(ctx))
	//	if  strings.Count(this.config.BasQuote.Addr, ".")==3  {
	err, _ := this.redirectCoinQuote(ctx, this.config.BasQuote.Addr, "/v1/coin/set", "", ctx.Request().Body, "POST")
	if err != nil {
		l4g.Error("redirectCoinQuote username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "redirectCoinQuote_ERROR:" + err.Error()})
		return err
	}
	ctx.Next()
	ctx.JSON(Response{Code: 0, Message: "ok"})
	l4g.Debug("deal UpdateIds username[%s] ok", utils.GetValueUserName(ctx))
	return nil
	//	}
	//	ipAddr, err := net.LookupHost(this.config.BasQuote.Addr)
	//	if err != nil {
	//		l4g.Error("Domain LookupHost[%s] err[%s]", this.config.BasQuote.Addr, err.Error())
	//		ctx.JSON(&Response{Code: -1, Message: "Domain_LookupHost_ERROR:" + err.Error()})
	//		return err
	//	}
	//	var netWait sync.WaitGroup
	//
	//	ff :=func(addr string, err *error) {
	//		netWait.Add(1)
	//		defer netWait.Done()
	//		*err, _ = this.redirectCoinQuote(ctx, addr, "/v1/set", "", ctx.Request().Body, "POST")
	//	}
	//	errArr := make([]error, len(ipAddr))
	//	for i:=0; i < len(ipAddr); i++ {//端口默认80,只设置成功一次就行了
	//		go ff(ipAddr[i],&errArr[i])
	//	}
	//	netWait.Wait()
	//	for i:=0; i < len(errArr); i++ {
	//		if errArr[i] != nil {
	//			l4g.Error("redirectCoinQuote username[%s][%s] err[%s]", utils.GetValueUserName(ctx), ipAddr[i], errArr[i].Error())
	//			ctx.JSON(&Response{Code: -1, Message: "redirectCoinQuote_ERROR:" + ipAddr[i]+" "+errArr[i].Error()})
	//		}
	//	}
	ctx.JSON(Response{Code: 0, Message: "ok"})
	l4g.Debug("deal UpdateIds username[%s] ok", utils.GetValueUserName(ctx))
	return nil
}
