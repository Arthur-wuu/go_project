package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/util"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var GLogin Login

type Login struct {
}


const(
	RsaBits1024 = 1024
	RsaBits2048 = 2048
	RsaEncodeLimit1024 = RsaBits1024 / 8 - 11
	RsaDecodeLimit1024 = RsaBits1024 / 8
	RsaEncodeLimit2048 = RsaBits2048 / 8 - 11
	RsaDecodeLimit2048 = RsaBits2048 / 8
)

const   (
	serverRsaUrl = "/wallet/api/security/handshake/server"
	clientRsaUrl = "/wallet/api/security/handshake/client"
	Login_Url = "/wallet/api/user/oauth/login"
	Transfer_Url = "/wallet/api/user/assets/pay/transfer"
	Refresh_Url = "/wallet/api/user/oauth/refresh_token"
	GetUid_Url = "/wallet/api/user/list_by_phone"
	Check_Url = "/wallet/api/user/assets/pay/confirm_safety"

)

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

var publicKey = []byte(`
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuW30rrCPsvjtXMtCEV7e
lJdQ81NC2r309zTItBx+0KOcvysUSs8lQTMWyONlLsM6RCZQoERUGuK+K+isyLNv
40Ng4DsbMCyehX+FrLsJ7O6tVjfHKB1OnnLqvOjfKToow7BU8uBZZgQTlyH7+QmX
3hx6hk2Zpzsz5Qv+Uqknp93DmP19OMCrcZubLg2laaAi2fUmBR2u6WWVXU4hRugS
xYTSQJHeXm2jcOA13XJsO5/RcJrZ8Xod81/6T0sHTt3Rpq/YAVldz/mMf+pjTmTc
7Ga3wf6X4Dfjop4ahwssn8KUGkZH0LVJYUsoTL6Z1XF2HLjuFk8gOgHF1QqcrAMX
lQIDAQAB
`)

//登录去拿token 为转账做准备
func (l Login )LoginBastionPay (phone, pwd string) (string,string, error){
	times :=time.Now().Unix()
	timeStr :=strconv.FormatInt(times,10)
	nonce:= "dhsjkf"

	//aes的cbc加密数据
	secretKey, err :=GetServerSecretKey()
	fmt.Println("[--- login start ---]:",string(secretKey))
	push_id := "1"
	data := "{\"timestamp\":\"" + timeStr + "\", \"nonce\":\"" + nonce + "\", \"phone\":\"" + phone + "\", \"password\":\"" + pwd + "\",\"push_id\":\"" + push_id + "\"}"
	encyData, err := utils.Aes128Encrypt([]byte(data), []byte(secretKey))

	//fmt.Println("****encyData**",string(encyData))
	//签名
	signStr := "url=/wallet/api/user/oauth/login&timestamp="+timeStr+"&nonce="+nonce+"&data="+string(encyData)+"&"+string(secretKey)
	//fmt.Println("***sign***",signStr)


	signString :=GetSHA256HashCode([]byte(signStr))
	//将数据加密
	reqBody, _ := json.Marshal(map[string]interface{}{
		"timestamp": timeStr,
		"nonce": nonce,
		"data": string(encyData),
		"signature": signString,
	})
	fmt.Println("**signString**",signString)
	//请求bastion
	result, err := SendToBastion(reqBody, Login_Url)
	if err != nil  {
		ZapLog().Error( "send req to bastion err", zap.Error(err))

		return "get token err1", "refresh token err1", err
	}
	fmt.Println("[** result login ** ] :",result.Signature,result.Nonce,result.Code,result.Timestamp)
	if result.Code == 1001  {
		return "get token err2","refresh token err2", err
	}


	//验证签名参数
	timestpStr := strconv.FormatInt(result.Timestamp,10)
	var verifySign string
	if result.EncryptData != "" {
		verifySign = "url="+Login_Url+"&timestamp="+timestpStr+"&nonce="+result.Nonce+"&data="+result.EncryptData+"&"+string(secretKey)
	}

	if result.EncryptData == "" {
		verifySign = "url="+Login_Url+"&timestamp="+timestpStr+"&nonce="+result.Nonce+"&"+string(secretKey)
	}
	//有的没有返回签名
	if result.Signature == ""{
		return "get token err3","refresh token err3", err
	}
	sign2 :=GetSHA256HashCode([]byte(verifySign))
	//fmt.Println("**two sign login **",result.Signature,sign2)

	if result.Signature == sign2 {
		if result.EncryptData != ""{
			//将加密数据解密，返回给H5
			decrpData, err := utils.Aes128Decrypt([]byte(result.EncryptData),[]byte(secretKey))
			if err != nil {
				ZapLog().Error( "decry return data err", zap.Error(err))

				return "get token err4","refresh token err4", err
			}

			var index int

			for k, v := range decrpData {
				if v == 0 {
					index = k; break;
				}
			}

			dataInfo := new(ResDataInfo)
			err = json.Unmarshal(decrpData[:index],dataInfo)

			if err != nil {
				ZapLog().Error( "unmarshal return data err", zap.Error(err))
				//返回token ******
			}
			return dataInfo.Token, dataInfo.RefreshToken,  nil
		}
		if result.EncryptData == ""{
			return "no data back","refresh token err", err
		}

	}
	if result.Signature != sign2 {
		ZapLog().Error( "sign not equal", zap.Error(err))
		return "get token err5","refresh token err5", err
	}

	return "get token err6","refresh token err6", err
}



