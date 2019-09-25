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
	l4g "github.com/alecthomas/log4go"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"strings"
)

func NewRedirectController(config *tools.Config) *RedirectController {
	bp := &RedirectController{
		config: config,
	}
	return bp
}

type RedirectController struct {
	config *tools.Config
}

type ResNotifyMsg struct {
	Err                 *int        `json:"err,omitempty"`
	ErrMsg              *string     `json:"errmsg,omitempty"`
	TemplateGroupList   interface{} `json:"templategrouplist,omitempty"`
	Templates           interface{} `json:"template,omitempty"`
	TemplateHistoryList interface{} `json:"templatehistorylist,omitempty"`
}

func (this *ResNotifyMsg) GetErr() int {
	if this.Err == nil {
		return 0
	}
	return *this.Err
}

func (this *ResNotifyMsg) GetErrMsg() string {
	if this.ErrMsg == nil {
		return ""
	}
	return *this.ErrMsg
}

func (bp *RedirectController) HandlerV1BasNotify(ctx iris.Context) {
	query := ctx.Request().URL.RawQuery
	newUrl := bp.config.Notify.Addr + ctx.Path()
	if len(query) != 0 {
		newUrl = newUrl + "?" + query
	}
	req, err := http.NewRequest(ctx.Method(), newUrl, ctx.Request().Body)
	if err != nil {
		l4g.Error("http NewRequest username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code(), Message: "NewRequest_ERROR:" + err.Error()})
		return
	}

	l4g.Debug("url:%s ", newUrl)

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l4g.Error("http Do username[%s] url[%s] err[%s]", utils.GetValueUserName(ctx), newUrl, err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "BasNotify_HTTP_DO_ERROR:" + err.Error()})
		return
	}
	if resp.StatusCode != 200 {
		l4g.Error("http Do Response username[%s] url[%s] err[%s]", utils.GetValueUserName(ctx), newUrl, resp.Status)
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "BasNotify_HTTP_RESPONSE_ERROR:" + resp.Status})
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l4g.Error("Body readAll username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "BASNOTIFY_READ_BODY_ERROR:" + err.Error()})
		return
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		l4g.Error("username[%s] admin.api[%s] response is null", utils.GetValueUserName(ctx), newUrl)
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "BASNOTIFY_REDIRECT_ERROR:response body content null"})
		return
	}
	if string(content) == "Not Found" {
		l4g.Error("username[%s] admin.api[%s] response is Not Found", utils.GetValueUserName(ctx), newUrl)
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "BASNOTIFY_REDIRECT_ERROR:response is Not Found"})
		return
	}
	l4g.Debug("BASNOTIFY response content[%s]", string(content))
	adminRes := new(ResNotifyMsg)
	if err := json.Unmarshal(content, adminRes); err != nil {
		l4g.Error("Unmarshal username[%s] content[%s] err[%s]", utils.GetValueUserName(ctx), string(content), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_DATA_UNPACK_ERROR.Code(), Message: "BASNOTIFY_REDIRECT_ERROR:response cannot Unmarshal, " + err.Error()})
		return
	}

	if adminRes.GetErr() != 0 {
		l4g.Error("BASNOTIFY username[%s] Response.Status.Code[%d] err[%s]", utils.GetValueUserName(ctx), adminRes.GetErr(), adminRes.GetErrMsg())
		ctx.JSON(&Response{Code: adminRes.GetErr(), Message: "BASNOTIFY_REDIRECT_ERROR: " + adminRes.GetErrMsg()})
		return
	}
	if adminRes.TemplateGroupList != nil {
		ctx.JSON(&Response{Code: adminRes.GetErr(), Message: adminRes.GetErrMsg(), Data: adminRes.TemplateGroupList})
		return
	}
	if adminRes.Templates != nil {
		ctx.JSON(&Response{Code: adminRes.GetErr(), Message: adminRes.GetErrMsg(), Data: adminRes.Templates})
		return
	}
	if adminRes.TemplateHistoryList != nil {
		ctx.JSON(&Response{Code: adminRes.GetErr(), Message: adminRes.GetErrMsg(), Data: adminRes.TemplateHistoryList})
		return
	}
	ctx.JSON(&Response{Code: adminRes.GetErr(), Message: adminRes.GetErrMsg()})
	//	l4g.Debug("deal HandleV1Admin username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
	return
}

