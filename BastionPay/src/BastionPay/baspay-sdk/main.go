package main

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/baspay-sdk/base"
	"BastionPay/baspay-sdk/go-sdk"
	"BastionPay/baspay-sdk/util"
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	//config := new(go_sdk.Config)
	//config.KeyPath = "baspay-sdk/go-sdk/pem"

	trans := new(go_sdk.PaySdk)
	//trans.Init()
	//err := trans.Init(config)
	//if err != nil {
	//
	//}
	uidAndNum := ReadCsv("1.csv")

	//a :="[[1,98],[2,78],[3,45]]"
	succNum := 0
	failNum := 0
	for j := 0; j < len(uidAndNum); j++ {
		time.Sleep(4 * time.Second)

		uid := uidAndNum[j][0]
		amount := uidAndNum[j][1]
		//fmt.Println("amount",amount)
		if amount == "0" {
			continue
		}

		//requestNo := time.Now().UnixNano()
		//requestNoStr := strconv.FormatInt(requestNo, 10)
		//status, err := TransferCoin("SHINE", "201808", token, uid, requestNoStr, amount)
		status, err := trans.Transfer(amount, "SHINE", "12", uid, "满一年的6个的shine", "http://test.com")

		if err != nil {
			ZapLog().Error("send req to bastion err", zap.Error(err))
			fmt.Println("transfer err")
		}

		if status == "succ" {
			//fmt.sPrintf("转给", uid, amount)
			s := fmt.Sprintf("转给 [%v] [%v] 个 SHINE", uid, amount)
			fmt.Println(s)
			succNum++
			//ZapLog().Info("",)
		}
		if status == "fail" {
			s := fmt.Sprintf("失败 [%v]", uid)
			fmt.Println(s)
			failNum++
		}
		if status == "unpay" {
			s := fmt.Sprintf("失败 [%v]", uid)
			fmt.Println(s)
			failNum++

		}
	}
	fmt.Println("success num", succNum)
	fmt.Println("fail num", failNum)
	////转账测试
	//status, _ := trans.Transfer("4998","SHINE", "10", "11", "测试11okg","http://test.com")
	//time.Sleep(time.Second * 5)
	//fmt.Println("status",status)

	//有效币种测试，传参检验是否存在，不传参数返回全部
	//res,_ := trans.AvailAssets("BTC")
	//res1,_ := trans.AvailAssets()
	//fmt.Println("res",res, res1)

	//创建二维码订单测试
	//res , _ := trans.CreateQrtrade("1","OKG","8","57","test","test","test","test","test","test",900)
	//fmt.Println("res",res)

	//创建sdk订单测试
	//res , _ := trans.CreateSdkTrade("1","OKG","8","57","test","test","test","test","test","test",900)
	//fmt.Println("res",res)

	//创建wap订单测试
	//res , _ := trans.CreateWapTrade("1","OKG","8","57","test","test","test","test","test","test",900)
	//fmt.Println("res",res)

	//查询订单交易信息
	//res , _ := trans.SearchTradeInfo("8","5aa32a38-4214-4e52-a384-5ad8e3b54fb0","test")
	//fmt.Println("res",res)

}

