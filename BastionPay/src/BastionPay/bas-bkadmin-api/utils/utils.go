package utils

import (
	"errors"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/alecthomas/log4go"
	"github.com/kataras/iris"
	"runtime/debug"
	"sync"
)

type (
	UserUtils struct {
	}
)

var (
	u    *UserUtils
	once sync.Once
)

func NewUtils() *UserUtils {
	once.Do(func() {
		u = &UserUtils{}
	})

	return u
}

func (u *UserUtils) GetValueUserInfo(ctx iris.Context) (*models.Account, error) {

	body := ctx.Values().Get(ctx.GetHeader("token"))
	if body == nil {
		return nil, errors.New("get user info is nil")
	}

	return body.(*models.Account), nil
}

func (u *UserUtils) GetValueUserId(ctx iris.Context) int64 {

	body := ctx.Values().Get(ctx.GetHeader("token"))
	if body == nil {
		return 0
	}

	return body.(*models.Account).Id
}

func GetValueUserName(ctx iris.Context) string {
	body := ctx.Values().Get(ctx.GetHeader("token"))
	if body == nil {
		return ""
	}
	if body == nil {
		return ""
	}
	acc, ok := body.(*models.Account)
	if !ok {
		return ""
	}
	return acc.Name
}

func PanicPrint() {
	if err := recover(); err != nil {
		log4go.Error("panic err[%v] stack[\n%v\n]", err, string(debug.Stack()))
	}
}
