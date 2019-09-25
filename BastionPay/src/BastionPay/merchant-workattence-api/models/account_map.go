package models

import (
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/db"
	"github.com/jinzhu/gorm"
)

type AccountMap struct {
	Id             *int    `json:"id,omitempty"        gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"`           //加上type:int(11)后AUTO_INCREMENT无效
	AccId          *int    `json:"account_id,omitempty"     gorm:"column:account_id;type:int(11);index:idx_account_id"`  //bastionpay account id
	StaffId        *string `json:"user_id,omitempty"     gorm:"column:user_id;type:varchar(50);index:idx_user_id"`       //钉钉用户ID
	MemberId       *int    `json:"member_id,omitempty"     gorm:"column:member_id;type:int(11)"`                         //钉钉用户ID
	CorpId         *string `json:"corp_id,omitempty"          gorm:"column:corp_id;type:varchar(100);index:idx_corp_id"` //企业ID
	Phone          *string `json:"phone,omitempty"      gorm:"column:phone;type:varchar(30)"`                            //员工手机号
	Name           *string `json:"name,omitempty"   gorm:"column:name;type:varchar(30)"`                                 //员工姓名
	DepartmentId   *int    `json:"department_id,omitempty" gorm:"column:department_id;type:int(11);default:0"`           //部门ID
	CompanyId      *int    `json:"company_id,omitempty" gorm:"column:company_id;type:int(11);default:0"`                 //公司ID
	CompanyName    *string `json:"company_name,omitempty"   gorm:"column:company_name;type:varchar(50)"`                 //公司名称
	DepartmentName *string `json:"department_name,omitempty"   gorm:"column:department_name;type:varchar(50)"`           //部门名称
	HiredAt        *int64  `json:"hired_at,omitempty" gorm:"column:hired_at;type:int(11)"`                               //入职时间
	DeleteFlag     *int    `json:"valid,omitempty"   gorm:"column:delete_flag;type:tinyint(3);default:1"`                //删除标志0（删除）| 1（正常）
	CreatedAt      *int64  `json:"created_at,omitempty" gorm:"column:created_at;type:int(11)"`
	UpdatedAt      *int64  `json:"updated_at,omitempty" gorm:"column:updated_at;type:int(11)"`
}

type DepartmentInfo struct {
	DepartmentId   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	Numbers        int    `json:"numbers,omitempty"`
	Score          int    `json:"score,omitempty"`
}

var CorpIdMap map[int]string

func (this *AccountMap) TableName() string {
	return "workattence_account_map"
}

func (this *AccountMap) BkParseAdd(p *api.BkAccountMapAdd) (*AccountMap, error) {
	//check had add
	act, err := this.GetByAccIdAndStaffId(*p.AccId, *p.StaffId)

	if act == nil && err == nil {
		atm := &AccountMap{
			AccId:   p.AccId,
			StaffId: p.StaffId,
			CorpId:  p.CorpId,
		}

		if p.Name != nil {
			atm.Name = p.Name
		}

		if p.Phone != nil {
			atm.Phone = p.Phone
		}

		err := db.GDbMgr.Get().Create(atm).Error
		if err != nil {
			return nil, err
		}

		actm := new(AccountMap)
		err = db.GDbMgr.Get().Where("id = ?", *atm.Id).Last(actm).Error

		if err != nil {
			return nil, err
		}

		return actm, nil
	}

	return act, nil
}

func (this *AccountMap) BkParseUpdate(p *api.BkAccountMapUpdate) (*AccountMap, error) {
	atm := &AccountMap{
		Id:      p.Id,
		AccId:   p.AccId,
		StaffId: p.StaffId,
		CorpId:  p.CorpId,
	}

	if p.Name != nil {
		atm.Name = p.Name
	}

	if p.Phone != nil {
		atm.Phone = p.Phone
	}

	err := db.GDbMgr.Get().Model(this).Where("id = ?", atm.Id).Update(atm).Error
	if err != nil {
		return nil, err
	}

	act := new(AccountMap)
	err = db.GDbMgr.Get().Model(act).Where("id = ?", atm.Id).Last(act).Error
	if err != nil {
		return nil, err
	}

	return act, nil
}

