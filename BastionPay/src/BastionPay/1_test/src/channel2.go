package main

import (
	"crypto"
	"net/url"


	//"crypto/sha256"
	"encoding/base64"
	//"net/url"
	"os"
	//"strings"

	//"strings"
	"crypto/rand"
	//"net/url"

	//"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	//"strings"

	//"time"
)
//
//func main() {
//	go func() {
//		time.Sleep(1 * time.Hour)
//	}()
//	c := make(chan int)
//	go func() {
//		for i := 0; i < 10; i = i + 1 {
//			c <- i
//		}
//		close(c)
//	}()
//	for i := range c {
//		fmt.Println(i)
//	}
//	fmt.Println("Finished")
//}

var h5PublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuW30rrCPsvjtXMtCEV7e
lJdQ81NC2r309zTItBx+0KOcvysUSs8lQTMWyONlLsM6RCZQoERUGuK+K+isyLNv
40Ng4DsbMCyehX+FrLsJ7O6tVjfHKB1OnnLqvOjfKToow7BU8uBZZgQTlyH7+QmX
3hx6hk2Zpzsz5Qv+Uqknp93DmP19OMCrcZubLg2laaAi2fUmBR2u6WWVXU4hRugS
xYTSQJHeXm2jcOA13XJsO5/RcJrZ8Xod81/6T0sHTt3Rpq/YAVldz/mMf+pjTmTc
7Ga3wf6X4Dfjop4ahwssn8KUGkZH0LVJYUsoTL6Z1XF2HLjuFk8gOgHF1QqcrAMX
lQIDAQAB
-----END PUBLIC KEY-----
`)

var h5PrivateKeys = []byte(`-----BEGIN RSA PRIVATE KEY-----
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

var testPrivateKeys = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCbZbLvSvijwn3W
xHrcqpEGoPwtBYwvvh6qlG4cN7qNAng352Y368WX4jhrptVzigCDZoE/sHVBbFEt
AvTMPsD6e4AQ9hW0QQxg2DCjSWbGnC+n1UvOSYtOSOzGrm1fxQcNlq+QCjQGqLU+
b13eo+AyenPc37wyvA4lEE/nEjG0aDxsowYoh8gAe3mBRAhsmZ+6rMTapXrHe7vD
S14FSz0g/AnWAHua+y7TNnERcETyR7EA/SIdfovaxcLvyj0OR5gdDf75vuaKBTIT
vmVWPIciabbY+ZU+4gvq/w50I0exCev76Mp6r1CCluDMoxEn2q/bQ6I+g21caZo6
nnBcdDqnAgMBAAECggEAdXr7PeFF/DfrftRnti/VGFfYjgjlpKps8LTqUbbn9/bo
AAuWwawjY+IImYo1UPYB0VVLXWUAIIfNDuRvQYInzrZTaX9BhVawDv8iNjAl3Pzz
IkUk3D3JbVPAfawc0AxaerFy5Mhx8J7W9u6m3syxkDf3JAKZexmk7+xXG/ArV6FN2
NvblT5mgosMTYFjYGsU1IKvnPRw3uHwLRCeTuT+HtxRyrtwp0eZ0qYBRYmhvODh/G
U8mQA77aj9h9bWvN6vav2f3L0cDJdlO5xAUDLVfqr1Z855QJwuuYrFsMW+RLd87h/
pCtVx8eku+4fYXYGDn7j+aHLwINTxQBS6j8hwUQKBgQDj7DLCYUFqTq/VNvdWo6Tx
17LsHcU2WJo17hsnH0oQMcFURcfJ/hAUrCGpkAI3YQrPvr6GZY7OFcVAGEoPAKadS
iyM9rCagWXyy2Co5r+uimQjnGNQmPw3YcLVZ28wBqGOCRid+zPj3Ir/I6F7mnsmPb
54KLEPJaxPKdM+WG3lqwKBgQCuilRccy2gyYY0Uf/4t6eSc9FD2Y1ZfJRVfGVdVR4
YPJDBc8F3GV4MzNDoOZBS36sbPuE3YDKt1ddln6+Ouky/Ss9Izx+SM7JDKbqglWA7
t/DNw3b7cAUlKXwMbjTDkQoVOzfOR9T5LoT9R2DSQasNJfRdt0fpkRUkBOOm+KxK9
QKBgHhBHR711V/TmG40jBeIS/TVy69MncrowKSHtofTuG4G8mwWTS1EARQHJdOjCa
hSaTPm/ftHBiuxzNredeSogUAn7I2Lcu5yK2oI6Dz1Ulky51bqonPZ4+kMiZGy+zU
pqn+YSQbBjUVCDYxELmVawnMQzLf1MEY/qEQ0WyJf4cv1AoGADx6CckOz5yKt0mhs
APJ/vIr1zKfSu7az7rfI3A3cfoL4kxlg3909rWQskIE0BEnFu6V1wuM9YJuOfgoYH
gf7T/K+A/OVK4f44CKEPRbTcDjdziUpcFxixbZTPYxqW6p7sh0gF2lXhIJIGNyPAY
eYtpncEiYnP49Gwoj941/VJOUCgYEAqXI37ZwLD9zz1qynVpoYKMzMOaBSjLgKLRk
2v+UHW9dah4NMyOWWeV03NH42y9opvTBfexnHzSNig4vfzUiXHwq3X9cl6E8EOzr0
UYLraNxFKwBzyjvMjXBk36D7eHuXNw/gHRkiSPkr6aWu56S2Ps9geHEc311YJKi6fbyK9pY=
-----END RSA PRIVATE KEY-----
`)

var testPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAm2Wy70r4o8J91sR63KqRB
qD8LQWML74eqpRuHDe6jQJ4N+dmN+vFl+I4a6bVc4oAg2aBP7B1QWxRLQL0zD7A+n
uAEPYVtEEMYNgwo0lmxpwvp9VLzkmLTkjsxq5tX8UHDZavkAo0Bqi1Pm9d3qPgMnpz
3N+8MrwOJRBP5xIxtGg8bKMGKIfIAHt5gUQIbJmfuqzE2qV6x3u7w0teBUs9IPwJ1g
B7mvsu0zZxEXBE8kexAP0iHX6L2sXC78o9DkeYHQ3++b7migUyE75lVjyHImm22PmV
PuIL6v8OdCNHsQnr++jKeq9QgpbgzKMRJ9qv20OiPoNtXGmaOp5wXHQ6pwIDAQAB
-----END PUBLIC KEY-----
`)


