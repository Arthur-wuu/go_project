package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	"github.com/alecthomas/log4go"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
)

func NewDownloadController(config *tools.Config) *DownloadController {
	bp := &DownloadController{
		config: config,
	}
	return bp
}

type DownloadController struct {
	config *tools.Config
}

func (this *DownloadController) GetStatus(ctx iris.Context) {
	if len(this.config.BasQuote.Addr) < 5 { //暂时把这个功能放在 quote-collect中
		log4go.Error("config not set, host[%s]", this.config.BasQuote.Addr)
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code()))
		return
	}

	url := this.config.BasQuote.Addr + "/v1/filetransfer/status/get" + "?" + ctx.Request().URL.RawQuery
	resp, err := http.Get(url)
	if err != nil {
		log4go.Error("http get[%s] err[%s]", url, err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code()))
		return
	}

	if resp.StatusCode != 200 {
		log4go.Error("http response[%s] err[%s]", resp.Status, err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log4go.Error("Body readAll  err[%s]", err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_SERVICE_UNKNOWN_ERROR.Code()))
		return
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		log4go.Error("download response is null url[%s]", url)
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}
	//if string(content) == "Not Found" {
	//	ZapLog().With(zap.String("url", url)).Error("download response is Not Found")
	//	ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
	//	return
	//}
	resMsg := &struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			QueryID   string `json:"query_id"`
			Status    int    `json:"status"`
			Data      string `json:"data"`
			UserParam string `json:"user_param,omitempty"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(content, resMsg); err != nil {
		log4go.Error("json Unmarshal[%s] err[%s]", string(content), err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_DATA_UNPACK_ERROR.Code()))
		return
	}
	if resMsg.Code != 0 {
		log4go.Error("download response code[%s][%v]", resMsg.Code, resMsg)
		ctx.JSON(common.NewErrorResponse(ctx, nil, fmt.Sprintf("%d-%s", resMsg.Code, resMsg.Message), apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	//ZapLog().With(zap.String("content", string(content))).Debug("download response content")

	ctx.JSON(common.NewSuccessResponse(ctx, resMsg.Data))
	ctx.Next()
}
