package models

import (
	"BastionPay/bas-userhelp/common"
	"BastionPay/bas-userhelp/db"
	"BastionPay/bas-userhelp/models/table"
	"github.com/jinzhu/gorm"
)

type (
	AppWhiteLabelAdd struct {
		//NameId        *int `valid:"optional" json:"name_id"`
		Name          *string `valid:"optional" json:"name"`
	}

	AppWhiteLabel struct {
		Id             *uint64 `valid:"required" json:"id"`
		//NameId        *int `valid:"optional" json:"name_id"`
		Name          *string `valid:"optional" json:"name"`
	}

	AppWhiteLabelList struct {
		//NameId        *int `valid:"optional" json:"name_id"`
		Id             *uint64 `valid:"optional" json:"id"`
		Name          *string `valid:"optional" json:"name"`

		Page           int64   `valid:"required" json:"page"`
		Size           int64   `valid:"optional" json:"size"`
	}
)

func CreateAppWhiteLabel() error {
	return db.GDbMgr.Get().AutoMigrate(&table.AppWhiteLabel{}).Error
}

func (this *AppWhiteLabelAdd) Add() (*table.AppWhiteLabel, error) {
	model := &table.AppWhiteLabel{
		Name:         this.Name,
		//NameId:       this.NameId,
	}

	err := db.GDbMgr.Get().Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (this *AppWhiteLabel) Update() (*table.AppWhiteLabel, error) {
	model := &table.AppWhiteLabel{
		Id :          this.Id,
		Name:         this.Name,
		//NameId:      this.NameId,
	}
	if err := db.GDbMgr.Get().Model(&table.AppWhiteLabel{}).Update(model).Error; err != nil {
		return nil, err
	}

	appVersion := &table.AppWhiteLabel{}
	if err := db.GDbMgr.Get().Where(&table.AppWhiteLabel{}).Last(appVersion).Error; err != nil {
		return nil, err
	}

	return appVersion, nil
}

func (this *AppWhiteLabel) GetBy(id *uint64) (*table.AppWhiteLabel, error) {
	appVersion := &table.AppWhiteLabel{}
	err := db.GDbMgr.Get().Where(&table.AppWhiteLabel{}).Where("id = ? ", id).Last(appVersion).Error;
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return appVersion, nil
}

func (this *AppWhiteLabelList) List() (*common.Result, error) {
	var list []*table.AppWhiteLabel

	query := db.GDbMgr.Get()

	model := &table.AppWhiteLabel{
		Id: this.Id,
		Name:      this.Name,
		//NameId:    this.NameId,
	}
	query = query.Where(model)

	return new(common.Result).PageQuery(query, &table.AppWhiteLabel{}, &list, this.Page, this.Size, nil, "")
}
