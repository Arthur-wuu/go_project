package common

import (
	"crypto/rand"
	"io"
)

func RandomDigit(length int) string {
	return parseDigitsToString(randomDigits(length))
}

// randomDigits returns a byte slice of the given length containing
// pseudorandom numbers in range 0-9. The slice can be used as a captcha
// solution.
func randomDigits(length int) []byte {
	return randomBytesMod(length, 10)
}

// randomBytes returns a byte slice of the given length read from CSPRNG.
func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

// randomBytesMod returns a byte slice of the given length, where each byte is
// a random number modulo mod.
func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r := randomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}
}

//parseDigitsToString parse randomDigits to normal string
func parseDigitsToString(bytes []byte) string {
	ssbb := make([]byte, len(bytes))
	for idx, by := range bytes {
		ssbb[idx] = by + '0'
	}
	return string(ssbb)
}
