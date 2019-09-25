/*表结构*/
package models

import (
	"github.com/BastionPay/bas-admin-api/common"
)

// 用户表
type User struct {
	common.Model

	//钱包服务uuid绑定
	Uuid string `gorm:"type:varchar(100)"`
	// 公司名称
	CompanyName string `gorm:"type:varchar(255)"`
	// 手机
	Phone string `gorm:"type:varchar(50)"`
	// 国家代码 国际电信联盟的国际电话区号(E.164)
	CountryCode string `gorm:"type:varchar(20)"`
	// 邮箱
	Email string `gorm:"type:varchar(320)"`
	// 密码
	Secret   Secret
	SecretID uint `gorm:"not null;unique"`
	// 谷歌验证
	Ga string `gorm:"type:varchar(32)"`

	// 用户等级
	VipLevel uint8 `gorm:"default:0"`
	// 国籍
	Citizenship string `gorm:"type:varchar(255)"`
	// 语言
	Language string `gorm:"type:varchar(20)"`
	// 时区
	Timezone string `gorm:"type:varchar(10)"`
	// 法定币种
	LegalCurrency string `gorm:"type:varchar(50)"`
	// 推荐人
	Recommended uint
	// 注册类型 email, phone
	RegistrationType string `gorm:"type:varchar(50)"`
	// 是否禁用
	Blocked bool `gorm:"default:0"`

	// apikey
	Apikey   Apikey
	Optional Optional
}

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
type UserLevel struct {
	common.Model
	// 等级
	Level uint8
	// 等级名称
	Name string `gorm:"type:varchar(255)"`
}

// apikey
type Apikey struct {
	common.Model
	UserId    uint
	Apikey    string `gorm:"type:varchar(255)"`
	Signature string `gorm:"type:varchar(255)"`
	Blocked   bool   `gorm:"default:0"`
}

// 自选配置
type Optional struct {
	common.Model
	UserId   uint
	Optional string `gorm:"type:text"`
}

// 登录日志
type LogLogin struct {
	common.Model
	UserId  uint
	Ip      string `gorm:"type:varchar(255)"`
	Country string `gorm:"type:varchar(255)"`
	City    string `gorm:"type:varchar(255)"`
	Device  string `gorm:"type:varchar(10)"`
}

// 操作日志
type LogOperation struct {
	common.Model
	UserId    uint
	Operation string `gorm:"type:varchar(255)"`
	Ip        string `gorm:"type:varchar(255)"`
	Country   string `gorm:"type:varchar(255)"`
	City      string `gorm:"type:varchar(255)"`
}

// 公告表
type NoticeInfo struct {
	ID        *uint  `gorm:"primary_key"`
	CreatedAt *int64 `gorm:"type:bigint(20)"`
	UpdatedAt *int64 `gorm:"type:bigint(20)"`
	OnlinedAt *int64 `gorm:"type:bigint(20)"`
	// 下线时间
	OfflinedAt *int64 `gorm:"type:bigint(20)" `
	// 语言
	Language *string `gorm:"type:varchar(20)" `
	// 置顶标志
	Focus *bool `gorm:"type:tinyint(1)" `
	Race  *bool `gorm:"type:tinyint(1)" `
	// 标题
	Title *string `gorm:"type:varchar(100)" `
	//
	Author *string `gorm:"type:varchar(20)" `
	// 摘要
	Abstract *string `gorm:"type:varchar(100)"`
	// 内容
	Content *string `gorm:"type:text"`

	IsRead *bool `gorm:"-"`
}

// 已读公告表
type NoticeRead struct {
	//多于2个主键的情况下，不能设置id，不然Save->PrimaryKeyZero返回id，且id判断为空，则认为没有主键
	Noticeid   *uint   `gorm:"primary_key"`
	Userid     *uint   `gorm:"primary_key"`
	Userkey    *string `gorm:"type:VARCHAR(20)"`
	CreatedAt  *int64  `gorm:"type:BIGINT(20)"`
	OnlinedAt  *int64  `gorm:"type:bigint(20)"`
	OfflinedAt *int64  `gorm:"type:bigint(20)"`
}

type NoticeUser struct {
	Id          *uint  `gorm:"type:INT(11)"`
	Userid      *uint  `gorm:"primary_key"`
	MaxNoticeId *uint  `gorm:"type:INT(11) "`
	CreatedAt   *int64 `gorm:"type:BIGINT(20)"`
	UpdatedAt   *int64 `gorm:"type:bigint(20)"`
}
