package models

import (
	"BastionPay/bas-userhelp/common"
	"BastionPay/bas-userhelp/db"
	"BastionPay/bas-userhelp/models/table"
)

type (
	UserHelpAdd struct {
		CountryCode *string `valid:"optional" json:"country_code"`
		Phone       *string `valid:"optional" json:"phone"`
		Email       *string `valid:"optional" json:"email"`
		Website     *string `valid:"optional" json:"website"`
		AppName     *string `valid:"optional" json:"app_name"`
		Remark      *string `valid:"optional" json:"remark"`
		Status      *int    `valid:"optional" json:"status"`
		Type        *int    `valid:"optional" json:"type"`
		Entrance    *string `valid:"optional" json:"entrance"`
		Name        *string `valid:"optional" json:"name"`
	}

	UserHelpUpdate struct {
		Id     *uint64 `valid:"required" json:"id"`
		Status *int    `valid:"required" json:"status"`
	}

	UserHelpList struct {
		UserHelpAdd

		Id             *uint64 `valid:"optional" json:"id"`
		StartCreatedAt *int64  `valid:"optional" json:"start_created_at"`
		EndCreatedAt   *int64  `valid:"optional" json:"end_created_at"`
		StartUpdatedAt *int64  `valid:"optional" json:"start_updated_at"`
		EndUpdatedAt   *int64  `valid:"optional" json:"end_updated_at"`
		Page           int64   `valid:"required" json:"page"`
		Size           int64   `valid:"optional" json:"size"`
	}
)

func Init() error {
	return db.GDbMgr.Get().AutoMigrate(&table.UserHelp{}).Error
}

func (this *UserHelpAdd) Add() (*table.UserHelp, error) {
	model := &table.UserHelp{
		CountryCode: this.CountryCode,
		Phone:       this.Phone,
		Email:       this.Email,
		Website:     this.Website,
		AppName:     this.AppName,
		Remark:      this.Remark,
		Type:        this.Type,
		Entrance:    this.Entrance,
		Name:        this.Name,
		Status:      this.Status,
	}

	err := db.GDbMgr.Get().Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (this *UserHelpUpdate) Update() (*table.UserHelp, error) {
	if err := db.GDbMgr.Get().Model(&table.UserHelp{}).Where("id = ?", this.Id).Update("status", this.Status).Error; err != nil {
		return nil, err
	}

	userHelp := &table.UserHelp{}
	if err := db.GDbMgr.Get().Where(&table.UserHelp{}).Find(userHelp).Error; err != nil {
		return nil, err
	}

	return userHelp, nil
}

func (this *UserHelpList) List() (*common.Result, error) {
	var list []*table.UserHelp

	query := db.GDbMgr.Get()

	model := &table.UserHelp{
		Id:          this.Id,
		CountryCode: this.CountryCode,
		Phone:       this.Phone,
		Email:       this.Email,
		Website:     this.Website,
		AppName:     this.AppName,
		Remark:      this.Remark,
		Type:        this.Type,
		Entrance:    this.Entrance,
		Name:        this.Name,
		Status:      this.Status,
	}
	query = query.Where(model)

	if this.StartCreatedAt != nil && this.EndCreatedAt != nil {
		query = query.Where("created_at BETWEEN ? AND ?", this.StartCreatedAt, this.EndCreatedAt)
	}
	if this.StartUpdatedAt != nil && this.EndUpdatedAt != nil {
		query = query.Where("updated_at BETWEEN ? AND ?", this.StartUpdatedAt, this.EndUpdatedAt)
	}

	query = query.Order("id desc")

	return new(common.Result).PageQuery(query, &table.UserHelp{}, &list, this.Page, this.Size, nil, "")
}
