package utils

import (
	"bytes"
	"crypto"
	"crypto/sha512"
	"io/ioutil"
	"testing"
)

var testDir = "/Users/henly.liu/gotestdir"
var priFile1024 = "/private_1024.pem"
var pubFile1024 = "/public_1024.pem"
var priFile2048 = "/private_2048.pem"
var pubFile2048 = "/public_2048.pem"

var priKey1024 []byte
var pubKey1024 []byte
var priKey2048 []byte
var pubKey2048 []byte

func TestRsaGen(t *testing.T) {
	var err error
	err = RsaGen(RsaBits1024, testDir+priFile1024, testDir+pubFile1024)
	if err != nil {
		t.FailNow()
	}
	priKey1024, err = ioutil.ReadFile(testDir + priFile1024)
	if err != nil {
		t.FailNow()
	}
	pubKey1024, err = ioutil.ReadFile(testDir + pubFile1024)
	if err != nil {
		t.FailNow()
	}

	err = RsaGen(RsaBits2048, testDir+priFile2048, testDir+pubFile2048)
	if err != nil {
		t.FailNow()
	}
	priKey2048, err = ioutil.ReadFile(testDir + priFile2048)
	if err != nil {
		t.FailNow()
	}
	pubKey2048, err = ioutil.ReadFile(testDir + pubFile2048)
	if err != nil {
		t.FailNow()
	}
}

func TestRsaEncrypt(t *testing.T) {
	var originData []byte
	for i := 0; i < 1024; i++ {
		originData = append(originData, byte(i))
	}

	cipherData1024, err := RsaEncrypt(originData, pubKey1024, RsaEncodeLimit1024)
	if err != nil {
		t.FailNow()
	}
	originData1024, err := RsaDecrypt(cipherData1024, priKey1024, RsaDecodeLimit1024)
	if err != nil {
		t.FailNow()
	}
	if bytes.Compare(originData1024, originData) != 0 {
		t.Error("rsa 1024 failed")
		t.FailNow()
	}

	cipherData2048, err := RsaEncrypt(originData, pubKey2048, RsaEncodeLimit2048)
	if err != nil {
		t.FailNow()
	}
	originData2048, err := RsaDecrypt(cipherData2048, priKey2048, RsaDecodeLimit2048)
	if err != nil {
		t.FailNow()
	}
	if bytes.Compare(originData2048, originData) != 0 {
		t.Error("rsa 1024 failed")
		t.FailNow()
	}
}

func TestRsaSign(t *testing.T) {
	t.Log("test rsa sign and verify...")
	var originData []byte
	for i := 0; i < 1024; i++ {
		originData = append(originData, byte(i))
	}

	var hashData []byte
	hs := sha512.New()
	hs.Write(originData)
	hashData = hs.Sum(nil)

	signData1024, err := RsaSign(crypto.SHA512, hashData, priKey1024)
	if err != nil {
		t.FailNow()
	}
	err = RsaVerify(crypto.SHA512, hashData, signData1024, pubKey1024)
	if err != nil {
		t.Error("rsa 1024 verify failed")
		t.FailNow()
	}
	err = RsaVerify(crypto.SHA512, hashData, signData1024, pubKey2048)
	if err == nil {
		t.Error("rsa 1024 verify failed")
		t.FailNow()
	}

	signData2048, err := RsaSign(crypto.SHA512, hashData, priKey2048)
	if err != nil {
		t.FailNow()
	}
	err = RsaVerify(crypto.SHA512, hashData, signData2048, pubKey2048)
	if err != nil {
		t.Error("rsa 2048 verify failed")
		t.FailNow()
	}
	err = RsaVerify(crypto.SHA512, hashData, signData2048, pubKey1024)
	if err == nil {
		t.Error("rsa 2048 verify failed")
		t.FailNow()
	}
}
