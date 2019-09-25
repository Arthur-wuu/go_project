package controllers

import (
	"github.com/kataras/iris"
	"BastionPay/bas-notify/models"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
)

type Sms struct{
	Controllers
}

func (this * Sms) Send(ctx iris.Context) {
	param := new(models.SmsMsg)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}
	if flag := this.ReqNotifyMsgIsValid(param); !flag {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	errCode, err := param.Send(true)
	if errCode != 0 {
		this.ExceptionSerive(ctx, errCode, err.Error())
		ZapLog().Error( "Send err", zap.Error(err))
		return
	}
	this.Response(ctx, nil)
}

func (this * Sms) Sends(ctx iris.Context) {
	param := make([]*models.SmsMsg, 0)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	res:= make([]*Response, 0, len(param))
	for i:=0; i < len(param); i++ {
		if flag := this.ReqNotifyMsgIsValid(param[i]); !flag {
			//this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
			ZapLog().Error( "param err", zap.Error(err))
			res =append(res, &Response{
				Err : apibackend.BASERR_INVALID_PARAMETER.Code(),
				ErrMsg: "param err",
			})
			continue
		}

		errCode, err := param[i].Send(true)
		if errCode != 0 {
			//this.ExceptionSerive(ctx, errCode, err.Error())
			ZapLog().Error( "send err", zap.Error(err))
			//continue
		}
		res =append(res, &Response{
			Err : errCode,
			ErrMsg: err.Error(),
		})

	}

	ctx.JSON(res)

}

func (this *Sms) ReqNotifyMsgIsValid(req *models.SmsMsg) bool {
	//if req.Recipient == nil {
	//	return false
	//}
	//if req.TempId != nil {
	//	return true
	//}
	//if req.TempAlias != nil {
	//	return true
	//}
	if (req.GroupId == nil) &&(req.GroupName == nil) /*&&(req.GroupAlias == nil)*/ {
		return false
	}
	if req.Lang == nil {
		return false
	}
	return true
}