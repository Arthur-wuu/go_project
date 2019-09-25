package models

import (
	backend "BastionPay/bas-api/apibackend/v1/backend"
	"github.com/jinzhu/gorm"
	"BastionPay/pay-user-merchant-api/db"
)

type User struct {
	Id    *int64       `json:"ID,omitempty"  gorm:"AUTO_INCREMENT:1;column:ID;primary_key;not null"` //加上type:int(11)后AUTO_INCREMENT无效
	//钱包服务uuid绑定
	NickName string `gorm:"column:NICK_NAME;type:varchar(32)"`
	// 公司名称
	Company string `gorm:"column:COUNTRY;type:varchar(20)"`
	// 手机
	Phone string `gorm:"column:PHONE;type:varchar(50)"`
	Country string `gorm:"column:COUNTRY;type:varchar(20)"`
	// 国家代码 国际电信联盟的国际电话区号(E.164)
	PhonneDistrict string `gorm:"column:PHONE_DISTRICT;type:varchar(20)"`
	// 谷歌验证
	LoginPassword string `gorm:"column:LOGIN_PASSWORD;type:varchar(32)"`
	// 用户等级
	Status uint8 `gorm:"column:STATUS;type:tinyint(4)"`
	// 语言
	Language string `gorm:"column:LANGUA;type:varchar(20)"`
	// 时区
	Channnel  string  `gorm:"column:CHANNEL;type:varchar(30)"`
	Ga        string  `gorm:"column:CHANNEL;type:varchar(30)"`
	Email string `gorm:"-"`
	VipLevel  uint8 `gorm:"-"`
	CreateTime       *int64          `json:"create_time,omitempty"        gorm:"column:CREATE_TIME;type:varchar(255)"`
	LastUpdateTime   *int64         `json:"last_update_time,omitempty"        gorm:"column:LAST_UPDATE_TIME;type:varchar(255)"`
}


func (this *User) GetById(id int64) (*User, error) {
	user := &User{Id: &id}
	err := db.GDbMgr.Get().Where(user).Last(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return user, err
	}

	return user, nil
}

//func (this *User) GetByUserkey(userKey string) (*User, error) {
//	user := &User{Uuid: userKey}
//	err := db.GDbMgr.Get().Where(user).Last(&user).Error
//	if err == gorm.ErrRecordNotFound {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, err
//	}
//	//if user == nil {
//	//	return nil, errors.New("USER_NOT_FIND")
//	//}
//	return user, nil
//}

func (this *User) GetByName(username string) (*User, error) {
	user := &User{}
	err := db.GDbMgr.Get().Where("email = ? OR phone = ?", username, username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil,nil
	}
	if err != nil {
		return nil, err
	}
	//if user == nil {
	//	return nil, errors.New("USER_NOT_FIND")
	//}
	return user, nil
}

func (this *User) GetByPhone(countryCode, phone string) (*User, error) {
	user := &User{
		PhonneDistrict:countryCode,
		Phone:phone,
	}
	err := db.GDbMgr.Get().Where(user).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil,nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetEmail(userId int64) (string, error) {
	user, err := u.GetById(userId)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (u *User) GetPhone(userId int64) (string, error) {
	user, err := u.GetById(userId)
	if err != nil {
		return "", err
	}
	return user.Phone, nil
}

func (u *User) GetPhoneRecipient(userId int64) (string, error) {
	user, err := u.GetById(userId)
	if err != nil {
		return "", err
	}
	return user.PhonneDistrict + user.Phone, nil
}

func (u *User) GetGa(userId int64) (string, error) {
	user, err := u.GetById(userId)
	if err != nil {
		return "", err
	}
	return user.Ga, nil
}

func (u *User) UserExisted(username string) (bool, error) {
	count := 0
	err := db.GDbMgr.Get().Model(&User{}).Where("email = ? OR phone = ?", username, username).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

//func (u *User) CompanyExisted(name string) (bool, error) {
//	count := 0
//	err := db.GDbMgr.Get().Model(&User{}).Where("company_name = ? ", name).Count(&count).Error
//	if err != nil {
//		return false, err
//	}
//	if count > 0 {
//		return true, nil
//	} else {
//		return false, nil
//	}
//}


func (u *User) SetLanguage(userId uint, language string) error {
	if err := db.GDbMgr.Get().Model(&User{}).Where("ID = ?", userId).Update("language", language).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) SetTimezone(userId uint, timezone string) error {
	if err := db.GDbMgr.Get().Model(&User{}).Where("ID = ?", userId).Update("timezone", timezone).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) SetEmail(userId uint, email string) error {
	if err := db.GDbMgr.Get().Model(&User{}).Where("ID = ?", userId).Update("email", email).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) SetPhone(userId uint, phone string, countryCode string) error {
	if err := db.GDbMgr.Get().Model(&User{}).Where("ID = ?", userId).Updates(&User{Phone: phone, PhonneDistrict: countryCode}).
		Error; err != nil {
		return err
	}

	return nil
}

func (u *User) ListUserCountByBasic(conds map[string]interface{}) (int, error) {
	var count int
	err := db.GDbMgr.Get().Model(&User{}).Where(conds).Count(&count).Error
	return count, err
}

func (u *User) ListUsersByBasic(offset, pagesize int, conds map[string]interface{}) ([]backend.UserBasic, error) {
	userBasics := make([]backend.UserBasic, 0)
	//users := make([]User, 0)
	//err := db.GDbMgr.Get().Model(&User{}).Where(conds).Select("id, uuid, phone, email, vip_level, registration_type, created_at, updated_at").Order("id desc").Offset(offset).Limit(pagesize).Find(&users).Error
	//if err != nil {
	//	return nil, err
	//}
	//for i := 0; i < len(users); i++ {
	//	basic := new(backend.UserBasic)
	//	basic.Id = int(*(users[i].Id))
	//	//if users[i].RegistrationType == "email" {
	//	//	basic.UserName = users[i].Email
	//	//} else {
	//		basic.UserName = users[i].Phone
	//	//}
	//	basic.UserKey = users[i].Uuid
	//	basic.Level = int(users[i].VipLevel)
	//	basic.CreateTime = *users[i].CreatedAt
	//	basic.UpdateTime = *users[i].UpdatedAt
	//	basic.UserEmail = users[i].Email
	//	basic.UserMobile = users[i].Phone
	//	userBasics = append(userBasics, *basic)
	//}
	return userBasics, nil
}

func (u *User) SetLevel(user_key string, vipLevel, valid int, expireAt int64) error {
	tx := db.GDbMgr.Get().Begin()
	if err := tx.Model(&User{}).Where("uuid = ?", user_key).Update("vip_level", vipLevel).Error; err != nil {
		tx.Rollback()
		return err
	}
	ulv := &UserLevel{
		UserKey: &user_key,
		Level:  &vipLevel,
		Valid:  &valid,
		ExpireAt: &expireAt,
	}
	if err := tx.Model(&UserLevel{}).Save(ulv).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) BindGa(userId uint, secret string) error {
	if err := db.GDbMgr.Get().Model(&User{}).Where("ID = ?", userId).Update("GA", secret).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) UnBindGa(userId uint) error {
	if err := db.GDbMgr.Get().Model(&User{}).Where("ID = ?", userId).Update("GA", nil).Error; err != nil {
		return err
	}

	return nil
}