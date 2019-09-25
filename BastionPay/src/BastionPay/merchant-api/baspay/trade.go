package baspay

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/base"
	"fmt"
	"strings"

	//"BastionPay/merchant-api/controllers"

	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/models"
	"BastionPay/merchant-api/util"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	// rsa encode/decode bytes length limited, according to secret key bits
	RsaBits1024        = 1024
	RsaBits2048        = 2048
	RsaEncodeLimit1024 = RsaBits1024/8 - 11
	RsaDecodeLimit1024 = RsaBits1024 / 8
	RsaEncodeLimit2048 = RsaBits2048/8 - 11
	RsaDecodeLimit2048 = RsaBits2048 / 8
)

var serverPubByte = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArdDMsRHlTzzf8rfryGo8
2NDQ6VntnD07ax+7CMsKAAlICv28NxLHPoWRZAl9dRhM/uWGpgOPs2sKDayilyyR
0gZ8NPIVU4AWmn4xnv5l4Vu5HND9DcIoyvHLCiel+Lj/6HcpUzlJ+GmJ6L0QO/PI
CPq4KyR24ggCfknzAfLi8DQ+LUGFOhiSnu1ta3z4rVeOIyy72thlGoN7aTxXSMe6
yTi1bshkmFLgHyOcM2vpx4Vhtfb7xfu77LkRQEwi2k4vIZozInp4s5UaVFstd/Zd
IM/hMlwKP5zv4caLhI6Op3PrG+/6McLhx3j4tRxZhc6IdfSpvzEqO7icD+oRa5Sd
DwIDAQAB
-----END PUBLIC KEY-----
`)

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

var h5PublicKey = []byte(`
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuW30rrCPsvjtXMtCEV7e
lJdQ81NC2r309zTItBx+0KOcvysUSs8lQTMWyONlLsM6RCZQoERUGuK+K+isyLNv
40Ng4DsbMCyehX+FrLsJ7O6tVjfHKB1OnnLqvOjfKToow7BU8uBZZgQTlyH7+QmX
3hx6hk2Zpzsz5Qv+Uqknp93DmP19OMCrcZubLg2laaAi2fUmBR2u6WWVXU4hRugS
xYTSQJHeXm2jcOA13XJsO5/RcJrZ8Xod81/6T0sHTt3Rpq/YAVldz/mMf+pjTmTc
7Ga3wf6X4Dfjop4ahwssn8KUGkZH0LVJYUsoTL6Z1XF2HLjuFk8gOgHF1QqcrAMX
lQIDAQAB
`)

var prikeyWeb = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIGkAgEAAh8baPUMsl76IuHGIiOA8yC7GHh8Kxvq8lsgozWuFFInAgMBAAECHwYy
VDPqn04tVJ1WWnBshpc+FPiSwPf85vINpJu3QwECEAOfs/tD8gcm/ylfqJ1POnsC
EAeQVPnbqfsNukGNQAxq/UUCEADwivluv7XNDcJLlGvdnDsCEADbjbcTADAWO2NJ
Z9TAoN0CEAKVLMnwrY3t+iFgIpZ4aMg=
-----END RSA PRIVATE KEY-----`)

var pubkeyWeb = []byte(`-----BEGIN PUBLIC KEY-----
MDowDQYJKoZIhvcNAQEBBQADKQAwJgIfG2j1DLJe+iLhxiIjgPMguxh4fCsb6vJb
IKM1rhRSJwIDAQAB
-----END PUBLIC KEY-----`)

var CommonPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEArQ+bLJlYR4dn0l8yRGdN2ZSFns9NYLdL/DBLB70twKeRu5nL
7cAm0SzIMj6YtXgY9itPwOup53xFTCtvN322YY9Cffb1oBkcIY2PmtyX1C7RbDpD
P7tzwZWdz+eP67be19EMqfYXyYoPmMnoL2sPlMY34L/Jgop+2xw8F6kMRSUyZ8WU
+N0OswYF2W9S++fLWe1wZ+37WaUtAiOe721CtIVyQF2B7SDdGHoJ5AxRjcSP+rd5
qck99iGF1aAT+MYDFsBDBg5jZl8z6y1XccOQ50EVaNW7Gt6E+7951rZFA5RgAZqO
0lJtKFeKC3TSIuWsDhioGJ6m0aLqK7SDGz3FhwIDAQABAoIBAQCOofQkl+X4XhMl
gawuUG4LS6utLfH6KlgH682K/VI+HF2yHpnCw8G6WIxPTOQTfH4mNaAvwotv7C45
DvtE4ul0EtyccPQUFV3oEYIwAmtoR4X3CzXtyxMmk6dTeOhXP4r+mJ81XUxRoOYl
6RLiMfzPg2b+Z1MvsfHHqMemQH+KZDjdxXLGazk1DNAQXt8V/5e1SJhVcpv48QRC
ZeZHVRQ3f/6FL66edpGpQ5/oLBxw5m4tHiTSRg4T+a7sO4cCoZapCcsxaKx+s9eJ
DYjIOxPTNweqVVRDjrJ+fXxXKefbYx2AddUCecUIacyPMhDbWs5x//+FJkCFKvyH
eefKxWxxAoGBAM9Vxi0HaMCHdUkzUh/fX0uV8+f1kdl2k4W0CycXQA5zt1h2f2ve
fVcAEff5YwKYSTWlVgySxqyOVtz/ellzcEKcHiDZadIeCDYgoa23v266mnVpPYqT
n8jFwHhyVrPD/S+iLz8PrpF7NPYkQUNn40RIxNrEjMvUjyzatqgv8sTrAoGBANWu
Y/6t40mcBkCoWy0UJBuujPskiFrR9GLXQZ0GWhlSC+KKCOnO8/suWCicQqD1CD51
Nalpk+E3+wbXI/qDUpRq5WBob+l4CWKugtHhgaUdw8c+naWK6xn9okmD0Gadd1qr
jtAiHludI9fkNmkVIzaUkK26uBwOK41g9jMW/ErVAoGAG8BbWkOXnc2DwVyBLYr0
cmWL1AxmjTj13fuPUpgmFskeTVTvET0igbacsRhMTFid0/RhZCVxOj+DGOmJMtfk
usWysqrnIxyp9LTBb0Mc+HE5o2WGuzmvNWxiqryDJmShSvLmaAZtU0OufxOzOJZ7
MPSchLuyLMYys8pCkJh6YikCgYB/YUJC5C4GB3jCupn/uW39AoUQgaq3WUmyUlfO
36Z+SabEGT1PBAv1xJ7RNrWRdgDAGucuYr3BGLoQTdgo0ng7+a1bV2a/astNhHJ/
40qBv8ih0fXwZWvZRpWj9Wwaf+xSpMqx0GUAgCCJ5oV5BxzCwLWumwx9zQSxdwfN
VPp5MQKBgQDCwEBWR4xzjbt5MoE8pggF0an8zRGv3Z6yP0vFP+q/I6ykCrfNxByP
uchEixbp9Oe43zDzw8r/dIcCPPRw44A7voWfrAQKhyUjMVVki+ovGSQBZ5I6n9/7
AV5pE9Mk3YF7sLozbtKGyZev1YxnV3poC1DS+hTbcJdLFPCj3TQC1A==
-----END RSA PRIVATE KEY-----
`)

type (
	Trade struct {
		MerchantTradeNo *string `valid:"optional" json:"merchant_trade_no,omitempty"`
		PayeeId         *string `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"required" json:"amount,omitempty"`
		ProductName     *string `valid:"required" json:"product_name,omitempty"`
		ProductDetail   *string `valid:"required" json:"product_detail,omitempty"`
		ExpireTime      *int64  `valid:"required" json:"expire_time,omitempty"`
		Remark          *string `valid:"required" json:"remark,omitempty"`
		ReturnUrl       *string `valid:"required" json:"return_url,omitempty"`
		ShowUrl         *string `valid:"required" json:"show_url,omitempty"`
		DeviceId        *string `valid:"optional" json:"device_id,omitempty"`
		Request
	}

	QrTrade struct {
		MerchantTradeNo *string `valid:"optional" json:"merchant_trade_no,omitempty"`
		PayeeId         *string `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"optional" json:"amount,omitempty"`
		Legal           *string `valid:"optional" json:"legal,omitempty"`
		LegalNum        *string `valid:"optional" json:"legal_num,omitempty"`
		ProductName     *string `valid:"required" json:"product_name,omitempty"`
		ProductDetail   *string `valid:"optional" json:"product_detail,omitempty"`
		ExpireTime      *int64  `valid:"optional" json:"expire_time,omitempty"`
		Remark          *string `valid:"optional" json:"remark,omitempty"`
		ReturnUrl       *string `valid:"optional" json:"return_url,omitempty"`
		ShowUrl         *string `valid:"optional" json:"show_url,omitempty"`
		Request
	}

	CoffeeQrTrade struct {
		MerchantTradeNo *string `valid:"optional" json:"merchant_trade_no,omitempty"`
		PayeeId         *string `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"optional" json:"amount,omitempty"`
		ProductName     *string `valid:"required" json:"product_name,omitempty"`
		ProductDetail   *string `valid:"optional" json:"product_detail,omitempty"`
		ExpireTime      *int64  `valid:"optional" json:"expire_time,omitempty"`
		Remark          *string `valid:"optional" json:"remark,omitempty"`
		ReturnUrl       *string `valid:"optional" json:"return_url,omitempty"`
		ShowUrl         *string `valid:"optional" json:"show_url,omitempty"`
		Request
	}

	PosTrade struct {
		PosMachineId  *string `valid:"required" json:"pos_machine_id,omitempty"`
		MerchantPosNo *string `valid:"required" json:"merchant_pos_no,omitempty"`
		PayeeId       *string `valid:"required" json:"payee_id,omitempty"`
		Assets        *string `valid:"required" json:"assets,omitempty"`
		Amount        *string `valid:"required" json:"amount,omitempty"`
		MerchantId    *string `valid:"required" json:"merchant_id,omitempty"`
		SignType      *string `valid:"required" json:"sign_type,omitempty"`
		Signature     *string `valid:"required" json:"signature,omitempty"`
		Timestamp     *string `valid:"required" json:"timestamp,omitempty"`
		NotifyUrl     *string `valid:"optional" json:"notify_url,omitempty"`
		PayVoucher    *string `valid:"required" json:"pay_voucher,omitempty"`
		ProductName   *string `valid:"required" json:"product_name,omitempty"`
		ProductDetail *string `valid:"optional" json:"product_detail,omitempty"`
		Remark        *string `valid:"optional" json:"remark,omitempty"`
	}

	PosOrders struct {
		TimeStamp     *string `valid:"required" json:"timestamp,omitempty"`
		BeginTime     *string `valid:"optional" json:"begin_time,omitempty"`
		EndTime       *string `valid:"optional" json:"end_time,omitempty"`
		MerchantId    *string `valid:"required" json:"merchant_id,omitempty"`
		MerchantPosNo *string `valid:"optional" json:"merchant_pos_no"`
		NotifyUrl     *string `valid:"optional" json:"notify_url,omitempty"`
		Page          *string `valid:"optional" json:"page,omitempty"`
		PageSize      *string `valid:"optional" json:"page_size,omitempty"`
		PosMachineId  *string `valid:"required" json:"pos_machine_id,omitempty"`
		TradeStatus   *string `valid:"optional" json:"trade_status,omitempty"`
		SignType      *string `valid:"required" json:"sign_type,omitempty"`
		Signature     *string `valid:"required" json:"signature,omitempty"`
	}

	ResTrade struct {
		MerchantTradeNo *string `json:"merchant_trade_no,omitempty"`
		TradeNo         *string `json:"trade_order_no,omitempty"`
		Response
	}

	TradeWebRes struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	TradeQrRes struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	CoffeeTradeQrRes struct {
		Code    int       `json:"code"`
		Message string    `json:"message"`
		Data    BaspayRes `json:"data"`
	}

	BaspayRes struct {
		Merchant_trade_no string `json:"merchant_trade_no,omitempty"`
		Qr_code           string `json:"qr_code,omitempty"`
	}

	//pos机
	PosRes struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    PosDataRes `json:"data"`
	}

	PosDataRes struct {
		Assets        string `json:"assets,omitempty"`
		Amount        string `json:"amount,omitempty"`
		MerchantPosNo string `json:"merchant_pos_no,omitempty"`
		PosNo         string `json:"pos_no,omitempty"`
		Status        string `json:"status,omitempty"`
		TimeStamp     string `json:"time_stamp,omitempty"`
		Uid           string `json:"uid,omitempty"`
		Legal         string `json:"legal,omitempty"`
		LegalAmount   string `json:"lagal_amount,omitempty"`
	}

	//查询pos订单返回值
	PosOrdersRes struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	TradeSdkRes struct {
		Code    int     `json:"code,omitempty"`
		Message string  `json:"message,omitempty"`
		Data    DataSdk `json:"data,omitempty"`
	}

	DataSdk struct {
		MerchantTradeNo string `json:"merchant_trade_no,omitempty"`
		TradeNo         string `json:"trade_order_no,omitempty"`
	}
)

func (this *Trade) Parse(p *api.Trade) *Trade {
	trade := &Trade{
		MerchantTradeNo: &p.MerchantTradeNo,
		PayeeId:         p.PayeeId,
		Assets:          p.Assets,
		Amount:          p.Amount,
		ProductName:     p.ProductName,
		ProductDetail:   p.ProductDetail,
		ExpireTime:      &p.ExpireTime,
		Remark:          p.Remark,
		ReturnUrl:       p.ReturnUrl,
		ShowUrl:         p.ShowUrl,
		DeviceId:        p.DeviceId,
	}
	trade.MerchantId = p.MerchantId
	trade.NotifyUrl = p.NotifyUrl
	return trade
}

func (this *QrTrade) QrParse(p *api.QrTrade) *QrTrade {
	qrTrade := &QrTrade{
		MerchantTradeNo: &p.MerchantTradeNo,
		PayeeId:         p.PayeeId,
		Assets:          p.Assets,
		Amount:          p.Amount,
		Legal:           p.Legal,
		LegalNum:        p.LegalNum,
		ProductName:     p.ProductName,
		ProductDetail:   p.ProductDetail,
		ExpireTime:      &p.ExpireTime,
		Remark:          p.Remark,
	}
	qrTrade.MerchantId = p.MerchantId
	qrTrade.NotifyUrl = p.NotifyUrl
	return qrTrade
}

func (this *CoffeeQrTrade) CoffeeQrParse(p *api.CoffeeQrTrade) *CoffeeQrTrade {
	qrTrade := &CoffeeQrTrade{
		MerchantTradeNo: &p.MerchantTradeNo,
		PayeeId:         p.PayeeId,
		Assets:          p.Assets,
		Amount:          p.Amount,
		ProductName:     p.ProductName,
		ProductDetail:   p.ProductDetail,
		ExpireTime:      &p.ExpireTime,
		Remark:          p.Remark,
	}
	qrTrade.MerchantId = p.MerchantId
	qrTrade.NotifyUrl = p.NotifyUrl
	return qrTrade
}

func (this *PosTrade) PosParse(p *api.PosTrade) *PosTrade {
	qrTrade := &PosTrade{
		MerchantPosNo: &p.MerchantPosNo,
		PayeeId:       p.PayeeId,
		Assets:        p.Assets,
		Amount:        p.Amount,
		Timestamp:     p.TimeStamp,
		PayVoucher:    p.PayVoucher,
		ProductDetail: p.ProductDetail,
		ProductName:   p.ProductName,
		Remark:        p.Remark,
		PosMachineId:  p.PosMachineId,
	}
	qrTrade.MerchantId = p.MerchantId
	qrTrade.NotifyUrl = p.NotifyUrl
	return qrTrade
}

func (this *PosOrders) PosOrderParse(p *api.PosOrders) *PosOrders {
	posOrders := &PosOrders{
		MerchantPosNo: p.MerchantPosNo,
		TimeStamp:     p.TimeStamp,
		BeginTime:     p.BeginTime,
		EndTime:       p.EndTime,
		MerchantId:    p.MerchantId,
		NotifyUrl:     p.NotifyUrl,
		Page:          p.Page,
		PageSize:      p.PageSize,
		TradeStatus:   p.TradeStatus,
		PosMachineId:  p.PosMachineId,
	}
	return posOrders
}

func (this *Trade) Send() (*TradeWebRes, error) {
	//先去数据库根据设备id找到收款人id
	payee_id, err := new(models.BkConfig).GetPayeeId(*this.DeviceId)
	if err != nil {
		ZapLog().Error("get payid err", zap.Error(err))
		return nil, err
	}

	payId := strconv.FormatInt(*payee_id, 10)

	//往baspay 创建交易订单
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	expireTime := "900"

	this.Request.SignType = &signType
	this.Request.Timestamp = &timeStamp
	this.Request.NotifyUrl = &notifyUrl

	//传进去的的金额限制八位
	amount := Decimal7(*this.Amount)
	this.Amount = &amount

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          payId,                //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		//	"sign_type": this.Request.SignType,			 //签名算法类型
		//	"signature": this.Request.Signature,		 //⽤户请求的签名串
		"timestamp":  this.Request.Timestamp, //发送请求的时间
		"notify_url": this.Request.NotifyUrl, //回调通知商户服务器器的地址
	})

	signStr := RequestBodyToSignStr(reqBodySign)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	//fmt.Println("***final sign***",finalSign)
	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          payId,                //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       this.ExpireTime,      //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   this.Request.SignType,   //签名算法类型
		"signature":   finalSign,               //⽤户请求的签名串
		"timestamp":   this.Request.Timestamp,  //发送请求的时间
		"notify_url":  this.Request.NotifyUrl,  //回调通知商户服务器器的地址
	})

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/create_wap_trade", bytes.NewBuffer(reqBody), "POST", nil) //map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"ab9dd65725876c597","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" }
	//fmt.Println("**request result**",string(result))
	ZapLog().Info("create wap order result:", zap.String("web result:", string(result)))
	if err != nil {
		ZapLog().Error("send message to pastionpay create trade err", zap.Error(err))
		return nil, err
	}

	//
	resTrade := new(TradeWebRes)
	json.Unmarshal(result, resTrade)
	//订单号扔到channal中
	//p := new(api.Merchant)
	//p.MerchantTradeNo = *this.MerchantTradeNo
	//p.MerchantId = *this.Request.MerchantId
	//comsumer.GTasker.SendOrderNo(p)

	return resTrade, nil
}

//创建二维码订单

func (this *QrTrade) SendQr() (*TradeQrRes, error) {
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl
	showUrl := config.GConfig.CallBack.ShowUrl
	returnUrl := config.GConfig.CallBack.ReturnUrl

	amount, err := GetAmount(*this.Legal, *this.LegalNum, *this.Assets)
	//fmt.Println("**amount",*amount)
	if err != nil {
		ZapLog().Error("get amount err", zap.Error(err))
		return nil, err
	}

	amountDec := Decimal(*amount)

	stringAmount := fmt.Sprintf("%v", amountDec)

	this.Amount = &stringAmount
	expireTime := "900"
	//merchartId := this.MerchantId

	this.Request.SignType = &signType
	this.Request.Timestamp = &timeStamp
	this.Request.NotifyUrl = &notifyUrl
	this.ReturnUrl = &returnUrl
	this.ShowUrl = &showUrl

	//传进去的的金额限制八位
	amounts := Decimal7(*this.Amount)
	this.Amount = &amounts

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		//"sign_type": this.Request.SignType,			 //签名算法类型
		//	"signature": this.Request.Signature,		 //⽤户请求的签名串
		"timestamp":  this.Request.Timestamp, //发送请求的时间
		"notify_url": this.Request.NotifyUrl, //回调通知商户服务器器的地址
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       this.ExpireTime,      //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   this.Request.SignType,   //签名算法类型
		"signature":   finalSign,               //⽤户请求的签名串
		"timestamp":   this.Request.Timestamp,  //发送请求的时间
		"notify_url":  this.Request.NotifyUrl,  //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/create_qr_trade", bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("send message to pastionpay create trade err", zap.Error(err))
		return nil, err
	}

	ZapLog().Info("create qr order result:", zap.String("qr result:", string(result)))

	resTradeQr := new(TradeQrRes)
	json.Unmarshal(result, resTradeQr)

	return resTradeQr, nil
}

//价格有折扣的可以选这个方法，atm 收手续费的
func (this *QrTrade) SendDiscountQr() (*TradeQrRes, error) {
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl
	showUrl := config.GConfig.CallBack.ShowUrl
	returnUrl := config.GConfig.CallBack.ReturnUrl

	amount, err := GetAmount(*this.Legal, *this.LegalNum, *this.Assets)

	cash2coin_fee, err := strconv.ParseFloat(config.GConfig.Fee.Cash2coin, 64)
	if err != nil {
		ZapLog().Error("string to float err", zap.Error(err))
		return nil, err
	}

	amounts := *amount * (1 + cash2coin_fee)

	//fmt.Println("**amount",*amount)
	if err != nil {
		ZapLog().Error("get amount err", zap.Error(err))
		return nil, err
	}

	amountDec := Decimal(amounts)

	stringAmount := fmt.Sprintf("%v", amountDec)

	this.Amount = &stringAmount
	expireTime := "900"
	//merchartId := this.MerchantId

	this.Request.SignType = &signType
	this.Request.Timestamp = &timeStamp
	this.Request.NotifyUrl = &notifyUrl
	this.ReturnUrl = &returnUrl
	this.ShowUrl = &showUrl

	//传进去的的金额限制八位
	amount7 := Decimal7(*this.Amount)
	this.Amount = &amount7

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		//"sign_type": this.Request.SignType,			 //签名算法类型
		//	"signature": this.Request.Signature,		 //⽤户请求的签名串
		"timestamp":  this.Request.Timestamp, //发送请求的时间
		"notify_url": this.Request.NotifyUrl, //回调通知商户服务器器的地址
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(CommonPrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       this.ExpireTime,      //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   this.Request.SignType,   //签名算法类型
		"signature":   finalSign,               //⽤户请求的签名串
		"timestamp":   this.Request.Timestamp,  //发送请求的时间
		"notify_url":  this.Request.NotifyUrl,  //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/create_qr_trade", bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("send message to pastionpay create trade err", zap.Error(err))
		return nil, err
	}

	ZapLog().Info("create qr order result:", zap.String("qr result:", string(result)))

	resTradeQr := new(TradeQrRes)
	json.Unmarshal(result, resTradeQr)

	return resTradeQr, nil
}

//coffee 二维码订单
func (this *CoffeeQrTrade) SendCoffeeQr() (*CoffeeTradeQrRes, error) {
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl
	showUrl := config.GConfig.CallBack.ShowUrl
	returnUrl := config.GConfig.CallBack.ReturnUrl

	//this.Amount = &stringAmount
	expireTime := "900"
	//merchartId := this.MerchantId

	this.Request.SignType = &signType
	this.Request.Timestamp = &timeStamp
	this.Request.NotifyUrl = &notifyUrl
	this.ReturnUrl = &returnUrl
	this.ShowUrl = &showUrl

	//传进去的的金额限制八位
	amounts := Decimal7(*this.Amount)
	this.Amount = &amounts

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		//"sign_type": this.Request.SignType,			 //签名算法类型
		//	"signature": this.Request.Signature,		 //⽤户请求的签名串
		"timestamp":  this.Request.Timestamp, //发送请求的时间
		"notify_url": this.Request.NotifyUrl, //回调通知商户服务器器的地址
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(CommonPrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       this.ExpireTime,      //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   this.Request.SignType,   //签名算法类型
		"signature":   finalSign,               //⽤户请求的签名串
		"timestamp":   this.Request.Timestamp,  //发送请求的时间
		"notify_url":  this.Request.NotifyUrl,  //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/create_qr_trade", bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("send message to pastionpay create trade err", zap.Error(err))
		return nil, err
	}

	ZapLog().Info("create qr order result:", zap.String("qr result:", string(result)))

	resTradeQr := new(CoffeeTradeQrRes)
	json.Unmarshal(result, resTradeQr)

	return resTradeQr, nil
}

//Pos机
func (this *PosTrade) PosSend() (*PosRes, error) {
	signType := "RSA"
	//timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	this.SignType = &signType
	this.NotifyUrl = &notifyUrl

	//传进去的的金额限制八位
	amounts := Decimal7(*this.Amount)
	this.Amount = &amounts

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"pos_machine_id":  this.PosMachineId,
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"payee_id":        this.PayeeId,       //商户BastionPay 用户Id
		"assets":          this.Assets,        //数字货币币种
		"amount":          this.Amount,        //数字货币数量
		"pay_voucher":     this.PayVoucher,
		"merchant_id":     this.MerchantId,    //商户在BabstionPay 注册的商户ID
		"product_name":    this.ProductName,   //订单标题
		"product_detail":  this.ProductDetail, //订单描述
		"remark":          this.Remark,        //交易易备注信息
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature  	 //⽤户请求的签名串
		"timestamp":  this.Timestamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"pos_machine_id":  this.PosMachineId,
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"payee_id":        this.PayeeId,       //商户BastionPay 用户Id
		"assets":          this.Assets,        //数字货币币种
		"amount":          this.Amount,        //数字货币数量
		"pay_voucher":     this.PayVoucher,
		"merchant_id":     this.MerchantId,    //商户在BabstionPay 注册的商户ID
		"product_name":    this.ProductName,   //订单标题
		"product_detail":  this.ProductDetail, //订单描述
		"remark":          this.Remark,        //交易易备注信息
		"sign_type":       this.SignType,      //签名算法类型
		"signature":       finalSign,          //⽤户请求的签名串
		"timestamp":       this.Timestamp,     //发送请求的时间
		"notify_url":      this.NotifyUrl,     //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/pos", bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("send message to pastionpay create pos trade err", zap.Error(err))
		return nil, err
	}

	ZapLog().Info("create pos order result:", zap.String("pos result:", string(result)))

	posRes := new(PosRes)
	json.Unmarshal(result, posRes)

	return posRes, nil
}

//Pos机订单查询
func (this *PosOrders) PosOrdersSend() (*PosOrdersRes, error) {
	signType := "RSA"
	//timeStamp := GetTimeStamp()
	//notifyUrl := config.GConfig.CallBack.NotifyUrl

	this.SignType = &signType
	//this.NotifyUrl = &notifyUrl

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"pos_machine_id":  this.PosMachineId,
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"page":            this.Page,          //商户BastionPay 用户Id
		"page_size":       this.PageSize,      //数字货币币种
		//"amount": this.Amount,		         	 //数字货币数量
		//"pay_voucher": this.PayVoucher,
		"merchant_id":  this.MerchantId,  //商户在BabstionPay 注册的商户ID
		"begin_time":   this.BeginTime,   //订单标题
		"end_time":     this.EndTime,     //订单描述
		"trade_status": this.TradeStatus, //交易易备注信息
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature  	     //⽤户请求的签名串
		"timestamp":  this.TimeStamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"pos_machine_id":  this.PosMachineId,
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"page":            this.Page,          //商户BastionPay 用户Id
		"page_size":       this.PageSize,      //数字货币币种
		//"amount": this.Amount,		         	 //数字货币数量
		//"pay_voucher": this.PayVoucher,
		"merchant_id":  this.MerchantId,  //商户在BabstionPay 注册的商户ID
		"begin_time":   this.BeginTime,   //订单标题
		"end_time":     this.EndTime,     //订单描述
		"trade_status": this.TradeStatus, //交易易备注信息
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature  	     //⽤户请求的签名串
		"timestamp":  this.TimeStamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
		"sign_type":  this.SignType,  //签名算法类型
		"signature":  finalSign,      //⽤户请求的签名串
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/pos_orders", bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("send message to pastionpay list pos orders trade err", zap.Error(err))
		return nil, err
	}

	ZapLog().Info("create pos order result:", zap.String("pos result:", string(result)))

	posOrdersRes := new(PosOrdersRes)
	json.Unmarshal(result, posOrdersRes)

	return posOrdersRes, nil
}

//创建sdk订单
func (this *Trade) SendSdk() (*TradeSdkRes, error) {

	//往baspay 创建交易订单
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	expireTime := "900"
	//merchartId := this.MerchantId

	this.Request.SignType = &signType
	this.Request.Timestamp = &timeStamp
	this.Request.NotifyUrl = &notifyUrl

	//传进去的的金额限制八位
	amounts := Decimal7(*this.Amount)
	this.Amount = &amounts

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		//	"sign_type": this.Request.SignType,			 //签名算法类型
		//	"signature": this.Request.Signature,		 //⽤户请求的签名串
		"timestamp":  this.Request.Timestamp, //发送请求的时间
		"notify_url": this.Request.NotifyUrl, //回调通知商户服务器器的地址
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	//fmt.Println("***signStr***",signStr)
	//sign, err := RsaSignWithSha1Hex(signStr, string(h5PrivateKey))

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
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       this.ExpireTime,      //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.Request.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   this.Request.SignType,   //签名算法类型
		"signature":   finalSign,               //⽤户请求的签名串
		"timestamp":   this.Request.Timestamp,  //发送请求的时间
		"notify_url":  this.Request.NotifyUrl,  //回调通知商户服务器器的地址
	})
	//fmt.Println("**reqBody**",string(reqBody))

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/create_trade", bytes.NewBuffer(reqBody), "POST", nil) //map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"ab9dd65725876c597","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" }
	//fmt.Println("**request result**",string(result))
	if err != nil {
		ZapLog().Error("send message to pastionpay create trade err", zap.Error(err))
		return nil, err
	}

	ZapLog().Info("create sdk order result:", zap.String("sdk result:", string(result)))

	//
	resTrade := new(TradeSdkRes)
	json.Unmarshal(result, resTrade)

	return resTrade, nil
}

func GetTimeStamp() string {
	times := time.Now().Format("2006-01-02 15:04:05")
	return times
}

type (
	AmountResponse struct {
		Code   int               `json:"err,omitempty"`
		Quotes []QuoteDetailInfo `json:"quotes,omitempty" doc:"币简称"`
	}

	AmountDetailInfo struct {
		Symbol     string       `json:"symbol,omitempty" doc:"币简称"`
		Id         int          `json:"id,omitempty" doc:"币id"`
		MoneyInfos []MoneyInfos `json:"detail,omitempty" doc:"币行情数据"`
	}

	MoneyInfos struct {
		Symbol             *string  `json:"symbol,omitempty" doc:"币简称"`
		Price              float64  `json:"price,omitempty" doc:"最新价"`
		Volume_24h         *float64 `json:"volume_24h,omitempty" doc:"24小时成交量"`
		Market_cap         *float64 `json:"market_cap,omitempty" doc:"总市值"`
		Percent_change_1h  *float64 `json:"percent_change_1h,omitempty" doc:"1小时涨跌幅"`
		Percent_change_24h *float64 `json:"percent_change_24h,omitempty" doc:"24小时涨跌幅"`
		Percent_change_7d  *float64 `json:"percent_change_7d,omitempty" doc:"7天涨跌幅"`
		Last_updated       *int64   `json:"last_updated,omitempty" doc:"最近更新时间"`
	}
)

//http://test-quote.rkuan.com
func GetAmount(legal, legalNum, assets string) (*float64, error) {

	coinAmountResult, err := base.HttpSend(config.GConfig.BastionpayUrl.QuoteUrl+"/api/v1/coin/exchange?from="+strings.ToLower(legal)+"&to="+assets+"&amount="+legalNum, nil, "GET", nil)
	if err != nil {
		ZapLog().Error("select quote info err", zap.Error(err))
		return nil, err
	}

	amountResponse := new(AmountResponse)
	json.Unmarshal(coinAmountResult, amountResponse)

	if len(amountResponse.Quotes) == 0 || len(amountResponse.Quotes[0].MoneyInfos) == 0 {
		return nil, err
	}
	return &amountResponse.Quotes[0].MoneyInfos[0].Price, err
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", value), 64)
	return value
}

func Decimal7(value string) string {
	if len(value) < 10 {
		return value
	}
	int := strings.IndexAny(value, ".")
	if len(value)-int <= 8 {
		return value
	}
	s2 := value[0 : int+9]
	return s2
}
