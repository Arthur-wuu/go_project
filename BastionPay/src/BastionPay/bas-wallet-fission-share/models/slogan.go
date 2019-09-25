package models

import (
	"BastionPay/bas-wallet-fission-share/api"
	"BastionPay/bas-wallet-fission-share/common"
	"BastionPay/bas-wallet-fission-share/db"
	"github.com/jinzhu/gorm"
)

type (
	Slogan struct {
		Id         *int    `json:"id,omitempty"         gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		ActivityId *int    `json:"activity_id,omitempty"  gorm:"column:activity_id;type:int(11)"`
		TitleId    *int    `json:"title_id,omitempty"  gorm:"column:src_uid;type:bigint(20)"`
		Text       *string `json:"text,omitempty"    gorm:"column:text;type:varchar(200)"`
		Table
	}
)

func (this *Slogan) TableName() string {
	return "fission_slogan"
}

func (this *Slogan) ParseAdd(p *api.SloganAdd) *Slogan {
	acty := &Slogan{
		ActivityId: p.ActivityId,
		TitleId:    p.TitleId,
		Text:       p.Text,
	}
	acty.Valid = p.Valid
	if acty.Valid == nil {
		acty.Valid = new(int)
		*acty.Valid = 1
	}
	return acty
}

func (this *Slogan) Parse(p *api.Slogan) *Slogan {
	acty := &Slogan{
		ActivityId: p.ActivityId,
		TitleId:    p.TitleId,
		Text:       p.Text,
	}
	acty.Valid = p.Valid
	return acty
}

func (this *Slogan) ParseList(p *api.SloganList) *Slogan {
	acty := &Slogan{
		ActivityId: p.ActivityId,
		TitleId:    p.TitleId,
		Text:       p.Text,
	}
	acty.Valid = p.Valid
	return acty
}

func (this *Slogan) Add() (*Slogan, error) {
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

func (this *Slogan) Update() (*Slogan, error) {
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

func (this *Slogan) List(page, size int64) (*common.Result, error) {
	var list []*Slogan
	query := db.GDbMgr.Get().Where(this)

	return new(common.Result).PageQuery(query, &Slogan{}, &list, page, size, nil, "")
}

func (this *Slogan) GetAll(activity_id int) ([]*Slogan, error) {
	var list []*Slogan
	err := db.GDbMgr.Get().Where("activity_id = ?", activity_id).Find(list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return list, err
}
