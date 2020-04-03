package controllers

import (
	"BastionPay/pay-user-merchant-api/common"
	//"BastionPay/pay-user-merchant-api/config"
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-api/apibackend/v1/backend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/pay-user-merchant-api/api"
	"BastionPay/pay-user-merchant-api/models"
	"encoding/json"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type InfoController struct {
	Controllers
}

func NewInfoController() *InfoController {
	return &InfoController{}
}

func (this *InfoController) GetInformation(ctx iris.Context) {
	var (
		bindGa bool
		userId uint
	)

	userId = common.GetUserIdFromCtx(ctx)

	user, err := new(models.User).GetById(int64(userId))
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Uint("userId", userId)).Error("GetUserById err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "SERVER_INTERNAL_ERROR")
		return
	}

	if user.Ga == "" {
		bindGa = false
	} else {
		bindGa = true
	}

	this.Response(ctx, &api.ResUserInfo{user.Id, user.CreateTime, common.SecretPhone(user.Phone), user.PhonneDistrict,
		common.SecretEmail(user.Email), bindGa, user.VipLevel, user.Country,
		user.Language, "", "",
		user.Company})

}

func (this *InfoController) GetInformationNoHide(ctx iris.Context) {
	var (
		bindGa bool
		userId uint
	)

	userId = common.GetUserIdFromCtx(ctx)

	user, err := new(models.User).GetById(int64(userId))
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Uint("userId", userId)).Error("GetUserById err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "SERVER_INTERNAL_ERROR")
		return
	}

	if user.Ga == "" {
		bindGa = false
	} else {
		bindGa = true
	}

	this.Response(ctx, &api.ResUserInfoNoHide{user.Id, user.Phone, user.PhonneDistrict, user.Email, bindGa})
}

func (this *InfoController) SetInformation(ctx iris.Context) {
	params := new(api.UserInfoSet)
	userId := common.GetUserIdFromCtx(ctx)

	err := ctx.ReadJSON(params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	if params.Language != "" {
		if err = new(models.User).SetLanguage(userId, params.Language); err != nil {
			ZapLog().With(zap.Error(err)).Error("SetLanguage err")
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "FAILED_TO_SET_LANGUAGE")
			return
		}
	}

	if params.Timezone != "" {
		if err = new(models.User).SetTimezone(userId, params.Timezone); err != nil {
			ZapLog().With(zap.Error(err)).Error("SetTimezone err")
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "FAILED_TO_SET_TIMEZONE")
			return
		}
	}

	this.Response(ctx, nil)
}

