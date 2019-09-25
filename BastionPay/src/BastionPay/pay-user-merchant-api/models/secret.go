package models

import (
	"BastionPay/pay-user-merchant-api/common"
	"regexp"
	"fmt"
	"errors"
	"BastionPay/pay-user-merchant-api/db"
)

var GSecretModel SecretModel

type SecretModel struct {

}

func (s *SecretModel) Verify(id uint, password string) (bool, error) {
	var (
		secret = Secret{}
		err    error
	)

	if err = db.GDbMgr.Get().First(&secret, id).Error; err != nil {
		return false, err
	}

	if common.NewHash(password).SetSalt(secret.Salt).SHA512().CipherText == secret.Secret {
		return true, nil
	}

	return false, nil
}

func (s *SecretModel) CreateSecret(password string) (uint, error) {
	if err:= s.vaildPasswd(password); err != nil {
		return 0, err
	}
	h := common.NewHash(password).AddSalt(128).SHA512()

	secret := Secret{
		Secret:    h.CipherText,
		Salt:      h.Salt,
		Algorithm: h.Algorithm,
	}

	if err := db.GDbMgr.Get().Save(&secret).Error; err != nil {
		return 0, err
	}

	return secret.ID, nil
}

func (s *SecretModel) vaildPasswd(password string) error {
	reg := regexp.MustCompile(`.{8,20}`)
	if !reg.MatchString(password) {
		return errors.New("length not in 8-20")
	}
	arr := []string{ `(?:[0-9]+)`, `(?:[a-zA-Z]+)`,  `(?:[^[:alnum:]]+)`}
	count := 0
	for i :=0; i < len(arr); i++ {
		reg := regexp.MustCompile(arr[i])
		flag := reg.MatchString(password)
		if flag {
			count++
		}
		fmt.Println(flag)
	}
	if count < 2 {
		return errors.New("at least more than 2 in dig num symbol ")
	}
	return nil
}