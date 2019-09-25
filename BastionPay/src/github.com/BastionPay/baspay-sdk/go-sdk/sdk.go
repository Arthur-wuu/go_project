package go_sdk

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	. "BastionPay/baspay-sdk/config"
	"BastionPay/bas-base/config"
)

const testCommonUrl    =  "https://test-openapi.bastionpay.io/"
const commonUrl        =  "https://open-api.bastionpay.com/"

const transUrl         =  "open-api/trade/transfer"
const availAssetsUrl   =  "open-api/trade/avail_assets"
const createQrUrl      =  "open-api/trade/create_qr_trade"
const createSdkUrl     =  "open-api/trade/create_trade"
const createWapUrl     =  "open-api/trade/create_wap_trade"
const tradeInfoUrl     =  "open-api/trade/info"
const posUrl           =  "open-api/trade/pos"
const posOrderUrl      =  "open-api/trade/pos_orders"


//商户10号的 11号的
var h5PrivateKey = []byte(`
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

//商户8号的
var h5PrivateKey111 = []byte(`
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


type Configs struct {
		KeyPath string   //`yaml:"keyPath"`
	}

var (
	//私钥
	prikey []byte

	//公钥
	pubkey []byte
)

type PaySdk struct{

}

//pem文件的路径
func (this* PaySdk) Init() error {
	//////加载pem文件，检测有效性
	////if err := loadRsaKeys(); err != nil {
	////	fmt.Println("err",err)
	////	fmt.Errorf("BastionPay Init: %s", err.Error())
	////	os.Exit(1)
	////}
	//
	////加载yaml文件
	laxFlag := config.NewLaxFlagDefault()
	cfgPath := laxFlag.String("conf_path", "config.yaml", "config path")
	logPath := laxFlag.String("log_path", "zap.conf", "log conf path")
	laxFlag.LaxParseDefault()
	fmt.Printf("command param: conf_path=%s, log_path=%s\n", *cfgPath, *logPath)
	LoadConfig(*cfgPath)

	return nil
}

//转账接口       数量，币种，商户id，收款人id
func (this* PaySdk) Transfer(amount, assets, merchant_id, payee_id, productName, notifyUrl string)( string, error) {
	param := new(TransferParam)
	times := time.Now().Local().Format("2006-01-02 15:04:05")

	param.Timestamp = &times
	param.MerchantTransNo = GenerateUuid()
	param.Amount = &amount
	param.Assets = &assets
	param.PayeeId = &payee_id
	param.MerchantId = &merchant_id
	param.ProductName = &productName
	param.NotifyUrl = &notifyUrl


	transferStatus, err := param.Send()
	if err != nil {
		return "fail", err
	}

	if transferStatus == "1" {
		return "unpay", err
	}

	if transferStatus == "3" {
		return "succ", err
	}

	if transferStatus == "4" {
		return "fail", err
	}

	return "fail", nil
}




// 加载数据
func loadRsaKeys(config *Configs) error {
	var err error

	private := fmt.Sprintf("%s/%s", strings.Trim(config.KeyPath, "/"), "private.pem")
	prikey, err = ioutil.ReadFile(private)
	if err != nil {
		return err
	}

	public := fmt.Sprintf("%s/%s", strings.Trim(config.KeyPath, "/"), "public.pem")
	pubkey, err = ioutil.ReadFile(public)
	if err != nil {
		return err
	}
	return nil
}



//有效币种接口 avail_assets
func (this* PaySdk) AvailAssets( arg ... interface{})( interface{}, error) {
	param := new(AvailAssetsParam)
	if len(arg) ==  1 {
		argParam := arg[0].(string)
		param.Assets = &argParam
		availAssetsRes, err := param.Send()
		if err != nil {
			return nil, err
		}
		return availAssetsRes, nil
	}

	if len(arg) ==  0 {
		availAssetsRes, err := param.Send()
		if err != nil {
			return nil, err
		}
		return availAssetsRes, nil
	}

	return nil, nil
}



//二维码下单接口       create_qr_trade
func (this* PaySdk) CreateQrtrade(amount, assets, merchant_id, payee_id, product_detail,product_name,remark,return_url,show_url, notifyUrl string , expire_time int64)( interface{}, error) {
	param := new(QrTrade)
	times := time.Now().Local().Format("2006-01-02 15:04:05")

	param.Timestamp = &times
	param.MerchantTradeNo = GenerateUuid()
	param.Amount = &amount
	param.Assets = &assets
	param.PayeeId = &payee_id
	param.MerchantId = &merchant_id
	param.ProductName = &product_name
	param.ProductDetail = &product_detail
	param.Remark = &remark
	param.ReturnUrl = &return_url
	param.ShowUrl = &show_url
	param.NotifyUrl = &notifyUrl
	param.ExpireTime = expire_time


	createQrRes, err := param.Send()
	if err != nil {
		return "fail", err
	}

	return createQrRes, nil
}





//Sdk下单接口       create sdk trade
func (this* PaySdk) CreateSdkTrade(amount, assets, merchant_id, payee_id, product_detail,product_name,remark,return_url,show_url, notifyUrl string , expire_time int64)( interface{}, error) {
	param := new(SdkTrade)
	times := time.Now().Local().Format("2006-01-02 15:04:05")

	param.Timestamp = &times
	param.MerchantTradeNo = GenerateUuid()
	param.Amount = &amount
	param.Assets = &assets
	param.PayeeId = &payee_id
	param.MerchantId = &merchant_id
	param.ProductName = &product_name
	param.ProductDetail = &product_detail
	param.Remark = &remark
	param.ReturnUrl = &return_url
	param.ShowUrl = &show_url
	param.NotifyUrl = &notifyUrl
	param.ExpireTime = expire_time


	createSdkRes, err := param.Send()
	if err != nil {
		return "fail", err
	}

	return createSdkRes, nil
}






//wap下单接口       create wap trade
func (this* PaySdk) CreateWapTrade(amount, assets, merchant_id, payee_id, product_detail,product_name,remark,return_url,show_url, notifyUrl string , expire_time int64)( interface{}, error) {
	param := new(WapTrade)
	times := time.Now().Local().Format("2006-01-02 15:04:05")

	param.Timestamp = &times
	param.MerchantTradeNo = GenerateUuid()
	param.Amount = &amount
	param.Assets = &assets
	param.PayeeId = &payee_id
	param.MerchantId = &merchant_id
	param.ProductName = &product_name
	param.ProductDetail = &product_detail
	param.Remark = &remark
	param.ReturnUrl = &return_url
	param.ShowUrl = &show_url
	param.NotifyUrl = &notifyUrl
	param.ExpireTime = expire_time


	createWapRes, err := param.Send()
	if err != nil {
		return "fail", err
	}

	return createWapRes, nil
}




//查询交易信息       5aa32a38-4214-4e52-a384-5ad8e3b54fb0
func (this* PaySdk) SearchTradeInfo ( merchant_id, merchant_trade_no, notifyUrl string )( interface{}, error) {
	param := new(TradeInfo)
	param.MerchantId = &merchant_id
	param.MerchantTradeNo = &merchant_trade_no
	param.NotifyUrl = &notifyUrl

	infoRes, err := param.Send()
	if err != nil {
		return "fail", err
	}

	return infoRes, nil
}



//pos支付接口       create_qr_trade
func (this* PaySdk) CreatePos (amount, assets, merchant_id, payee_id,pay_voucher,pos_merchine_id, product_detail,product_name,remark, notifyUrl string )( interface{}, error) {
	param := new(PosTrade)
	times := time.Now().Local().Format("2006-01-02 15:04:05")

	param.Timestamp = &times
	param.MerchantPosNo = GenerateUuid()
	param.PosMachineId = &pos_merchine_id
	param.PayVoucher = &pay_voucher
	param.Amount = &amount
	param.Assets = &assets
	param.PayeeId = &payee_id
	param.MerchantId = &merchant_id
	param.ProductName = &product_name
	param.ProductDetail = &product_detail
	param.Remark = &remark
	param.NotifyUrl = &notifyUrl


	createQrRes, err := param.Send()
	if err != nil {
		return "fail", err
	}

	return createQrRes, nil
}




//pos支付订单列表接口
func (this* PaySdk) PosOrderList (begin_time,end_time, merchant_id,merchant_pos_no, pos_merchine_id, page,page_size,trade_status, notifyUrl string )( interface{}, error) {
	param := new(PosOrderList)
	times := time.Now().Local().Format("2006-01-02 15:04:05")

	param.TimeStamp = &times
	param.MerchantPosNo = &merchant_pos_no
	param.PosMachineId = &pos_merchine_id
	param.MerchantId = &merchant_id
	param.NotifyUrl = &notifyUrl
	param.BeginTime = &begin_time
	param.EndTime = &end_time
	param.Page = &page
	param.PageSize = &page_size
	param.TradeStatus = &trade_status


	createQrRes, err := param.Send()
	if err != nil {
		return "fail", err
	}

	return createQrRes, nil
}





















