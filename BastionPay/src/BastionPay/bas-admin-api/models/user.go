package models

import (
	"errors"
	backend "github.com/BastionPay/bas-api/apibackend/v1/backend"
	"github.com/jinzhu/gorm"
)

type UserModel struct {
	conn *gorm.DB
}

func NewUserModel(conn *gorm.DB) *UserModel {
	// 增加company_name字段
	conn.AutoMigrate(&User{})
	return &UserModel{conn: conn}
}

func (u *UserModel) GetUserById(userId uint) (*User, error) {
	user := &User{}
	err := u.conn.Find(user, userId).Error
	if err != nil {
		return user, err
	}

	if user == nil {
		return nil, errors.New("USER_NOT_FIND")
	}

	return user, nil
}

func (u *UserModel) GetUserByUserkey(userKey string) (*User, error) {
	user := &User{}
	err := u.conn.Debug().Where("uuid = ? ", userKey).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("USER_NOT_FIND")
	}
	return user, nil
}

func (u *UserModel) GetUserByName(username string) (*User, error) {
	user := &User{}
	err := u.conn.Debug().Where("email = ? OR phone = ?", username, username).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("USER_NOT_FIND")
	}
	return user, nil
}

func (u *UserModel) GetEmail(userId uint) (string, error) {
	user, err := u.GetUserById(userId)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (u *UserModel) GetPhone(userId uint) (string, error) {
	user, err := u.GetUserById(userId)
	if err != nil {
		return "", err
	}
	return user.Phone, nil
}

func (u *UserModel) GetPhoneRecipient(userId uint) (string, error) {
	user, err := u.GetUserById(userId)
	if err != nil {
		return "", err
	}
	return user.CountryCode + user.Phone, nil
}

func (u *UserModel) GetGa(userId uint) (string, error) {
	user, err := u.GetUserById(userId)
	if err != nil {
		return "", err
	}
	return user.Ga, nil
}

func (u *UserModel) UserExisted(username string) (bool, error) {
	count := 0
	err := u.conn.Model(&User{}).Where("email = ? OR phone = ?", username, username).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *UserModel) CompanyExisted(name string) (bool, error) {
	count := 0
	err := u.conn.Model(&User{}).Where("company_name = ? ", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *UserModel) CreateUser(companyName string, email string, phone string, countryCode string, password uint, citizenship string,
	language string, timezone string, registrationType, uuid string) (uint, error) {
	data := User{
		CompanyName:      companyName,
		Email:            email,
		Phone:            phone,
		CountryCode:      countryCode,
		SecretID:         password,
		VipLevel:         0,
		Citizenship:      citizenship,
		Language:         language,
		Timezone:         timezone,
		RegistrationType: registrationType,
		Uuid:             uuid,
	}

	err := u.conn.Save(&data).Error
	if err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (u *UserModel) UpdateSecretId(userId uint, secretId uint) error {
	if err := u.conn.Model(&User{}).Where("id = ?", userId).Update("secret_id", secretId).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserModel) BindGa(userId uint, secret string) error {
	if err := u.conn.Model(&User{}).Where("id = ?", userId).Update("ga", secret).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserModel) UnBindGa(userId uint) error {
	if err := u.conn.Model(&User{}).Where("id = ?", userId).Update("ga", nil).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserModel) SetLanguage(userId uint, language string) error {
	if err := u.conn.Model(&User{}).Where("id = ?", userId).Update("language", language).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserModel) SetTimezone(userId uint, timezone string) error {
	if err := u.conn.Model(&User{}).Where("id = ?", userId).Update("timezone", timezone).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserModel) SetEmail(userId uint, email string) error {
	if err := u.conn.Model(&User{}).Where("id = ?", userId).Update("email", email).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserModel) SetPhone(userId uint, phone string, countryCode string) error {
	if err := u.conn.Model(&User{}).Where("id = ?", userId).Updates(&User{Phone: phone, CountryCode: countryCode}).
		Error; err != nil {
		return err
	}

	return nil
}

func (u *UserModel) ListUserCountByBasic(conds map[string]interface{}) (int, error) {
	var count int
	err := u.conn.Model(&User{}).Where(conds).Count(&count).Error
	return count, err
}

func (u *UserModel) ListUsersByBasic(offset, pagesize int, conds map[string]interface{}) ([]backend.UserBasic, error) {
	userBasics := make([]backend.UserBasic, 0)
	users := make([]User, 0)
	err := u.conn.Model(&User{}).Where(conds).Select("id, uuid, phone, email, vip_level, registration_type, created_at, updated_at").Order("id desc").Offset(offset).Limit(pagesize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(users); i++ {
		basic := new(backend.UserBasic)
		basic.Id = int(users[i].ID)
		if users[i].RegistrationType == "email" {
			basic.UserName = users[i].Email
		} else {
			basic.UserName = users[i].Phone
		}
		basic.UserKey = users[i].Uuid
		basic.Level = int(users[i].VipLevel)
		basic.CreateTime = users[i].CreatedAt
		basic.UpdateTime = users[i].UpdatedAt
		basic.UserEmail = users[i].Email
		basic.UserMobile = users[i].Phone
		userBasics = append(userBasics, *basic)
	}
	return userBasics, nil
}
