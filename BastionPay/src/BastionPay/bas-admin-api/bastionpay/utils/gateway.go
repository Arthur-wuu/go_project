package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

	return bytes.Join(s, []byte("")), nil
}

// Rsa decode cipher data with pem format private key
func RsaDecrypt(cipherData []byte, priKey []byte, limit int) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("decode private key error")
	}

	priInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// we need decode by sections according limited bytes
	length := len(cipherData)

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
		srcData := cipherData[offset:offsetto]

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
