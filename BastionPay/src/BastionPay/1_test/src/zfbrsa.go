package encoding

import (
	"crypto"
	"crypto/rand"
	//"crypto/rsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

func RSAEncryptPKCS1(plaintext, key []byte) ([]byte, error) {
	pub, err := ParsePKCS1PublicKey(key)
	if err != nil {
		return nil, err
	}

	var data = packageData(plaintext, pub.N.BitLen()/8-11)
	var cipherData = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, pub, d)
		if e != nil {
			return nil, e
		}
		cipherData = append(cipherData, c...)
	}

	return cipherData, nil
}

func RSAEncryptPKCS1WithKey(plaintext []byte, key *rsa.PublicKey) ([]byte, error) {
	var data = packageData(plaintext, key.N.BitLen()/8-11)
	var cipherData = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		cipherData = append(cipherData, c...)
	}

	return cipherData, nil
}

func RSADecryptPKCS1(ciphertext, key []byte) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	var data = packageData(ciphertext, pri.PublicKey.N.BitLen()/8)
	var plainData = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, pri, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func RSADecryptPKCS1WithKey(ciphertext []byte, key *rsa.PrivateKey) ([]byte, error) {
	var data = packageData(ciphertext, key.PublicKey.N.BitLen()/8)
	var plainData = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func SignPKCS1v15(src, key []byte, hash crypto.Hash) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return SignPKCS1v15WithKey(src, pri, hash)
}

func SignPKCS1v15WithKey(src []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func VerifyPKCS1v15(src, sig, key []byte, hash crypto.Hash) error {
	pub, err := ParsePKCS1PublicKey(key)
	if err != nil {
		return err
	}
	return VerifyPKCS1v15WithKey(src, sig, pub, hash)
}

func VerifyPKCS1v15WithKey(src, sig []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sig)
}


func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key error")
	}

	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func ParsePKCS1PublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key error")
	}

	return key, err
}