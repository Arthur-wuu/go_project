package baspay

import (
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/util"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/base"

	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type(
	//退款
	RefundTrade struct{
		MerchantId                *string   `valid:"required" json:"merchant_id,omitempty"`
		MerchantRefundNo          *string   `valid:"required" json:"merchant_refund_no"`
		NotifyUrl                 *string   `valid:"optional" json:"notify_url,omitempty"`
		OriginalMerchantTradeNo   *string   `valid:"required" json:"original_merchant_trade_no,omitempty"`
		Remark                    *string   `valid:"optional" json:"remark,omitempty"`
		SignType                  *string   `valid:"required" json:"sign_type,omitempty"`
		Signature                 *string   `valid:"required"json:"signature,omitempty"`
		Timestamp                 *string   `valid:"required"json:"timestamp,omitempty"`
	}

	RefundRes struct{
		Code       *string        `json:"code"`
		Data       Datas          `json:"data"`
		Message    *string        `json:"message"`
	}

	Datas struct{
		OriginalMerchantTradeNo  *string    `json:"original_merchant_trade_no"`
		Status                   *string    `json:"status"`
	}
)

func (this * RefundTrade) Parse(f *api.RefundTrade) *RefundTrade {
	return &RefundTrade{
		MerchantId : f.MerchantId,
		MerchantRefundNo: f.MerchantRefundNo,
		NotifyUrl: f.NotifyUrl,
		OriginalMerchantTradeNo: f.OriginalMerchantTradeNo,
		Remark: f.Remark,
	}
}




func (this * RefundTrade) Send() (*RefundRes, error){
	var merchant_id string
	merchant_id=*this.MerchantId

	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	this.SignType = &signType
	this.Timestamp = &timeStamp
	this.NotifyUrl = &notifyUrl

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_id": this.MerchantId,
		"merchant_refund_no": this.MerchantRefundNo,
		"original_merchant_trade_no": this.OriginalMerchantTradeNo,
		"remark": this.Remark,
		"timestamp": this.Timestamp,
		"notify_url": this.NotifyUrl,
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	//fmt.Println("***signStr***",signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	if merchant_id == "8" {
		sha1.SetPriKey(h5PrivateKey)
	}

	if merchant_id == "1" || merchant_id == "2" || merchant_id == "3" || merchant_id == "4"{
		sha1.SetPriKey(CommonPrivateKey)
	}

	finalSign, err := sha1.Sign(signStr)
	//baseSign := utils.Base64Encode([]byte(finalSign))

	if err != nil {
		ZapLog().Error( "sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_id": this.MerchantId,
		"merchant_refund_no": this.MerchantRefundNo,
		"original_merchant_trade_no": this.OriginalMerchantTradeNo,
		"remark": this.Remark,
		"timestamp": this.Timestamp,
		"notify_url": this.NotifyUrl,
		"sign_type": this.SignType,
		"signature": finalSign,
	})

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/refund", bytes.NewBuffer(reqBody),"POST", nil) //map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"ab9dd65725876c597","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" }
	if err != nil {
		ZapLog().Error( "send message to refund err", zap.Error(err))
		return nil, err
	}
	fmt.Println("**trade info result**",string(result))

	resRefund := new(RefundRes)
	json.Unmarshal(result, resRefund)
	fmt.Println("**result resRefunf**", resRefund)

	return resRefund, nil

}


