package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/pay-user-merchant-api/api"
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/config"
	"BastionPay/pay-user-merchant-api/models"
	"fmt"
	"github.com/bugsnag/bugsnag-go/errors"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

var (
	VerificationType = []string{
		common.VerificationTypeEmail,
		common.VerificationTypeSms,
		common.VerificationTypeCaptcha,
		common.VerificationTypeGa,
	}
	OperatingType = []string{
		"login",
		"register",
		"forget_password",
		"withdrawal",
		"withdrawal_address",
		"trading",
		"bind_ga",
		"unbind_ga",
		"bind_email",
		"bind_phone",
		"rebind_phone",
	}
	phoneSmsLimiter   *common.BusLimiter
	ipSmsLimiter      *common.BusLimiter
	phoneEmailLimiter *common.BusLimiter
	ipEmailLimiter    *common.BusLimiter
)

type VerificationController struct {
	Controllers
}

func NewVerificationController() *VerificationController {
	return &VerificationController{}
}

func (this *VerificationController) Send(ctx iris.Context) {
	var (
		id           string
		vType        string
		operating    string
		captcha      string
		recipient    string
		captchaToken string
		ga           string
		err          error
	)

	user, _ := common.GetUserIdFromCtxUnsafe(ctx)

	vType = ctx.Params().Get("type")
	operating = ctx.URLParam("operating")
	recipient = ctx.URLParam("recipient")
	captchaToken = ctx.URLParam("captcha_token")

	if !common.InArray(VerificationType, vType) {
		ZapLog().Error("unsopported VerificationType", zap.String("vType", vType))
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "UNSUPPORTED_VERIFICATION_TYPE")
		return
	}

	if !common.InArray(OperatingType, operating) {
		ZapLog().Error("unsopported OperatingType", zap.String("operating", operating))
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "UNSUPPORTED_OPERATING_TYPE")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{user.UserId == 0 && (vType == common.VerificationTypeEmail || vType == common.VerificationTypeSms) && recipient == "",
			"recipient"},
		{(operating == "bind_email" || operating == "bind_phone") && vType == common.VerificationTypeSms && recipient == "",
			"recipient"},
		{user.UserId == 0 && (vType == common.VerificationTypeEmail || vType == common.VerificationTypeSms) && captchaToken == "",
			"captcha_token"}})
	if pv != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS")
		ZapLog().Error("CheckParams MISSING_PARAMETERS err", zap.String("error", pv.ErrMsg))
		return
	}

	// 如果在未登录情况下请求短信或邮件，需要检测图形验证码，防止暴力请求接口
	if user.UserId == 0 && (vType == common.VerificationTypeEmail || vType == common.VerificationTypeSms) {
		captchaPass, err := common.NewVerification(operating, common.VerificationTypeCaptcha).Check(captchaToken, 0, "")
		if err != nil || !captchaPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
			}
			this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_CAPTCHA_TOKEN_AUTHENTICATION")
			return
		}
	}

	// 获得接收者
	if user.UserId != 0 && operating != "bind_email" && operating != "bind_phone" {
		if vType == common.VerificationTypeEmail {
			recipient, err = new(models.User).GetEmail(int64(user.UserId))
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("GetEmail err")
				this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "CAN_NOT_GET_EMAIL")
				return
			}
		} else if vType == common.VerificationTypeSms {
			recipient, err = new(models.User).GetPhoneRecipient(int64(user.UserId))
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("GetPhoneRecipient err")
				this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "CAN_NOT_GET_PHONE")
				return
			}
		}
	}

	// 获得ga token
	if vType == common.VerificationTypeGa {
		if user.UserId == 0 {
			ZapLog().Info("AUTHENTICATION_FAILED err")
			this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), "AUTHENTICATION_FAILED")
			return
		}
		ga, err = new(models.User).GetGa(int64(user.UserId))
		if err != nil || ga == "" {
			ZapLog().Info("CAN_NOT_GET_GA err")
			this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "CAN_NOT_GET_GA")
			return
		}
	}

	// 生成验证操作
	verification := common.NewVerification(operating, vType)
	switch vType {
	case common.VerificationTypeEmail:
		//if recipient == "" {
		//	ctx.JSON(common.NewResponse(ctx).Error(common.ResponseErrorParams).
		//		SetMsgWithParams("MISSING_PARAMETERS", "recipient"))
		//	return
		//}
		id, err = this.sendEmail(ctx, verification, user.UserId, recipient)
		break
	case common.VerificationTypeSms:
		//if recipient == "" {
		//	ctx.JSON(common.NewResponse(ctx).Error(common.ResponseErrorParams).
		//		SetMsgWithParams("MISSING_PARAMETERS", "recipient"))
		//	return
		//}
		id, err = this.sendSms(ctx, verification, user.UserId, recipient)
		break
	case common.VerificationTypeCaptcha:
		id, captcha, err = this.sendCaptcha(ctx, verification, user.UserId)
		break
	case common.VerificationTypeGa:
		id, err = this.sendGa(ctx, verification, user.UserId, ga)
		break
	}
	if err != nil {
		return
	}

	if vType == common.VerificationTypeCaptcha {
		this.Response(ctx, &api.ResVerificationSend{&id, &captcha})
	} else {
		this.Response(ctx, &api.ResVerificationSend{&id, nil})
	}
}

