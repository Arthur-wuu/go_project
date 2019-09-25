package models

import (
	"BastionPay/bas-wallet-fission-share/api"
	"BastionPay/bas-wallet-fission-share/common"
	"BastionPay/bas-wallet-fission-share/db"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type (
	Activity struct {
		Id           *int             `json:"id,omitempty"        gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		Name         *string          `json:"name,omitempty"      gorm:"column:name;type:varchar(30)"`
		Sponsor      *string          `json:"sponsor,omitempty"   gorm:"column:sponsor;type:varchar(30)"`
		Merchant     *string          `json:"merchant,omitempty"   gorm:"column:merchant;type:varchar(30)"`
		Country      *string          `json:"country,omitempty"   gorm:"column:country;type:varchar(30)"`
		Language     *string          `json:"language,omitempty"   gorm:"column:language;type:varchar(10)"`
		Mode         *int             `json:"mode,omitempty"   gorm:"column:mode;type:int(11)"`
		Symbol       *string          `json:"symbol,omitempty"   gorm:"column:symbol;type:varchar(30)"`
		RemainCoin   *decimal.Decimal `json:"remain_coin,omitempty"   gorm:"column:remain_coin;type:decimal(38,12)"`
		TotalCoin    *decimal.Decimal `json:"total_coin,omitempty"   gorm:"column:total_coin;type:decimal(38,12)"`
		TotalRed     *int64           `json:"total_red,omitempty"   gorm:"column:total_red;type:bigint(20)"`
		RemainRed    *int64           `json:"remain_red,omitempty"   gorm:"column:remain_red;type:bigint(20)"`
		RedExpire    *int64           `json:"red_expire,omitempty"   gorm:"column:red_expire;type:bigint(20)"`
		RobRedExpire *int64           `json:"rob_red_expire,omitempty"   gorm:"column:rob_red_expire;type:bigint(20)"`
		TotalRob     *int64           `json:"total_rob,omitempty"   gorm:"column:total_rob;type:bigint(20)"`
		OnAt         *int64           `json:"on_at,omitempty"   gorm:"column:on_at;type:bigint(20)"`
		OffAt        *int64           `json:"off_at,omitempty"   gorm:"column:off_at;type:bigint(20)"`
		Precision    *int32           `json:"precision,omitempty"  gorm:"column:precision;type:int(11)"`
		Table
	}
)

func (this *Activity) TableName() string {
	return "fission_activity"
}

func (this *Activity) ParseAdd(p *api.ActivityAdd) *Activity {
	acty := &Activity{
		Name:         p.Name,
		Sponsor:      p.Sponsor,
		Merchant:     p.Merchant,
		Country:      p.Country,
		Language:     p.Language,
		Mode:         p.Mode,
		Symbol:       p.Symbol,
		RemainCoin:   p.TotalCoin,
		TotalCoin:    p.TotalCoin,
		TotalRed:     p.TotalRed,
		RemainRed:    p.TotalRed,
		RedExpire:    p.RedExpire,
		RobRedExpire: p.RobRedExpire,
		TotalRob:     p.TotalRob,
		OnAt:         p.OnAt,
		OffAt:        p.OffAt,
		Precision:    p.Precision,
	}
	acty.Valid = p.Valid
	if acty.Valid == nil {
		acty.Valid = new(int)
		*acty.Valid = 1
	}
	return acty
}

func (this *Activity) Parse(p *api.Activity) *Activity {
	acty := &Activity{
		Id:           p.Id,
		Name:         p.Name,
		Sponsor:      p.Sponsor,
		Merchant:     p.Merchant,
		Country:      p.Country,
		Language:     p.Language,
		Mode:         p.Mode,
		Symbol:       p.Symbol,
		TotalCoin:    p.TotalCoin,
		TotalRed:     p.TotalRed,
		RedExpire:    p.RedExpire,
		RobRedExpire: p.RobRedExpire,
		TotalRob:     p.TotalRob,
		OnAt:         p.OnAt,
		OffAt:        p.OffAt,
		Precision:    p.Precision,
	}
	acty.Valid = p.Valid
	return acty
}

func (this *Activity) ParseList(p *api.ActivityList) *Activity {
	acty := &Activity{
		Name:         p.Name,
		Sponsor:      p.Sponsor,
		Merchant:     p.Merchant,
		Country:      p.Country,
		Language:     p.Language,
		Mode:         p.Mode,
		Symbol:       p.Symbol,
		RemainCoin:   p.TotalCoin,
		TotalCoin:    p.TotalCoin,
		TotalRed:     p.TotalRed,
		RemainRed:    p.TotalRed,
		RedExpire:    p.RedExpire,
		RobRedExpire: p.RobRedExpire,
		OnAt:         p.OnAt,
		OffAt:        p.OffAt,
	}
	acty.Valid = p.Valid
	return acty
}

func (this *Activity) Unique() (bool, error) {
	count := 0
	err := db.GDbMgr.Get().Model(this).Where("name = ?", this.Name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, err
}

func (this *Activity) Add() (*Activity, error) {
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

func (this *Activity) Update() (*Activity, error) {
	err := db.GDbMgr.Get().Model(this).Update(this).Error
	if err != nil {
		return nil, err
	}
	err = db.GDbMgr.Get().Where("id = ?", *this.Id).Last(this).Error
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *Activity) TxUpdate(tx *gorm.DB) (*Activity, error) {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	err := tx.Model(this).Update(this).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (this *Activity) List(page, size int64) (*common.Result, error) {
	var list []*Activity
	query := db.GDbMgr.Get().Where(this)

	return new(common.Result).PageQuery(query, &Activity{}, &list, page, size, nil, "")
}

//func (this *Activity) ListForFront(page, size int64) (*common.Result, error) {
//	var list []*Activity
//	query := db.GDbMgr.Get().Where(this).Select("name", "sponsor", "mode", "symbol", )
//
//	return new(common.Result).PageQuery(query, &Activity{}, &list, page, size, nil, "")
//}

func (this *Activity) Get() (*Activity, error) {
	err := db.GDbMgr.Get().Where("id = ?", *this.Id).Last(this).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return this, err
}

func (this *Activity) TxGetById(tx *gorm.DB, id int) (*Activity, error) {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	err := tx.Where("id = ?", id).Last(this).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return this, err
}

func (this *Activity) ExistById(id int) (bool, error) {
	count := 0
	err := db.GDbMgr.Get().Model(this).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (this *Activity) ReduceRed(id int) (*Activity, error) {
	err := db.GDbMgr.Get().Model(this).Where("id = ?", id).Update("remain_red", gorm.Expr("remain_red - 1")).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return this, err
}

//func (this *Activity) TxReduceRed(tx *gorm.DB, id int,redNum int64, coin float64) (*Activity, error) {
//	if tx == nil {
//		tx = db.GDbMgr.Get()
//	}
//	err := tx.Model(this).Where("id = ?", id).Update("remain_red", gorm.Expr("remain_red - 1")).Error
//	if err == gorm.ErrRecordNotFound {
//		return nil,nil
//	}
//	return this,err
//}
