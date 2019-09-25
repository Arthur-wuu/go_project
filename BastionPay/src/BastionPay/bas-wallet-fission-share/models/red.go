package models

import (
	"BastionPay/bas-wallet-fission-share/api"
	"BastionPay/bas-wallet-fission-share/common"
	"BastionPay/bas-wallet-fission-share/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

type (
	Red struct {
		Id         *int             `json:"id,omitempty"         gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		ActivityId *int             `json:"activity_id,omitempty"  gorm:"column:activity_id;type:int(11)"`
		Symbol     *string          `json:"symbol,omitempty"   gorm:"column:symbol;type:varchar(30)"`
		TotalRob   *int64           `json:"total_rob,omitempty"   gorm:"column:total_rob;type:bigint(20)"` //红包领取个数
		RemainRob  *int64           `json:"remain_rob,omitempty"   gorm:"column:remain_rob;type:bigint(20)"`
		TotalCoin  *decimal.Decimal `json:"total_coin,omitempty"   gorm:"column:total_coin;type:decimal(38,12)"`
		RemainCoin *decimal.Decimal `json:"remain_coin,omitempty"   gorm:"column:remain_coin;type:decimal(38,12)"`
		OrderId    *string          `json:"order_id,omitempty"   gorm:"column:order_id;type:varchar(100)"`
		Mode       *int             `json:"mode,omitempty"   gorm:"column:mode;type:int(11)"`
		SrcUid     *int64           `json:"src_uid,omitempty"    gorm:"column:src_uid;type:bigint(20)"`
		ExpireAt   *int64           `json:"expire_at,omitempty"  gorm:"column:expire_at;type:bigint(20)"`
		RobExpire  *int64           `json:"rob_expire,omitempty"  gorm:"column:rob_expire;type:bigint(20)"`
		Precision  *int32           `json:"precision,omitempty"  gorm:"column:precision;type:int(11)"`
		Table
	}
)

func (this *Red) TableName() string {
	return "fission_red"
}

func (this *Red) ParseAdd(p *api.RedAdd, acty *Activity) *Red {
	r := &Red{
		ActivityId: p.ActivityId,
		Symbol:     p.Symbol,
		TotalRob:   p.TotalSize,
		RemainRob:  p.TotalSize,
		TotalCoin:  p.TotalCoin,
		RemainCoin: p.TotalCoin,
		Mode:       p.Mode,
		SrcUid:     p.SrcUid,
		ExpireAt:   p.ExpireAt,
		RobExpire:  p.RobExpire,
		OrderId:    p.OrderId,
	}
	r.Valid = p.Valid
	if r.Valid == nil {
		r.Valid = new(int)
		*r.Valid = 1
	}
	r.ExpireAt = new(int64)
	r.RobExpire = acty.RobRedExpire
	*r.ExpireAt = time.Now().Unix() + *acty.RedExpire
	r.Symbol = acty.Symbol
	r.Mode = acty.Mode
	r.TotalCoin = new(decimal.Decimal)
	if *acty.RemainRed == 1 {
		*r.TotalCoin = *acty.RemainCoin
	} else {
		*r.TotalCoin = acty.TotalCoin.Div(decimal.NewFromFloat(float64(*acty.TotalRed))).Truncate(*acty.Precision)
	}
	fmt.Println(" remian=", *acty.RemainCoin)
	r.RemainCoin = r.TotalCoin
	*acty.RemainCoin = acty.RemainCoin.Sub(*r.TotalCoin)
	*acty.RemainRed--
	r.TotalRob = acty.TotalRob
	r.RemainRob = acty.TotalRob
	r.Precision = acty.Precision
	fmt.Println("totoal=coin=", *r.TotalCoin, " remian=", *acty.RemainCoin)
	return r
}

func (this *Red) ParseList(p *api.RedList) *Red {
	r := &Red{
		Id:         p.Id,
		ActivityId: p.ActivityId,
		Symbol:     p.Symbol,
		Mode:       p.Mode,
		SrcUid:     p.SrcUid,
		ExpireAt:   p.ExpireAt,
		OrderId:    p.OrderId,
	}
	r.Valid = p.Valid
	return r
}

func (this *Red) TxUnique(tx *gorm.DB) (bool, error) {
	if this.OrderId == nil || len(*this.OrderId) == 0 {
		return true, nil
	}
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	count := 0
	err := tx.Model(this).Where("activity_id = ? and order_id = ?", this.ActivityId, this.OrderId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, err
}

func (this *Red) Add() (*Red, error) {
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

func (this *Red) TxAdd(tx *gorm.DB) (*Red, error) {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	err := tx.Create(this).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (this *Red) Update() (*Red, error) {
	err := db.GDbMgr.Get().Model(this).Update(this).Error
	if err != nil {
		return nil, err
	}
	err = db.GDbMgr.Get().Model(this).Where("id = ?", *this.Id).Last(this).Error
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *Red) TxUpdate(tx *gorm.DB) (*Red, error) {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	err := tx.Model(this).Update(this).Error
	if err != nil {
		return nil, err
	}
	err = tx.Model(this).Where("id = ?", *this.Id).Last(this).Error
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *Red) GetById(id int) (*Red, error) {
	red := new(Red)
	err := db.GDbMgr.Get().Where("id = ?", id).Last(red).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return red, nil
}

func (this *Red) TxGetById(tx *gorm.DB, id int) (*Red, error) {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	red := new(Red)
	err := tx.Where("id = ?", id).Last(red).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return red, nil
}

func (this *Red) List(page, size int64) (*common.Result, error) {
	var list []*Red
	query := db.GDbMgr.Get().Where(this)

	return new(common.Result).PageQuery(query, &Robber{}, &list, page, size, nil, "")
}