type SHAwithRSA struct {
	privateKey *rsa.PrivateKey
}

func (this *SHAwithRSA) SetPriKey(pkey []byte) {
	block, _ := pem.Decode(pkey)
	fmt.Println("block.block block",block)
	if block == nil {
		fmt.Println("pem.Decode err")
		return
	}

	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	//private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return
	}
	//base64.StdEncoding.EncodeToString([]byte(private))

	this.privateKey = private
	//fmt.Println("private:",private)

	return
}

func (this *SHAwithRSA) Sign(data string) ([]byte, error) {

	h := crypto.Hash.New(crypto.SHA1)
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	fmt.Println("this.privateKey：", this.privateKey)
	signature, err := rsa.SignPKCS1v15(nil, this.privateKey, crypto.SHA1, hashed)
	fmt.Println("signature:", signature )
	//signature, err := rsa.SignPSS(nil, this.privateKey, crypto.SHA1, hashed,nil)
	if err != nil {
		fmt.Println("Error from signing: %s\n", err)
		return nil, err
	}
	sig :=base64.StdEncoding.EncodeToString(signature)
	fmt.Println("base64 sign :", sig )

	urls := url.QueryEscape(sig)
	fmt.Println("url encoding :", urls )

	//strings, _ := url.QueryUnescape(urls)
	//fmt.Println("signature byte  tt:", strings )
	//
	////t,_ := url.ParseQuery(urls)
	//
	//byte, _ :=base64.StdEncoding.DecodeString(strings)
	//
	//fmt.Println("signature byte  ttttt:", byte )
	return signature, nil
}

//dUtZK1FsMWVaTVp1T0toSnJhSGlKSVYzSWJVeHp5RS80S2E4UHpETzg1ZkZNTjl6


