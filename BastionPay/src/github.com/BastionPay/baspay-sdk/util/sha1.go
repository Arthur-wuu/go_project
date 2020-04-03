package utils

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type SHAwithRSA struct {
	privateKey *rsa.PrivateKey
}

func (this *SHAwithRSA) SetPriKey(pkey []byte) {
	block, _ := pem.Decode(pkey)
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

	this.privateKey = private

	return
}

func (this *SHAwithRSA) Sign(data string) (string, error) {

	h := crypto.Hash.New(crypto.SHA1)
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(nil, this.privateKey, crypto.SHA1, hashed)

	//signature, err := rsa.SignPSS(nil, this.privateKey, crypto.SHA1, hashed,nil)
	if err != nil {
		fmt.Println("Error from signing: %s\n", err)
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(signature)

	sReplace1 := strings.Replace(s, "\n", "", -1)
	sReplace2 := strings.Replace(sReplace1, "[", "", -1)
	sReplace3 := strings.Replace(sReplace2, "]", "", -1)
	sReplace4 := strings.Replace(sReplace3, "\r", "", -1)

	ss := url.QueryEscape(sReplace4)

	return ss, nil
}

func VerifySign(signingPubKey, data, sign []byte) {
	block, _ := pem.Decode(signingPubKey)
	if block == nil {

	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {

	}

	rsaPubKey, ok := key.(*rsa.PublicKey)
	if !ok {

	}

	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA1, data, sign)
	if err != nil {

	}
}

//---------------AES加密  解密--------------------

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
	origData = ZeroPadding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, origData)

	return string(Base64Encode(crypted)), nil
}

func Aes128Decrypt(crypted, key []byte) ([]byte, error) {
	crypted1, _ := Base64Decode(crypted)

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
	blockMode := cipher.NewCBCDecrypter(block, IV[:blockSize])
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

var coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

func RsaDecrypt3(cipherData []byte, priKey []byte) ([]byte, error) {
	clipherBase64, err := Base64Decode(cipherData)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("decode private key error")
	}

	priInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	dstData, err := rsa.DecryptPKCS1v15(rand.Reader, priInterface, clipherBase64)

	if err != nil {
		return nil, err
	}
	return dstData, nil
}

// Rsa encode origin data with pem format public key
func RsaEncrypt(originData []byte, pubKey []byte, limit int) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("decode public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	// we need encode by sections according limited bytes
	length := len(originData)

	cnt := length / limit
	if length%limit != 0 {
		cnt += 1
	}
	s := make([][]byte, cnt)

	index := 0
	offset := 0
	for offset < length {
		offsetto := 0
		if length-offset > limit {
			offsetto = offset + limit
		} else {
			offsetto = length
		}
		srcData := originData[offset:offsetto]

		dstData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, srcData)
		if err != nil {
			return nil, err
		}
		s[index] = dstData

		index++
		offset = offsetto
	}
	baseCode := Base64Encode(bytes.Join(s, []byte("")))

	return baseCode, nil
	//return bytes.Join(s, []byte("")), nil
}
