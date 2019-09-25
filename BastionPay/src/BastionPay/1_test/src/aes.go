package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const (
	key = "2018201820182018"
	iv  = "1234567887654321"
)

func main() {
	//
	content := []byte("url=/wallet/api/user/oauth/reset_login_password&timestamp=1551868437&nonce=dhsjkf&data=ekEmr+qshaAOxJixoZnpKk0t6lSkYdZJhjYBBqa3kyewrlu4qjSx2Dxpfzgtb69Rj3bENvi1AhHDiLIjPFqFQTBv4W4GEo+x6mc8UHCAQ4cRXuvcPT2p374i51zKpNoT&d9-9178-23224d464225424362556446")

	str := GetSHA256HashCode(content)

	fmt.Println("str:",str)

}

func AesEncrypt(encodeStr string, key []byte) (string, error) {
	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encodeBytes = PKCS5Padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

//func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	//填充
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//
//	return append(ciphertext, padtext...)
//}

func AesDecrypt(decodeStr string, key []byte) ([]byte, error) {
	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//func PKCS5UnPadding(origData []byte) []byte {
//	length := len(origData)
//	unpadding := int(origData[length-1])
//	return origData[:(length - unpadding)]
//}



func Padding(plainText []byte,blockSize int) []byte{
	//计算要填充的长度
	n:= blockSize-len(plainText)%blockSize
	//对原来的明文填充n个n
	temp:=bytes.Repeat([]byte{byte(n)},n)
	plainText=append(plainText,temp...)
	return plainText
}

var coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}
func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}


func AesEncrypt3(origData, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {

		return nil, err

	}

	blockSize := block.BlockSize()

	origData = PKCS5Padding(origData, blockSize)

	// origData = ZeroPadding(origData, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	crypted := make([]byte, len(origData))

	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以

	// crypted := origData

	blockMode.CryptBlocks(crypted, origData)
	cyBase64 := Base64Encode(crypted)
	return cyBase64, nil

	return crypted, nil

}

func AesDecrypt3(crypted, key []byte) ([]byte, error) {
	cryData ,_ := Base64Decode(crypted)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cryData))
	// origData := crypted
	blockMode.CryptBlocks(origData, cryData)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func GetSHA256HashCode(message []byte)string{
	//方法一：
	//创建一个基于SHA256算法的hash.Hash接口的对象

	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串

	code := hex.EncodeToString(bytes)

	//返回哈希值
	return code

	////方法二：
	//bytes2:=sha256.Sum256(message)//计算哈希值，返回一个长度为32的数组
	//hashcode2:=hex.EncodeToString(bytes2[:])//将数组转换成切片，转换成16进制，返回字符串
	//return hashcode2
}