func main() {
	const  RsaEncodeLimit2048 = 2048
	sha := new(SHAwithRSA)
	sha.SetPriKey(h5PrivateKeys)
                              //amount=0.0001&assets=WC&expire_time=1000&merchant_id=13&merchant_trade_no=1&notify_url=http://example.com/notify_url&payee_id=100&product_detail=product detail&product_name=product name&remark=remark&return_url=http://example.com/return_url&show_url=http://example.com/show_url&timestamp=2019-04-16 23:59:59
	sign, err := sha.Sign("amount=0.0001&assets=WC&expire_time=1000&merchant_id=13&merchant_trade_no=1&notify_url=http://example.com/notify_url&payee_id=100&product_detail=product detail&product_name=product name&remark=remark&return_url=http://example.com/return_url&show_url=http://example.com/show_url&timestamp=2019-04-16 23:59:59")
	//sign, err := sha.Sign("amount=0.0001&assets=WC&expire_time=1000&merchant_id=13&merchant_trade_no=1&notify_url=http://example.com/notify_url&payee_id=100&product_detail=product detail&product_name=product name&remark=remark&return_url=http://example.com/return_url&show_url=http://example.com/show_url&timestamp=2019-04-16 23:59:59")
	fmt.Println("url sign:",sign,err)


	//VerifySign(h5PublicKey, []byte("amount=0.0001&assets=WC&expire_time=1000&merchant_id=1&merchant_trade_no=1&notify_url=http://example.com/notify_url&payee_id=1&product_detail=product detail&product_name=product name&remark=remark&return_url=http://example.com/return_url&show_url=http://example.com/show_url&sign_type=RSA&signature=string&timestamp=2019-04-15 23:59:59"), []byte(sign))



	dec := 0.0004323273872
	//stringAmount := strconv.FormatFloat(0.9999999999999, 'E', -1, 64)
	s :=fmt.Sprintf("%v",dec)

	fmt.Println("ss",s)

	//GenRsaKey(RsaEncodeLimit2048)
	//str := "0086"
	//	str = "+"+str[1:]

	//var  i int64 = 2333222

	//f := i / 1000

	//fmt.Println("b",f)
}

//base64 sign : run21G/iFK7cz781ufQjq1rtaf48/yz+nwi0yOFhrIU5+JyJlvDbVR+u6gWL94G7c7EpgeHwh/nj3Hf5UXW8/+jnyrx+16q2HftkBtYLYx7ySQEWGQESpOogWMIPUo/eS26E5OK3okmShs4o0XGmzWtYb24wadrxny6/PSfLaa/oHVbVdp8fPZ+DbiOKsflHBE+xyVfGoCj5lgkMUclOCR1lbYUhz4IF2qIZ6gJRmnbVFkSuX3si7X7EdmZm0RJbu2DRovGKn9FrbnbQU8dQvSEFKrnseWu9SdzNt2hu98G019MpmXvPKMXUDQiV90IrFJWGlbAAy3Pw0GJKy0RbXw==
//   url sign: run21G%2FiFK7cz781ufQjq1rtaf48%2Fyz%2Bnwi0yOFhrIU5%2BJyJlvDbVR%2Bu6gWL94G7c7EpgeHwh%2Fnj3Hf5UXW8%2F%2Bjnyrx%2B16q2HftkBtYLYx7ySQEWGQESpOogWMIPUo%2FeS26E5OK3okmShs4o0XGmzWtYb24wadrxny6%2FPSfLaa%2FoHVbVdp8fPZ%2BDbiOKsflHBE%2BxyVfGoCj5lgkMUclOCR1lbYUhz4IF2qIZ6gJRmnbVFkSuX3si7X7EdmZm0RJbu2DRovGKn9FrbnbQU8dQvSEFKrnseWu9SdzNt2hu98G019MpmXvPKMXUDQiV90IrFJWGlbAAy3Pw0GJKy0RbXw%3D%3D <nil>




func Base64Encode11(src []byte) []byte {
	var coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	return []byte(coder.EncodeToString(src))
}
func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func VerifySign (signingPubKey, signStr,sign []byte) {
	block, _ := pem.Decode(signingPubKey)
	if block == nil {
		fmt.Println("url sign11111:")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("22222err:",err)
	}

	rsaPubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		fmt.Println("rsaPubKey:",rsaPubKey)
	}

	c := crypto.Hash.New(crypto.SHA1)
	c.Write(signStr)
	digest := c.Sum(nil)

	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA1, digest, sign)
	if err != nil {
		fmt.Println("vvvv:",err)
	}
}
