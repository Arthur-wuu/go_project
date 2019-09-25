package controllers

import (
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	"github.com/asaskevich/govalidator"
)

type (
	Controllers struct {
		UserId int64
		Config *tools.Config
	}

	Response struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const SEARCHSIZE = 10000

var (
	Tools    *common.Tools
	User     *utils.UserUtils
	Validate *govalidator.Validator
)

func init() {
	Tools = common.New()
	User = utils.NewUtils()
}

func NewResponse(code int, msg string) *Response {
	return &Response{
		Code:    code,
		Message: msg,
	}
}

func (r *Response) SetLimitResult(result interface{}, total interface{}, page interface{}) *Response {
	r.Data = struct {
		Total interface{} `json:"total"`
		Page  interface{} `json:"page"`
		Data  interface{} `json:"data"`
	}{total, page, result}
	return r
}
