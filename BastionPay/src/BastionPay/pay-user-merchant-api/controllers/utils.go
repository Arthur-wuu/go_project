package controllers

import (
	"BastionPay/pay-user-merchant-api/api"
	"fmt"
	"github.com/kataras/iris"
	"math/rand"
	"time"
)

func GetAppUserInfo(ctx iris.Context) (*api.BkLoginUserInfo, error) {
	appClaimsInterface := ctx.Values().Get(CONST_CTX_APPUSERINFO)
	if appClaimsInterface == nil {
		return nil, fmt.Errorf("nil " + CONST_CTX_APPUSERINFO)
	}

	return appClaimsInterface.(*api.BkLoginUserInfo), nil
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
