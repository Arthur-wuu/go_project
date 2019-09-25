package common

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
)

const (
	// 成功
	ResponseSuccess = 0
	// 系统错误
	ResponseError = 1
	// 访问超限
	ResponseErrorLimiter = 2
	// 参数错误
	ResponseErrorParams = 2000

	// 鉴权错误
	ResponseErrorToken   = 3000
	ResponseErrorCaptcha = 3001
	ResponseErrorEmail   = 3002
	ResponseErrorPhone   = 3003
	ResponseErrorGa      = 3004

	// 签名验证失败
	ResponseErrorSignatureExpired = 3005
	// 来源IP不合法
	ResponseErrorIllegalIP = 4001
)

// This is default response struct
// swagger:response Response
type Response struct {
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

func NewResponse(ctx iris.Context) *Response {
	return &Response{ctx: ctx}
}

func NewSuccessResponse(ctx iris.Context, result interface{}) *Response {
	r := &Response{ctx: ctx}

	r.Success()
	r.SetMsg("SUCCESS")
	r.Result = result

	return r
}

func NewErrorResponse(ctx iris.Context, result interface{}, msg string, code int) *Response {
	r := &Response{ctx: ctx}

	r.Error(code)
	r.SetMsg(msg)
	r.Result = result

	return r
}

func (r *Response) Success() *Response {
	r.Status.Code = ResponseSuccess
	return r
}

func (r *Response) Error(code int) *Response {
	r.Status.Code = code
	return r
}

func (r *Response) SetMsg(msg string) *Response {
	r.Status.Msg = i18n.Translate(r.ctx, msg)
	if r.Status.Msg == "" {
		r.Status.Msg = msg
	}
	return r
}

func (r *Response) SetMsgWithParams(msg string, params ...interface{}) *Response {
	r.Status.Msg = i18n.Translate(r.ctx, msg, params...)
	if r.Status.Msg == "" {
		r.Status.Msg = msg
	}
	return r
}

func (r *Response) SetResult(result interface{}) *Response {
	r.Result = result
	return r
}

func (r *Response) SetLimitResult(result interface{}, total interface{}, page interface{}) *Response {
	r.Result = struct {
		Total interface{} `json:"total"`
		Page  interface{} `json:"page"`
		Data  interface{} `json:"data"`
	}{total, page, result}
	return r
}
