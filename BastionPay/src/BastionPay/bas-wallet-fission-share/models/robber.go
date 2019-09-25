package models

import (
	"BastionPay/bas-wallet-fission-share/api"
	"BastionPay/bas-wallet-fission-share/common"
	"BastionPay/bas-wallet-fission-share/db"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

type (
	Robber struct {
		Id           *int             `json:"id,omitempty"         gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		RedId        *int             `json:"red_id,omitempty"     gorm:"column:red_id;type:int(11)"`
		CountryCode  *string          `json:"country_code,omitempty"   gorm:"column:country_code;type:varchar(10)"`
		Phone        *string          `json:"phone,omitempty"      gorm:"column:phone;type:varchar(30)"`
		Symbol       *string          `json:"symbol,omitempty"   gorm:"column:symbol;type:varchar(30)"`
		Coin         *decimal.Decimal `json:"coin,omitempty"   gorm:"column:coin;type:decimal(38,12)"`
		SrcUrl       *string          `json:"src_url,omitempty"    gorm:"column:src_url;type:varchar(200)"`
		SrcUid       *int64           `json:"src_uid,omitempty"    gorm:"column:src_uid;type:bigint(20)"`
		ExpireAt     *int64           `json:"expire_at,omitempty"  gorm:"column:expire_at;type:bigint(20)"`
		TransferFlag *int             `json:"transfer_flag,omitempty"  gorm:"column:transfer_flag;type:int(11);default 0"`
		Table
	}
)

func (this *Robber) TableName() string {
	return "fission_robber"
}

func (this *Robber) ParseAdd(p *api.RobberAdd, red *Red) *Robber {
	this.RedId = p.RedId
	this.CountryCode = p.CountryCode
	this.Phone = p.Phone
	this.Coin = p.Coin
	this.Symbol = red.Symbol
	this.SrcUrl = p.SrcUrl
	this.SrcUid = p.SrcUid
	this.ExpireAt = new(int64)
	*this.ExpireAt = time.Now().Unix() + *red.RobExpire
	this.Valid = new(int)
	*this.Valid = 1
	this.TransferFlag = new(int)
	*this.TransferFlag = 0
	return this
}

func (this *Robber) ParseList(p *api.RobberList) *Robber {
	this.RedId = p.RedId
	this.CountryCode = p.CountryCode
	this.Phone = p.Phone
	this.Coin = p.Coin
	this.Symbol = p.Symbol
	this.SrcUrl = p.SrcUrl
	this.SrcUid = p.SrcUid
	this.ExpireAt = p.ExpireAt

	return this
}

func (this *Robber) Add() (*Robber, error) {
	err := db.GDbMgr.Get().Create(this).Error
	if err != nil {
		return nil, err
	}
	err = db.GDbMgr.Get().Where("id = ?", *this.Id).Last(this).Error
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *Robber) TxAdd(tx *gorm.DB) (*Robber, error) {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	err := tx.Create(this).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (this *Robber) GetBy(CountryCode, phone string) (*Robber, error) {
	robber := new(Robber)
	err := db.GDbMgr.Get().Where("country_code = ? and phone = ? ", CountryCode, phone).Last(robber).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return robber, nil
}

func (this *Robber) TxGetBy(tx *gorm.DB, redId int, CountryCode, phone string) (*Robber, error) {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	robber := new(Robber)
	err := tx.Where("red_id = ? and country_code = ? and phone = ? ", redId, CountryCode, phone).Last(robber).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return robber, nil
}

func (this *Robber) List(page, size int64) (*common.Result, error) {
	var list []*Robber
	query := db.GDbMgr.Get().Where(this)

	return new(common.Result).PageQuery(query, &Robber{}, &list, page, size, nil, "")
}
