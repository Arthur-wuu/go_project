package sms

import (
	"BastionPay/bas-notify/common"
	"BastionPay/bas-notify/config"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bytes"
	"encoding/hex"
)

const timeTemplate1 = "20060102150405" //常规类型

type ReqYunTongXin struct {
	To         string   `json:"to"`
	AppId      string   `json:"appId"`
	TemplateId string   `json:"templateId"`
	Datas      []string `json:"datas"`
}

type ResYunTongXin struct {
	StatusCode    string `json:"statusCode"`
	SmsMessageSid string `json:"smsMessageSid"`
	DateCreated   string `json:"dateCreated"`
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func (this *SmsMgr) DirectSendYunTongXin(tempId, appId, toPhones string, params []string) error {
	timeStr := time.Now().In(time.FixedZone("UTC", 8*3600)).Format(timeTemplate1)

	sig := md5V(config.GConfig.YunTongXun.AccountSid + config.GConfig.YunTongXun.AuthToken + timeStr)
	ulrStr := config.GConfig.YunTongXun.Url + "/2013-12-26/Accounts/" + config.GConfig.YunTongXun.AccountSid + "/SMS/TemplateSMS?sig=" + sig
	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Content-Type"] = "application/json;charset=utf-8"

	header["Authorization"] = base64.StdEncoding.EncodeToString([]byte(config.GConfig.YunTongXun.AccountSid + ":" + timeStr))

	req := &ReqYunTongXin{
		To:         toPhones,
		AppId:      appId,
		TemplateId: tempId,
		Datas:      params,
	}
	reqEds, err := json.Marshal(req)
	if err != nil {
		return err
	}

	header["Content-Length"] = fmt.Sprintf("%d", len(reqEds))

	resData, err := common.HttpSend(ulrStr, bytes.NewReader(reqEds), "POST", header)
	if err != nil {
		return err
	}

	res := new(ResYunTongXin)
	if err := json.Unmarshal(resData, res); err != nil {
		return err
	}
	if res.StatusCode != "000000" {
		return fmt.Errorf("%s-%s-%s", res.StatusCode, res.SmsMessageSid)
	}
	return nil
}
