package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

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
