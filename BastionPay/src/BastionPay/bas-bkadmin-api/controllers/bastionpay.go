package controllers

import (
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-api/apibackend/v1/backend"
	"BastionPay/bas-bkadmin-api/models"
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/bastionpay"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	"github.com/BastionPay/bas-tools/sdk.notify.mail"
	l4g "github.com/alecthomas/log4go"
	"github.com/kataras/iris"
	"strings"
)

type BastionPayController struct {
	supportedFunction map[string]string
	mConfig           *tools.Config
}

func NewBastionPayController(config *tools.Config) *BastionPayController {
	bp := &BastionPayController{
		supportedFunction: make(map[string]string),
		mConfig:           config,
	}

	bp.supportedFunction = make(map[string]string)
	for _, path := range config.WalletPaths {
		index := strings.LastIndex(path, "/")

		relativePath := path[0:index]
		function := path[index+1:]

		bp.supportedFunction[function] = relativePath
	}

	err := sdk_notify_mail.GNotifySdk.Init(bp.mConfig.Notify.Addr, "bas-bkadmin-api")
	if err != nil {
		l4g.Error("mailNotify sdk init err[%s]", err.Error())
	}
	l4g.Info("notify sdk init ok [%s]", bp.mConfig.Notify.Addr)
	return bp
}

