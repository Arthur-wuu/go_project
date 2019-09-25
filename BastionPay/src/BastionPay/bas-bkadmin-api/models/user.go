package models

import (
	"encoding/json"
	"github.com/BastionPay/bas-api/admin"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	l4g "github.com/alecthomas/log4go"
	"github.com/bugsnag/bugsnag-go/errors"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
)

func NewUserModel(config *tools.Config) *UserModel {
	// 增加company_name字段
	return &UserModel{mAdminUrl: config.BasAmin.Url}
}

const ConstAdminGetUserInfoPath = "/v1/inner/user/account/getuserinfo?userkey="

type UserModel struct {
	mAdminUrl string
}

func (this *UserModel) GetUserInfo(ctx iris.Context, userkey string) ([]admin.UserDetailInfo, error) {

	adminUrl := this.mAdminUrl + ConstAdminGetUserInfoPath + userkey
	resp, err := http.Get(adminUrl)
	if resp.StatusCode != 200 {
		l4g.Error("admin.api[%s] response status[%s]", adminUrl, resp.Status)
		return nil, errors.New(resp.Status, 0)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if len(content) == 0 {
		l4g.Error(" admin.api[%s] response is null", adminUrl)
		return nil, errors.New("content empty", 0)
	}
	if string(content) == "Not Found" {
		l4g.Error("admin.api[%s] response is Not Found", adminUrl)
		return nil, errors.New("Not Found", 0)
	}
	l4g.Debug("admin.api response content[%s]", string(content))
	adminRes := new(admin.AdminResponse)
	if err := json.Unmarshal(content, adminRes); err != nil {
		l4g.Error("Unmarshal  content[%s] err[%s]", string(content), err.Error())
		return nil, err
	}

	if adminRes.Status.Code != 0 {
		l4g.Error("admin.api Response.Status.Code[%s] err[%s]", adminRes.Status.Code, adminRes.Status.Msg)
		return nil, err
	}

	resStr, ok := adminRes.Result.(string)
	if !ok {
		switch v := adminRes.Result.(type) {
		default:
			l4g.Error("type:%v", v)
		}
		return nil, errors.New("type err", 0)
	}

	resAddr := make([]admin.UserDetailInfo, 0)
	if err := json.Unmarshal([]byte(resStr), &resAddr); err != nil {
		l4g.Error("json Unmarshal err:%s", err.Error())
		return nil, err
	}

	return resAddr, nil
}
