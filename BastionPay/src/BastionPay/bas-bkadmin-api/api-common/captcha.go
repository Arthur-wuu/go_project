package common

import (
	"github.com/mojocn/base64Captcha"
)

const (
	CaptchaTypeDigit     = "digit"
	CaptchaTypeCharacter = "character"
)

var configCharacter = base64Captcha.ConfigCharacter{
	Height: 48,
	Width:  168,
	// CaptchaModeNumber:数字,
	// CaptchaModeAlphabet:字母,
	// CaptchaModeArithmetic:算术,
	// CaptchaModeNumberAlphabet:数字字母混合.
	Mode:               base64Captcha.CaptchaModeNumberAlphabet,
	ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
	ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
	IsUseSimpleFont:    false,
	IsShowHollowLine:   true,
	IsShowNoiseDot:     false,
	IsShowNoiseText:    false,
	IsShowSlimeLine:    false,
	IsShowSineLine:     false,
	CaptchaLen:         6,
}

var configDigit = base64Captcha.ConfigDigit{
	Height:     48,
	Width:      168,
	MaxSkew:    0.7,
	DotCount:   80,
	CaptchaLen: 6,
}

type Captcha struct {
	Type    string
	Id      string
	Value   string
	Captcha string
}

func NewCaptcha(id string, t string) *Captcha {
	return &Captcha{Id: id, Type: t}
}

func (c *Captcha) Generate() *Captcha {
	var (
		config interface{}
	)
	switch c.Type {
	case CaptchaTypeDigit:
		config = configDigit
	case CaptchaTypeCharacter:
		config = configCharacter
	default:
		config = configDigit
	}

	_, digitCap := base64Captcha.GenerateCaptcha(c.Id, config)

	verifyValue := ""

	switch digitCap.(type) {
	case *base64Captcha.Audio:
		verifyValue = digitCap.(*base64Captcha.Audio).VerifyValue
	case *base64Captcha.CaptchaImageDigit:
		verifyValue = digitCap.(*base64Captcha.CaptchaImageDigit).VerifyValue
	case *base64Captcha.CaptchaImageChar:
		verifyValue = digitCap.(*base64Captcha.CaptchaImageChar).VerifyValue
	}

	base64Png := base64Captcha.CaptchaWriteToBase64Encoding(digitCap)

	c.Value = verifyValue
	c.Captcha = base64Png
	return c
}
