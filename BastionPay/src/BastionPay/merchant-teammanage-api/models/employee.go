package models

import (
	"BastionPay/merchant-teammanage-api/api"
	"BastionPay/merchant-teammanage-api/common"
	"BastionPay/merchant-teammanage-api/db"
	"github.com/jinzhu/gorm"
)

type (
	Employee struct {
		Id           *int64  `json:"id,omitempty"        gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		CompanyId    *int64  `json:"company_id,omitempty"   gorm:"column:company_id;type:bigint(20)"`
		DepartmentId *int64  `json:"department_id,omitempty"   gorm:"column:department_id;type:bigint(20)"`
		Name         *string `json:"name,omitempty"     gorm:"column:name;type:varchar(40)"`
		Sex          *int    `json:"sex,omitempty"   gorm:"column:sex;type:tinyint(2)"`
		HiredAt      *int64  `json:"hired_at,omitempty"   gorm:"column:hired_at;type:bigint(20)"`
		BirthAt      *int64  `json:"birth_at,omitempty"   gorm:"column:birth_at;type:bigint(20)"`
		EeNo         *string `json:"ee_no,omitempty"     gorm:"column:ee_no;type:varchar(70)"`
		MechNo       *string `json:"mech_no,omitempty"     gorm:"column:mech_no;type:varchar(70)"`
		BpUid        *int64  `json:"bp_uid,omitempty"   gorm:"column:bp_uid;type:bigint(20)"`
		Phone        *string `json:"phone,omitempty"     gorm:"column:phone;type:varchar(25)"`
		Wchat        *string `json:"wchat,omitempty"     gorm:"column:wchat;type:varchar(30)"`
		Email        *string `json:"email,omitempty"     gorm:"column:email;type:varchar(40)"`
		RegularAt    *int64  `json:"regular_at,omitempty"   gorm:"column:regular_at;type:bigint(20)"`
		Table
	}
)

func (this *Employee) TableName() string {
	return "teammanage_employee"
}

func (this *Employee) ParseAdd(p *api.EmployeeAdd) *Employee {
	c := &Employee{
		CompanyId:    p.CompanyId,
		DepartmentId: p.DepartmentId,
		Name:         p.Name,
		Sex:          p.Sex,
		HiredAt:      p.HiredAt,
		BirthAt:      p.BirthAt,
		RegularAt:    p.RegularAt,
		EeNo:         p.EeNo,
		BpUid:        p.BpUid,
		MechNo:       p.MechNo,
		Phone:        p.Phone,
		Wchat:        p.Wchat,
		Email:        p.Email,
	}
	c.Valid = p.Vaild
	if c.Valid == nil {
		c.Valid = new(int)
		*c.Valid = 1
	}
	return c
}

func (this *Employee) Parse(p *api.Employee) *Employee {
	c := &Employee{
		Id:           p.Id,
		CompanyId:    p.CompanyId,
		DepartmentId: p.DepartmentId,
		Name:         p.Name,
		Sex:          p.Sex,
		HiredAt:      p.HiredAt,
		BirthAt:      p.BirthAt,
		RegularAt:    p.RegularAt,
		EeNo:         p.EeNo,
		BpUid:        p.BpUid,
		MechNo:       p.MechNo,
		Phone:        p.Phone,
		Wchat:        p.Wchat,
		Email:        p.Email,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Employee) FtParseGet(p *api.FtEmployeeGet) *Employee {
	c := &Employee{
		Id:           p.Id,
		CompanyId:    p.CompanyId,
		DepartmentId: p.DepartmentId,
		Name:         p.Name,
		Sex:          p.Sex,
		BpUid:        p.BpUid,
		EeNo:         p.EeNo,
		MechNo:       p.MechNo,
	}
	c.Valid = p.Vaild
	if c.Valid == nil {
		c.Valid = new(int)
		*c.Valid = 1
	}
	return c
}

func (this *Employee) FtParseGets(p *api.FtEmployeeGets) *Employee {
	c := &Employee{
		CompanyId:    p.CompanyId,
		DepartmentId: p.DepartmentId,
	}
	c.Valid = p.Vaild
	if c.Valid == nil {
		c.Valid = new(int)
		*c.Valid = 1
	}
	return c
}

func (this *Employee) ParseList(p *api.EmployeeList) *Employee {
	c := &Employee{
		CompanyId:    p.CompanyId,
		DepartmentId: p.DepartmentId,
		Name:         p.Name,
		Sex:          p.Sex,
		HiredAt:      p.HiredAt,
		BirthAt:      p.BirthAt,
		EeNo:         p.EeNo,
		BpUid:        p.BpUid,
		MechNo:       p.MechNo,
		Phone:        p.Phone,
		Wchat:        p.Wchat,
		Email:        p.Email,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Employee) FtParseList(p *api.FtEmployeeList) *Employee {
	c := &Employee{
		CompanyId:    p.CompanyId,
		DepartmentId: p.DepartmentId,
		Name:         p.Name,
		Sex:          p.Sex,
		HiredAt:      p.HiredAt,
		BirthAt:      p.BirthAt,
		EeNo:         p.EeNo,
		BpUid:        p.BpUid,
		MechNo:       p.MechNo,
	}
	c.Valid = p.Vaild
	return c
}

func (this *Employee) ExistByDepartmentId(departmentId int64) (bool, error) {
	count := 0
	err := db.GDbMgr.Get().Where("department_id = ?", departmentId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (this *Employee) Add() (*Employee, error) {
	err := db.GDbMgr.Get().Create(this).Error
	if err != nil {
		return nil, err
	}
	acty := new(Employee)
	err = db.GDbMgr.Get().Where("id = ?", *this.Id).Last(acty).Error
	if err != nil {
		return nil, err
	}
	return acty, nil
}

func (this *Employee) Get() (*Employee, error) {
	acty := new(Employee)
	err := db.GDbMgr.Get().Where(this).Last(acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Employee) Del() error {
	err := db.GDbMgr.Get().Where("id = ?", *this.Id).Delete(&Employee{}).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (this *Employee) Gets() ([]*Employee, error) {
	var acty []*Employee
	err := db.GDbMgr.Get().Where(this).Find(&acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Employee) Update() (*Employee, error) {
	if err := db.GDbMgr.Get().Updates(this).Error; err != nil {
		return nil, err
	}
	acty := new(Employee)
	err := db.GDbMgr.Get().Where("id = ?", *this.Id).Last(acty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return acty, err
}

func (this *Employee) ListWithConds(page, size int64, condPair []*SqlPairCondition) (*common.Result, error) {
	var list []*Employee
	query := db.GDbMgr.Get().Where(this)
	for i := 0; i < len(condPair); i++ {
		if condPair[i] == nil {
			continue
		}
		query = query.Where(condPair[i].Key, condPair[i].Value)
	}
	query = query.Order("valid desc").Order("id")

	return new(common.Result).PageQuery(query, &Employee{}, &list, page, size, nil, "")
}
