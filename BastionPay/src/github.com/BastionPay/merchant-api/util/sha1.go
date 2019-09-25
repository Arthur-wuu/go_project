package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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
	s :=base64.StdEncoding.EncodeToString(signature)

	sReplace1 :=strings.Replace(s, "\n","",-1)
	sReplace2 :=strings.Replace(sReplace1, "[","",-1)
	sReplace3 :=strings.Replace(sReplace2, "]","",-1)
	sReplace4 :=strings.Replace(sReplace3, "\r","",-1)

	ss := url.QueryEscape(sReplace4)

	return ss, nil
}


func VerifySign (signingPubKey, data,sign []byte) {
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