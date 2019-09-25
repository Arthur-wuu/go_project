package controllers

import (
	"BastionPay/bas-notify/common"
	"github.com/kataras/iris/context"
	//"gopkg.exa.center/blockshine-ex/api-article/config"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/models/table"
	"go.uber.org/zap"
	"runtime/debug"
)

type (
	Controllers struct {
	}

	Response struct {
		Err                 int         `json:"err"`
		ErrMsg              string      `json:"errmsg"`
		TemplateGroupList   interface{} `json:"templategrouplist,omitempty"`
		Templates           interface{} `json:"template,omitempty"`
		TemplateHistoryList interface{} `json:"templatehistorylist,omitempty"`
	}
)

var (
	Tools *common.Tools
)

func init() {
	Tools = common.New()
}

func (c *Controllers) ResponseGroupList(
	ctx context.Context,
	data interface{}) {

	ctx.JSON(
		Response{
			Err:               0,
			ErrMsg:            "Success",
			TemplateGroupList: data,
		})
}

func (c *Controllers) ResponseTemps(
	ctx context.Context,
	data interface{}) {

	ctx.JSON(
		Response{
			Err:       0,
			ErrMsg:    "Success",
			Templates: data,
		})
}

func (c *Controllers) ResponseHistoryList(
	ctx context.Context,
	data interface{}) {

	ctx.JSON(
		Response{
			Err:                 0,
			ErrMsg:              "Success",
			TemplateHistoryList: data,
		})
}

func (c *Controllers) ExceptionSerive(
	ctx context.Context,
	code int,
	message string) {

	ctx.JSON(
		Response{
			Err:    code,
			ErrMsg: message,
		})
}

func (c *Controllers) Response(
	ctx context.Context,
	data interface{}) {

	ctx.JSON(
		Response{
			Err:    0,
			ErrMsg: "Success",
		})
}

type TemplateGroupList struct {
	Total_lines    int                    `json:"total_lines"`
	Page_index     int                    `json:"page_index"`
	Max_disp_lines int                    `json:"max_disp_lines"`
	TemplateGroups []*table.TemplateGroup `json:"templategroup,omitempty"`
}

type TemplateHistoryList struct {
	Total_lines      int              `json:"total_lines"`
	Page_index       int              `json:"page_index"`
	Max_disp_lines   int              `json:"max_disp_lines"`
	TemplateHistorys []*table.History `json:"templatehistory,omitempty"`
}

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}
