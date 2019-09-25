package common

import (
	"testing"
)

func TestHash_MD5(t *testing.T) {
	h := NewHash("HelloWorld").MD5()
	t.Log(h)
}

func TestHash_SHA1(t *testing.T) {
	h := NewHash("HelloWorld").SHA1()
	t.Log(h)
}

func TestHash_SHA256(t *testing.T) {
	h := NewHash("HelloWorld").SHA256()
	t.Log(h)
}

func TestHash_SHA512(t *testing.T) {
	h := NewHash("HelloWorld").SHA512()
	t.Log(h)
}
