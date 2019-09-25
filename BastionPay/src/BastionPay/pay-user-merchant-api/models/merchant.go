package models

import (
	"BastionPay/pay-user-merchant-api/db"
	"github.com/jinzhu/gorm"
	"BastionPay/pay-user-merchant-api/api"
)

type Merchant struct{
	ID               *int64        `json:"id,omitempty"            gorm:"AUTO_INCREMENT:1;column:ID;primary_key;not null"` //加上type:int(11)后AUTO_INCREMENT无效
	MerchantId       *string        `json:"merchant_id,omitempty"        gorm:"column:MERCHANT_ID;type:varchar(255)"`
	MerchantName     *string        `json:"merchant_name,omitempty"        gorm:"column:MERCHANT_NAME;type:varchar(255)"`
	NotifyUrl        *string       `json:"notify_url,omitempty"        gorm:"column:NOTIFY_URL;type:varchar(255)"`
	SignType         *string       `json:"sign_type,omitempty"        gorm:"column:SIGN_TYPE;type:varchar(255)"`
	SignKey          *string       `json:"sign_key,omitempty"        gorm:"column:SIGN_KEY;type:varchar(2000)"`
	PayeeId          *int64       `json:"payee_id,omitempty"        gorm:"column:PAYEE_ID;type:int(11)"`
	LanguageType     *string        `json:"language_type,omitempty"        gorm:"column:LANGUAGE_TYPE;type:varchar(255)"`
	LegalCurrency    *string        `json:"legal_currency,omitempty"        gorm:"column:LEGAL_CURRENCY;type:varchar(255)"`
	Contact          *string       `json:"contact,omitempty"        gorm:"column:CONTACT;type:varchar(255)"`
	ContactPhone     *string        `json:"contact_phone,omitempty"        gorm:"column:CONTACT_PHONE;type:varchar(255)"`
	ContactEmail     *string         `json:"contact_email,omitempty"        gorm:"column:CONTACT_EMAIL;type:varchar(255)"`
	Country          *string        `json:"country,omitempty"        gorm:"column:COUNTRY;type:varchar(255)"`
	CreateTime       *int64          `json:"create_time,omitempty"        gorm:"column:CREATE_TIME;type:bigint(20)"`
	LastUpdateTime   *int64         `json:"last_update_time,omitempty"        gorm:"column:LAST_UPDATE_TIME;type:bigint(20)"`
}

func (this *Merchant) ParseAdd(mer *api.MerchantAdd)  *Merchant {
	//this.MerchantId = mer.MerchantId
	this.MerchantName = mer.MerchantName
	this.NotifyUrl = mer.NotifyUrl
	this.SignType = mer.SignType
	this.SignKey   =  mer.SignKey
	this.LanguageType = mer.LanguageType
	this.LegalCurrency = mer.LegalCurrency
	this.Contact =      mer.Contact
	this.ContactPhone =  mer.ContactPhone
	this.ContactEmail =  mer.ContactEmail
	this.Country  =      mer.Country
	this.CreateTime     = mer.CreateTime
	this.LastUpdateTime = mer.LastUpdateTime
	return this
}

func (this *Merchant) Parse(mer *api.Merchant) *Merchant {
	this.MerchantId = mer.MerchantId
	this.MerchantName = mer.MerchantName
	this.NotifyUrl = mer.NotifyUrl
	this.SignType = mer.SignType
	this.SignKey   =  mer.SignKey
	this.LanguageType = mer.LanguageType
	this.LegalCurrency = mer.LegalCurrency
	this.Contact =      mer.Contact
	this.ContactPhone =  mer.ContactPhone
	this.ContactEmail =  mer.ContactEmail
	this.Country  =      mer.Country
	this.CreateTime     = mer.CreateTime
	this.LastUpdateTime = mer.LastUpdateTime
	return this
}

func (this *Merchant) Unique() (bool,error) {
	var count int
	if err := db.GDbMgr.Get().Model(this).Where("MERCHANT_ID = ？or PAYEE_ID = ?", this.MerchantId, this.PayeeId).Count(&count).Error; err != nil {
		return false,err
	}
	return count == 0, nil
}

func (this *Merchant) Add() (*Merchant, error) {
	if err := db.GDbMgr.Get().Create(this).Error; err != nil {
		return nil,err
	}
	if this.ID == nil || *this.ID == 0 {
		return nil,nil
	}
	if err := db.GDbMgr.Get().Model(this).Where("ID = ?", *this.ID).Last(this).Error; err != nil {
		return nil,err
	}
	return this,nil
}

func (this *Merchant)GetByPayeeId(payeeId int64) (*Merchant, error) {
	 err := db.GDbMgr.Get().Model(this).Where("PAYEE_ID = ?", payeeId).Last(this).Error
	 if err == gorm.ErrRecordNotFound {
	 	return nil,nil
	 }
	 return this,err
}

func (this *Merchant)Update() (*Merchant, error) {
	err := db.GDbMgr.Get().Model(this).Where("PAYEE_ID = ?", this.PayeeId).Updates(this).Error
	if err != nil {
		return nil,err
	}
	err = db.GDbMgr.Get().Model(this).Where("ID = ?", *this.ID).Last(this).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return this,err
}

//func (this *Merchant)List() ([]*Merchant, error) {
//	err := db.GDbMgr.Get().Model(this).Updates(this).Error
//	if err != nil {
//		return nil,err
//	}
//	err = db.GDbMgr.Get().Model(this).Where("ID = ?", *this.ID).Last(this).Error
//	if err == gorm.ErrRecordNotFound {
//		return nil, nil
//	}
//	return this,err
//}