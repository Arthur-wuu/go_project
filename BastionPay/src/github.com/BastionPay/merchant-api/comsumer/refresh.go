package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var GRefresh Refresh

type Refresh struct {

}

//刷新token
func (l Refresh ) RefreshToken (token, refreshToken string ) (string,string,int64, error){
	times :=time.Now().Unix()
	timeStr :=strconv.FormatInt(times,10)
	nonce:= "dhsjkf"

	//aes的cbc加密数据
	secretKey, err :=GetServerSecretKey()
	fmt.Println("[--- refresh token start ---]:",string(secretKey))

	data, _ := json.Marshal(map[string]interface{}{
		"timestamp": timeStr,
		"nonce": nonce,
		"token": token,
		"refresh_token": refreshToken,
	})
	encyData, err := utils.Aes128Encrypt([]byte(data), []byte(secretKey))

	//fmt.Println("****encyData**",string(encyData))
	//签名
	signStr := "url=/wallet/api/user/oauth/refresh_token&timestamp="+timeStr+"&nonce="+nonce+"&data="+string(encyData)+"&"+string(secretKey)
	//fmt.Println("***sign***",signStr)


	signString :=GetSHA256HashCode([]byte(signStr))
	//将数据加密
	reqBody, _ := json.Marshal(map[string]interface{}{
		"timestamp": timeStr,
		"nonce": nonce,
		"data": string(encyData),
		"signature": signString,
	})
	//fmt.Println("**signString**",signString)
	//请求bastion
	result, err := SendToBastion(reqBody, Refresh_Url)
	if err != nil  {
		ZapLog().Error( "send req to bastion err", zap.Error(err))

		return "","",0, err
	}
	fmt.Println("[** result refresh token **] :",result.Signature,result.Nonce,result.Code,result.Timestamp)
	if result.Code == 1001  {
		return "","", 1001, err
	}
	if result.Code == 1014  {
		return "","", 1014,  err
	}


	//验证签名参数
	timestpStr := strconv.FormatInt(result.Timestamp,10)
	var verifySign string
	if result.EncryptData != "" {
		verifySign = "url="+Refresh_Url+"&timestamp="+timestpStr+"&nonce="+result.Nonce+"&data="+result.EncryptData+"&"+string(secretKey)
	}

	if result.EncryptData == "" {
		verifySign = "url="+Refresh_Url+"&timestamp="+timestpStr+"&nonce="+result.Nonce+"&"+string(secretKey)
	}
	//有的没有返回签名
	if result.Signature == ""{
		return "","", 0, err
	}
	sign2 :=GetSHA256HashCode([]byte(verifySign))
	//fmt.Println("**two sign refresh**",result.Signature,sign2)

	if result.Signature == sign2 {
		if result.EncryptData != ""{
			//将加密数据解密，返回给H5
			decrpData, err := utils.Aes128Decrypt([]byte(result.EncryptData),[]byte(secretKey))
			if err != nil {
				ZapLog().Error( "decry return data err", zap.Error(err))
				return "","", 0, err
			}
			var index int

			for k, v := range decrpData {
				if v == 0 {
					index = k; break;
				}
			}

			dataInfo := new(RefreshDataInfo)
			err = json.Unmarshal(decrpData[:index],dataInfo)

			if err != nil {
				ZapLog().Error( "unmarshal return data err", zap.Error(err))

				return "","", 0, err      //返回转账状态 1，处理中，2，成功，3，失败******
			}
			return dataInfo.Token, dataInfo.Refresh_token, 0, nil
		}
		if result.EncryptData == ""{
			return "","", 0, err
		}

	}
	if result.Signature != sign2 {
		ZapLog().Error( "sign not equal", zap.Error(err))
		return "","", 0, err
	}

	return "","", 0, err
}



type RefreshDataInfo struct {
	Refresh_token   string  `json:"refresh_token,omitempty"`
	Token           string  `json:"token,omitempty"`
}