func (this *VerificationController) sendEmail(ctx iris.Context, verification *common.Verification, userId uint, recipient string) (string, error) {
	muchFlag, rateStr, err := this.overLimitMail(ctx, userId)
	if err != nil {
		ZapLog().With(zap.String("recipient", recipient), zap.Error(err), zap.Uint("userId", userId), zap.String("ip", this.getRemoteAddr(ctx))).Error("overLimitMail err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "SEND_FAILED")
		return "", err
	}
	if muchFlag {
		ZapLog().With(zap.String("recipient", recipient), zap.Uint("userId", userId), zap.String("ip", this.getRemoteAddr(ctx)), zap.String("timelimit", rateStr)).Error("overLimitMail send too much email")
		this.ExceptionSerive(ctx, apibackend.BASERR_OPERATE_FREQUENT.Code(), fmt.Sprintf("Send_Too_Much_Email_%s", rateStr))
		return "", errors.New("send too much email", 0)
	}
	language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
	id, err := verification.
		GenerateEmail(userId, recipient, config.GConfig.BasNotify.VerifyCodeMailTmp, language)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetString err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), "SEND_FAILED")
		return "", err
	}
	return id, nil
}

func (this *VerificationController) sendSms(ctx iris.Context, verification *common.Verification, userId uint, recipient string) (string, error) {
	muchFlag, rateStr, err := this.overLimitSms(ctx, userId)
	if err != nil {
		ZapLog().With(zap.String("recipient", recipient), zap.Error(err), zap.Uint("userId", userId), zap.String("ip", this.getRemoteAddr(ctx))).Error("overLimitSms err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "SEND_FAILED")
		return "", err
	}
	if muchFlag {
		ZapLog().With(zap.String("recipient", recipient), zap.Uint("userId", userId), zap.String("ip", this.getRemoteAddr(ctx)), zap.String("timelimit", rateStr)).Error("overLimitSMS send too much sms")
		this.ExceptionSerive(ctx, apibackend.BASERR_OPERATE_FREQUENT.Code(), fmt.Sprintf("Send_Too_Much_SMS_%s", rateStr))
		return "", errors.New("send too much sms", 0)
	}
	language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
	id, err := verification.
		GenerateSms(userId, recipient, config.GConfig.BasNotify.VerifyCodeSmsTmp, language)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetString err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), "SEND_FAILED")
		return "", err
	}
	return id, nil
}

func (this *VerificationController) sendGa(ctx iris.Context, verification *common.Verification, userId uint, secret string) (string, error) {
	id, err := verification.GenerateGA(userId, secret)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GenerateGA err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "GET_GA_FAILED")
		return "", err
	}
	return id, nil
}

