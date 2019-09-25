package models

import (
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/db"
	"github.com/shopspring/decimal"
	"time"
)

type (
	RubbishClassify struct {
		Id             *int             `json:"id,omitempty"        	gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"`
		RecordId       *int             `json:"record_id,omitempty"     gorm:"column:record_id;type:int(11);default:0"`
		UserId         *int             `json:"user_id,omitempty"       gorm:"column:user_id;type:int(11)"`
		CompanyId      *int             `json:"company_id,omitempty"	gorm:"column:company_id;type:int(11)"`
		DepartmentId   *int             `json:"department_id,omitempty" gorm:"column:department_id;type:int(11)"`
		CompanyName    *string          `json:"company_name,omitempty"   gorm:"column:company_name;type:varchar(50)"`       //公司名称
		DepartmentName *string          `json:"department_name,omitempty"   gorm:"column:department_name;type:varchar(50)"` //部门名称
		Name           *string          `json:"name,omitempty"			gorm:"column:name;type:varchar(30)"`
		BpUid          *int             `json:"bp_uid,omitempty" 		gorm:"column:bp_uid;type:int(11)"`
		Symbol         *string          `json:"symbol,omitempty"   gorm:"column:symbol;type:varchar(30)"`
		Coin           *decimal.Decimal `json:"coin,omitempty"   gorm:"column:coin;type:decimal(38,12)"`
		Score          *int             `json:"score,omitempty"  gorm:"column:score;type:tinyint(3);default:0"`                 // 0 差 1 良 2 优
		TransferFlag   *int             `json:"transfer_flag,omitempty"  gorm:"column:transfer_flag;type:tinyint(3);default:0"` //这个主要是给让第三方转账用的，第三方转账完了可以设置下，避免重复转账
		CreatedAt      *int64           `json:"created_at,omitempty" gorm:"column:created_at;type:int(11)"`
		UpdatedAt      *int64           `json:"updated_at,omitempty" gorm:"column:updated_at;type:int(11)"`
	}

	RClassifyChan struct {
		Coin       decimal.Decimal
		Symbol     string
		MerchantId int
		Times      int
		AwardId    int
		AccountId  int
	}
)

var RcChan chan RClassifyChan

func (this *RubbishClassify) TableName() string {
	return "rubbish_classify_award"
}

func (this *RubbishClassify) AddNew(p *api.StaffListData, coin decimal.Decimal, symbol string, score int) (*RubbishClassify, error) {
	rcd := &RubbishClassify{
		UserId: p.Id,
		BpUid:  p.BpUid,
		Symbol: &symbol,
		Coin:   &coin,
		Score:  &score,
	}

	if p.CompanyId != nil {
		rcd.CompanyId = p.CompanyId
	}

	if p.Name != nil {
		rcd.Name = p.Name
	}

	if p.DepartmentId != nil {
		rcd.DepartmentId = p.DepartmentId
	}

	if p.CompanyName != nil {
		rcd.CompanyName = p.CompanyName
	}

	if p.DepartmentName != nil {
		rcd.DepartmentName = p.DepartmentName
	}

	err := db.GDbMgr.Get().Create(rcd).Error

	if err != nil {
		return nil, err
	}

	return rcd, nil
}

func (this *RubbishClassify) AddNewByAccount(act AccountMap, coin decimal.Decimal, symbol string, score int, rId int) (*RubbishClassify, error) {
	rcd := &RubbishClassify{
		UserId:         act.MemberId,
		BpUid:          act.AccId,
		CompanyName:    act.CompanyName,
		DepartmentId:   act.DepartmentId,
		DepartmentName: act.DepartmentName,
		CompanyId:      act.CompanyId,
		Name:           act.Name,
		Symbol:         &symbol,
		Coin:           &coin,
		Score:          &score,
		RecordId:       &rId,
	}

	err := db.GDbMgr.Get().Create(rcd).Error

	if err != nil {
		return nil, err
	}

	return rcd, nil
}

func (this *RubbishClassify) ResendRecord() {
	var awds []*RubbishClassify
	//get half hour ago unsend successful and attempt less ten times records
	err := db.GDbMgr.Get().
		Where("transfer_flag = ? and created_at <= ? and created_at >= ?", 0, time.Now().Unix()-1800, time.Now().Unix()-12000).
		Limit(100).
		Order("created_at asc").
		Find(&awds).Error

	if err == nil && awds != nil {
		for _, awd := range awds {
			RcChan <- RClassifyChan{
				Coin:       *awd.Coin,
				Symbol:     *awd.Symbol,
				MerchantId: config.GConfig.Company.RubbishClassify.MerchantId,
				Times:      1,
				AwardId:    *awd.Id,
				AccountId:  *awd.BpUid,
			}
		}
	}
}

func (this *RubbishClassify) SetTransferSuccess(id int) error {
	err := db.GDbMgr.Get().Model(RubbishClassify{}).Where("id = ?", id).Update("transfer_flag", 1).Error

	if err != nil {
		return err
	}

	return nil
}

func (this *RubbishClassify) BkParseList(p *api.BkRubbishClassifyAwardList) *RubbishClassify {
	rbr := &RubbishClassify{
		UserId:       p.UserId,
		CompanyId:    p.CompanyId,
		DepartmentId: p.DepartmentId,
		Name:         p.Name,
		BpUid:        p.BpUid,
		Score:        p.Score,
		TransferFlag: p.TransferFlag,
	}

	return rbr
}

func (this *RubbishClassify) List(page, size int64) (*common.Result, error) {
	var list []*RubbishClassify
	query := db.GDbMgr.Get().Where(this).Order("created_at desc")

	return new(common.Result).PageQuery(query, &RubbishClassify{}, &list, page, size, nil, "")
}

func (this *RubbishClassify) AwardListByRecordId(rId int, page, size int64) (*common.Result, error) {
	var list []*RubbishClassify
	query := db.GDbMgr.Get().Where("record_id = ?", rId).Order("created_at desc")

	return new(common.Result).PageQuery(query, &RubbishClassify{}, &list, page, size, nil, "")
}
