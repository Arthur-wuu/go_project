package ga

import (
	"encoding/base64"
	"github.com/dgryski/dgoogauth"
	"math/rand"
	"net/url"
	"rsc.io/qr"
	"time"
)

type GA struct {
	URI         string `json:"uri"`
	Secret      string `json:"secret"`
	Image       string `json:"image"`
	Account     string `json:"account"`
	CompanyName string `json:"-"`
	Code        string `valid:"required" json:"code"`
}

func NewGA(companyName string) *GA {
	return &GA{
		CompanyName: companyName,
	}
}

func (g *GA) Verify(secret string, value string) (bool, error) {
	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	val, err := otpc.Authenticate(value)
	if err != nil {
		return false, err
	}

	return val, nil
}

func (g *GA) Generate(account string) error {
	var err error

	g.Account = account

	if err = g.generateSecret(); err != nil {
		return err
	}

	if err = g.generateURI(); err != nil {
		return err
	}

	if err = g.generateImage(); err != nil {
		return err
	}

	return nil
}

func (g *GA) generateSecret() error {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	g.Secret = string(result)

	return nil
}

func (g *GA) generateURI() error {
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		return err
	}

	URL.Path += "/" + url.PathEscape(g.CompanyName) + ":" + url.PathEscape(g.Account)

	params := url.Values{}
	params.Add("secret", g.Secret)
	params.Add("issuer", g.CompanyName)

	URL.RawQuery = params.Encode()

	g.URI = URL.String()

	return nil
}

func (g *GA) generateImage() error {
	code, err := qr.Encode(g.URI, qr.Q)
	if err != nil {
		return err
	}

	b := code.PNG()
	g.Image = "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)

	return nil
}
