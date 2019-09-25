package account

import (
	"github.com/BastionPay/bas-bkadmin-api/models"
)

type (
	Login struct {
		Name     string `valid:"required" json:"name"`
		Password string `valid:"required,length(6|50)" json:"password"`
	}

	UserLogin interface {
		GetUserInfoByEmail(email string) (*models.Account, error)
		GetUserInfoByMobile(mobile string) (*models.Account, error)
		GetUserInfoByName(name string) (*models.Account, error)
	}
)

func (this *Login) GetUserInfoByEmail(email string) (*models.Account, error) {
	account := &models.Account{}

	err := models.DB.Where("email = ? AND valid = ?", email, "0").First(account).Error
	return account, err
}

func (this *Login) GetUserInfoByMobile(mobile string) (*models.Account, error) {
	account := &models.Account{}

	err := models.DB.Where("mobile = ? AND valid = ?", mobile, "0").First(account).Error
	return account, err
}

func (this *Login) GetUserInfoByName(name string) (*models.Account, error) {
	account := &models.Account{}

	err := models.DB.Where("name = ? AND valid = ?", name, "0").First(account).Error
	return account, err
}
