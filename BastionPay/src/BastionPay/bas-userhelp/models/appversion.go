package models

import (
	"BastionPay/bas-userhelp/common"
	"BastionPay/bas-userhelp/db"
	"BastionPay/bas-userhelp/models/table"
	"fmt"
)

type (
	AppVersionSet struct {
		Name         *string `valid:"optional" json:"name"`
		LabelId      *uint64 `valid:"required" json:"label_id" `
		Language     *string `valid:"optional" json:"language"`
		ShowMode     *int    `valid:"optional" json:"show_mode"`
		Version      *string `valid:"required" json:"version"`
		Url          *string `valid:"optional" json:"url"`
		Instructions *string `valid:"optional" json:"instructions"`
		SysType      *string `valid:"required" json:"sys_type"`
		UpgradeMode  *int    `valid:"optional" json:"upgrade_mode"`
		UpgradeAt    *int64  `valid:"optional" json:"upgrade_at"`
	}

	AppVersionGet struct {
		//NameId        *int     `valid:"optional" json:"name_id"`
		Name     *string `valid:"optional" json:"name"`
		Language *string `valid:"optional" json:"language"`
		SysType  *string `valid:"required" json:"sys_type"`
	}

	AppVersionUpdate struct {
		Id           *uint64 `valid:"required" json:"id"`
		LabelId      *uint64 `valid:"optional" json:"label_id" `
		Language     *string `valid:"optional" json:"language"`
		ShowMode     *int    `valid:"optional" json:"show_mode"`
		Version      *string `valid:"optional" json:"version"`
		Url          *string `valid:"optional" json:"url"`
		Instructions *string `valid:"optional" json:"instructions"`
		SysType      *string `valid:"optional" json:"sys_type"`
		UpgradeMode  *int    `valid:"optional" json:"upgrade_mode"`
		UpgradeAt    *int64  `valid:"optional" json:"upgrade_at"`
	}

	AppVersionList struct {
		LabelId        *uint64 `valid:"optional" json:"label_id" `
		Language       *string `valid:"optional" json:"language"`
		ShowMode       *int    `valid:"optional" json:"show_mode"`
		Version        *string `valid:"optional" json:"version"`
		Url            *string `valid:"optional" json:"url"`
		Instructions   *string `valid:"optional" json:"instructions"`
		SysType        *string `valid:"optional" json:"sys_type"`
		UpgradeMode    *int    `valid:"optional" json:"upgrade_mode"`
		UpgradeAt      *int64  `valid:"optional" json:"upgrade_at"`
		Id             *uint64 `valid:"optional" json:"id"`
		StartCreatedAt *int64  `valid:"optional" json:"start_created_at"`
		EndCreatedAt   *int64  `valid:"optional" json:"end_created_at"`
		StartUpdatedAt *int64  `valid:"optional" json:"start_updated_at"`
		EndUpdatedAt   *int64  `valid:"optional" json:"end_updated_at"`
		Page           int64   `valid:"required" json:"page"`
		Size           int64   `valid:"optional" json:"size"`
	}
)

func CreateAppvsionTab() error {
	return db.GDbMgr.Get().AutoMigrate(&table.AppVersion{}).Error
}

func (this *AppVersionSet) Set() (*table.AppVersion, error) {
	model := &table.AppVersion{
		Name:         this.Name,
		LabelId:      this.LabelId,
		Version:      this.Version,
		Url:          this.Url,
		Instructions: this.Instructions,
		Language:     this.Language,
		ShowMode:     this.ShowMode,
		UpgradeMode:  this.UpgradeMode,
		SysType:      this.SysType,
		UpgradeAt:    this.UpgradeAt,
	}

	err := db.GDbMgr.Get().Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

//
//func (this *AppVersionGet) Get() (*table.AppVersion, error) {
//	model := &table.AppVersion{
//		LabelId:         this.LabelId,
//		Version:      this.Version,
//		Url:          this.Url,
//		Instructions: this.Instructions,
//		Language:     this.Language,
//		ShowMode:     this.ShowMode,
//		UpgradeMode:  this.UpgradeMode,
//		SysType:      this.SysType,
//	}
//
//	err := db.GDbMgr.Get().Last(model).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return model, nil
//}

func (this *AppVersionGet) GetForFront(name, SysType *string) (*table.AppVersion, error) {
	model := new(table.AppVersion)
	err := db.GDbMgr.Get().Where("name = ? and sys_type = ?", name, SysType).Last(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (this *AppVersionUpdate) Update() (*table.AppVersion, error) {
	model := &table.AppVersion{
		Id:           this.Id,
		LabelId:      this.LabelId,
		Version:      this.Version,
		Url:          this.Url,
		Instructions: this.Instructions,
		Language:     this.Language,
		ShowMode:     this.ShowMode,
		UpgradeMode:  this.UpgradeMode,
		SysType:      this.SysType,
	}
	if err := db.GDbMgr.Get().Model(&table.AppVersion{}).Where("id = ?", this.Id).Update(model).Error; err != nil {
		return nil, err
	}

	appVersion := &table.AppVersion{}
	if err := db.GDbMgr.Get().Where(&table.AppVersion{}).Find(appVersion).Error; err != nil {
		return nil, err
	}

	return appVersion, nil
}

func (this *AppVersionSet) RowsAffectNumUpdate() (int64, error) {
	model := &table.AppVersion{
		Version:      this.Version,
		Url:          this.Url,
		Instructions: this.Instructions,
		Language:     this.Language,
		ShowMode:     this.ShowMode,
		UpgradeMode:  this.UpgradeMode,
		SysType:      this.SysType,
	}
	fmt.Println("go here up ****")
	newDB := db.GDbMgr.Get().Model(&table.AppVersion{}).Update(model)
	fmt.Println("go here up 222****")

	err := db.GDbMgr.Get().Model(&table.AppVersion{}).Create(model).Error

	return newDB.RowsAffected, err
}

func (this *AppVersionList) List() (*common.Result, error) {
	var list []*table.AppVersion

	query := db.GDbMgr.Get()

	model := &table.AppVersion{
		LabelId:      this.LabelId,
		Version:      this.Version,
		Url:          this.Url,
		Instructions: this.Instructions,
		Language:     this.Language,
		ShowMode:     this.ShowMode,
		UpgradeMode:  this.UpgradeMode,
		SysType:      this.SysType,
	}
	query = query.Where(model)

	if this.StartCreatedAt != nil && this.EndCreatedAt != nil {
		query = query.Where("created_at BETWEEN ? AND ?", this.StartCreatedAt, this.EndCreatedAt)
	}
	if this.StartUpdatedAt != nil && this.EndUpdatedAt != nil {
		query = query.Where("updated_at BETWEEN ? AND ?", this.StartUpdatedAt, this.EndUpdatedAt)
	}

	query = query.Order("created_at desc")

	return new(common.Result).PageQuery(query, &table.AppVersion{}, &list, this.Page, this.Size, nil, "")
}
