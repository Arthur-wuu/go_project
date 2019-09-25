package controllers

import (
	"BastionPay/bas-quote/common"
	"BastionPay/bas-quote/config"
	"BastionPay/bas-quote/db"
	//"fmt"
	. "BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
)

type (
	Inspection struct {
	}
    Response struct {
	   Code    int         `json:"code"`
	   Message interface{} `json:"message"`
	   Data    interface{} `json:"data"`
})

func (this *Inspection) VerifyIsBlock(ctx iris.Context) {
	termBlock := common.NewTermBlocker(&db.GRedis, config.GConfig.TermBlock,  ctx.Method()+ctx.Path(), "", ctx)
	isBlock, err := termBlock.IsBlock()
	if err !=nil {
		ZapLog().With(zap.Error(err)).Error("termBlocker isBlock err")
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: apibackend.BASERR_DATABASE_ERROR.Desc()})
		return
	}else if isBlock {
		ctx.JSON(Response{Code: apibackend.BASERR_BLOCK_ACCOUNT.Code(), Message: apibackend.BASERR_BLOCK_ACCOUNT.Desc()})
		return
	}

	tbResper, err := termBlock.Done(false)
	//fmt.Println("tbResper ",tbResper.OpenFlag,tbResper.OnBlock)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("termBlocker done ")
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: apibackend.BASERR_DATABASE_ERROR.Desc()})
		return
	}
	//if !tbResper.OpenFlag {
	//	ctx.JSON(Response{Code: 1, Message: "user is locked not open flag "})
	//	return
	//}

	if  tbResper.OnBlock {
		ctx.JSON(Response{Code: apibackend.BASERR_OPERATE_FREQUENT.Code(), Message: apibackend.BASERR_OPERATE_FREQUENT.Desc()})
		return
	}

	ctx.Next()
}







