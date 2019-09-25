package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-api/apibackend"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func NewDownloadController(config *config.Config) *DownloadController {
	bp := &DownloadController{
		config: config,
	}
	return bp
}

type DownloadController struct {
	config *config.Config
}

func (this *DownloadController) GetStatus(ctx iris.Context) {
	if len(this.config.Bas_quote.Addr) < 5 { //暂时把这个功能放在 quote-collect中
		ZapLog().With(zap.String("host", this.config.Bas_quote.Addr)).Info("config not set")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code()))
		return
	}

	url := this.config.Bas_quote.Addr + "/v1/filetransfer/status/get" + "?" + ctx.Request().URL.RawQuery
	resp, err := http.Get(url)
	if err != nil {
		ZapLog().With(zap.String("url", url), zap.Error(err)).Error("http get err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code()))
		return
	}

	if resp.StatusCode != 200 {
		ZapLog().With(zap.String("status", resp.Status)).Error("http response err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Body readAll  err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_SERVICE_UNKNOWN_ERROR.Code()))
		return
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		ZapLog().With(zap.String("url", url)).Error("download response is null")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}
	if string(content) == "Not Found" {
		ZapLog().With(zap.String("url", url)).Error("download response is Not Found")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}
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
		ZapLog().With(zap.String("content", string(content)), zap.Error(err)).Error("json Unmarshal err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_DATA_UNPACK_ERROR.Code()))
		return
	}
	if resMsg.Code != 0 {
		ZapLog().With(zap.String("url", url)).Error("download response is Not Found")
		ctx.JSON(common.NewErrorResponse(ctx, nil, fmt.Sprintf("%d-%s", resMsg.Code, resMsg.Message), apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	ZapLog().With(zap.String("content", string(content))).Debug("download response content")

	ctx.JSON(common.NewSuccessResponse(ctx, resMsg.Data))
}
