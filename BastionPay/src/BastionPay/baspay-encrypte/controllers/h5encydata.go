package controllers

//
//import (
//	"BastionPay/bas-api/apibackend"
//	"BastionPay/baspay-encrypte/comsumer"
//	"BastionPay/baspay-encrypte/util"
//	"encoding/json"
//	"fmt"
//	"github.com/kataras/iris"
//	"go.uber.org/zap"
//	"io/ioutil"
//	"sort"
//	"strconv"
//	//"go.uber.org/zap"
//	. "BastionPay/bas-base/log/zap"
//)
//
//var GH5EncyData H5EncyData
//type H5EncyData struct {
//	Controllers
//}
////const(
////	RsaBits1024 = 1024
////	RsaBits2048 = 2048
////	RsaEncodeLimit1024 = RsaBits1024 / 8 - 11
////	RsaDecodeLimit1024 = RsaBits1024 / 8
////	RsaEncodeLimit2048 = RsaBits2048 / 8 - 11
////	RsaDecodeLimit2048 = RsaBits2048 / 8
////)
//
////const   (
////	serverRsaUrl = "/wallet/api/security/handshake/server"
////	clientRsaUrl = "/wallet/api/security/handshake/client"
////)
//
//var h5PrivateKey = []byte(`
//-----BEGIN RSA PRIVATE KEY-----
//MIIEpAIBAAKCAQEAuW30rrCPsvjtXMtCEV7elJdQ81NC2r309zTItBx+0KOcvysU
//Ss8lQTMWyONlLsM6RCZQoERUGuK+K+isyLNv40Ng4DsbMCyehX+FrLsJ7O6tVjfH
//KB1OnnLqvOjfKToow7BU8uBZZgQTlyH7+QmX3hx6hk2Zpzsz5Qv+Uqknp93DmP19
//OMCrcZubLg2laaAi2fUmBR2u6WWVXU4hRugSxYTSQJHeXm2jcOA13XJsO5/RcJrZ
//8Xod81/6T0sHTt3Rpq/YAVldz/mMf+pjTmTc7Ga3wf6X4Dfjop4ahwssn8KUGkZH
//0LVJYUsoTL6Z1XF2HLjuFk8gOgHF1QqcrAMXlQIDAQABAoIBAAbgGth3V3ytWi+8
//oaB/QgWEbs324l21+WVJIb/75n/Z8S/tav0zHRCD0m50cAYB2yE3wnmOzweILMD4
//CaAPdzCYCPmHd4Sbwuz0Q4KaFM4iM28k9k94drfNTwLeJ5ghve0zlpLC2s27jPZn
//ni4nuufTan+cVPwsZ8owXae2+e9xa1qy1gJBveXT26kGPbo6q06zK0kfcwKs5cYG
//o12tnvBzDv2h2TVoXujhADKNWNIAOkTmOk3VSAb+bTiNCipep7SbgeToYnX5rqx/
//SHIFjFlwxMYqiadG8Ysq2hK0bg8HMMvK7L2PNQ/Ge4B5aXXJHxF4m+e+2EjL3uJl
//VD2stGkCgYEA71oyBzPmbm4u/fhZsedwtCdHnfkQMFMje14sQiPbPonlWxo417cD
//pwcPU8bd/s5X8saGim+1p3WShnY1+wHH/VF5aiFqwt/K69c3ofAfrmbmXbyN2gJJ
//umuAR7Hb8Urd2teSova4LC4Fc7sSvtLXPle3eWaoD0mO1sxlQ37pSDcCgYEAxlOi
//TNFB4uEbIJfmOqzqXywzibOxMasgjZLKEu+JE6Vk8y3YMzCfrlKl/2wr/I+DrhWx
//yX7cf0/SfJgox0D9K23yh/Uu72QNRMZpVSRUGjCurtWvILRcD0mdzFkHijC7cfDE
//3XqMHRjFrRgAWlt01YVojbjYJ/F8cudVNyXWYJMCgYEAuDrofvrHxwAwU3OxNmo6
//KbCCQ2nNuCSGDxMxZcdLnhtt2m2YixFnUkzw0z8i6FnTAB8mt6+8VqT8n1qlugpo
//8OahWbtW/aBcBKOnQpIdEJRLhKL5XHCeZ0sPdh/Edzl1AlkjmSPmJrtVnvrDNvX6
//jxXdNyh4+ytXMqYo24b38IkCgYBkZ9UEJPDBRwuvzZcuX3psYnlZHpL3vVZGtmj9
//ey2ft51LDAunptdArvEBRidivtmAmdUfWM2S2ruKfpIuhjVl9kzSDgwMAFBDYFvV
//UgYOGFVniCEYYpc02iU8XlpV2OQdBDL2meMzm+YAAuWy2RhmPRs4nLs6RaSmm31l
//5Q8KZwKBgQCVp0hRdSphZBjcUEAJyBHDcf3D4f0nJx3CLvmNg6TuoT+5hendCi2E
//dArT9yVYco35HJd0nxBm/UOzEptNUuAeX41AVeN9o0UYoJomt/l07cyY77HtrnIb
//ZpqTdTAmvK7JaZooBF736k+gTMmX0qdzfmoJSDNqHZaoq4tmuojZZg==
//-----END RSA PRIVATE KEY-----
//`)
//
//var h5PublicKey = []byte(`
//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuW30rrCPsvjtXMtCEV7e
//lJdQ81NC2r309zTItBx+0KOcvysUSs8lQTMWyONlLsM6RCZQoERUGuK+K+isyLNv
//40Ng4DsbMCyehX+FrLsJ7O6tVjfHKB1OnnLqvOjfKToow7BU8uBZZgQTlyH7+QmX
//3hx6hk2Zpzsz5Qv+Uqknp93DmP19OMCrcZubLg2laaAi2fUmBR2u6WWVXU4hRugS
//xYTSQJHeXm2jcOA13XJsO5/RcJrZ8Xod81/6T0sHTt3Rpq/YAVldz/mMf+pjTmTc
//7Ga3wf6X4Dfjop4ahwssn8KUGkZH0LVJYUsoTL6Z1XF2HLjuFk8gOgHF1QqcrAMX
//lQIDAQAB
//`)
//
////
//
////var privateKey, publicKey []byte
////
////func (ency *EncyData) Init() {
////	var err error
////
////	//生成公钥和私钥的文件
////	err = utils.GenRsaKey(RsaBits2048)
////	if err != nil {
////		ZapLog().Error( "gen pub and pri pem file err", zap.Error(err))
////	}
////
////	publicKey, err = ioutil.ReadFile("public.pem")
////	if err != nil {
////		os.Exit(-1)
////	}
////	privateKey, err = ioutil.ReadFile("private.pem")
////	if err != nil {
////		os.Exit(-1)
////	}
////}
//
////H5 得到服务端处理后的数据
//func (this *H5EncyData) H5Get(ctx iris.Context) {
//	//param := new(comsumer.ReqEncyData)
//
//	requestParams := make(map[string]interface{},0)
//
//	data, err :=ioutil.ReadAll(ctx.Request().Body)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("Body readAll  err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//		return
//	}
//
//	err = json.Unmarshal(data, requestParams)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("requestbody to requestParams err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//		return
//	}
//
//	//将param的key排序，
//	keysSort := make([]string, 0)
//	for k, _ := range requestParams{
//		keysSort = append(keysSort, k)
//	}
//	sort.Strings(keysSort)
//
//	//拼接签名字符串
//	signH5Str := ""
//	for _, v := range keysSort {
//		signH5Str = v+"="+requestParams[v].(string)+"&"
//	}
//	signH5Str = signH5Str[0:len(signH5Str)-1]
//
//
//
//
//
//
//
//
//
//
//
//	fmt.Println("*data*",string(data))
//	//时间戳，nonce值 是拿用户请求传进来的
//	//analysisData := new(Data)
//	//json.Unmarshal(data, analysisData)
//	//
//	//times := analysisData.Timestamp
//	//timeStr :=strconv.FormatInt(1551868437,10)
//	timeStr := "1551868437"
//	nonce:= "dhsjkf"
//
//	//aes的cbc加密数据
//	secretKey, err :=GetServerSecretKey()
//	//secretKey := "d9-9178-23224d464225424362556446"
//	fmt.Println("***secretKey***",secretKey)
//	encyData, err := utils.Aes128Encrypt(data, []byte(secretKey))
//
//	//fmt.Println("***encyData***",string(encyData))
//	//签名
//	signStr := "url="+ctx.Path()+"&timestamp="+timeStr+"&nonce="+nonce+"&data="+string(encyData)+"&"+string(secretKey)
//	fmt.Println("***sign***",signStr)
//
//	//hash := sha256.New()
//	//hash.Write([]byte(signStr))
//	//sign1 := hash.Sum(nil)
//	signString :=GetSHA256HashCode([]byte(signStr))
//	//将数据加密
//	reqBody, _ := json.Marshal(map[string]interface{}{
//		"timestamp": timeStr,
//		"nonce": nonce,
//		"data": string(encyData),
//		"signature": signString,
//	})
//	fmt.Println("**signString**",signString)
//	//请求bastion
//	result, err := comsumer.SendToBastion(reqBody, ctx.Path())
//	if err != nil  {
//		ZapLog().Error( "send req to bastion err", zap.Error(err))
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//		return
//	}
//
//	if result.Code != 0  {
//		this.Response(ctx, "no data back")
//		return
//	}
//
//	//验证签名参数
//	timestpStr := strconv.FormatInt(result.Timestamp,10)
//	var verifySign string
//	verifySign = "url="+ctx.Path()+"&timestamp="+timestpStr+"&nonce="+result.Nonce+"&data="+result.EncryptData+"&"+string(secretKey)
//
//	if result.EncryptData == "" {
//		verifySign = "url="+ctx.Path()+"&timestamp="+timestpStr+"&nonce="+result.Nonce+"&"+string(secretKey)
//	}
//
//	sign2 :=GetSHA256HashCode([]byte(verifySign))
//	fmt.Println("**two sign**",result.Signature,sign2)
//	if result.Signature == sign2 {
//		this.Response(ctx, "nil data back")
//		return
//	}
//	if result.Signature != sign2 {
//		ZapLog().Error( "sign not equal", zap.Error(err))
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//		return
//	}
//
//	//将加密数据解密，返回给H5
//	decrpData, err := utils.Aes128Decrypt([]byte(result.EncryptData),[]byte(secretKey))
//	if err != nil {
//		ZapLog().Error( "decry return data err", zap.Error(err))
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//		return
//	}
//	dataInfo := new(comsumer.ResEncyData)
//	err = json.Unmarshal(decrpData,dataInfo)
//	if err != nil {
//		ZapLog().Error( "unmarshal return data err", zap.Error(err))
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//		return
//	}
//
//	//if dataInfo.Timestamp != result.Timestamp || dataInfo.Nonce != result.Nonce {
//	//	ZapLog().Error( "timestamp or nonce not equal", zap.Error(err))
//	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//	//	return
//	//}
//	this.Response(ctx, dataInfo.EncryptData)
//}
//
//
////到服务端拿RSA的公钥
////func GetServerRsaPub() (string, error) {
////	baspayRes := new(Body)
////	bastionPayRes, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+serverRsaUrl,nil,"POST",map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"32132149078","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" })
////	if err != nil {
////		ZapLog().Error( "send message to pastionpay get rsa pub err", zap.Error(err))
////		return "", err
////	}
////	json.Unmarshal(bastionPayRes, baspayRes)
////	serverPub := baspayRes.Data
////
////	return serverPub, nil
////}
//
////生成本地的公私钥
////func GenRsaPairKey()(string, string, error){
////	pri, pub, err :=utils.MakeSSHKeyPair()
////	if err != nil {
////		ZapLog().Error( "gen rsa pub pri err", zap.Error(err))
////		return "", "", err
////	}
////	return pri, pub, nil
////}
//
//
////上传客户端 公钥 并申请SecretKey
////func GetServerSecretKey() ([]byte, error) {
////	baspayRes := new(Body)
////	//加密因子
////	//secretKeyRamdom :=GenRandomString(32)
////
////	body :=[]byte("client_public_key="+string(publicKey)+"&secrect_key_ramdom=35321243")
////
////	fmt.Println("****body body**",string(body))
////	//拿服务端的公钥 加密body
////	//serverPub, err := GetServerRsaPub()
////	//fmt.Println("****serverPub pubkey",serverPub)
////	//if err != nil {
////	//	ZapLog().Error( "get server pub err", zap.Error(err))
////	//}
////
////	serverPubByte := []byte(`
////-----BEGIN PUBLIC KEY-----
////MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArdDMsRHlTzzf8rfryGo8
////2NDQ6VntnD07ax+7CMsKAAlICv28NxLHPoWRZAl9dRhM/uWGpgOPs2sKDayilyyR
////0gZ8NPIVU4AWmn4xnv5l4Vu5HND9DcIoyvHLCiel+Lj/6HcpUzlJ+GmJ6L0QO/PI
////CPq4KyR24ggCfknzAfLi8DQ+LUGFOhiSnu1ta3z4rVeOIyy72thlGoN7aTxXSMe6
////yTi1bshkmFLgHyOcM2vpx4Vhtfb7xfu77LkRQEwi2k4vIZozInp4s5UaVFstd/Zd
////IM/hMlwKP5zv4caLhI6Op3PrG+/6McLhx3j4tRxZhc6IdfSpvzEqO7icD+oRa5Sd
////DwIDAQAB
////-----END PUBLIC KEY-----
////`)
////
////	encyBody, err := utils.RsaEncrypt(body, serverPubByte, RsaEncodeLimit2048)
////
////	if err != nil {
////		ZapLog().Error( "rsa encybody  err")
////	}
////	fmt.Println("****encyBody** 加密body，**",string(encyBody))
////
////	//请求服务端的 得到SecretKey的密文
////	bastionPayRes, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+clientRsaUrl, bytes.NewBuffer(encyBody),"POST",map[string]string{"Client":"1", "DeviceType":"1", "DeviceName":"huawei","DeviceId":"32132149078","Version":"1.1.0","Content-Type":"application/json;charset=UTF-8" })
////	fmt.Println("****bastionPayRes** ，**",string(bastionPayRes))
////	if err != nil {
////		ZapLog().Error( "send message to pastionpay get secret key err", zap.Error(err))
////		return nil, err
////	}
////	json.Unmarshal(bastionPayRes, baspayRes)
////	fmt.Println("****baspayRes** ，**",baspayRes.Data)
////	//客户端私钥解密 得到SecretKey的密文
////	secretKey, err  := utils.RsaDecrypt3([]byte(baspayRes.Data), privateKey)
////
////	fmt.Println("****secretKey** ，**",string(secretKey))
////	return secretKey, nil
////}
//
//
////
////type Data struct{
////	Timestamp  int64   `json:"timestamp,omitempty"`
////	Nonce      string  `json:"nonce,omitempty"`
////}
////
////type Body struct{
////	Code     int         `json:"code,omitempty"`
////	Data     string      `json:"data,omitempty"`
////	Message  string      `json:"message,omitempty"`
////}
////
//////生成随机字符串
////func GenRandomString(l int) string{
////	str := "123456789"
////	bt := []byte(str)
////	result := []byte{}
////	r := rand.New(rand.NewSource(time.Now().UnixNano()))
////	for i := 0; i < l; i++ {
////		result = append(result, bt[r.Intn(len(bt))])
////	}
////	return  string(result)
////}
////
////func GetSHA256HashCode(message []byte)string{
////
////	hash := sha256.New()
////	//输入数据
////	hash.Write(message)
////	//计算哈希值
////	bytes := hash.Sum(nil)
////	//将字符串编码为16进制格式,返回字符串
////	code := hex.EncodeToString(bytes)
////	//返回哈希值
////	return code
////
////}