func (this *InfoController) BindEmail(ctx iris.Context) {

	params := new(api.BindEmail)

	userId := common.GetUserIdFromCtx(ctx)

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Email == "", "email"},
		{params.EmailToken == "", "email_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
		return
	}

	// 检测邮箱验证码
	tokenPass, err := common.NewVerification("bind_email", common.VerificationTypeEmail).
		Check(params.EmailToken, userId, params.Email)
	if err != nil || !tokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION")
		return
	}

	user, err := new(models.User).GetById(int64(userId))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "GET_USER_ERROR")
		return
	}
	if user.Email != "" {
		ZapLog().Error("YOU_HAVE_ALREADY_BIND_EMAIL err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "YOU_HAVE_ALREADY_BIND_EMAIL")
		return
	}

	// 检测邮箱是否绑定过
	exists, err := new(models.User).UserExisted(params.Email)
	if err != nil || exists {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("UserExisted err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "THE_EMAIL_HAS_ALREADY_EXISTED")
		return
	}

	if err = new(models.User).SetEmail(userId, params.Email); err != nil {
		ZapLog().With(zap.Error(err)).Error("SetEmail err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "FAILED_TO_BIND_THE_EMAIL")
		return
	}

	this.Response(ctx, nil)
	//ctx.Next()
}

func (this *InfoController) BindPhone(ctx iris.Context) {
	params := new(api.BindPhone)
	userId := common.GetUserIdFromCtx(ctx)

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Phone == "", "phone"},
		{params.CountryCode == "", "country_code"},
		{params.SmsToken == "", "sms_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
		return
	}

	// 检测手机
	tokenPass, err := common.NewVerification("bind_phone", common.VerificationTypeSms).
		Check(params.SmsToken, userId, params.CountryCode+params.Phone)
	if err != nil || !tokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_SMS_TOKEN_AUTHENTICATION")
		return
	}

	user, err := new(models.User).GetById(int64(userId))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "GET_USER_ERROR")
		return
	}
	if user.Phone != "" {
		ZapLog().Error("YOU_HAVE_ALREADY_BIND_PHONE err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code(), "YOU_HAVE_ALREADY_BIND_PHONE")
		return
	}

	// 检测手机是否绑定过
	exists, err := new(models.User).UserExisted(params.Phone)
	if err != nil || exists {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("UserExisted err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "THE_PHONE_HAS_ALREADY_EXISTED")
		return
	}

	if err = new(models.User).SetPhone(userId, params.Phone, params.CountryCode); err != nil {
		ZapLog().With(zap.Error(err)).Error("SetPhone err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "FAILED_TO_BIND_THE_PHONE")
		return
	}

	this.Response(ctx, nil)
	//ctx.Next()
}

func (this *InfoController) RebindPhone(ctx iris.Context) {

	params := new(api.RebindPhone)
	userId := common.GetUserIdFromCtx(ctx)

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Phone == "", "phone"},
		{params.CountryCode == "", "country_code"},
		{params.OldSmsToken == "", "old_sms_token"},
		{params.NewSmsToken == "", "new_sms_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
		return
	}

	user, err := new(models.User).GetById(int64(userId))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "GET_USER_ERROR")
		return
	}
	if user.Phone == "" {
		ZapLog().Error("YOU_HAVE_NOT_BIND_PHONE err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code(), "YOU_HAVE_NOT_BIND_PHONE")
		return
	}
	if user.Phone == params.Phone {
		ZapLog().Error("PLEASE_ENTER_A_DIFFERENT_PHONE err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_DATA_SAME.Code(), "PLEASE_ENTER_A_DIFFERENT_PHONE")
		return
	}
	//if user.Email == "" {
	//	ZapLog().Error("YOU_HAVE_NOT_BIND_EMAIL err")
	//	ctx.JSON(common.NewErrorResponse(ctx, nil, "YOU_HAVE_NOT_BIND_EMAIL", common.ResponseError))
	//	return
	//}

	pv = common.CheckParams(&[]common.ParamsVerification{
		{user.Email != "" && params.EmailToken == "", "email_token"},
		{user.Ga != "" && params.GaToken == "", "ga_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
		return
	}

	// 检测邮箱验证码
	if user.Email != "" {
		tokenPass, err := common.NewVerification("rebind_phone", common.VerificationTypeEmail).
			Check(params.EmailToken, userId, user.Email)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
			}
			this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION")
			return
		}
	}

	// 检测GA验证码
	if user.Ga != "" {
		tokenPass, err := common.NewVerification("rebind_phone", common.VerificationTypeGa).
			Check(params.GaToken, userId, user.Ga)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
			}
			this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_GA_TOKEN_AUTHENTICATION")
			return
		}
	}

	// 检测旧手机
	oldTokenPass, err := common.NewVerification("rebind_phone", common.VerificationTypeSms).
		Check(params.OldSmsToken, userId, user.PhonneDistrict+user.Phone)
	if err != nil || !oldTokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_OLD_SMS_TOKEN_AUTHENTICATION")
		return
	}

	// 检测新手机
	newTokenPass, err := common.NewVerification("bind_phone", common.VerificationTypeSms).
		Check(params.NewSmsToken, userId, params.CountryCode+params.Phone)
	if err != nil || !newTokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_NEW_SMS_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
		return
	}

	// 检测手机是否绑定过
	exists, err := new(models.User).UserExisted(params.Phone)
	if err != nil || exists {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("UserExisted err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "THE_PHONE_HAS_ALREADY_EXISTED")
		return
	}

	if err = new(models.User).SetPhone(userId, params.Phone, params.CountryCode); err != nil {
		ZapLog().With(zap.Error(err)).Error("SetPhone err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "FAILED_TO_BIND_THE_PHONE")
		return
	}

	this.Response(ctx, nil)
	//ctx.Next()
}

//func (this *InfoController) GetUserInfo(ctx iris.Context) {
//	userkeys := ctx.FormValue("userkey")
//	userkeysNoSpace := strings.Replace(userkeys, " ", "", len(userkeys))
//	userkeyArr := strings.Split(userkeysNoSpace, ",")
//
//	resArr := make([]*api.ResUserInfo, 0, len(userkeyArr))
//	for i := 0; i < len(userkeyArr); i++ {
//		if len(userkeyArr[i]) == 0 {
//			continue
//		}
//		user, err := new(models.User).GetByUserkey(userkeyArr[i])
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("GetUserByUserkey err")
//			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
//			return
//		}
//		if user == nil {
//			continue
//		}
//
//		bindGa := true
//		if user.Ga == "" {
//			bindGa = false
//		}
//
//		res := &api.ResUserInfo{
//			Id: user.Id,
//			CreatedAt: user.CreatedAt,
//			Phone: user.Phone,
//			CountryCode: user.CountryCode,
//			Email: user.Email,
//			BindGa: bindGa,
//			VipLevel: user.VipLevel,
//			Citizenship: user.Citizenship,
//			Language: user.Language,
//			Timezone: user.Timezone,
//			RegistrationType: user.RegistrationType,
//			UserKey: userkeyArr[i],
//			CompanyName: user.CompanyName,
//		}
//
//
//		ZapLog().With(zap.Int("userid", *user.Id)).Info("")
//		resArr = append(resArr, res)
//	}
//
//
//	content, err := json.Marshal(resArr)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("Marshal err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATA_PACK_ERROR.Code(), "Json_marshal_ERROR")
//		return
//	}
//	this.Response(ctx, string(content))
//}

func (this *InfoController) Listusers(ctx iris.Context) {
	reqUserList := new(api.UserList)

	err := ctx.ReadJSON(reqUserList)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	var (
		pageNum   = reqUserList.MaxDispLines
		totalLine = reqUserList.TotalLines
		pageIndex = reqUserList.PageIndex

		beginIndex = 0
	)

	if totalLine == 0 {
		//totalLine, err = db.ListUserCount()
		totalLine, err = new(models.User).ListUserCountByBasic(reqUserList.Condition)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("ListUserCountByBasic err")
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "SERVER_INTERNAL_ERROR")
			return
		}
	}

	if pageNum < 1 || pageNum > 100 {
		pageNum = 50
	}

	beginIndex = pageNum * (pageIndex - 1)

	//ackUserList, err := db.ListUsers(beginIndex, pageNum)
	userBasics, err := new(models.User).ListUsersByBasic(beginIndex, pageNum, reqUserList.Condition)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ListUsersByBasic err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "SERVER_INTERNAL_ERROR")
		return
	}

	ackUserList := new(backend.AckUserList)

	ackUserList.PageIndex = pageIndex
	ackUserList.MaxDispLines = pageNum
	ackUserList.TotalLines = totalLine
	ackUserList.Data = userBasics

	// to ack
	dataAck, err := json.Marshal(ackUserList)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("jsonMarshal err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATA_PACK_ERROR.Code(), "JSON_MUSHAL_ERROR")
		return
	}

	this.Response(ctx, string(dataAck))
}
