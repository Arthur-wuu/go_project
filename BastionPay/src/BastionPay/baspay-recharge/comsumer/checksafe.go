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

type CheckSafe struct {
}

//安全校验
func (l CheckSafe) Check(assets, token, payee_id string) (int64, error) {
	times := time.Now().Unix()
	timeStr := strconv.FormatInt(times, 10)
	nonce := "dhsjkf"
	//aes的cbc加密数据
	secretKey, err := GetServerSecretKey()
	fmt.Println("*** secretKey check***", string(secretKey))

	data, _ := json.Marshal(map[string]interface{}{
		"timestamp": timeStr,
		"nonce":     nonce,
		"assets":    assets,
		"token":     token,
		"payee_id":  payee_id,
	})
	fmt.Println("*** data check**", string(data))
	encyData, err := utils.Aes128Encrypt([]byte(data), []byte(secretKey))

	fmt.Println("****encyData check**", string(encyData))
	//签名
	signStr := "url=/wallet/api/user/assets/pay/confirm_safety&timestamp=" + timeStr + "&nonce=" + nonce + "&data=" + string(encyData) + "&" + string(secretKey)

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
	result, err := SendToBastion(reqBody, Check_Url)
	if err != nil {
		ZapLog().Error("send req to bastion err", zap.Error(err))

		return 1, err
	}
	fmt.Println("result :", result.Signature, result.Nonce, result.Code, result.Timestamp)
	if result.Code == 1001 {
		return 1, err
	}

	if result.Code == 0 {
		return 0, nil
	}

	//验证签名参数
	timestpStr := strconv.FormatInt(result.Timestamp, 10)
	var verifySign string
	if result.EncryptData != "" {
		verifySign = "url=" + Check_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&data=" + result.EncryptData + "&" + string(secretKey)
	}

	if result.EncryptData == "" {
		verifySign = "url=" + Check_Url + "&timestamp=" + timestpStr + "&nonce=" + result.Nonce + "&" + string(secretKey)
	}
	//有的没有返回签名
	if result.Signature == "" {
		return 5001, err
	}
	sign2 := GetSHA256HashCode([]byte(verifySign))
	fmt.Println("**two sign transfer**", result.Signature, sign2)

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