func (this *VerificationController) sendCaptcha(ctx iris.Context, verification *common.Verification, userId uint) (string, string, error) {
	id, captcha, err := verification.GenerateCaptcha(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GenerateCaptcha err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "GET_CAPTCHA_FAILED")
		return "", "", err
	}
	return id, captcha, nil
}

func (this *VerificationController) Verification(ctx iris.Context) {

	params := new(api.Verification)
	user, _ := common.GetUserIdFromCtxUnsafe(ctx)
	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Id == "", "id"},
		{params.Value == "", "value"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS")
		return
	}

	verification := common.NewVerification("", "")

	bol, err := verification.Verify(params.Id, user.UserId, params.Value, params.Recipient)
	if err != nil || bol == false {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("Verify err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INCORRECT_VERIFYCODE.Code(), "FAILED_TO_VERIFICATION")
		return
	}

	this.Response(ctx, nil)
}

func (this *VerificationController) overLimitSms(ctx iris.Context, userId uint) (bool, string, error) {
	if userId == 0 {
		limitFlag, rateStr, err := ipSmsLimiter.Check(this.getRemoteAddr(ctx))
		if err != nil {
			ZapLog().With(zap.Error(err), zap.String("rateStr", rateStr)).Error("redis err")
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
			return true, rateStr, err
		}
		if limitFlag {
			ZapLog().With(zap.Error(err)).Error("ip sms limit err")
			this.ExceptionSerive(ctx, apibackend.BASERR_OPERATE_FREQUENT.Code(), apibackend.BASERR_OPERATE_FREQUENT.Desc())
			return true, rateStr, err
		}
		return false, rateStr, nil
	}
	//已登录用户
	limitFlag, rateStr, err := phoneSmsLimiter.Check(fmt.Sprintf("%d", userId))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return true, rateStr, err
	}
	if limitFlag {
		ZapLog().With(zap.Error(err)).Error("phone sms limit err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OPERATE_FREQUENT.Code(), apibackend.BASERR_OPERATE_FREQUENT.Desc())
		return true, rateStr, err
	}
	return false, rateStr, nil
}

func (this *VerificationController) overLimitMail(ctx iris.Context, userId uint) (bool, string, error) {
	if userId == 0 {
		limitFlag, rateStr, err := ipEmailLimiter.Check(this.getRemoteAddr(ctx))
		if err != nil {
			ZapLog().With(zap.Error(err), zap.String("rateStr", rateStr)).Error("redis err")
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
			return true, rateStr, err
		}
		if limitFlag {
			ZapLog().With(zap.Error(err)).Error("ip email limit err")
			this.ExceptionSerive(ctx, apibackend.BASERR_OPERATE_FREQUENT.Code(), apibackend.BASERR_OPERATE_FREQUENT.Desc())
			return true, rateStr, err
		}
		return false, rateStr, nil
	}
	//已登录用户
	limitFlag, rateStr, err := phoneEmailLimiter.Check(fmt.Sprintf("%d", userId))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return true, rateStr, err
	}
	if limitFlag {
		ZapLog().With(zap.Error(err)).Error("phone email limit err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OPERATE_FREQUENT.Code(), apibackend.BASERR_OPERATE_FREQUENT.Desc())
		return true, rateStr, err
	}
	return false, rateStr, nil
}

func (this *VerificationController) getRemoteAddr(ctx iris.Context) string {
	ZapLog().With(zap.String("X-Forwarded-For", ctx.GetHeader("X-Forwarded-For")), zap.String("remoteAddr", ctx.RemoteAddr())).Debug("remote addr")
	ip := ctx.GetHeader("X-Forwarded-For")
	if len(ip) < 3 {
		ip = ctx.RemoteAddr()
	}
	ipArr := strings.Split(ip, ",")
	if len(ipArr) > 1 {
		ip = ipArr[0]
	}
	ipArr2 := strings.Split(ip, ":")
	if len(ipArr2) > 1 {
		ip = ipArr2[0]
	}
	ip = strings.TrimSpace(ip)
	return ip
}
