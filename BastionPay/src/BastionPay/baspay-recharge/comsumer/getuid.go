package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/baspay-recharge/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type GetUid struct {
}

//刷新token
func (l GetUid) GetUidByPhone(phone, code, token string) (string, string, error) {
	times := time.Now().Unix()
	timeStr := strconv.FormatInt(times, 10)
	nonce := "dhsjkf"

	//aes的cbc加密数据
	secretKey, err := GetServerSecretKey()
	fmt.Println("[--- get uid by phone start ---]", string(secretKey))

	data, _ := json.Marshal(map[string]interface{}{
		"timestamp":         timeStr,
		"nonce":             nonce,
		"phone":             phone,
		"token":             token,
		"phone_district_no": code,
	})
	encyData, err := utils.Aes128Encrypt([]byte(data), []byte(secretKey))

	//fmt.Println("****GetUidByPhone Data**",string(data))
	//签名
	signStr := "url=/wallet/api/user/list_by_phone&timestamp=" + timeStr + "&nonce=" + nonce + "&data=" + string(encyData) + "&" + string(secretKey)

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
	result, err := SendToBastion(reqBody, GetUid_Url)
	if err != nil {
		ZapLog().Error("send req to bastion err", zap.Error(err))

		return "", "", err
	}
	fmt.Println("[** get id result ** ]:", result.Signature, result.Nonce, result.Code, result.Timestamp)
	if result.Code == 1001 {
		return "", "", err
	}
	if result.Code == 1014 {
		return "", "", err
	}

	//验证签名参数
	timestpStr := strconv.FormatInt(result.Timestamp, 10)
	var verifySign string
	if result.EncryptData != "" {
		verifySign = "url=" + GetUid_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&data=" + result.EncryptData + "&" + string(secretKey)
	}

	if result.EncryptData == "" {
		verifySign = "url=" + GetUid_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&" + string(secretKey)
	}
	//有的没有返回签名
	if result.Signature == "" {
		return "", "", err
	}
	sign2 := GetSHA256HashCode([]byte(verifySign))
	//fmt.Println("**two sign GetUidByPhone**",result.Signature,sign2)

	if result.Signature == sign2 {
		if result.EncryptData != "" {
			//将加密数据解密，返回给H5
			decrpData, err := utils.Aes128Decrypt([]byte(result.EncryptData), []byte(secretKey))
			if err != nil {
				ZapLog().Error("decry return data err", zap.Error(err))
				return "", "", err
			}
			var index int

			for k, v := range decrpData {
				if v == 0 {
					index = k
					break
				}
			}

			dataInfo := new(GetUidDataInfo)
			err = json.Unmarshal(decrpData[:index], dataInfo)

			if err != nil {
				ZapLog().Error("unmarshal return data err", zap.Error(err))
				return dataInfo.UserId, dataInfo.RegistTime, err //返回转账状态 1，处理中，2，成功，3，失败******
			}
			return dataInfo.UserId, dataInfo.RegistTime, nil
		}
		if result.EncryptData == "" {
			return "", "", err
		}
	}
	if result.Signature != sign2 {
		ZapLog().Error("sign not equal", zap.Error(err))
		return "", "", err
	}
	return "", "", err
}

type GetUidDataInfo struct {
	UserId     string `json:"user_id,omitempty"`
	RegistTime string `json:"regist_time,omitempty"`
	//Token           string  `json:"token,omitempty"`
}
