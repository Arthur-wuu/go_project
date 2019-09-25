package models

import (
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/db"
	"github.com/shopspring/decimal"
	"time"
)

type StaffMotivation struct {
	Id             *int             `json:"id,omitempty"        	gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"`
	UserId         *int             `json:"user_id,omitempty"       gorm:"column:user_id;type:int(11)"`
	CompanyId      *int             `json:"company_id,omitempty"	gorm:"column:company_id;type:int(11)"`
	DepartmentId   *int             `json:"department_id,omitempty" gorm:"column:department_id;type:int(11)"`
	CompanyName    *string          `json:"company_name,omitempty"   gorm:"column:company_name;type:varchar(50)"`       //公司名称
	DepartmentName *string          `json:"department_name,omitempty"   gorm:"column:department_name;type:varchar(50)"` //部门名称
	Name           *string          `json:"name,omitempty"			gorm:"column:name;type:varchar(30)"`
	HiredAt        *int64           `json:"hired_at,omitempty"		gorm:"column:hired_at;type:int(11)"`
	BpUid          *int             `json:"bp_uid,omitempty" 		gorm:"column:bp_uid;type:int(11)"`
	Valid          *int             `json:"valid,omitempty"			gorm:"column:valid;type:tinyint(3)"`
	Symbol         *string          `json:"symbol,omitempty"   gorm:"column:symbol;type:varchar(30)"`
	Coin           *decimal.Decimal `json:"coin,omitempty"   gorm:"column:coin;type:decimal(38,12)"`
	TransferFlag   *int             `json:"transfer_flag,omitempty"  gorm:"column:transfer_flag;type:tinyint(3);default:0"` //这个主要是给让第三方转账用的，第三方转账完了可以设置下，避免重复转账
	CreatedAt      *int64           `json:"created_at,omitempty" gorm:"column:created_at;type:int(11)"`
	UpdatedAt      *int64           `json:"updated_at,omitempty" gorm:"column:updated_at;type:int(11)"`
}

type StaffChan struct {
	Coin       decimal.Decimal
	Symbol     string
	MerchantId int
	Times      int
	AwardId    int
	AccountId  int
}

type StaffMotivationAwardResponse struct {
	UserId         int             `json:"user_id,omitempty"`
	DepartmentName string          `json:"department_name"`
	Name           string          `json:"name"`
	Coin           decimal.Decimal `json:"coin,omitempty"`
	HiredAt        int64           `json:"hired_at,omitempty"`
	Years          string          `json:"years,omitempty"`
}

var SfChan chan StaffChan

func (this *StaffMotivation) TableName() string {
	return "staff_motivation_record"
}

func (this *StaffMotivation) AddNew(p *api.StaffListData, coin decimal.Decimal, symbol string) (*StaffMotivation, error) {
	rcd := &StaffMotivation{
		UserId: p.Id,
		BpUid:  p.BpUid,
		Valid:  p.Valid,
		Symbol: &symbol,
		Coin:   &coin,
	}

	if p.CompanyId != nil {
		rcd.CompanyId = p.CompanyId
	}

	if p.HiredAt != nil {
		rcd.HiredAt = p.HiredAt
	}

	if p.Valid != nil {
		rcd.Valid = p.Valid
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

func (this *StaffMotivation) ResendRecord() {
	var awds []*StaffMotivation
	//get half hour ago unsend successful and attempt less ten times records
	err := db.GDbMgr.Get().
		Where("transfer_flag = ? and created_at <= ? and created_at >= ?", 0, time.Now().Unix()-1800, time.Now().Unix()-12000).
		Limit(100).
		Order("created_at asc").
		Find(&awds).Error

	if err == nil && awds != nil {
		for _, awd := range awds {
			SfChan <- StaffChan{
				Coin:       *awd.Coin,
				Symbol:     *awd.Symbol,
				MerchantId: config.GConfig.Company.ServiceAward.MerchantId,
				Times:      1,
				AwardId:    *awd.Id,
				AccountId:  *awd.BpUid,
			}
		}
	}
}

func (this *StaffMotivation) SetTransferSuccess(id int) error {
	err := db.GDbMgr.Get().Model(StaffMotivation{}).Where("id = ?", id).Update("transfer_flag", 1).Error

	if err != nil {
		return err
	}

	return nil
}

func (this *StaffMotivation) GetDayAwardList(btime, page, size int64) (*common.Result, error) {
	var list []*StaffMotivationAwardResponse
	query := db.GDbMgr.Get().Select("user_id,name,department_name,hired_at,coin").Where("created_at between ? and ?", btime, btime+86400).Order("created_at desc")

	return new(common.Result).PageQueryScan(query, &StaffMotivation{}, &list, page, size, nil, "")
}

func (this *StaffMotivation) GetTotalAwardList(btime, page, size int64) (*common.Result, error) {
	var list []*StaffMotivationAwardResponse
	query := db.GDbMgr.Get().Select("user_id,any_value(department_name) as department_name ,any_value(name) as name,any_value(hired_at) as hired_at,sum(coin) as coin").Where("created_at < ?", btime+86400).Group("user_id").Order("user_id asc")

	return new(common.Result).PageQueryScan(query, &StaffMotivation{}, &list, page, size, nil, "")
}
