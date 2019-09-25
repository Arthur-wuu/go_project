package controllers

import (
	"BastionPay/bas-api/apibackend"
	"encoding/json"
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-admin-api/models"
	"github.com/BastionPay/bas-api/admin"
	"github.com/BastionPay/bas-api/apibackend/v1/backend"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

type InfoController struct {
	redis     *common.Redis
	db        *gorm.DB
	config    *config.Config
	userModel *models.UserModel
}

func NewInfoController(redis *common.Redis, db *gorm.DB, config *config.Config) *InfoController {
	models.GlobalBasModel.Init(config)
	return &InfoController{
		redis: redis,
		db:    db, config: config,
		userModel: models.NewUserModel(db),
	}
}

func (i *InfoController) GetInformation(ctx iris.Context) {
	var (
		bindGa bool
		userId uint
	)

	userId = common.GetUserIdFromCtx(ctx)

	user, err := i.userModel.GetUserById(userId)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Uint("userId", userId)).Error("GetUserById err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	if user.Ga == "" {
		bindGa = false
	} else {
		bindGa = true
	}

	userAllStatus, err := models.GlobalBasModel.GetUserAllAccountStatus(user.Uuid)
	if err != nil { //为了兼容，暂时不处理
		ZapLog().With(zap.Error(err), zap.String("userkey", user.Uuid)).Error("getUserAuditeStatus err")
		//		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", common.ResponseError))
		//		return
	}

	remainLvLimits, err := common.GetAllLevelRemainLimits(i.redis, i.config.LevelPathLimits, user.Uuid, 0)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.String("userkey", user.Uuid)).Error("GetAllLevelRemainLimits err")
	}

	ctx.JSON(common.NewSuccessResponse(ctx, &struct {
		Id               uint        `json:"id"`
		CreatedAt        int64       `json:"created_at"`
		Phone            string      `json:"phone"`
		CountryCode      string      `json:"country_code"`
		Email            string      `json:"email"`
		BindGa           bool        `json:"bind_ga"`
		VipLevel         uint8       `json:"vip_level"`
		Citizenship      string      `json:"citizenship"`
		Language         string      `json:"language"`
		Timezone         string      `json:"timezone"`
		RegistrationType string      `json:"registration_type"`
		UserKey          string      `json:"user_key"`
		CompanyName      string      `json:"company_name"`
		AuditeStatus     uint        `json:"audite_status"`
		TransferStatus   uint        `json:"transfer_status"`
		LevelLimits      interface{} `json:"level_limits,omitempty"`
	}{user.ID, user.CreatedAt, common.SecretPhone(user.Phone), user.CountryCode,
		common.SecretEmail(user.Email), bindGa, user.VipLevel, user.Citizenship,
		user.Language, user.Timezone, user.RegistrationType,
		user.Uuid, user.CompanyName, userAllStatus.AuditeStatus, userAllStatus.TransferStatus, remainLvLimits}))
}

func (i *InfoController) GetInformationNoHide(ctx iris.Context) {
	var (
		bindGa bool
		userId uint
	)

	userId = common.GetUserIdFromCtx(ctx)

	user, err := i.userModel.GetUserById(userId)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Uint("userId", userId)).Error("GetUserById err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	if user.Ga == "" {
		bindGa = false
	} else {
		bindGa = true
	}

	ctx.JSON(common.NewSuccessResponse(ctx, &struct {
		Id          uint   `json:"id"`
		Phone       string `json:"phone"`
		CountryCode string `json:"country_code"`
		Email       string `json:"email"`
		BindGa      bool   `json:"bind_ga"`
		UserKey     string `json:"user_key"`
	}{user.ID, user.Phone, user.CountryCode, user.Email, bindGa, user.Uuid}))
}

