package baspay

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/util"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

type (
	TradeInfo struct {
		MerchantTradeNo *string `json:"merchant_trade_no,omitempty"`
		MerchantId      *string `json:"merchant_id,omitempty"`
		SignType        *string `json:"sign_type,omitempty"`
		Signature       *string `json:"signature,omitempty"`
		Timestamp       *string `json:"timestamp,omitempty"`
		NotifyUrl       *string `json:"notify_url,omitempty"`
	}

	ResTradeInfo struct {
		Code    *int        `json:"code"`
		Data    interface{} `json:"data"`
		Message *string     `json:"message"`
	}

	ResTradeCoffee struct {
		Code    *int    `json:"code"`
		Data    Data    `json:"data"`
		Message *string `json:"message"`
	}

	Data struct {
		MerchantTradeNo *string `json:"merchant_trade_no"`
		Assets          *string `json:"assets"`
		Amount          *string `json:"amount"`
		Status          *int64  `json:"status"`
		TradeNo         *string `json:"trade_no"`
	}
)

func (this *TradeInfo) Parse(f *api.TradeInfo) *TradeInfo {
	return &TradeInfo{
		MerchantTradeNo: f.MerchantTradeNo,
		MerchantId:      f.MerchantId,
	}
}

func (this *TradeInfo) Send() (interface{}, error) {
	var merchant_id string
	merchant_id = *this.MerchantId
	//往baspay 查询订单状态
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	this.SignType = &signType
	this.Timestamp = &timeStamp
	this.NotifyUrl = &notifyUrl

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo,
		"merchant_id":       this.MerchantId,
		//	"sign_type": this.Request.SignType,
		//	"signature": this.Request.Signature,
		"timestamp":  this.Timestamp,
		"notify_url": this.NotifyUrl,
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("***signStr***", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	if merchant_id == "8" {
		sha1.SetPriKey(h5PrivateKey)
	}

	if merchant_id == "1" || merchant_id == "2" || merchant_id == "3" || merchant_id == "4" {
		sha1.SetPriKey(CommonPrivateKey)
	}

	finalSign, err := sha1.Sign(signStr)
	//baseSign := utils.Base64Encode([]byte(finalSign))

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo,
		"merchant_id":       this.MerchantId,
		"sign_type":         this.SignType,
		"signature":         finalSign,
		"timestamp":         this.Timestamp,
		"notify_url":        this.NotifyUrl,
	})

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/info", bytes.NewBuffer(reqBody), "POST", nil) //map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"ab9dd65725876c597","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" }
	if err != nil {
		ZapLog().Error("send message to get trade info err", zap.Error(err))
		return nil, err
	}
	fmt.Println("**trade info result**", string(result))

	resTradeInfo := new(ResTradeInfo)
	json.Unmarshal(result, resTradeInfo)
	//tempAmount := Decimal8(*resTradeInfo.Data.Amount)
	//resTradeInfo.Data.Amount = &tempAmount

	fmt.Println("**result2**", resTradeInfo.Data)

	return resTradeInfo.Data, nil

}

//coffee 查询订单状态
func (this *TradeInfo) SendCoffee() (string, error) {
	//往baspay 查询订单状态
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	this.SignType = &signType
	this.Timestamp = &timeStamp
	this.NotifyUrl = &notifyUrl

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo,
		"merchant_id":       this.MerchantId,
		//	"sign_type": this.Request.SignType,
		//	"signature": this.Request.Signature,
		"timestamp":  this.Timestamp,
		"notify_url": this.NotifyUrl,
	})

	signStr := RequestBodyToSignStr(reqBodySign)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(CommonPrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo,
		"merchant_id":       this.MerchantId,
		"sign_type":         this.SignType,
		"signature":         finalSign,
		"timestamp":         this.Timestamp,
		"notify_url":        this.NotifyUrl,
	})

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/info", bytes.NewBuffer(reqBody), "POST", nil) //map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"ab9dd65725876c597","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" }
	if err != nil {
		ZapLog().Error("send message to get trade info err", zap.Error(err))
		return "", err
	}
	//fmt.Println("**result**",string(result))

	resTradeCoffee := new(ResTradeCoffee)
	err = json.Unmarshal(result, resTradeCoffee)
	if err != nil {
		ZapLog().Error("json unmarshal err", zap.Error(err))
		return "", err
	}

	if resTradeCoffee.Data.Status == nil {
		return "1", err
	}
	ZapLog().Sugar().Infof("status %v", *resTradeCoffee.Data.Status)

	statusSucc3 := strconv.FormatInt(*resTradeCoffee.Data.Status, 10)

	return statusSucc3, nil

}

func Decimal8(value string) string {
	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return ""
	}
	f64, _ := strconv.ParseFloat(fmt.Sprintf("%.9f", float), 64)
	s2 := strconv.FormatFloat(f64, 'g', -1, 64) //float64
	s2 = s2[:len(s2)-1]
	return s2
}
