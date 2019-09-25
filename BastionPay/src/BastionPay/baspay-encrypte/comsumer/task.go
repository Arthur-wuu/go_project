package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/baspay-encrypte/base"
	"BastionPay/baspay-encrypte/config"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type ReqEncyData struct {
	//Data   string  `json:"data"`
}

type ResEncyData struct {
	Timestamp   int64  `json:"timestamp,omitempty"`
	Nonce       string `json:"nonce,omitempty"`
	EncryptData string `json:"data,omitempty"`
	Signature   string `json:"signature,omitempty"`
	Code        int64  `json:"code,omitempty"`
}

func (d *ReqEncyData) SendToBastion(body []byte, url string, header map[string][]string) (*ResEncyData, error) {
	resEncyData := new(ResEncyData)
	fmt.Println("***send to bastion url", config.GConfig.BastionpayUrl.Bastionurl+url)
	bastionPayRes, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+url, bytes.NewBuffer(body), "POST", header)
	if err != nil {
		ZapLog().Error("send message to pastionpay get res ency data err", zap.Error(err))
		return nil, err
	}
	json.Unmarshal(bastionPayRes, resEncyData)

	return resEncyData, nil
}

func SendToBastion(body []byte, url string, header map[string][]string) (*ResEncyData, error) {
	resEncyData := new(ResEncyData)
	fmt.Println("***send to bastion url", config.GConfig.BastionpayUrl.Bastionurl+url)
	bastionPayRes, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+url, bytes.NewBuffer(body), "POST", header)
	if err != nil {
		ZapLog().Error("send message to pastionpay get res ency data err", zap.Error(err))
		return nil, err
	}
	json.Unmarshal(bastionPayRes, resEncyData)

	return resEncyData, nil
}