func (i *InfoController) SetInformation(ctx iris.Context) {
	var (
		userId uint
		params = struct {
			Language string `json:"language"`
			Timezone string `json:"timezone"`
		}{}
		err error
	)
	userId = common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	if params.Language != "" {
		if err = i.userModel.SetLanguage(userId, params.Language); err != nil {
			ZapLog().With(zap.Error(err)).Error("SetLanguage err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_SET_LANGUAGE", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	if params.Timezone != "" {
		if err = i.userModel.SetTimezone(userId, params.Timezone); err != nil {
			ZapLog().With(zap.Error(err)).Error("SetTimezone err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_SET_TIMEZONE", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
}

func (i *InfoController) BindEmail(ctx iris.Context) {
	var (
		userId uint
		params = struct {
			EmailToken string `json:"email_token"`
			Email      string
		}{}
		err error
	)
	userId = common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Email == "", "email"},
		{params.EmailToken == "", "email_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 检测邮箱验证码
	tokenPass, err := common.NewVerification(i.redis, "bind_email", common.VerificationTypeEmail).
		Check(params.EmailToken, userId, params.Email)
	if err != nil || !tokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
		return
	}

	user, err := i.userModel.GetUserById(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_USER_ERROR", common.ResponseError))
		return
	}
	if user.Email != "" {
		ZapLog().Error("YOU_HAVE_ALREADY_BIND_EMAIL err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "YOU_HAVE_ALREADY_BIND_EMAIL", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	// 检测邮箱是否绑定过
	exists, err := i.userModel.UserExisted(params.Email)
	if err != nil || exists {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("UserExisted err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "THE_EMAIL_HAS_ALREADY_EXISTED", apibackend.BASERR_OBJECT_EXISTS.Code()))
		return
	}

	basBody, _ := json.Marshal(map[string]interface{}{
		"user_email": params.Email,
	})

	_, _, err = bastionpay.CallApi(user.Uuid, string(basBody), "/v1/account/updatecontact")
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("bastionpay.CallApi err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_ERR", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	if err = i.userModel.SetEmail(userId, params.Email); err != nil {
		ZapLog().With(zap.Error(err)).Error("SetEmail err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_BIND_THE_EMAIL", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
	ctx.Next()
}

func (i *InfoController) BindPhone(ctx iris.Context) {
	var (
		userId uint
		params = struct {
			SmsToken    string `json:"sms_token"`
			Phone       string `json:"phone"`
			CountryCode string `json:"country_code"`
		}{}
		err error
	)
	userId = common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Phone == "", "phone"},
		{params.CountryCode == "", "country_code"},
		{params.SmsToken == "", "sms_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 检测手机
	tokenPass, err := common.NewVerification(i.redis, "bind_phone", common.VerificationTypeSms).
		Check(params.SmsToken, userId, params.CountryCode+params.Phone)
	if err != nil || !tokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_SMS_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
		return
	}

	user, err := i.userModel.GetUserById(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_USER_ERROR", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}
	if user.Phone != "" {
		ZapLog().Error("YOU_HAVE_ALREADY_BIND_PHONE err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "YOU_HAVE_ALREADY_BIND_PHONE", apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code()))
		return
	}

	// 检测手机是否绑定过
	exists, err := i.userModel.UserExisted(params.Phone)
	if err != nil || exists {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("UserExisted err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "THE_PHONE_HAS_ALREADY_EXISTED", apibackend.BASERR_OBJECT_EXISTS.Code()))
		return
	}

	basBody, _ := json.Marshal(map[string]interface{}{
		"user_mobile":  params.Phone,
		"country_code": params.CountryCode,
	})
	_, _, err = bastionpay.CallApi(user.Uuid, string(basBody), "/v1/account/updatecontact")
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("bastionpay.CallApi err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_ERR", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	if err = i.userModel.SetPhone(userId, params.Phone, params.CountryCode); err != nil {
		ZapLog().With(zap.Error(err)).Error("SetPhone err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_BIND_THE_PHONE", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
	ctx.Next()
}

func (i *InfoController) RebindPhone(ctx iris.Context) {
	var (
		userId uint
		params = struct {
			EmailToken  string `json:"email_token"`
			GaToken     string `json:"ga_token"`
			OldSmsToken string `json:"old_sms_token"`
			NewSmsToken string `json:"new_sms_token"`
			Phone       string `json:"phone"`
			CountryCode string `json:"country_code"`
		}{}
		err error
	)
	userId = common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Phone == "", "phone"},
		{params.CountryCode == "", "country_code"},
		{params.OldSmsToken == "", "old_sms_token"},
		{params.NewSmsToken == "", "new_sms_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	user, err := i.userModel.GetUserById(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_USER_ERROR", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}
	if user.Phone == "" {
		ZapLog().Error("YOU_HAVE_NOT_BIND_PHONE err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "YOU_HAVE_NOT_BIND_PHONE", apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code()))
		return
	}
	if user.Phone == params.Phone {
		ZapLog().Error("PLEASE_ENTER_A_DIFFERENT_PHONE err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "PLEASE_ENTER_A_DIFFERENT_PHONE", apibackend.BASERR_OBJECT_DATA_SAME.Code()))
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
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 检测邮箱验证码
	if user.Email != "" {
		tokenPass, err := common.NewVerification(i.redis, "rebind_phone", common.VerificationTypeEmail).
			Check(params.EmailToken, userId, user.Email)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
				//				glog.Error(err.Error())
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
			return
		}
	}

	// 检测GA验证码
	if user.Ga != "" {
		tokenPass, err := common.NewVerification(i.redis, "rebind_phone", common.VerificationTypeGa).
			Check(params.GaToken, userId, user.Ga)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
				//				glog.Error(err.Error())
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_GA_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
			return
		}
	}

	// 检测旧手机
	oldTokenPass, err := common.NewVerification(i.redis, "rebind_phone", common.VerificationTypeSms).
		Check(params.OldSmsToken, userId, user.CountryCode+user.Phone)
	if err != nil || !oldTokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_OLD_SMS_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
		return
	}

	// 检测新手机
	newTokenPass, err := common.NewVerification(i.redis, "bind_phone", common.VerificationTypeSms).
		Check(params.NewSmsToken, userId, params.CountryCode+params.Phone)
	if err != nil || !newTokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_NEW_SMS_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
		return
	}

	// 检测手机是否绑定过
	exists, err := i.userModel.UserExisted(params.Phone)
	if err != nil || exists {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("UserExisted err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "THE_PHONE_HAS_ALREADY_EXISTED", apibackend.BASERR_OBJECT_EXISTS.Code()))
		return
	}

	basBody, _ := json.Marshal(map[string]interface{}{
		"user_mobile":  params.Phone,
		"country_code": params.CountryCode,
	})
	_, _, err = bastionpay.CallApi(user.Uuid, string(basBody), "/v1/account/updatecontact")
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("bastionpay.CallApi err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_ERR", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	if err = i.userModel.SetPhone(userId, params.Phone, params.CountryCode); err != nil {
		ZapLog().With(zap.Error(err)).Error("SetPhone err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_BIND_THE_PHONE", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
	ctx.Next()
}

func (this *InfoController) GetUserInfo(ctx iris.Context) {
	userkeys := ctx.FormValue("userkey")
	userkeysNoSpace := strings.Replace(userkeys, " ", "", len(userkeys))
	userkeyArr := strings.Split(userkeysNoSpace, ",")

	resArr := make([]*admin.UserDetailInfo, 0)
	for i := 0; i < len(userkeyArr); i++ {
		if len(userkeyArr[i]) == 0 {
			continue
		}
		user, err := this.userModel.GetUserByUserkey(userkeyArr[i])
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserByUserkey err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR:"+err.Error(), apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
			return
		}

		bindGa := true
		if user.Ga == "" {
			bindGa = false
		}

		res := new(admin.UserDetailInfo)
		res.Id = user.ID
		ZapLog().With(zap.Uint("userid", user.ID)).Info("")
		//		glog.Info(user.ID)
		res.CreatedAt = user.CreatedAt
		res.Phone = user.Phone
		res.CountryCode = user.CountryCode
		res.Email = user.Email
		res.BindGa = bindGa
		res.VipLevel = user.VipLevel
		res.Citizenship = user.Citizenship
		res.Language = user.Language
		res.Timezone = user.Timezone
		res.RegistrationType = user.RegistrationType
		res.UserKey = user.Uuid
		res.CompanyName = user.CompanyName

		resArr = append(resArr, res)
	}

	content, err := json.Marshal(resArr)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "Json_marshal_ERROR:"+err.Error(), apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}
	ctx.JSON(common.NewSuccessResponse(ctx, string(content)))
}

func (this *InfoController) Listusers(ctx iris.Context) {
	reqUserList := struct {
		TotalLines   int `json:"total_lines" doc:"总数,0：表示首次查询"`
		PageIndex    int `json:"page_index" doc:"页索引,1开始"`
		MaxDispLines int `json:"max_disp_lines" doc:"页最大数，100以下"`

		Condition map[string]interface{} `json:"condition" doc:"条件查询"`
	}{}

	err := ctx.ReadJSON(&reqUserList)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
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
		totalLine, err = this.userModel.ListUserCountByBasic(reqUserList.Condition)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("ListUserCountByBasic err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	if pageNum < 1 || pageNum > 100 {
		pageNum = 50
	}

	beginIndex = pageNum * (pageIndex - 1)

	//ackUserList, err := db.ListUsers(beginIndex, pageNum)
	userBasics, err := this.userModel.ListUsersByBasic(beginIndex, pageNum, reqUserList.Condition)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ListUsersByBasic err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
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
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, string(dataAck)))
}
