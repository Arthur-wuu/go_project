/*表结构*/
package models

import (
	"BastionPay/pay-user-merchant-api/common"
)

// 用户表
//type User struct {
//	common.Model
//
//	//钱包服务uuid绑定
//	Uuid string `gorm:"type:varchar(100)"`
//	// 公司名称
//	CompanyName string `gorm:"type:varchar(255)"`
//	// 手机
//	Phone string `gorm:"type:varchar(50)"`
//	// 国家代码 国际电信联盟的国际电话区号(E.164)
//	CountryCode string `gorm:"type:varchar(20)"`
//	// 邮箱
//	Email string `gorm:"type:varchar(320)"`
//	// 密码
//	Secret   Secret
//	SecretID uint `gorm:"not null;unique"`
//	// 谷歌验证
//	Ga string `gorm:"type:varchar(32)"`
//
//	// 用户等级
//	VipLevel uint8 `gorm:"default:0"`
//	// 国籍
//	Citizenship string `gorm:"type:varchar(255)"`
//	// 语言
//	Language string `gorm:"type:varchar(20)"`
//	// 时区
//	Timezone string `gorm:"type:varchar(10)"`
//	// 法定币种
//	LegalCurrency string `gorm:"type:varchar(50)"`
//	// 推荐人
//	Recommended uint
//	// 注册类型 email, phone
//	RegistrationType string `gorm:"type:varchar(50)"`
//	// 是否禁用
//	Blocked bool `gorm:"default:0"`
//
//	// apikey
//	Apikey   Apikey
//	Optional Optional
//}

// 密码表
type Secret struct {
	common.Model
	// 密文
	Secret string `gorm:"type:varchar(255)"`
	// 盐
	Salt string `gorm:"type:varchar(255)"`
	// 加密方式
	Algorithm string `gorm:"type:varchar(20)"`
}

// 用户等级表
//type UserLevel struct {
//	common.Model
//	// 等级
//	Level uint8
//	// 等级名称
//	Name string `gorm:"type:varchar(255)"`
//}

// apikey
type Apikey struct {
	common.Model
	UserId    uint
	Apikey    string `gorm:"type:varchar(255)"`
	Signature string `gorm:"type:varchar(255)"`
	Blocked   bool   `gorm:"default:0"`
}

//自选配置
type Optional struct {
	common.Model
	UserId   uint
	Optional string `gorm:"type:text"`
}

// 登录日志
//type LogLogin struct {
//	common.Model
//	UserId  uint
//	Ip      string `gorm:"type:varchar(255)"`
//	Country string `gorm:"type:varchar(255)"`
//	City    string `gorm:"type:varchar(255)"`
//	Device  string `gorm:"type:varchar(10)"`
//}

// 操作日志
//type LogOperation struct {
//	common.Model
//	UserId    uint
//	Operation string `gorm:"type:varchar(255)"`
//	Ip        string `gorm:"type:varchar(255)"`
//	Country   string `gorm:"type:varchar(255)"`
//	City      string `gorm:"type:varchar(255)"`
//}

// 公告表

// 已读公告表

type (
	UserLevel struct {
		Id        *int64  `json:"id,omitempty"                gorm:"AUTO_INCREMENT:1;column:id"`
		UserKey   *string `json:"user_key,omitempty"           gorm:"column:user_key;type:varchar(100);primary_key"`
		ExpireAt  *int64  `json:"expire_at,omitempty"         gorm:"column:expire_at;type:bigint(20)"`
		Level     *int    `json:"level,omitempty"           gorm:"column:level;type:int(11)"`
		Valid     *int    `json:"valid,omitempty" gorm:"column:valid"`
		CreatedAt *int64  `json:"created_at,omitempty" gorm:"column:created_at"`
		UpdatedAt *int64  `json:"updated_at,omitempty" gorm:"column:updated_at"`
		Author    *string `json:"author,omitempty" gorm:"column:author"`
	}
)

func (this *UserLevel) TableName() string {
	return "user_level"
}