//v2-->v1,暂时不能透传，未使用
func (bp *RedirectController) HandlerV2Gateway(ctx iris.Context) {
	l4g.Debug("start deal HandlerV2Gateway username[%s]", utils.GetValueUserName(ctx))
	newPath := strings.Trim(ctx.Path(), "\r")
	pathArr := strings.Split(newPath, "/")
	if len(pathArr) < 3 {
		l4g.Error("Path[%s] username[%s] err", newPath, utils.GetValueUserName(ctx))
		ctx.JSON(&Response{Code: apibackend.BASERR_UNSUPPORTED_METHOD.Code(), Message: "PATH_ERROR"})
		return
	}
	pathArr[1] = "v1"
	newPath = strings.Join(pathArr, "/")

	params := &struct {
		SubUserKey string      `json:"subuserkey"`
		Message    interface{} `json:"message"`
	}{}
	err := ctx.ReadJSON(&params)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err[%s] ", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	reqMsg, err := json.Marshal(params.Message)
	if err != nil {
		l4g.Error("Marshal username[%s] params[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: err.Error()})
		return
	}

	l4g.Debug("HandlerV2Gateway[%s] req content[%d]", utils.GetValueUserName(ctx), len(reqMsg))
	status, res, err := bastionpay.CallApi(params.SubUserKey, string(reqMsg), newPath)
	bp.doneHook(ctx, status.Err == 0, reqMsg, params.SubUserKey)
	if err != nil {
		l4g.Error("redirect GateWay username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "GATEWAY_REQUEST_ERRORS:" + err.Error()})
		return
	}
	ctx.JSON(&Response{Data: string(res)})
	l4g.Debug("deal HandlerV2Gateway username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
}

func (this *RedirectController) HandleV1Admin(ctx iris.Context) {
	// 不用ctx.Redirect()，因为Response结构不一样，还得有业务处理
	l4g.Debug("start deal HandleV1Admin username[%s]", utils.GetValueUserName(ctx))
	newPathArr := make([]string, 0)
	pathArr := strings.Split(ctx.Path(), "/")
	for i := 0; i < len(pathArr); i++ {
		if i == 3 {
			continue
		}
		if i == 2 {
			newPathArr = append(newPathArr, "inner", "user")
		}
		newPathArr = append(newPathArr, pathArr[i])
	}
	newPath := strings.Join(newPathArr, "/")
	adminUrl := this.config.BasAmin.Url + newPath
	req, err := http.NewRequest(ctx.Method(), adminUrl, ctx.Request().Body)
	if err != nil {
		l4g.Error("http NewRequest username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_CONFIG_ERROR.Code(), Message: "NewRequest_ERROR:" + err.Error()})
		return
	}

	l4g.Debug("url:%s ", adminUrl)

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l4g.Error("http Do username[%s] url[%s] err[%s]", utils.GetValueUserName(ctx), adminUrl, err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "ADMIN.API_HTTP_DO_ERROR:" + err.Error()})
		return
	}

	if resp.StatusCode != 200 {
		l4g.Error("http Do Response username[%s] url[%s] err[%s]", utils.GetValueUserName(ctx), adminUrl, resp.Status)
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "ADMIN.API_HTTP_DO_ERROR:" + resp.Status})
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l4g.Error("Body readAll username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "ADMIN.API_READ_BODY_ERROR:" + "From admin.api, " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		l4g.Error("username[%s] admin.api[%s] response is null", utils.GetValueUserName(ctx), adminUrl)
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "ADMIN.API_REDIRECT_ERROR: From admin.api response body content null"})
		return
	}
	if string(content) == "Not Found" {
		l4g.Error("username[%s] admin.api[%s] response is Not Found", utils.GetValueUserName(ctx), adminUrl)
		ctx.JSON(&Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "ADMIN.API_REDIRECT_ERROR:From admin.api response is Not Found"})
		return
	}
	l4g.Debug("admin.api response content[%s]", string(content))
	adminRes := new(AdminResponse)
	if err := json.Unmarshal(content, adminRes); err != nil {
		l4g.Error("Unmarshal username[%s] content[%s] err[%s]", utils.GetValueUserName(ctx), string(content), err.Error())
		ctx.JSON(&Response{Code: apibackend.BASERR_DATA_UNPACK_ERROR.Code(), Message: "ADMIN.API_REDIRECT_ERROR:From admin.api response cannot Unmarshal, " + err.Error()})
		return
	}

	if adminRes.Status.Code != 0 {
		l4g.Error("admin.api username[%s] Response.Status.Code[%d] err[%s]", utils.GetValueUserName(ctx), adminRes.Status.Code, adminRes.Status.Msg)
		ctx.JSON(&Response{Code: adminRes.Status.Code, Message: "ADMIN.API_REDIRECT_ERROR:From admin.api, " + adminRes.Status.Msg, Data: adminRes.Result})
		return
	}
	if (len(pathArr) > 0) && !strings.Contains(pathArr[len(pathArr)-1], "get") {
		ctx.Next()
	}

	ctx.JSON(&Response{Code: adminRes.Status.Code, Message: adminRes.Status.Msg, Data: adminRes.Result})
	l4g.Debug("deal HandleV1Admin username[%s] ok", utils.GetValueUserName(ctx))
	return
}