//
//
//func  main() {
//	//登录bastionPay的用户登录
//	token, refreshtoken, err := LoginBastionPay("+086_16666666666","BastionPay666")
//	if  err != nil {
//		ZapLog().Sugar().Errorf("get token err [%v]", err)
//		return
//	}
//	fmt.Println("token", token, refreshtoken)
//
//	uidAndNum := ReadCsv("1.csv")
//
//	//a :="[[1,98],[2,78],[3,45]]"
//	succNum := 0
//	failNum := 0
//	for j:=0; j <len(uidAndNum); j ++ {
//		time.Sleep(2 * time.Second)
//
//		uid := uidAndNum[j][0]
//		amount := uidAndNum[j][1]
//		//fmt.Println("amount",amount)
//		if amount == "0" {
//			continue
//		}
//
//		requestNo  := time.Now().UnixNano()
//		requestNoStr := strconv.FormatInt(requestNo, 10)
//		status ,err := TransferCoin("SHINE","201808",token,uid,requestNoStr, amount)
//
//		if err != nil  {
//			ZapLog().Error( "send req to bastion err", zap.Error(err))
//			fmt.Println("transfer err")
//		}
//
//		if status == 2 {
//			//fmt.sPrintf("转给", uid, amount)
//			s := fmt.Sprintf("转给 [%v] [%v] 个 SHINE",uid,amount)
//			fmt.Println(s)
//			succNum ++
//			//ZapLog().Info("",)
//		}
//		if status == 1 {
//			s := fmt.Sprintf("失败 [%v]",uid)
//			fmt.Println(s)
//			failNum++
//		}
//		if status == 3 {
//			s := fmt.Sprintf("失败 [%v]",uid)
//			fmt.Println(s)
//			failNum++
//
//		}
//	}
//	fmt.Println("success num", succNum)
//	fmt.Println("fail num", failNum)
//
//}

const (
	RsaBits1024        = 1024
	RsaBits2048        = 2048
	RsaEncodeLimit1024 = RsaBits1024/8 - 11
	RsaDecodeLimit1024 = RsaBits1024 / 8
	RsaEncodeLimit2048 = RsaBits2048/8 - 11
	RsaDecodeLimit2048 = RsaBits2048 / 8
)

const (
	serverRsaUrl = "/wallet/api/security/handshake/server"
	clientRsaUrl = "/wallet/api/security/handshake/client"
	Login_Url    = "/wallet/api/user/oauth/login"
	Transfer_Url = "/wallet/api/user/assets/pay/transfer"
	Refresh_Url  = "/wallet/api/user/oauth/refresh_token"
	GetUid_Url   = "/wallet/api/user/list_by_phone"
	Check_Url    = "/wallet/api/user/assets/pay/confirm_safety"
)

