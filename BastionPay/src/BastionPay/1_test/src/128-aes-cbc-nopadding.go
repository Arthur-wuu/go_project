package main
import(
	"crypto/aes"
	//"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	//"time"

	//"time"

	//"fmt"
	//"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"bytes"
	"os"
)

const (
	sKey        = "d9-9178-23224d46"
	ivParameter     = "4225424362556446"

)

//加密
//func PswEncrypt(src string)(string){
//	key := []byte(sKey)
//	iv := []byte(ivParameter)
//
//	result, err := Aes128Encrypt([]byte(src), key, iv)
//	if err != nil {
//		panic(err)
//	}
//	return result
//}
////解密
//func PswDecrypt(src string)(string) {
//
//	key := []byte(sKey)
//	iv := []byte(ivParameter)
//
//	var result []byte
//	var err error
//
//	result,err=base64.RawStdEncoding.DecodeString(src)
//	if err != nil {
//		panic(err)
//	}
//	origData, err := Aes128Decrypt(result, key, iv)
//	if err != nil {
//		panic(err)
//	}
//	return string(origData)
//
//}
func Aes128Encrypt(origData, key []byte) (string, error) {
	if key == nil || len(key) != 32 {
		return "", nil
	}
	key = key[0:16]
	IV := key[16:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = ZeroPadding(origData,blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, origData)

	return string(Base64Encode1(crypted)), nil
}
var coder1 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

func Base64Encode1(src []byte) []byte {
	return []byte(coder1.EncodeToString(src))
}
func Base64Decode1(src []byte) ([]byte, error) {
	return coder1.DecodeString(string(src))
}
//
//func Base58ss (src []byte) ([]byte, error) {
//	//fmt.Println(base.(src))
//	//s , _ := Base58(src)
//	//fmt.Println("string", string(s))
//	//return Base58(src)
//}


func Aes128Decrypt(crypted, key []byte) ([]byte, error) {
	crypted1,_ := Base64Decode1(crypted)

	if key == nil || len(key) != 32 {
		return nil, nil
	}
	key = key[0:16]
	IV := key[16:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block,IV[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted1)
	//origData = ZeroPadding(origData,blockSize[:])
	return origData, nil
}


func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
const(
	// rsa encode/decode bytes length limited, according to secret key bits
	RsaBits1024 = 1024
	RsaBits2048 = 2048
	RsaEncodeLimit1024 = RsaBits1024 / 8 - 11
	RsaDecodeLimit1024 = RsaBits1024 / 8
	RsaEncodeLimit2048 = RsaBits2048 / 8 - 11
	RsaDecodeLimit2048 = RsaBits2048 / 8
)

//func main(){
//	//fmt.Println(len([]byte("d9-9178-23224d464225424362556446")))
//	//encodingString,_ := Aes128Encrypt([]byte("123456"),[]byte("d9-9178-23224d464225424362556446"))
//	//fmt.Println("AES-128-C",encodingString)
//	//decodingString,_ := Aes128Decrypt([]byte(encodingString),[]byte("d9-9178-23224d464225424362556446"));
//	//fmt.Println("AES-128-C",string(decodingString))
//	//err := RsaGen(RsaEncodeLimit1024, "/Users/sywu/GoProject/test/", "/Users/sywu/GoProject/test/")
//	//err := GenRsaKeys(RsaBits2048)
//	//
//	//fmt.Println("err",err)
//
//}


//package main

//import (
//"fmt"
//"time"
//)

var POOL = 100

func groutine1(p chan int) {

	for i := 1; i <= POOL; i++ {
		p <- i
		if i%2 == 1 {
			fmt.Println("groutine-1:", i)
		}
	}
}

func groutine2(p chan int) {

	for i := 1; i <= POOL; i++ {
		<-p
		if i%2 == 0 {
			fmt.Println("groutine-2:", i)
		}
	}
}

func main() {
	//msg := make(chan int)
	//go groutine2(msg)
	//go groutine1(msg)
	//
	//
	//time.Sleep( time.Second * 12)
	//Base58([]byte("ss"))
	err := GenRsaKeys(RsaBits2048)
	fmt.Println(err)

}


func GenRsaKeys(bits int) error {
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
