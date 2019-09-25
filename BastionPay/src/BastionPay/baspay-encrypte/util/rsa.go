package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

const (
	// rsa encode/decode bytes length limited, according to secret key bits
	RsaBits1024        = 1024
	RsaBits2048        = 2048
	RsaEncodeLimit1024 = RsaBits1024/8 - 11
	RsaDecodeLimit1024 = RsaBits1024 / 8
	RsaEncodeLimit2048 = RsaBits2048/8 - 11
	RsaDecodeLimit2048 = RsaBits2048 / 8
)

// Generate Rsa secret key with pem format
func RsaGen(bits int, priPath string, pubPath string) error {
	// generate pri
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	file, err := os.Create(priPath)
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	// build public
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create(pubPath)
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

//网上找的 加密方式
func RsaEncrypt2(originData []byte, pubKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, originData)
}

func RsaDecrypt2(cipherData []byte, priKey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherData)
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

// Rsa decode cipher data with pem format private key
func RsaDecrypt(cipherData []byte, priKey []byte, limit int) ([]byte, error) {
	cDate, err := Base64Decode(cipherData)
	if err != nil {
		fmt.Println("decode err ", err)
		return nil, err
	}
	fmt.Println("****Base64Decode  ", cDate)
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("decode private key error")
	}

	priInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// we need decode by sections according limited bytes
	length := len(cDate)

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
		srcData := cDate[offset:offsetto]

		dstData, err := rsa.DecryptPKCS1v15(rand.Reader, priInterface, srcData)
		if err != nil {
			return nil, err
		}
		s[index] = dstData

		index++
		offset = offsetto
	}

	return bytes.Join(s, []byte("")), nil
}

// Rsa signature hash data with pem format private key
func RsaSign(hash crypto.Hash, hashData []byte, priKey []byte) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("decode private key error")
	}

	priInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.SignPKCS1v15(rand.Reader, priInterface, hash, hashData)
}

// Rsa verify signature data with pem public key and hash data
func RsaVerify(hash crypto.Hash, hashData []byte, signData []byte, pubKey []byte) error {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return errors.New("decode public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	pub := pubInterface.(*rsa.PublicKey)

	return rsa.VerifyPKCS1v15(pub, hash, hashData, signData)
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

func MakeSSHKeyPair() (string, string, error) {

	pkey, pubkey, err := GenerateKey(2048)
	if err != nil {
		return "", "", err
	}

	pub, err := EncodeSSHKey(pubkey)
	if err != nil {
		return "", "", err
	}
	fmt.Println("", string(EncodePrivateKey(pkey)))
	fmt.Println("", string(pub))

	//s :=fmt.Sprintf("privateKey=[%s]\n pubKey=[%s]",string(EncodePrivateKey(pkey)),string(pub))
	//fmt.Println(s)
	return string(EncodePrivateKey(pkey)), string(pub), nil

}

func GenerateKey(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return private, &private.PublicKey, nil

}

func EncodePrivateKey(private *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(private),
		Type:  "RSA PRIVATE KEY",
	})
}

func EncodeSSHKey(public *rsa.PublicKey) ([]byte, error) {
	publicKey, err := ssh.NewPublicKey(public)
	if err != nil {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(publicKey), nil
}

var coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

//自己改一下limit的，测试下
func RsaEncrypt3(originData []byte, pubKey []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("decode public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	dstData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, originData)
	if err != nil {
		return nil, err
	}
	//base64 加密一下
	base64EnCode := Base64Encode(dstData)

	return base64EnCode, nil

}

//自己改一下limit的，测试下
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