//商户8号的
var privateKey = []byte(`
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

//商户8号的
var publicKey = []byte(`
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuW30rrCPsvjtXMtCEV7e
lJdQ81NC2r309zTItBx+0KOcvysUSs8lQTMWyONlLsM6RCZQoERUGuK+K+isyLNv
40Ng4DsbMCyehX+FrLsJ7O6tVjfHKB1OnnLqvOjfKToow7BU8uBZZgQTlyH7+QmX
3hx6hk2Zpzsz5Qv+Uqknp93DmP19OMCrcZubLg2laaAi2fUmBR2u6WWVXU4hRugS
xYTSQJHeXm2jcOA13XJsO5/RcJrZ8Xod81/6T0sHTt3Rpq/YAVldz/mMf+pjTmTc
7Ga3wf6X4Dfjop4ahwssn8KUGkZH0LVJYUsoTL6Z1XF2HLjuFk8gOgHF1QqcrAMX
lQIDAQAB
`)

//商户10号的
var privateKey10 = []byte(`
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

//商户10号的
var publicKey10 = []byte(`
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArQ+bLJlYR4dn0l8yRGdN2
ZSFns9NYLdL/DBLB70twKeRu5nL7cAm0SzIMj6YtXgY9itPwOup53xFTCtvN322YY
9Cffb1oBkcIY2PmtyX1C7RbDpDP7tzwZWdz+eP67be19EMqfYXyYoPmMnoL2sPlMY
34L/Jgop+2xw8F6kMRSUyZ8WU+N0OswYF2W9S++fLWe1wZ+37WaUtAiOe721CtIVy
QF2B7SDdGHoJ5AxRjcSP+rd5qck99iGF1aAT+MYDFsBDBg5jZl8z6y1XccOQ50EVa
NW7Gt6E+7951rZFA5RgAZqO0lJtKFeKC3TSIuWsDhioGJ6m0aLqK7SDGz3FhwIDAQAB
`)

//登录去拿token 为转账做准备
func LoginBastionPay(phone, pwd string) (string, string, error) {
	times := time.Now().Unix()
	timeStr := strconv.FormatInt(times, 10)
	nonce := "dhsjkf"

	//aes的cbc加密数据
	secretKey, err := GetServerSecretKey()
	fmt.Println("[--- login start ---]:", string(secretKey))
	push_id := "1"
	data := "{\"timestamp\":\"" + timeStr + "\", \"nonce\":\"" + nonce + "\", \"phone\":\"" + phone + "\", \"password\":\"" + pwd + "\",\"push_id\":\"" + push_id + "\"}"
	encyData, err := utils.Aes128Encrypt([]byte(data), []byte(secretKey))

	//fmt.Println("****encyData**",string(encyData))
	//签名
	signStr := "url=/wallet/api/user/oauth/login&timestamp=" + timeStr + "&nonce=" + nonce + "&data=" + string(encyData) + "&" + string(secretKey)
	//fmt.Println("***sign***",signStr)

	signString := GetSHA256HashCode([]byte(signStr))
	//将数据加密
	reqBody, _ := json.Marshal(map[string]interface{}{
		"timestamp": timeStr,
		"nonce":     nonce,
		"data":      string(encyData),
		"signature": signString,
	})
	fmt.Println("**timeStr**", timeStr)
	fmt.Println("**nonce**", nonce)
	fmt.Println("**data**", string(encyData))
	fmt.Println("**signString**", signString)
	//请求bastion
	result, err := SendToBastion(reqBody, Login_Url)
	if err != nil {
		ZapLog().Error("send req to bastion err", zap.Error(err))

		return "get token err1", "refresh token err1", err
	}
	fmt.Println("[** result login ** ] :", result.Signature, result.Nonce, result.Code, result.Timestamp)
	if result.Code == 1001 {
		return "get token err2", "refresh token err2", err
	}

	//验证签名参数
	timestpStr := strconv.FormatInt(result.Timestamp, 10)
	var verifySign string
	if result.EncryptData != "" {
		verifySign = "url=" + Login_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&data=" + result.EncryptData + "&" + string(secretKey)
	}

	if result.EncryptData == "" {
		verifySign = "url=" + Login_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&" + string(secretKey)
	}
	//有的没有返回签名
	if result.Signature == "" {
		return "get token err3", "refresh token err3", err
	}
	sign2 := GetSHA256HashCode([]byte(verifySign))
	//fmt.Println("**two sign login **",result.Signature,sign2)

	if result.Signature == sign2 {
		if result.EncryptData != "" {
			//将加密数据解密，返回给H5
			decrpData, err := utils.Aes128Decrypt([]byte(result.EncryptData), []byte(secretKey))
			if err != nil {
				ZapLog().Error("decry return data err", zap.Error(err))

				return "get token err4", "refresh token err4", err
			}

			var index int

			for k, v := range decrpData {
				if v == 0 {
					index = k
					break
				}
			}

			dataInfo := new(ResDataInfo)
			err = json.Unmarshal(decrpData[:index], dataInfo)

			if err != nil {
				ZapLog().Error("unmarshal return data err", zap.Error(err))
				//返回token ******
			}
			return dataInfo.Token, dataInfo.RefreshToken, nil
		}
		if result.EncryptData == "" {
			return "no data back", "refresh token err", err
		}

	}
	if result.Signature != sign2 {
		ZapLog().Error("sign not equal", zap.Error(err))
		return "get token err5", "refresh token err5", err
	}

	return "get token err6", "refresh token err6", err
}

//上传客户端 公钥 并申请SecretKey
func GetServerSecretKey() ([]byte, error) {
	baspayRes := new(Body)
	//加密因子
	//secretKeyRamdom :=GenRandomString(32)

	body := []byte("client_public_key=" + string(publicKey) + "&secrect_key_ramdom=35321243")

	//fmt.Println("****body body**",string(body))
	//拿服务端的公钥 加密body
	//serverPub, err := GetServerRsaPub()
	//fmt.Println("****serverPub pubkey",serverPub)
	//if err != nil {
	//	ZapLog().Error( "get server pub err", zap.Error(err))
	//}

	serverPubByte := []byte(`
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

	encyBody, err := utils.RsaEncrypt(body, serverPubByte, RsaEncodeLimit2048)

	if err != nil {
		ZapLog().Error("rsa encybody  err")
	}
	//fmt.Println("****encyBody** 加密body，**",string(encyBody))

	//请求服务端的 得到SecretKey的密文
	bastionPayRes, err := base.HttpSendSer("https://user-api.bastionpay.io"+clientRsaUrl, bytes.NewBuffer(encyBody), "POST", map[string]string{"Client": "1", "DeviceType": "1", "DeviceName": "huawei", "DeviceId": "2585b79fdsgls7aab76c9dd", "Version": "1.1.0", "Content-Type": "application/json;charset=UTF-8"}) //
	//fmt.Println("****bastionPayRes** ，**",string(bastionPayRes))
	if err != nil {
		ZapLog().Error("send message to pastionpay get secret key err", zap.Error(err))
		return nil, err
	}
	json.Unmarshal(bastionPayRes, baspayRes)
	//fmt.Println("****baspayRes** ，**",baspayRes.Data)
	//客户端私钥解密 得到SecretKey的密文
	secretKey, err := utils.RsaDecrypt3([]byte(baspayRes.Data), privateKey)

	//fmt.Println("****secretKey** ，**",string(secretKey))
	return secretKey, nil
}

func SendToBastion(body []byte, url string) (*ResEncyData, error) {
	resEncyData := new(ResEncyData)
	//fmt.Println("[** send to bastion url** ]:",config.GConfig.BastionpayUrl.Bastionurl+url)
	bastionPayRes, err := base.HttpSendSer("https://user-api.bastionpay.io"+url, bytes.NewBuffer(body), "POST", map[string]string{"Client": "1", "DeviceType": "1", "DeviceName": "huawei", "DeviceId": "2585b79fdsgls7aab76c9dd", "Version": "1.1.0", "Content-Type": "application/json;charset=UTF-8"})
	if err != nil {
		ZapLog().Error("send message to pastionpay get res ency data err", zap.Error(err))
		return nil, err
	}
	json.Unmarshal(bastionPayRes, resEncyData)

	return resEncyData, nil
}

//
//func SendToBastionLoginPool (body []byte  ,url string, deviceId string ) (*ResEncyData, error) {
//	resEncyData := new(ResEncyData)
//	fmt.Println("[** device id** ]:",deviceId)
//	bastionPayRes, err := base.HttpSendSer(config.GConfig.BastionpayUrl.Bastionurl+url, bytes.NewBuffer(body),"POST", map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"2585b79fdsgls7aab76c9dd","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" })
//	if err != nil {
//		ZapLog().Error( "send message to pastionpay get res ency data err", zap.Error(err))
//		return nil, err
//	}
//	json.Unmarshal(bastionPayRes, resEncyData)
//	//fmt.Println("[** send to bastion resEncyData** ]:",resEncyData.Signature, resEncyData.Code, resEncyData.EncryptData)
//	return resEncyData , nil
//}

func GetSHA256HashCode(message []byte) string {

	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	code := hex.EncodeToString(bytes)
	//返回哈希值
	return code

}

type Data struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
}

type Body struct {
	Code    int    `json:"code,omitempty"`
	Data    string `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResEncyData struct {
	Timestamp   int64  `json:"timestamp,omitempty"`
	Nonce       string `json:"nonce,omitempty"`
	EncryptData string `json:"data,omitempty"`
	Signature   string `json:"signature,omitempty"`
	Code        int64  `json:"code,omitempty"`
}

type ResDataInfo struct {
	//Nick              string  `json:"nick,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	//Country           string  `json:"country,omitempty"`
	//UserId            int64   `json:"user_id,omitempty"`
	//PhoneDistrict     string  `json:"phone_district,omitempty"`
	//Language          string  `json:"language,omitempty"`
	//Currency          string  `json:"currency,omitempty"`
	//Expire_in         int64   `json:"expire_in,omitempty"`
	//Nonce             string  `json:"nonce,omitempty"`
	Token string `json:"token,omitempty"`
	//Timestamp         int64   `json:"timestamp,omitempty"`
}

func TransferCoin(assets, pay_password, token, payee_id, request_no, amount string) (int64, error) {
	times := time.Now().Unix()
	timeStr := strconv.FormatInt(times, 10)
	nonce := "dhsjkf"
	remark := "lucky package"
	//aes的cbc加密数据
	secretKey, err := GetServerSecretKey()
	//fmt.Println("[--- transfer start ---]:",string(secretKey))

	data, _ := json.Marshal(map[string]interface{}{
		"timestamp":    timeStr,
		"nonce":        nonce,
		"assets":       assets,
		"amount":       amount,
		"request_no":   request_no,
		"remark":       remark,
		"pay_password": pay_password,
		"token":        token,
		"payee_id":     payee_id,
	})
	//fmt.Println("[**transfer data**]",string(data))
	encyData, err := utils.Aes128Encrypt([]byte(data), []byte(secretKey))

	//fmt.Println("****encyData**",string(encyData))
	//签名
	signStr := "url=/wallet/api/user/assets/pay/transfer&timestamp=" + timeStr + "&nonce=" + nonce + "&data=" + string(encyData) + "&" + string(secretKey)

	signString := GetSHA256HashCode([]byte(signStr))
	//将数据加密
	reqBody, _ := json.Marshal(map[string]interface{}{
		"timestamp": timeStr,
		"nonce":     nonce,
		"data":      string(encyData),
		"signature": signString,
	})
	//fmt.Println("**signString**",signString)
	//请求bastion
	result, err := SendToBastion(reqBody, Transfer_Url)
	if err != nil {
		ZapLog().Error("send req to bastion err", zap.Error(err))

		return 1, err
	}
	fmt.Println("[** result transfer** ] :", result.Signature, result.Nonce, result.Code, result.Timestamp)
	if result.Code == 1001 {
		return 1, err
	}

	//验证签名参数
	timestpStr := strconv.FormatInt(result.Timestamp, 10)
	var verifySign string
	if result.EncryptData != "" {
		verifySign = "url=" + Transfer_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&data=" + result.EncryptData + "&" + string(secretKey)
	}

	if result.EncryptData == "" {
		verifySign = "url=" + Transfer_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&" + string(secretKey)
	}
	//有的没有返回签名
	if result.Signature == "" {
		return 5001, err
	}
	sign2 := GetSHA256HashCode([]byte(verifySign))
	//fmt.Println("[**two sign transfer**]",result.Signature,sign2)

	if result.Signature == sign2 {
		if result.EncryptData != "" {
			//将加密数据解密，返回给H5
			decrpData, err := utils.Aes128Decrypt([]byte(result.EncryptData), []byte(secretKey))
			if err != nil {
				ZapLog().Error("decry return data err", zap.Error(err))
				return 5002, err
			}

			var index int

			for k, v := range decrpData {
				if v == 0 {
					index = k
					break
				}
			}

			dataInfo := new(TransResDataInfo)
			err = json.Unmarshal(decrpData[:index], dataInfo)

			if err != nil {
				ZapLog().Error("unmarshal return data err", zap.Error(err))
				return 5003, err
				//返回转账状态 1，处理中，2，成功，3，失败******
			}
			return dataInfo.Status, nil
		}
		if result.EncryptData == "" {
			return 5004, err
		}
	}
	if result.Signature != sign2 {
		ZapLog().Error("sign not equal", zap.Error(err))
		return 5005, err
	}
	return 5006, err
}

type TransResDataInfo struct {
	Status int64 `json:"status,omitempty"`
}

func ReadCsv(name string) [][]string {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("open err[%v]", err)
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	allRecord := make([][]string, 0)
	//var newRecord []string
	for i := 0; true; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Read %v err[%v]", i, err)
			return nil
		}
		for j := 0; j < len(record); j++ {
			record[j] = strings.Replace(record[j], " ", "", -1)
			record[j] = strings.Replace(record[j], "\r", "", -1)
			record[j] = strings.Replace(record[j], "\n", "", -1)
		}
		allRecord = append(allRecord, record)
	}
	return allRecord
}