func (this *AccountMap) BkParseList(p *api.BkAccountMapList) *AccountMap {
	acty := &AccountMap{
		Id:         p.Id,
		Name:       p.Name,
		AccId:      p.AccId,
		Phone:      p.Phone,
		DeleteFlag: p.DeleteFlag,
		StaffId:    p.StaffId,
		CorpId:     p.CorpId,
		CreatedAt:  p.CreatedAt,
	}

	return acty
}

func (this *AccountMap) List(page, size int64) (*common.Result, error) {
	var list []*AccountMap
	query := db.GDbMgr.Get().Where(this).Order("created_at desc")

	return new(common.Result).PageQuery(query, &AccountMap{}, &list, page, size, nil, "")
}

func (this *AccountMap) GetById(id int) (*AccountMap, error) {
	acty := new(AccountMap)
	err := db.GDbMgr.Get().Where("id = ? ", id).Last(acty).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return acty, err
}

func (this *AccountMap) GetByAccIdAndStaffId(AccId int, StaffId string) (*AccountMap, error) {
	act := new(AccountMap)
	err := db.GDbMgr.Get().Where("account_id = ? and user_id = ?", AccId, StaffId).Find(&act).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return act, err
}

func (this *AccountMap) GetByStaffId(StaffId string) (*AccountMap, error) {
	act := new(AccountMap)
	err := db.GDbMgr.Get().Where("user_id = ?", StaffId).Find(&act).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return act, err
}

func (this *AccountMap) GetAllAccounts(corpId string) []*AccountMap {
	var acts []*AccountMap
	err := db.GDbMgr.Get().Where("delete_flag = ? and corp_id=?", 1, corpId).Find(&acts).Error

	if err != nil {
		return nil
	}

	return acts
}

func (this *AccountMap) AddOrUpdate(p *api.StaffListData) error {
	acr, err := this.GetByAccIdAndStaffId(*p.BpUid, *p.EeNo)

	if err != nil {
		return err
	}

	atm := &AccountMap{
		AccId:    p.BpUid,
		StaffId:  p.EeNo,
		MemberId: p.Id,
	}

	if p.CompanyId != nil {
		corpId := CorpIdMap[*p.CompanyId]
		atm.CompanyId = p.CompanyId
		atm.CorpId = &corpId
	}

	if p.HiredAt != nil {
		atm.HiredAt = p.HiredAt
	}

	if p.Valid != nil {
		atm.DeleteFlag = p.Valid
	}

	if p.Name != nil {
		atm.Name = p.Name
	}

	if p.DepartmentId != nil {
		atm.DepartmentId = p.DepartmentId
	}

	if p.CompanyName != nil {
		atm.CompanyName = p.CompanyName
	}

	if p.DepartmentName != nil {
		atm.DepartmentName = p.DepartmentName
	}

	if acr != nil {
		err := db.GDbMgr.Get().Model(AccountMap{}).Where("id = ?", acr.Id).Update(atm).Error

		if err != nil {
			return err
		}
	} else {
		err := db.GDbMgr.Get().Create(atm).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func (this *AccountMap) GetDepartmentInfo() ([]*DepartmentInfo, error) {
	var depart []*DepartmentInfo
	err := db.GDbMgr.Get().
		Model(AccountMap{}).
		Select("department_id,any_value(department_name) as department_name,count(1) as numbers").
		Where("delete_flag=?", 1).
		Group("department_id").Scan(&depart).Error

	if err != nil {
		return nil, err
	}

	return depart, nil
}

func (this *AccountMap) GetValidEmployerCount() (int, error) {
	count := 0
	err := db.GDbMgr.Get().Model(this).Select("id").Where("delete_flag=?", 1).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, err
}

func (this *AccountMap) GetAccountsByDepartId(departId int) []*AccountMap {
	var acts []*AccountMap
	err := db.GDbMgr.Get().Where("delete_flag = ? and department_id = ?", 1, departId).Find(&acts).Error

	if err != nil {
		return nil
	}

	return acts
}