func (bp *RedirectController) doneHook(ctx iris.Context, status bool, msg []byte, subUserKey string) {
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
			userModel := models.NewUserModel(bp.config)
			userInfos, err := userModel.GetUserInfo(ctx, reqFrozeUser.UserKey)
			if err != nil {
				l4g.Error("GetUserInfo[%s] err[%s]", reqFrozeUser.UserKey, err.Error())
				return
			}
			if len(userInfos) == 0 {
				l4g.Error("userkey[%s] nofind userInfo in admin.api err", reqFrozeUser.UserKey)
			}
			for i := 0; i < len(userInfos); i++ {
				fileSuffix := fmt.Sprintf("updatefrozen_%d_%s", reqFrozeUser.IsFrozen, userInfos[i].Language)
				l4g.Info("start deal userkey[%s] updatefrozen[%s] phone[%s] mail[%s]", reqFrozeUser.UserKey, fileSuffix, userInfos[i].Phone, userInfos[i].Email)
				var err error
				if len(userInfos[i].Phone) != 0 {
					err = models.GlobalNotifyMgr.DirectSms(bp.config.Notify.TmplateDir, fileSuffix, userInfos[i].CountryCode+userInfos[i].Phone)
					if err != nil {
						l4g.Error("UserKey[%s] Send Sms[%s] err[%s]", reqFrozeUser.UserKey, userInfos[i].CountryCode+userInfos[i].Phone, err.Error())
					}
				}
				if len(userInfos[i].Email) != 0 {
					err := models.GlobalNotifyMgr.DirectMail(bp.config.Notify.TmplateDir, fileSuffix, userInfos[i].Email, bp.config.Notify.SrcEmail)
					if err != nil {
						l4g.Error("UserKey[%s] Send Mail[%s][%s] err[%s]", reqFrozeUser.UserKey, userInfos[i].Email, bp.config.Notify.SrcEmail, err.Error())
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
			if !(reqUserAuditeStatus.AuditeStatus == backend.AUDITE_Status_Deny) {
				l4g.Info("userkey[%s] audite[%d] no need notify", subUserKey, reqUserAuditeStatus.AuditeStatus)
				return
			}
			userModel := models.NewUserModel(bp.config)
			userInfos, err := userModel.GetUserInfo(ctx, subUserKey)
			if err != nil {
				l4g.Error("GetUserInfo[%s] err[%s]", subUserKey, err.Error())
				return
			}
			if len(userInfos) == 0 {
				l4g.Error("subuserkey[%s] nofind userInfo in admin.api err", subUserKey)
			}
			for i := 0; i < len(userInfos); i++ {
				fileSuffix := fmt.Sprintf("updateaudite_%d_%s", reqUserAuditeStatus.AuditeStatus, userInfos[i].Language)
				l4g.Info("start deal userkey[%s] updateaudite[%s] phone[%s] mail[%s]", subUserKey, fileSuffix, userInfos[i].Phone, userInfos[i].Email)
				var err error
				if len(userInfos[i].Phone) != 0 {
					err = models.GlobalNotifyMgr.DirectSms(bp.config.Notify.TmplateDir, fileSuffix, userInfos[i].CountryCode+userInfos[i].Phone)
					if err != nil {
						l4g.Error("UserKey[%s] Send Sms[%s] err[%s]", subUserKey, userInfos[i].CountryCode+userInfos[i].Phone, err.Error())
					}
				}
				if len(userInfos[i].Email) != 0 {
					err := models.GlobalNotifyMgr.DirectMail(bp.config.Notify.TmplateDir, fileSuffix, userInfos[i].Email, bp.config.Notify.SrcEmail)
					if err != nil {
						l4g.Error("UserKey[%s] Send Mail[%s][%s] err[%s]", subUserKey, userInfos[i].Email, bp.config.Notify.SrcEmail, err.Error())
					}
				}
				l4g.Info("deal userkey[%s] updateaudite ok", subUserKey)

			}
		}()
	}
}
