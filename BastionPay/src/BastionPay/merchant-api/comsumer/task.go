package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"sort"

	"BastionPay/merchant-api/device"
	//"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"

	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/util"
)

var h5PrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAuW30rrCPsvjtXMtCEV7elJdQ81NC2r309zTItBx+0KOcvysU
Ss8lQTMWyONlLsM6RCZQoERUGuK+K+isyLNv40Ng4DsbMCyehX+FrLsJ7O6tVjfH
KB1OnnLqvOjfKToow7BU8uBZZgQTlyH7+QmX3hx6hk2Zpzsz5Qv+Uqknp93DmP19
OMCrcZubLg2laaAi2fUmBR2u6WWVXU4hRugSxYTSQJHeXm2jcOA13XJsO5/RcJrZ
8Xod81/6T0sHTt3Rpq/YAVldz/mMf+pjTmTc7Ga3wf6X4Dfjop4ahwssn8KUGkZH
0LVJYUsoTL6Z1XF2HLjuFk8gOgHF1QqcrAMXlQIDAQABAoIBAAbgGth3V3ytWi+8
oaB/QgWEbs324l21+WVJIb/75n/Z8S/tav0zHRCD0m50cAYB2yE3wnmOzweILMD4
CaAPdzCYCPmHd4Sbwuz0Q4KaFM4iM28k9k94drfNTwLeJ5ghve0zlpLC2s27jPZn
ni4nuufTan+cVPwsZ8owXae2+e9xa1qy1gJBveXT26kGPbo6q06zK0kfcwKs5cYG
o12tnvBzDv2h2TVoXujhADKNWNIAOkTmOk3VSAb+bTiNCipep7SbgeToYnX5rqx/
SHIFjFlwxMYqiadG8Ysq2hK0bg8HMMvK7L2PNQ/Ge4B5aXXJHxF4m+e+2EjL3uJl
VD2stGkCgYEA71oyBzPmbm4u/fhZsedwtCdHnfkQMFMje14sQiPbPonlWxo417cD
pwcPU8bd/s5X8saGim+1p3WShnY1+wHH/VF5aiFqwt/K69c3ofAfrmbmXbyN2gJJ
umuAR7Hb8Urd2teSova4LC4Fc7sSvtLXPle3eWaoD0mO1sxlQ37pSDcCgYEAxlOi
TNFB4uEbIJfmOqzqXywzibOxMasgjZLKEu+JE6Vk8y3YMzCfrlKl/2wr/I+DrhWx
yX7cf0/SfJgox0D9K23yh/Uu72QNRMZpVSRUGjCurtWvILRcD0mdzFkHijC7cfDE
3XqMHRjFrRgAWlt01YVojbjYJ/F8cudVNyXWYJMCgYEAuDrofvrHxwAwU3OxNmo6
KbCCQ2nNuCSGDxMxZcdLnhtt2m2YixFnUkzw0z8i6FnTAB8mt6+8VqT8n1qlugpo
8OahWbtW/aBcBKOnQpIdEJRLhKL5XHCeZ0sPdh/Edzl1AlkjmSPmJrtVnvrDNvX6
jxXdNyh4+ytXMqYo24b38IkCgYBkZ9UEJPDBRwuvzZcuX3psYnlZHpL3vVZGtmj9
ey2ft51LDAunptdArvEBRidivtmAmdUfWM2S2ruKfpIuhjVl9kzSDgwMAFBDYFvV
UgYOGFVniCEYYpc02iU8XlpV2OQdBDL2meMzm+YAAuWy2RhmPRs4nLs6RaSmm31l
5Q8KZwKBgQCVp0hRdSphZBjcUEAJyBHDcf3D4f0nJx3CLvmNg6TuoT+5hendCi2E
dArT9yVYco35HJd0nxBm/UOzEptNUuAeX41AVeN9o0UYoJomt/l07cyY77HtrnIb
ZpqTdTAmvK7JaZooBF736k+gTMmX0qdzfmoJSDNqHZaoq4tmuojZZg==
-----END RSA PRIVATE KEY-----
`)

type Data struct {
	MerchantTradeNo *string `json:"merchant_trade_no"`
	Assets          *string `json:"assets"`
	Amount          *string `json:"amount"`
	Status          *string `json:"status"`
	TradeNo         *string `json:"trade_no"`
}

var GTasker Tasker

type Tasker struct {
	mMerchantOrderChan chan *api.Merchant
}

func (this *Tasker) Start() {
	this.mMerchantOrderChan = make(chan *api.Merchant, 100)
	for i := 0; i < 10; i++ {
		go this.sChanWorker()
	}
}

func (this *Tasker) sChanWorker() {
	//查cache
	for {
		time.Sleep(time.Second * 2)
		orderNo := <-this.mMerchantOrderChan
		fmt.Println("开始拿chan 的数据...", orderNo.MerchantTradeNo, orderNo.MerchantId)

		//查询订单状态
		statusData, err := Send(orderNo.MerchantTradeNo, orderNo.MerchantId)
		if err != nil {
			ZapLog().Error("send message to pastionpay select order status err")
			continue
		}
		status := statusData["status"]

		//err = json.Unmarshal(statusData.([]byte), data)
		//if err != nil {
		//	ZapLog().Error( "unmarshal err", zap.Error(err))
		//	return
		//}

		fmt.Println("order status", status)

		if status == "3" {
			bool, err := new(models.Trade).UpdateRowsAffected(orderNo.MerchantTradeNo, models.EUM_TRANSFER_STATUS_SUCCESS)
			if err != nil {
				ZapLog().Error(" update by tradeNo err, notify end", zap.Error(err))
				//ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
				continue
			}
			if bool == false {
				ZapLog().Error("update 0 rows, notify end", zap.Error(err))
				//ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
				continue
			}
			//err := new(models.Trade).UpdateTransferStatusByTradeNo(nil, orderNo.MerchantTradeNo, models.EUM_TRANSFER_STATUS_SUCCESS)
			//if err != nil {
			//	ZapLog().Error( "succ update by tradeNo err, notify end", zap.Error(err))
			//	//ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
			//	return
			//}
			deviceId, gameCoin, err := new(models.Trade).GetDeviceId(nil, orderNo.MerchantTradeNo)
			if gameCoin == "" {
				continue
			}

			if err != nil {
				ZapLog().Error("get device id err, notify end", zap.Error(err))
				//ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
				continue
			}

			//获得机器id对应的机器
			device.GDeviceMgr.Init()
			devices := device.GDeviceMgr.Get(deviceId)

			err = devices.Send(gameCoin)
			fmt.Println("***send coin over ***", status)
			if err != nil {
				ZapLog().Error("send coin to machine fail, notify end", zap.Error(err))
				//ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
				continue
			}
		} else {
			//扔回去
			//fmt.Println("order status ！=3 " , status)
			this.mMerchantOrderChan <- orderNo
			continue
		}

	}
}

func (this *Tasker) SendOrderNo(msg *api.Merchant) {
	fmt.Println("数据进入chan 待消费中....", msg)
	this.mMerchantOrderChan <- msg
}

func Send(MerchantTradeNo, MerchantId string) (map[string]interface{}, error) {
	//往baspay 查询订单状态
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	//this.SignType = &signType
	//this.Timestamp = &timeStamp
	//this.NotifyUrl = &notifyUrl

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": MerchantTradeNo,
		"merchant_id":       MerchantId,
		//	"sign_type": this.Request.SignType,
		//	"signature": this.Request.Signature,
		"timestamp":  timeStamp,
		"notify_url": notifyUrl,
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	//fmt.Println("***signStr***",signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)
	//baseSign := utils.Base64Encode([]byte(finalSign))

	//fmt.Println("***final sign***",finalSign)
	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": MerchantTradeNo,
		"merchant_id":       MerchantId,
		"sign_type":         signType,
		"signature":         finalSign,
		"timestamp":         timeStamp,
		"notify_url":        notifyUrl,
	})

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/info", bytes.NewBuffer(reqBody), "POST", nil) //map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"ab9dd65725876c597","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" }
	if err != nil {
		ZapLog().Error("send message to get trade info err", zap.Error(err))
		return nil, err
	}
	//fmt.Println("**result**",string(result))

	resTradeInfo := new(ResTradeInfo)
	json.Unmarshal(result, resTradeInfo)

	//fmt.Println("**result2**", resTradeInfo.Data)

	return *resTradeInfo.Data, nil

}

func GetTimeStamp() string {
	times := time.Now().Format("2006-01-02 15:04:05")
	return times
}

type ResTradeInfo struct {
	Code    *int                    `json:"code"`
	Data    *map[string]interface{} `json:"data"`
	Message *string                 `json:"message"`
}

func RequestBodyToSignStr(body []byte) string {
	//fmt.Println("**body**",string(body))
	requestParams := make(map[string]string, 0)

	err := json.Unmarshal(body, &requestParams)
	//fmt.Println("**requestParams**",requestParams)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("requestbody to requestParams err")
		return ""
	}
	//将param的key排序，
	keysSort := make([]string, 0)
	for k, _ := range requestParams {
		keysSort = append(keysSort, k)
	}
	sort.Strings(keysSort)
	//fmt.Println("**keysSort**",keysSort)
	//拼接签名字符串
	signH5Str := ""
	for i := 0; i < len(keysSort); i++ {
		signH5Str += keysSort[i] + "=" + requestParams[keysSort[i]] + "&"
	}
	signH5Str = signH5Str[0 : len(signH5Str)-1]
	//fmt.Println("**signH5Str**",signH5Str)
	return signH5Str
}
