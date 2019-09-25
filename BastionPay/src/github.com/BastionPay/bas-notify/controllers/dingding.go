package controllers

import (
	"github.com/kataras/iris"
	"BastionPay/bas-notify/models"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/config"
)

type DingDing struct{
	Controllers
}

func (this * DingDing) Send(ctx iris.Context) {
	param := new(models.DingDingMsg)

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

func (this * DingDing) GetQuns(ctx iris.Context) {
	quns := make([]string, 0, len(config.GConfig.DingDing))
	for i:=0; i < len(config.GConfig.DingDing); i++ {
		quns = append(quns, config.GConfig.DingDing[i].QunName)
	}
	this.Response(ctx, quns)
}


//
//func (this * DingDing) Sends(ctx iris.Context) {
//	param := make([]*models.SmsMsg, 0)
//
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
//		ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//
//	res:= make([]*Response, 0, len(param))
//	for i:=0; i < len(param); i++ {
//		if flag := this.ReqNotifyMsgIsValid(param[i]); !flag {
//			//this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
//			ZapLog().Error( "param err", zap.Error(err))
//			res =append(res, &Response{
//				Err : apibackend.BASERR_INVALID_PARAMETER.Code(),
//				ErrMsg: "param err",
//			})
//			continue
//		}
//
//		errCode, err := param[i].Send(true)
//		if errCode != 0 {
//			//this.ExceptionSerive(ctx, errCode, err.Error())
//			ZapLog().Error( "send err", zap.Error(err))
//			//continue
//		}
//		res =append(res, &Response{
//			Err : errCode,
//			ErrMsg: err.Error(),
//		})
//
//	}
//
//	ctx.JSON(res)
//
//}

func (this *DingDing) ReqNotifyMsgIsValid(req *models.DingDingMsg) bool {
	//if req.Recipient == nil {
	//	return false
	//}
	//if req.TempId != nil {
	//	return true
	//}
	//if req.TempAlias != nil {
	//	return true
	//}
	//if (req.GroupId == nil) &&(req.GroupName == nil)&&(req.GroupAlias == nil) {
	//	return false
	//}
	if req.Lang == nil {
		return false
	}
	return true
}