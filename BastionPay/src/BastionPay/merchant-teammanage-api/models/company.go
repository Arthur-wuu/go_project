package models

import (
	"BastionPay/merchant-teammanage-api/api"
	"BastionPay/merchant-teammanage-api/common"
	"BastionPay/merchant-teammanage-api/db"
	"github.com/jinzhu/gorm"
)

type (
	Company struct {
		Id      *int64  `json:"id,omitempty"        gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		Name    *string `json:"name,omitempty"      gorm:"column:name;type:varchar(50);not null"`
		Address *string `json:"address,omitempty"  gorm:"column:address;type:varchar(50)"`
		Table
	}
)

func (this *Company) TableName() string {
	return "teammanage_company"
}

func (this *Company) ParseAdd(p *api.CompanyAdd) *Company {
	c := &Company{
		Name:    p.Name,
		Address: p.Address,
	}
	c.Valid = p.Vaild
	if c.Valid == nil {
		c.Valid = new(int)
		*c.Valid = 1
	}
	return c
}

func (this *Company) Parse(p *api.Company) *Company {
	c := &Company{
		Id:      p.Id,
		Name:    p.Name,
		Address: p.Address,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Company) ParseList(p *api.CompanyList) *Company {
	c := &Company{
		Name:    p.Name,
		Address: p.Address,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Company) Add() (*Company, error) {
	err := db.GDbMgr.Get().Create(this).Error
	if err != nil {
		return nil, err
	}
	acty := new(Company)
	err = db.GDbMgr.Get().Where("id = ?", *this.Id).Last(acty).Error
	if err != nil {
		return nil, err
	}
	return acty, nil
}

func (this *Company) Get() (*Company, error) {
	acty := new(Company)
	err := db.GDbMgr.Get().Where(this).Last(acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Company) Del() error {
	err := db.GDbMgr.Get().Where("id = ?", *this.Id).Delete(&Company{}).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (this *Company) Update() (*Company, error) {
	if err := db.GDbMgr.Get().Updates(this).Error; err != nil {
		return nil, err
	}
	acty := new(Company)
	err := db.GDbMgr.Get().Where("id = ?", *this.Id).Last(acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Company) ListWithConds(page, size int64, condPair []*SqlPairCondition) (*common.Result, error) {
	var list []*Company
	query := db.GDbMgr.Get().Where(this)
	for i := 0; i < len(condPair); i++ {
		if condPair[i] == nil {
			continue
		}
		query = query.Where(condPair[i].Key, condPair[i].Value)
	}
	query = query.Order("valid desc").Order("id")

	return new(common.Result).PageQuery(query, &Company{}, &list, page, size, nil, "")
}
