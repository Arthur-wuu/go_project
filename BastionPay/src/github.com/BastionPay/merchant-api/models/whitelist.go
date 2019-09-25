package models

import (
	//"BastionPay/merchant-api/common"
	//"BastionPay/merchant-api/common"
	"BastionPay/merchant-api/db"
	"errors"
	//"strconv"
)


type (
	WhiteList struct {
		Id            *uint64     `json:"id"              gorm:"AUTO_INCREMENT;primary_key;column:id"`
		CreatedAt     *int64      `json:"created_at"      gorm:"column:created_at;type:bigint(20)"`
		UpdatedAt     *int64      `json:"updated_at"      gorm:"column:updated_at;type:bigint(20)"`

		WhiteUrl       *string    `json:"white_url"       gorm:"column:white_url;type:varchar(80)"`
		//Symbol         *string    `json:"symbol"          gorm:"column:symbol;type:varchar(10)"`
		//Amount         *string    `json:"amount"          gorm:"column:amount;type:varchar(10)"`
		//DeviceType     *string    `json:"device_type"     gorm:"column:device_type;type:varchar(10)"`
		//SupportCoin    *string    `json:"support_coin"    gorm:"column:support_coin;type:varchar(90)"`
		//PayeeId        *int64     `json:"payee_id"        gorm:"column:payee_id;type:bigint(20)"`
		//ReturnUrl      *string    `json:"return_url"      gorm:"column:return_url;type:varchar(100)"`
		//ShowUrl        *string    `json:"show_url"        gorm:"column:show_url;type:varchar(100)"`

	}
)

func (this *WhiteList) TableName() string {
	return "white_list"
}

type (
	WhiteListGet struct {
		DeviceNo    *string `valid:"optional" json:"white_url"`
	}

)

func InitWhiteList() error {
	return db.GDbMgr.Get().AutoMigrate(&WhiteList{}).Error
}

func (this *WhiteListGet) Get() (*string ,error) {
	model := &WhiteList{}
	if err := db.GDbMgr.Get().Last(model).Error; err != nil {
		return model.WhiteUrl,  err
	}
	if model == nil {
		return nil, errors.New("WhiteList not find")
	}
	return model.WhiteUrl, nil
}
//
//func (this *BkConfigAdd) Add()  (error) {
//	//bkConfig := &BkConfigAdd{}
//	model := &DeviceConfig{
//		DeviceNo:    this.DeviceNo,
//		DeviceType:  this.DeviceType,
//		Amount:      this.Amount,
//		Symbol:      this.Symbol,
//		SupportCoin: this.SupportCoin,
//		PayeeId:     this.PayeeId,
//		ReturnUrl:   this.ReturnUrl,
//		ShowUrl:     this.ShowUrl,
//	}
//	err := db.GDbMgr.Get().Create(model).Error
//	if err != nil {
//		return  err
//	}
//	return nil
//}
//
//func (this *BkConfigUpdate) Update() ( error) {
//	model := &DeviceConfig{
//		Id:  		 this.Id,
//		DeviceNo:    this.DeviceNo,
//		DeviceType:  this.DeviceType,
//		Amount:      this.Amount,
//		Symbol:      this.Symbol,
//		SupportCoin: this.SupportCoin,
//		PayeeId:     this.PayeeId,
//		ReturnUrl:   this.ReturnUrl,
//		ShowUrl:     this.ShowUrl,
//	}
//	if err := db.GDbMgr.Get().Model(&DeviceConfig{}).Where("id = ?", this.Id).Updates(model).Error; err != nil {
//		return  err
//	}
//	return nil
//}

//
//
//func (this *BkConfigGet) GetCoinList() (*string ,error) {
//	model := &DeviceConfig{}
//	if err := db.GDbMgr.Get().Find(model,"device_no = ? ",this.DeviceNo ).Error; err != nil {
//		return nil,  err
//	}
//	if model == nil {
//		return nil, errors.New("bkConfig not find")
//	}
//	return model.SupportCoin, nil
//}
//
//func  GetPrice(no string ) (*float64 ,error) {
//	model := &DeviceConfig{}
//	if err := db.GDbMgr.Get().Find(model,"device_no = ? ", no ).Error; err != nil {
//		return nil,  err
//	}
//	if model == nil {
//		return nil, errors.New("bkConfig not find")
//	}
//
//	amount, err1 := strconv.ParseFloat(*model.Amount, 64)
//	if err1 != nil {
//		return nil, errors.New("bkConfig not find11")
//	}
//	return  &amount, nil
//}
//
//func (this *BkConfig) GetPayeeId(no string ) (*int64 ,error) {
//	model := &DeviceConfig{}
//	if err := db.GDbMgr.Get().Find(model,"device_no = ? ", no ).Error; err != nil {
//		return nil,  err
//	}
//	if model == nil {
//		return nil, errors.New("bkConfig not find")
//	}
//	return  model.PayeeId, nil
//}
//
//
//func (this *BkConfigList) List() (*common.Result, error) {
//	var list []*DeviceConfig
//
//	query := db.GDbMgr.Get()
//
//	model := &DeviceConfig{
//		DeviceNo:    this.DeviceNo,
//		DeviceType:  this.DeviceType,
//		Amount:      this.Amount,
//		Symbol:      this.Symbol,
//		SupportCoin: this.SupportCoin,
//		PayeeId:     this.PayeeId,
//		ReturnUrl:   this.ReturnUrl,
//		ShowUrl:     this.ShowUrl,
//	}
//	query = query.Where(model)
//
//	return new(common.Result).PageQuery(query, &DeviceConfig{}, &list, this.Page, this.Size, nil, "")
//}
//
//func (this *BkConfigDel) Delete() error {
//	bkConfig := &DeviceConfig{}
//	if err := db.GDbMgr.Get().Delete(bkConfig,"device_no = ? ",this.DeviceNo).Error; err != nil {
//		return   err
//	}
//	return nil
//}
//
