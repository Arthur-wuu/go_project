package models

import (
	//"BastionPay/merchant-api/common"
	"BastionPay/merchant-api/common"
	"BastionPay/merchant-api/db"
	"errors"
)


type (
	MerchantConfig struct {
		Id                  *uint64    `json:"id"                gorm:"AUTO_INCREMENT;primary_key;column:ID;type:bigint(20)"`
		MerchantId          *string    `json:"merchant_id"       gorm:"column:MERCHANT_ID;type:varchar(255)"`
		MerchantName        *string    `json:"merchant_name"     gorm:"column:MERCHANT_NAME;type:varchar(255)"`

		NotifyUrl           *string    `json:"notify_url"        gorm:"column:NOTIFY_URL;type:varchar(255)"`
		SignType            *string    `json:"sign_type"         gorm:"column:SIGN_TYPE;type:varchar(255)"`
		SignKey             *string    `json:"sign_key"          gorm:"column:SIGN_KEY;type:varchar(2000)"`
		PayeeId             *int64     `json:"payee_id"          gorm:"column:PAYEE_ID;type:bigint(20)"`
		LanguageType        *string    `json:"language_type"     gorm:"column:LANGUAGE_TYPE;type:varchar(255)"`
		LegalCurrency       *string    `json:"legal_currency"    gorm:"column:LEGAL_CURRENCY;type:varchar(255)"`
		Contact             *string    `json:"contact"           gorm:"column:CONTACT;type:varchar(255)"`
		ContactPhone        *string    `json:"contact_phone"     gorm:"column:CONTACT_PHONE;type:varchar(255)"`
		ContactEmail        *string    `json:"contact_email"     gorm:"column:CONTACT_EMAIL;type:varchar(255)"`
		Country   		    *string    `json:"country"           gorm:"column:COUNTRY;type:varchar(255)"`
		CreatedAt           *int64     `json:"created_at"        gorm:"column:CREATED_DATE;type:bigint(20)"`
		UpdatedAt           *int64     `json:"updated_at"        gorm:"column:LAST_UPDATE_TIME;type:bigint(20)"`
	}
)

func (this *MerchantConfig) TableName() string {
	return "MERCHANT"
}

type (
	MerchantAdd  struct {
	//	Id                  *int64  `valid:"optional" json:"id"`
		MerchantId          *string  `valid:"required" json:"merchant_id"`
		MerchantName        *string  `valid:"optional" json:"merchant_name"`
		NotifyUrl           *string  `valid:"required" json:"notify_url"`
		SignType            *string  `valid:"required" json:"sign_type"`
		SignKey             *string  `valid:"required" json:"sign_key"`
		PayeeId             *int64   `valid:"required" json:"payee_id"`
		LanguageType        *string  `valid:"optional" json:"language_type"`
		LegalCurrency       *string  `valid:"optional" json:"legal_currency"`
		Contact             *string  `valid:"optional" json:"contact"`
		ContactPhone        *string  `valid:"optional" json:"contact_phone"`
		ContactEmail        *string  `valid:"optional" json:"contact_email"`
		Country             *string  `valid:"optional" json:"country"`
	}

	MerchantUpdate struct {
		Id                  *uint64  `valid:"required" json:"id"`
		MerchantId          *string  `valid:"required" json:"merchant_id"`
		MerchantName        *string  `valid:"optional" json:"merchant_name"`
		NotifyUrl           *string  `valid:"optional" json:"notify_url"`
		SignType            *string  `valid:"optional" json:"sign_type"`
		SignKey             *string  `valid:"optional" json:"sign_key"`
		PayeeId             *int64   `valid:"optional" json:"payee_id"`
		LanguageType        *string  `valid:"optional" json:"language_type"`
		LegalCurrency       *string  `valid:"optional" json:"legal_currency"`
		Contact             *string  `valid:"optional" json:"contact"`
		ContactPhone        *string  `valid:"optional" json:"contact_phone"`
		ContactEmail        *string  `valid:"optional" json:"contact_email"`
		Country             *string  `valid:"optional" json:"country"`
	}
	MerchantGet struct {
		MerchantId          *string  `valid:"required" json:"merchant_id"`
	}
	MerchantDel struct {
		MerchantId          *string  `valid:"required" json:"merchant_id"`
	}
	MerchantList struct {
		Id                  *uint64  `valid:"optional" json:"id"`
		MerchantId          *string  `valid:"optional" json:"merchant_id"`
		MerchantName        *string  `valid:"optional" json:"merchant_name"`
		NotifyUrl           *string  `valid:"optional" json:"notify_url"`
		SignType            *string  `valid:"optional" json:"sign_type"`
		SignKey             *string  `valid:"optional" json:"sign_key"`
		PayeeId             *int64   `valid:"optional" json:"payee_id"`
		LanguageType        *string  `valid:"optional" json:"language_type"`
		LegalCurrency       *string  `valid:"optional" json:"legal_currency"`
		Contact             *string  `valid:"optional" json:"contact"`
		ContactPhone        *string  `valid:"optional" json:"contact_phone"`
		ContactEmail        *string  `valid:"optional" json:"contact_email"`
		Country             *string  `valid:"optional" json:"country"`

		Page            	int64    `valid:"required" json:"page"`
		Size          		int64    `valid:"optional" json:"size"`
	}
)


