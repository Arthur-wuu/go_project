package controllers

import (
	"BastionPay/bas-api/apibackend"
	"fmt"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-admin-api/models"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/bugsnag/bugsnag-go/errors"
	"github.com/jinzhu/gorm"
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
)

type VerificationController struct {
	redis     *common.Redis
	db        *gorm.DB
	config    *config.Config
	userModel *models.UserModel
}

func NewVerificationController(redis *common.Redis, db *gorm.DB, config *config.Config) *VerificationController {
	return &VerificationController{redis: redis, db: db, config: config, userModel: models.NewUserModel(db)}
}

func (v *VerificationController) Send(ctx iris.Context) {
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
		ctx.JSON(common.NewErrorResponse(ctx, nil, "UNSUPPORTED_VERIFICATION_TYPE", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	if !common.InArray(OperatingType, operating) {
		ctx.JSON(common.NewErrorResponse(ctx, nil, "UNSUPPORTED_OPERATING_TYPE", apibackend.BASERR_INVALID_PARAMETER.Code()))
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
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 如果在未登录情况下请求短信或邮件，需要检测图形验证码，防止暴力请求接口
	if user.UserId == 0 && (vType == common.VerificationTypeEmail || vType == common.VerificationTypeSms) {
		captchaPass, err := common.NewVerification(v.redis, operating, common.VerificationTypeCaptcha).Check(captchaToken, 0, "")
		if err != nil || !captchaPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_CAPTCHA_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
			return
		}
	}

	// 获得接收者
	if user.UserId != 0 && operating != "bind_email" && operating != "bind_phone" {
		if vType == common.VerificationTypeEmail {
			recipient, err = v.userModel.GetEmail(user.UserId)
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("GetEmail err")
				ctx.JSON(common.NewErrorResponse(ctx, nil, "CAN_NOT_GET_EMAIL", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
				return
			}
		} else if vType == common.VerificationTypeSms {
			recipient, err = v.userModel.GetPhoneRecipient(user.UserId)
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("GetPhoneRecipient err")
				ctx.JSON(common.NewErrorResponse(ctx, nil, "CAN_NOT_GET_PHONE", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
				return
			}
		}
	}

	// 获得ga token
	if vType == common.VerificationTypeGa {
		if user.UserId == 0 {
			ZapLog().Info("AUTHENTICATION_FAILED err")
			ctx.JSON(common.NewErrorResponse(ctx, nil, "AUTHENTICATION_FAILED", apibackend.BASERR_UNKNOWN_BUG.Code()))
			return
		}
		ga, err = v.userModel.GetGa(user.UserId)
		if err != nil || ga == "" {
			ZapLog().Info("CAN_NOT_GET_GA err")
			ctx.JSON(common.NewErrorResponse(ctx, nil, "CAN_NOT_GET_GA", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
			return
		}
	}

	// 生成验证操作
	verification := common.NewVerification(v.redis, operating, vType)
	switch vType {
	case common.VerificationTypeEmail:
		//if recipient == "" {
		//	ctx.JSON(common.NewResponse(ctx).Error(common.ResponseErrorParams).
		//		SetMsgWithParams("MISSING_PARAMETERS", "recipient"))
		//	return
		//}
		id, err = v.sendEmail(ctx, verification, user.UserId, recipient)
		break
	case common.VerificationTypeSms:
		//if recipient == "" {
		//	ctx.JSON(common.NewResponse(ctx).Error(common.ResponseErrorParams).
		//		SetMsgWithParams("MISSING_PARAMETERS", "recipient"))
		//	return
		//}
		id, err = v.sendSms(ctx, verification, user.UserId, recipient)
		break
	case common.VerificationTypeCaptcha:
		id, captcha, err = v.sendCaptcha(ctx, verification, user.UserId)
		break
	case common.VerificationTypeGa:
		id, err = v.sendGa(ctx, verification, user.UserId, ga)
		break
	}
	if err != nil {
		return
	}

	if vType == common.VerificationTypeCaptcha {
		ctx.JSON(common.NewSuccessResponse(ctx, struct {
			Id      string `json:"id"`
			Captcha string `json:"captcha"`
		}{id, captcha}))
	} else {
		ctx.JSON(common.NewSuccessResponse(ctx, struct {
			Id string `json:"id"`
		}{id}))
	}
}

func (v *VerificationController) sendEmail(ctx iris.Context, verification *common.Verification, userId uint, recipient string) (string, error) {
	muchFlag, expireAt, err := v.overLimitMail(ctx, userId)
	if err != nil {
		ZapLog().With(zap.String("recipient", recipient), zap.Error(err), zap.Uint("userId", userId), zap.String("ip", v.getRemoteAddr(ctx))).Error("overLimitMail err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_DATABASE_ERROR.Code()).
			SetMsg("SEND_FAILED"))
		return "", err
	}
	if muchFlag {
		ZapLog().With(zap.String("recipient", recipient), zap.Uint("userId", userId), zap.String("ip", v.getRemoteAddr(ctx)), zap.Int("timelimit", expireAt)).Error("overLimitMail send too much email")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_OPERATE_FREQUENT.Code()).
			SetMsg(fmt.Sprintf("Send_Too_Much_Email_%ds", expireAt)))
		return "", errors.New("send too much email", 0)
	}
	language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
	id, err := verification.
		GenerateEmail(userId, recipient, v.config.Bas_notify.VerifyCodeMailTmp, language)
	if err != nil {
		//		glog.Error(err.Error())
		ZapLog().With(zap.Error(err)).Error("GetString err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SEND_FAILED", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return "", err
	}
	return id, nil
}

func (v *VerificationController) sendSms(ctx iris.Context, verification *common.Verification, userId uint, recipient string) (string, error) {
	muchFlag, expireAt, err := v.overLimitSms(ctx, userId)
	if err != nil {
		ZapLog().With(zap.String("recipient", recipient), zap.Error(err), zap.Uint("userId", userId), zap.String("ip", v.getRemoteAddr(ctx))).Error("overLimitSms err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_DATABASE_ERROR.Code()).
			SetMsg("SEND_FAILED"))
		return "", err
	}
	if muchFlag {
		ZapLog().With(zap.String("recipient", recipient), zap.Uint("userId", userId), zap.String("ip", v.getRemoteAddr(ctx)), zap.Int("timelimit", expireAt)).Error("overLimitSMS send too much sms")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_OPERATE_FREQUENT.Code()).
			SetMsg(fmt.Sprintf("Send_Too_Much_SMS_%ds", expireAt)))
		return "", errors.New("send too much sms", 0)
	}
	language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
	id, err := verification.
		GenerateSms(userId, recipient, v.config.Bas_notify.VerifyCodeSmsTmp, language)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetString err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SEND_FAILED", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return "", err
	}
	return id, nil
}

func (v *VerificationController) sendGa(ctx iris.Context, verification *common.Verification, userId uint, secret string) (string, error) {
	id, err := verification.GenerateGA(userId, secret)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GenerateGA err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_GA_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return "", err
	}
	return id, nil
}

func (v *VerificationController) sendCaptcha(ctx iris.Context, verification *common.Verification, userId uint) (string, string, error) {
	id, captcha, err := verification.GenerateCaptcha(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GenerateCaptcha err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_CAPTCHA_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return "", "", err
	}
	return id, captcha, nil
}

func (v *VerificationController) Verification(ctx iris.Context) {
	var (
		params = struct {
			Id        string `json:"id"`
			Value     string `json:"value"`
			Recipient string `json:"recipient"`
		}{}
		err error
	)

	user, _ := common.GetUserIdFromCtxUnsafe(ctx)
	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Id == "", "id"},
		{params.Value == "", "value"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	verification := common.NewVerification(v.redis, "", "")

	bol, err := verification.Verify(params.Id, user.UserId, params.Value, params.Recipient)
	if err != nil || bol == false {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("Verify err")
			//			glog.Error(err.Error())
		}

		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_VERIFICATION", apibackend.BASERR_ADMIN_INCORRECT_VERIFYCODE.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
}

func (this *VerificationController) overLimitSms(ctx iris.Context, userId uint) (bool, int, error) {
	if userId == 0 {
		ip := this.getRemoteAddr(ctx)
		for i := 1; i < len(this.config.Limits.IpSms); i += 2 {
			arr := this.config.Limits.IpSms
			if arr[i] <= 0 || arr[i-1] <= 0 {
				continue
			}
			key := fmt.Sprintf("limits_sms_ip_%d_%s", arr[i], ip)
			ZapLog().With(zap.String("key", key)).Info("overLimitSms")
			muchFlag, err := this.overLimit(key, arr[i], arr[i-1])
			if err != nil {
				return true, 0, err
			}
			if muchFlag {
				return true, arr[i], nil
			}
		}
		return false, 0, nil
	}
	//已登录用户
	for i := 1; i < len(this.config.Limits.IdSms); i += 2 {
		arr := this.config.Limits.IdSms
		if arr[i] <= 0 || arr[i-1] <= 0 {
			continue
		}
		key := fmt.Sprintf("limits_sms_id_%d_%d_", arr[i], userId)
		ZapLog().With(zap.String("key", key)).Info("overLimitSms")
		muchFlag, err := this.overLimit(key, arr[i], arr[i-1])
		if err != nil {
			return true, 0, err
		}
		if muchFlag {
			return true, arr[i], nil
		}
	}
	return false, 0, nil
}

func (this *VerificationController) overLimitMail(ctx iris.Context, userId uint) (bool, int, error) {
	if userId == 0 {
		ip := this.getRemoteAddr(ctx)
		for i := 1; i < len(this.config.Limits.IpMail); i += 2 {
			arr := this.config.Limits.IpMail
			if arr[i] <= 0 || arr[i-1] <= 0 {
				continue
			}
			key := fmt.Sprintf("limits_mail_ip_%d_%s", arr[i], ip)
			ZapLog().With(zap.String("key", key)).Debug("overLimitMail")
			muchFlag, err := this.overLimit(key, arr[i], arr[i-1])
			if err != nil {
				return true, 0, err
			}
			if muchFlag {
				return true, arr[i], nil
			}
		}
		return false, 0, nil
	}
	//已登录用户
	for i := 1; i < len(this.config.Limits.IdMail); i += 2 {
		arr := this.config.Limits.IdMail
		if arr[i] <= 0 || arr[i-1] <= 0 {
			continue
		}
		key := fmt.Sprintf("limits_mail_id_%d_%d_", arr[i], userId)
		ZapLog().With(zap.String("key", key)).Info("overLimitMail")
		muchFlag, err := this.overLimit(key, arr[i], arr[i-1])
		if err != nil {
			return true, 0, err
		}
		if muchFlag {
			return true, arr[i], nil
		}
	}
	return false, 0, nil
}

func (this *VerificationController) overLimit(key string, expireTime, limitNum int) (bool, error) {
	if expireTime <= 0 || limitNum <= 0 {
		return false, nil
	}
	result, err := this.redis.Do("INCR", key)
	if err != nil {
		//		ZapLog().With(zap.Error(err)).Error("redis INCR err")
		return true, err
	}
	if result == nil {
		return true, errors.New("redis Do result is nil", 0)
	}
	res := int(result.(int64))
	if res == 1 {
		_, err = this.redis.Do("EXPIRE", key, expireTime)
		if err != nil {
			return true, err
		}
	}
	if res > limitNum {
		ZapLog().With(zap.String("key", key), zap.Int("res", res), zap.Int("limit", limitNum)).Error("overlimit")
		return true, nil
	}
	return false, nil
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
