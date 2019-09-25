package common

import (
	"BastionPay/bas-tv-proxy/api"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-tv-proxy/base"
	"runtime/debug"
	"unicode"
	"regexp"
)

func  PackPushMsg(reqster base.Requester, apiMsg interface{}) ([]byte, error) {
	apiResponse := api.NewResponse(reqster.GetQid())
	apiResponse.SetCounter(reqster.GetCounter())

	content, err := apiResponse.Marshal(apiMsg)
	if err != nil {
		ZapLog().Error("apiResponse Marshal err", zap.Error(err))
		return nil,err
	}
	return content, nil
}

func  PackResMsg(reqster base.Requester, errCode int32, apiMsg interface{}) ([]byte, error) {
	apiResponse := api.NewResponse(reqster.GetQid())
	apiResponse.SetErr(errCode)
	content, err := apiResponse.Marshal(apiMsg)
	if err != nil {
		ZapLog().Error("apiResponse Marshal err", zap.Error(err))
		return nil,err
	}
	return content, nil
}

func CtxJson(ctx iris.Context, qid string, apiMsg *api.MSG) error {
	response := api.NewResponse(qid)
	content, err := response.Marshal(apiMsg)
	if err != nil {
		ZapLog().Error("apiResponse Marshal err", zap.Error(err))
		return err
	}
	_, err =  ctx.Write(content)
	if err != nil {
		ZapLog().Error("ctx.Write err", zap.Error(err))
	}
	return err
}

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