//未使用
func (bp *BastionPayController) Admin(ctx iris.Context) {
	params := struct {
		Function string      `json:"function"`
		Message  interface{} `json:"message"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		l4g.Error(err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	//user, err := User.GetValueUserInfo(ctx)
	//if err != nil {
	//	ctx.JSON(Response{Code: -1, Message: err.Error()})
	//	return
	//}

	params.Function = strings.Trim(params.Function, "\r")
	relativePath, ok := bp.supportedFunction[params.Function]
	if !ok {
		l4g.Error("NOT SUPPORT FUNCTION:(%s)", params.Function)
		ctx.JSON(Response{Code: apibackend.BASERR_UNSUPPORTED_METHOD.Code(), Message: fmt.Sprintf("NOT SUPPORT FUNCTION:(%s)", params.Function)})
		return
	}

	reqMsg, err := json.Marshal(params.Message)
	if err != nil {
		l4g.Error(err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: err.Error()})
		return
	}

	_, res, err := bastionpay.CallApi("", string(reqMsg), relativePath+"/"+params.Function)
	if err != nil {
		l4g.Error(err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Code: apibackend.BASERR_SUCCESS.Code(), Data: string(res)})
	ctx.Next()
}

func (bp *BastionPayController) AdminByFunction(ctx iris.Context) {
	l4g.Debug("start deal AdminByFunction username[%s]", utils.GetValueUserName(ctx))
	var (
		path         string
		relativePath string
		function     string

		params struct {
			SubUserKey string      `json:"subuserkey"`
			Message    interface{} `json:"message"`
		}
	)

	path = ctx.Path()
	index := strings.LastIndex(path, "/")
	if index == -1 {
		l4g.Error(fmt.Errorf("Path is error format"))
		ctx.JSON(Response{Code: apibackend.BASERR_UNSUPPORTED_METHOD.Code(), Message: fmt.Errorf("Path is error format")})
		return
	}

	relativePath = path[0:index]
	function = path[index+1:]

	err := ctx.ReadJSON(&params)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err[%s] ", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	//user, err := User.GetValueUserInfo(ctx)
	//if err != nil {
	//	ctx.JSON(Response{Code: -1, Message: err.Error()})
	//	return
	//}

	function = strings.Trim(function, "\r")
	relativePath, ok := bp.supportedFunction[function]
	if !ok {
		l4g.Error("supportedFunction username[%s] FUNCTION[%s] err", utils.GetValueUserName(ctx), function)
		ctx.JSON(Response{Code: apibackend.BASERR_UNSUPPORTED_METHOD.Code(), Message: fmt.Sprintf("NOT SUPPORT FUNCTION:(%s)", function)})
		return
	}

	reqMsg, err := json.Marshal(params.Message)
	if err != nil {
		l4g.Error("Marshal username[%s] params[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: err.Error()})
		return
	}

	status, res, err := bastionpay.CallApi(params.SubUserKey, string(reqMsg), relativePath+"/"+function)
	if err != nil {
		l4g.Error("bastionpay.CallApi username[%s] param[%v] reqMsg[%s] path[%s] err[%s]", utils.GetValueUserName(ctx), params, string(reqMsg), relativePath+"/"+function, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}
	bp.doneHook(ctx, status.Err == 0, reqMsg, params.SubUserKey)

	ctx.JSON(Response{Code: 0, Data: string(res)})

	l4g.Debug("deal AdminByFunction username[%s] ok, result[%d]", utils.GetValueUserName(ctx), len(res))
	ctx.Next()
}

func (bp *BastionPayController) doneHook(ctx iris.Context, status bool, msg []byte, subUserKey string) {
	//这个依赖配置接口，接口变了这个也得变化
	if strings.Contains(ctx.Path(), "updatefrozen") && status {
		l4g.Debug("doneHook updatefrozen start")
		go func() {
			reqFrozeUser := new(backend.ReqFrozenUser)
			err := json.Unmarshal(msg, reqFrozeUser)
			if err != nil {
				l4g.Error("username[%s] userkey[%s] Unmarshal backend.ReqFrozenUser err[%s] ", utils.GetValueUserName(ctx), reqFrozeUser.UserKey, err.Error())
				return
			}
			//reqFrozeUser := new(backend.ReqFrozenUser)
			//err := ctx.ReadJSON(reqFrozeUser)
			//if err != nil {
			//	l4g.Error("Context[%s] ReadJSON err[%s] ", utils.GetValueUserName(ctx), err.Error())
			//	return
			//}
			userModel := models.NewUserModel(bp.mConfig)
			userInfos, err := userModel.GetUserInfo(ctx, reqFrozeUser.UserKey)
			if err != nil {
				l4g.Error("GetUserInfo[%s] err[%s]", reqFrozeUser.UserKey, err.Error())
				return
			}
			if len(userInfos) == 0 {
				l4g.Error("userkey[%s] nofind userInfo in admin.api err", reqFrozeUser.UserKey)
			}

			for i := 0; i < len(userInfos); i++ {
				fileSuffix := fmt.Sprintf("updatefrozen_%d", reqFrozeUser.IsFrozen)
				l4g.Info("start deal userkey[%s] updatefrozen[%s] phone[%s] mail[%s]", reqFrozeUser.UserKey, fileSuffix, userInfos[i].Phone, userInfos[i].Email)
				var err error
				if len(userInfos[i].Phone) != 0 {
					//	err = models.GlobalNotifyMgr.DirectSms(bp.mConfig.Notify.TmplateDir, fileSuffix, userInfos[i].CountryCode+userInfos[i].Phone)
					err = sdk_notify_mail.GNotifySdk.SendSmsByGroupName(fileSuffix, userInfos[i].Language, []string{userInfos[i].CountryCode + userInfos[i].Phone}, nil)
					if err != nil {
						l4g.Error("UserKey[%s] Send Sms[%s] err[%s]", reqFrozeUser.UserKey, userInfos[i].CountryCode+userInfos[i].Phone, err.Error())
					}
				}
				if len(userInfos[i].Email) != 0 {
					//	err := models.GlobalNotifyMgr.DirectMail(bp.mConfig.Notify.TmplateDir, fileSuffix, userInfos[i].Email, bp.mConfig.Notify.SrcEmail)
					err = sdk_notify_mail.GNotifySdk.SendMailByGroupName(fileSuffix, userInfos[i].Language, []string{userInfos[i].Email}, nil)
					if err != nil {
						l4g.Error("UserKey[%s] Send Mail[%s][%s] err[%s]", reqFrozeUser.UserKey, userInfos[i].Email, bp.mConfig.Notify.SrcEmail, err.Error())
					}
				}
				l4g.Info("deal userkey[%s] updatefrozen ok", reqFrozeUser.UserKey)

			}
		}()
	}

	if strings.Contains(ctx.Path(), "updateaudite") && status {
		l4g.Debug("doneHook updateaudite start")
		go func() {
			reqUserAuditeStatus := new(backend.ReqUserAuditeStatus)
			err := json.Unmarshal(msg, reqUserAuditeStatus)
			if err != nil {
				l4g.Error("userkey[%s] Unmarshal backend.ReqFrozenUser err[%s] ", subUserKey, err.Error())
				return
			}
			// backend.AUDITE_Status_Pass  通过 不需要发短信
			if !(reqUserAuditeStatus.AuditeStatus == backend.AUDITE_Status_Deny) {
				l4g.Info("userkey[%s] audite[%d] no need notify", subUserKey, reqUserAuditeStatus.AuditeStatus)
				return
			}
			userModel := models.NewUserModel(bp.mConfig)
			userInfos, err := userModel.GetUserInfo(ctx, subUserKey)
			if err != nil {
				l4g.Error("GetUserInfo[%s] err[%s]", subUserKey, err.Error())
				return
			}
			if len(userInfos) == 0 {
				l4g.Error("subuserkey[%s] nofind userInfo in admin.api err", subUserKey)
			}

			for i := 0; i < len(userInfos); i++ {
				fileSuffix := fmt.Sprintf("audite_status_%d", reqUserAuditeStatus.AuditeStatus)
				l4g.Info("start deal userkey[%s] updateaudite[%s] phone[%s] mail[%s]", subUserKey, fileSuffix, userInfos[i].Phone, userInfos[i].Email)
				var err error
				if len(userInfos[i].Phone) != 0 {
					//err = models.GlobalNotifyMgr.DirectSms(bp.mConfig.Notify.TmplateDir, fileSuffix, userInfos[i].CountryCode+userInfos[i].Phone)
					err = sdk_notify_mail.GNotifySdk.SendSmsByGroupName(fileSuffix, userInfos[i].Language, []string{userInfos[i].CountryCode + userInfos[i].Phone}, nil)
					if err != nil {
						l4g.Error("UserKey[%s] Send Sms[%s] err[%s]", subUserKey, userInfos[i].CountryCode+userInfos[i].Phone, err.Error())
					}
				}
				if len(userInfos[i].Email) != 0 {
					//err := models.GlobalNotifyMgr.DirectMail(bp.mConfig.Notify.TmplateDir, fileSuffix, userInfos[i].Email, bp.mConfig.Notify.SrcEmail)
					err = sdk_notify_mail.GNotifySdk.SendMailByGroupName(fileSuffix, userInfos[i].Language, []string{userInfos[i].Email}, nil)
					if err != nil {
						l4g.Error("UserKey[%s] Send Mail[%s][%s] err[%s]", subUserKey, userInfos[i].Email, bp.mConfig.Notify.SrcEmail, err.Error())
					}
				}
				l4g.Info("deal userkey[%s] updateaudite ok", subUserKey)

			}
		}()
	}
	if strings.HasSuffix(ctx.Path(), "adminupdateprofile") && status {
		go func() {
			req := new(backend.ReqUserUpdateProfile)
			if err := json.Unmarshal(msg, req); err != nil {
				l4g.Error("userkey[%s] Unmarshal backend.ReqUserUpdateProfile err[%s] ", subUserKey, err.Error())
				return
			}
			// backend.AUDITE_Status_Pass  通过 不需要发短信
			if len(req.PublicKey) == 0 && len(req.SourceIP) == 0 && len(req.CallbackUrl) == 0 {
				l4g.Info("userkey[%s] adminupdateprofile no need notify", subUserKey)
				return
			}
			userModel := models.NewUserModel(bp.mConfig)
			userInfos, err := userModel.GetUserInfo(ctx, subUserKey)
			if err != nil {
				l4g.Error("GetUserInfo[%s] err[%s]", subUserKey, err.Error())
				return
			}
			if len(userInfos) == 0 {
				l4g.Error("subuserkey[%s] nofind userInfo in admin.api err", subUserKey)
				return
			}
			temlateName := "admin_update_profile"
			for i := 0; i < len(userInfos); i++ {
				l4g.Info("start deal userkey[%s] temlateName[%s] phone[%s] mail[%s]", subUserKey, temlateName, userInfos[i].Phone, userInfos[i].Email)
				zhParams := ""
				enParams := ""
				if len(req.PublicKey) != 0 {
					zhParams += "公钥，"
					enParams += "public key,"
				}
				if len(req.SourceIP) != 0 {
					zhParams += "IP白名单，"
					enParams += "while list,"
				}
				if len(req.CallbackUrl) != 0 {
					zhParams += "回调地址，"
					enParams += "callback address,"
				}
				zhParams = strings.TrimRight(zhParams, "，")
				enParams = strings.TrimRight(enParams, ",")
				Params := zhParams
				lang := userInfos[i].Language
				if (lang != "en-US") && (lang != "zh-CN") {
					lang = "en-US"
				}
				if lang == "en-US" {
					Params = enParams
				}

				//发送给管理员的
				err = sdk_notify_mail.GNotifySdk.SendSmsByGroupName(temlateName, "zh-CN", nil, map[string]interface{}{"key1": zhParams})
				if err != nil {
					l4g.Error("admin SendSmsByGroupName[%s][%s] err[%s]", userInfos[i].CompanyName, subUserKey, err.Error())
					//发送给管理员 失败，可能是 手机号错了，不能返回，要继续
				}

				//发送给用户的
				if (len(userInfos[i].CountryCode) != 0) && (len(userInfos[i].Phone) != 0) {
					err = sdk_notify_mail.GNotifySdk.SendSmsByGroupName(temlateName, lang, []string{userInfos[i].CountryCode + userInfos[i].Phone}, map[string]interface{}{"key1": Params}, false)
					if err != nil {
						l4g.Error("user SendSmsByGroupName[%s][%s][%s] err[%s]", userInfos[i].CompanyName, subUserKey, userInfos[i].CountryCode+userInfos[i].Phone, err.Error())
					}
				}
				if len(userInfos[i].Email) != 0 {
					err = sdk_notify_mail.GNotifySdk.SendMailByGroupName(temlateName, lang, []string{userInfos[i].Email}, map[string]interface{}{"key1": Params})
					if err != nil {
						l4g.Error("user SendMailByGroupName[%s][%s][%s] err[%s]", userInfos[i].CompanyName, subUserKey, userInfos[i].Email, err.Error())
					}
				}
				l4g.Info("deal userkey[%s] updateaudite ok", subUserKey)
			}
		}()
	}
}