//上传客户端 公钥 并申请SecretKey
func GetServerSecretKey() ([]byte, error) {
	baspayRes := new(Body)
	//加密因子
	//secretKeyRamdom :=GenRandomString(32)

	body :=[]byte("client_public_key="+string(publicKey)+"&secrect_key_ramdom=35321243")

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
		ZapLog().Error( "rsa encybody  err")
	}
	//fmt.Println("****encyBody** 加密body，**",string(encyBody))

	//请求服务端的 得到SecretKey的密文
	bastionPayRes, err := base.HttpSendSer(config.GConfig.BastionpayUrl.BastionUser+clientRsaUrl, bytes.NewBuffer(encyBody),"POST", map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":config.GConfig.Login.DeviceId,"Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" })//
	//fmt.Println("****bastionPayRes** ，**",string(bastionPayRes))
	if err != nil {
		ZapLog().Error( "send message to pastionpay get secret key err", zap.Error(err))
		return nil, err
	}
	json.Unmarshal(bastionPayRes, baspayRes)
	//fmt.Println("****baspayRes** ，**",baspayRes.Data)
	//客户端私钥解密 得到SecretKey的密文
	secretKey, err  := utils.RsaDecrypt3([]byte(baspayRes.Data), privateKey)

	//fmt.Println("****secretKey** ，**",string(secretKey))
	return secretKey, nil
}


func SendToBastion (body []byte  ,url string ) (*ResEncyData, error) {
	resEncyData := new(ResEncyData)
	fmt.Println("[** send to bastion url** ]:",config.GConfig.BastionpayUrl.Bastionurl+url)
	bastionPayRes, err := base.HttpSendSer(config.GConfig.BastionpayUrl.BastionUser+url, bytes.NewBuffer(body),"POST", map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":config.GConfig.Login.DeviceId,"Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" })
	if err != nil {
		ZapLog().Error( "send message to pastionpay get res ency data err", zap.Error(err))
		return nil, err
	}
	json.Unmarshal(bastionPayRes, resEncyData)

	return resEncyData , nil
}


func SendToBastionLoginPool (body []byte  ,url string, deviceId string ) (*ResEncyData, error) {
	resEncyData := new(ResEncyData)
	fmt.Println("[** device id** ]:",deviceId)
	bastionPayRes, err := base.HttpSendSer(config.GConfig.BastionpayUrl.BastionUser+url, bytes.NewBuffer(body),"POST", map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":deviceId,"Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" })
	if err != nil {
		ZapLog().Error( "send message to pastionpay get res ency data err", zap.Error(err))
		return nil, err
	}
	json.Unmarshal(bastionPayRes, resEncyData)
	//fmt.Println("[** send to bastion resEncyData** ]:",resEncyData.Signature, resEncyData.Code, resEncyData.EncryptData)
	return resEncyData , nil
}


func GetSHA256HashCode(message []byte)string{

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

//type Data struct{
//	Timestamp  int64   `json:"timestamp,omitempty"`
//	Nonce      string  `json:"nonce,omitempty"`
//}

type Body struct{
	Code     int         `json:"code,omitempty"`
	Data     string      `json:"data,omitempty"`
	Message  string      `json:"message,omitempty"`
}

type ResEncyData struct {
	Timestamp   int64  `json:"timestamp,omitempty"`
	Nonce       string `json:"nonce,omitempty"`
	EncryptData string `json:"data,omitempty"`
	Signature   string `json:"signature,omitempty"`
	Code        int64 `json:"code,omitempty"`
}

type ResDataInfo struct {
	//Nick              string  `json:"nick,omitempty"`
	RefreshToken      string  `json:"refresh_token,omitempty"`
	//Country           string  `json:"country,omitempty"`
	//UserId            int64   `json:"user_id,omitempty"`
	//PhoneDistrict     string  `json:"phone_district,omitempty"`
	//Language          string  `json:"language,omitempty"`
	//Currency          string  `json:"currency,omitempty"`
	//Expire_in         int64   `json:"expire_in,omitempty"`
	//Nonce             string  `json:"nonce,omitempty"`
	Token             string  `json:"token,omitempty"`
	//Timestamp         int64   `json:"timestamp,omitempty"`
}



