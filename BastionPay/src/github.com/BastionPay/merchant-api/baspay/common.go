package baspay

import (
	. "BastionPay/bas-base/log/zap"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"go.uber.org/zap"
	"sort"
)

type(
	Request struct{
		MerchantId *string  `json:"merchant_id,omitempty"`
		SignType   *string  `json:"sign_type,omitempty"`
		Signature  *string  `json:"signature,omitempty"`
		Timestamp  *string  `json:"timestamp,omitempty"`
		NotifyUrl  *string  `json:"notify_url,omitempty"`
	}

	Response struct{
		Code       int       `json:"code,omitempty"`
		Message    string    `json:"message,omitempty"`
		Signature  string    `json:"signature,omitempty"`
	}
)


func RequestBodyToSignStr (body []byte) (string){
	requestParams := make(map[string]string,0)

	err := json.Unmarshal(body, &requestParams)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("requestbody to requestParams err")
		return ""
	}
	//将param的key排序，
	keysSort := make([]string, 0)
	for k, _ := range requestParams{
		keysSort = append(keysSort, k)
	}
	sort.Strings(keysSort)
	//拼接签名字符串
	signH5Str := ""
	for i:=0; i<len(keysSort); i++ {
		signH5Str += keysSort[i]+"="+requestParams[keysSort[i]]+"&"
	}
	signH5Str = signH5Str[0:len(signH5Str)-1]
	return  signH5Str
}



// 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func RsaSignWithSha1Hex(data string, prvKey string) (string, error) {
	block, _ := pem.Decode([]byte(prvKey))
	if block == nil { // 失败情况
		fmt.Println("En err")
	}

	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	fmt.Println("signature***",signature)

	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := hex.EncodeToString(signature)
	return out, nil
}


//  验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func RsaVerySignWithSha1Base64(originalData, signData, pubKey string) error{
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	public, _ := base64.StdEncoding.DecodeString(pubKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

func Sha1(data string) []byte {

	// stage1Hash = SHA1(password)
	crypt := sha1.New()
	crypt.Write([]byte(data))
	stage1 := crypt.Sum(nil)

	return stage1
}