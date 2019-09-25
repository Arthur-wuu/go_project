package sms

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/config"
	sns "BastionPay/bas-tools/sdk.aws.sns"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/saintpete/twilio-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var GSmsMgr SmsMgr

type SmsMgr struct {
	mAwsSdk *sns.SnsSdk // aws
	mTwlSdk *twilio.Client
}

func (this *SmsMgr) Init() error {
	awsConfig := config.GConfig.Aws
	this.mAwsSdk = sns.NewSnsSdk(awsConfig.SnsRegion, awsConfig.Accesskeyid, awsConfig.Accesskey, awsConfig.Accesstoken)
	twConfig := &config.GConfig.Twilio
	this.mTwlSdk = twilio.NewClient(twConfig.Sid, twConfig.Token, nil)
	return nil
}

func (this *SmsMgr) DirectSendAws(body, phone string, senderId *string) error {
	if err := this.mAwsSdk.Send(body, phone, senderId); err != nil {
		return err
	}
	return nil
}

type ReqChuanglanMsg struct {
	Account string  `json:"account"`
	Pwd     string  `json:"password"`
	Msg     string  `json:"msg"`
	Report  string  `json:"report"`
	Phone   *string `json:"phone,omitempty"`
	Params  *string `json:"params,omitempty"`
}

type ResChuanglanMsg struct {
	Code       string `json:"code"`
	ErrorMsg   string `json:"errorMsg"`
	MsgId      string `json:"msgId"`
	Time       string `json:"time"`
	FailNum    string `json:"failNum"`
	successNum string `json:"successNum"`
}

//成功数量
func (this *SmsMgr) DirectSendChuanglan(body string, phones []string, params []string) (num int, err error) {
	if len(phones) == 0 {
		return 0, nil
	}
	if len(params) == 0 {
		phonesStr := strings.Join(phones, ",")
		num, err = this.sendToChuanglan(body, phonesStr, "")
	} else {
		newParams := ""
		for i := 0; i < len(phones); i++ {
			newParams += phones[i]
			for j := 0; j < len(params); j++ {
				newParams += "," + params[j]
			}
			newParams += ";"
		}
		phonesStr := strings.Join(phones, ",")
		newParams = strings.TrimRight(newParams, ";")
		num, err = this.sendToChuanglan(body, phonesStr, newParams)
	}
	if (num == 0) && (err == nil) {
		return len(phones), err
	}
	return num, err
}

func (this *SmsMgr) sendToChuanglan(body, phone string, params string) (int, error) {
	req := new(ReqChuanglanMsg)
	req.Account = config.GConfig.ChuangLan.Account
	req.Pwd = config.GConfig.ChuangLan.Pwd
	req.Report = "true"
	req.Msg = url.QueryEscape(body)
	req.Msg = body
	if len(params) != 0 {
		req.Params = new(string)
		*req.Params = params
	}
	req.Phone = new(string)
	*req.Phone = phone

	bytesData, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}
	ZapLog().Info("sendToChuanglan", zap.String("json", string(bytesData)))
	reader := bytes.NewReader(bytesData)
	url := config.GConfig.ChuangLan.Url
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("%d_%s", resp.StatusCode, resp.Status)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	resMsg := new(ResChuanglanMsg)
	if err = json.Unmarshal(respBytes, resMsg); err != nil {
		return 0, err
	}
	if resMsg.Code != "0" {
		return 0, fmt.Errorf("%s_%s_%s", resMsg.Code, resMsg.ErrorMsg, resMsg.MsgId)
	}
	ZapLog().Info("sendToChuanglan", zap.String("res", resMsg.Code+" "+resMsg.ErrorMsg+" "+resMsg.successNum+" "+resMsg.FailNum))
	num, _ := strconv.Atoi(resMsg.successNum)
	return num, nil
}

//一旦包含中文字符，短信最大字符数50个左右(官网67个)，超过发不出短信来，测试账户
func (this *SmsMgr) DirectSendTwl(body string, phone string) error {
	if (len(phone) > 4) && (phone[0] != '+') {
		phone = "+" + phone[2:]
	}
	if _, err := this.mTwlSdk.Messages.SendMessage(config.GConfig.Twilio.FromPhone, phone, body, nil); err != nil {
		return err
	}
	return nil
}