func InitMerchantConfig() error {
	return db.GDbMgr.Get().AutoMigrate(&MerchantConfig{}).Error
}

func (this *MerchantAdd) Add()  (error) {
	model := &MerchantConfig{
		MerchantId:    this.MerchantId,
		MerchantName:  this.MerchantName,
		NotifyUrl:     this.NotifyUrl,
		SignType:      this.SignType,
		SignKey:       this.SignKey,
		PayeeId:       this.PayeeId,
		LanguageType:  this.LanguageType,
		LegalCurrency: this.LegalCurrency,
		Contact:       this.Contact,
		ContactPhone:  this.ContactPhone,
		ContactEmail:  this.ContactEmail,
		Country:       this.Country,
	}
	err := db.GDbMgr.Get().Create(model).Error
	if err != nil {
		return  err
	}
	return nil
}

func (this *MerchantUpdate) Update() ( error) {
	model := &MerchantConfig{
		Id: 		   this.Id,
		MerchantId:    this.MerchantId,
		MerchantName:  this.MerchantName,
		NotifyUrl:     this.NotifyUrl,
		SignType:      this.SignType,
		SignKey:       this.SignKey,
		PayeeId:       this.PayeeId,
		LanguageType:  this.LanguageType,
		LegalCurrency: this.LegalCurrency,
		Contact:       this.Contact,
		ContactPhone:  this.ContactPhone,
		ContactEmail:  this.ContactEmail,
		Country:       this.Country,
	}
	if err := db.GDbMgr.Get().Model(&MerchantConfig{}).Where("id = ?", this.Id).Updates(model).Error; err != nil {
		return  err
	}
	return nil
}


func (this *MerchantGet) Get() (*MerchantConfig ,error) {
	model := &MerchantConfig{}
	if err := db.GDbMgr.Get().Find(model,"MERCHANT_ID = ? ",this.MerchantId ).Error; err != nil {
		return model,  err
	}
	if model == nil {
		return nil, errors.New("bkConfig not find")
	}
	return model, nil
}


func (this *MerchantList) List() (*common.Result, error) {
	var list []*MerchantConfig

	query := db.GDbMgr.Get()

	model := &MerchantConfig{
		Id:            this.Id,
		MerchantId:    this.MerchantId,
		MerchantName:  this.MerchantName,
		NotifyUrl:     this.NotifyUrl,
		SignType:      this.SignType,
		SignKey:       this.SignKey,
		PayeeId:       this.PayeeId,
		LanguageType:  this.LanguageType,
		LegalCurrency: this.LegalCurrency,
		Contact:       this.Contact,
		ContactPhone:  this.ContactPhone,
		ContactEmail:  this.ContactEmail,
		Country:       this.Country,
	}
	query = query.Where(model)

	//query = query.Where("MERCHANT_ID = ? ",this.MerchantId)

	return new(common.Result).PageQuery(query, &MerchantConfig{}, &list, this.Page, this.Size, nil, "")
}

func (this *MerchantDel) Delete() error {
	bkConfig := &MerchantConfig{}
	if err := db.GDbMgr.Get().Delete(bkConfig,"MERCHANT_ID = ? ",this.MerchantId).Error; err != nil {
		return   err
	}
	return nil
}

