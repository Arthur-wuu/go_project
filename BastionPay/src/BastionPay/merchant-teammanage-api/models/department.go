package models

import (
	"BastionPay/merchant-teammanage-api/api"
	"BastionPay/merchant-teammanage-api/common"
	"BastionPay/merchant-teammanage-api/db"
	"github.com/jinzhu/gorm"
)

type (
	Department struct {
		Id        *int64  `json:"id,omitempty"         gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		CompanyId *int64  `json:"company_id,omitempty"   gorm:"column:company_id;type:bigint(20)"`
		ParentId  *int64  `json:"parent_id,omitempty"   gorm:"column:parent_id;type:bigint(20)"`
		Name      *string `json:"name,omitempty"        gorm:"column:name;type:varchar(50)"`
		EpeNum    *int    `json:"employee_num,omitempty"  gorm:"column:employee_num;type:int(11);default 0"`
		Table
	}
)

func (this *Department) TableName() string {
	return "teammanage_department"
}

func (this *Department) ParseAdd(p *api.DepartmentAdd) *Department {
	c := &Department{
		Name:      p.Name,
		CompanyId: p.CompanyId,
		EpeNum:    p.EpeNum,
		ParentId:  p.ParentId,
	}
	c.Valid = p.Vaild
	if c.Valid == nil {
		c.Valid = new(int)
		*c.Valid = 1
	}
	return c
}

func (this *Department) Parse(p *api.Department) *Department {
	c := &Department{
		Id:        p.Id,
		Name:      p.Name,
		CompanyId: p.CompanyId,
		EpeNum:    p.EpeNum,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Department) FtParseGet(p *api.FtDepartmentGet) *Department {
	c := &Department{
		Id:        p.Id,
		Name:      p.Name,
		CompanyId: p.CompanyId,
		EpeNum:    p.EpeNum,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Department) FtParseList(p *api.FtDepartmentList) *Department {
	c := &Department{
		Name:      p.Name,
		CompanyId: p.CompanyId,
		EpeNum:    p.EpeNum,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Department) ParseList(p *api.DepartmentList) *Department {
	c := &Department{
		Name:      p.Name,
		CompanyId: p.CompanyId,
		EpeNum:    p.EpeNum,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Department) Add() (*Department, error) {
	err := db.GDbMgr.Get().Create(this).Error
	if err != nil {
		return nil, err
	}
	acty := new(Department)
	err = db.GDbMgr.Get().Where("id = ?", *this.Id).Last(acty).Error
	if err != nil {
		return nil, err
	}
	return acty, nil
}

func (this *Department) ExistByCompanyId(CompanyId int64) (bool, error) {
	count := 0
	err := db.GDbMgr.Get().Where("company_id = ?", CompanyId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (this *Department) Get() (*Department, error) {
	acty := new(Department)
	err := db.GDbMgr.Get().Where(this).Last(acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Department) Del() error {
	err := db.GDbMgr.Get().Where("id = ?", *this.Id).Delete(&Department{}).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (this *Department) GetsByCompany(cId int64) ([]*Department, error) {
	var acty []*Department
	err := db.GDbMgr.Get().Where("company_id = ?", cId).Find(&acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Department) Update() (*Department, error) {
	if err := db.GDbMgr.Get().Updates(this).Error; err != nil {
		return nil, err
	}
	acty := new(Department)
	err := db.GDbMgr.Get().Where("id = ?", *this.Id).Last(acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Department) ListWithConds(page, size int64, condPair []*SqlPairCondition) (*common.Result, error) {
	var list []*Department
	query := db.GDbMgr.Get().Where(this)
	for i := 0; i < len(condPair); i++ {
		if condPair[i] == nil {
			continue
		}
		query = query.Where(condPair[i].Key, condPair[i].Value)
	}
	query = query.Order("valid desc").Order("id")

	return new(common.Result).PageQuery(query, &Department{}, &list, page, size, nil, "")
}
