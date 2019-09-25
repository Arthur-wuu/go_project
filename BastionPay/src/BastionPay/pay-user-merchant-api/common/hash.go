package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"time"
)

func GetRandomString(strLen int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < strLen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

type Hash struct {
	ClearText  string
	CipherText string
	Algorithm  string
	Salt       string
}

func NewHash(clearText string) *Hash {
	return &Hash{ClearText: clearText}
}

func (h *Hash) AddSalt(strLen int) *Hash {
	h.Salt = GetRandomString(strLen)
	return h
}

func (h *Hash) SetSalt(salt string) *Hash {
	h.Salt = salt
	return h
}

func (h *Hash) MD5() *Hash {
	ctx := md5.New()
	ctx.Write([]byte(h.ClearText + h.Salt))
	h.CipherText = hex.EncodeToString(ctx.Sum(nil))
	h.Algorithm = "MD5"
	return h
}

func (h *Hash) SHA1() *Hash {
	ctx := sha1.New()
	ctx.Write([]byte(h.ClearText + h.Salt))
	h.CipherText = hex.EncodeToString(ctx.Sum(nil))
	h.Algorithm = "SHA1"
	return h
}

func (h *Hash) SHA256() *Hash {
	ctx := sha256.New()
	ctx.Write([]byte(h.ClearText + h.Salt))
	h.CipherText = hex.EncodeToString(ctx.Sum(nil))
	h.Algorithm = "SHA256"
	return h
}

func (h *Hash) SHA512() *Hash {
	ctx := sha512.New()
	ctx.Write([]byte(h.ClearText + h.Salt))
	h.CipherText = hex.EncodeToString(ctx.Sum(nil))
	h.Algorithm = "SHA512"
	return h
}
