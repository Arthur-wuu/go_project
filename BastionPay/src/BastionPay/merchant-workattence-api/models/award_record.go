package models

import (
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/db"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

type (
	AwardRecord struct {
		Id             *int             `json:"id,omitempty"        gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"`          //加上type:int(11)后AUTO_INCREMENT无效
		CheckinId      *int             `json:"checkin_id,omitempty"     gorm:"column:checkin_id;type:int(11);index:idx_checkin_id"` //考情记录id
		AccId          *int             `json:"account_id,omitempty"     gorm:"column:account_id;type:int(11);index:idx_account_id"` //bastionpay account id
		MemberId       *int             `json:"member_id,omitempty"     gorm:"column:member_id;type:int(11)"`
		DepartmentId   *int             `json:"department_id,omitempty" gorm:"column:department_id;type:int(11);default:0"`     //部门ID
		StaffId        *string          `json:"user_id,omitempty"     gorm:"column:user_id;type:varchar(50);index:idx_user_id"` //钉钉用户ID
		CompanyId      *int             `json:"company_id,omitempty" gorm:"column:company_id;type:int(11);default:0"`           //公司ID
		CompanyName    *string          `json:"company_name,omitempty"   gorm:"column:company_name;type:varchar(50)"`           //公司名称
		DepartmentName *string          `json:"department_name,omitempty"   gorm:"column:department_name;type:varchar(50)"`     //部门名称
		Name           *string          `json:"name,omitempty"   gorm:"column:name;type:varchar(30)"`                           //员工姓名
		Symbol         *string          `json:"symbol,omitempty"   gorm:"column:symbol;type:varchar(30)"`
		Coin           *decimal.Decimal `json:"coin,omitempty"   gorm:"column:coin;type:decimal(38,12)"`
		TransferFlag   *int             `json:"transfer_flag,omitempty"  gorm:"column:transfer_flag;type:tinyint(3);default:0"` //这个主要是给让第三方转账用的，第三方转账完了可以设置下，避免重复转账
		Type           *int             `json:"type,omitempty"  gorm:"column:type;type:tinyint(3);default:1"`                   //奖励类型，1考勤，2加班
		CreatedAt      *int64           `json:"created_at,omitempty" gorm:"column:created_at;type:int(11)"`
		UpdatedAt      *int64           `json:"updated_at,omitempty" gorm:"column:updated_at;type:int(11)"`
	}

	SendChan struct {
		Coin       decimal.Decimal
		Symbol     string
		MerchantId int
		Times      int
		AwardId    int
		AccountId  int
	}

	ResponseOvertimeList struct {
		Firsttime      string          `json:"firsttime"`
		Lasttime       string          `json:"lasttime"`
		CheckinId      int             `json:"checkin_id,omitempty"`
		Name           string          `json:"name"`
		MemberId       int             `json:"member_id"`
		DepartmentName string          `json:"department_name"`
		Symbol         string          `json:"symbol,omitempty"`
		Coin           decimal.Decimal `json:"coin,omitempty"`
		Duration       int64           `json:"duration"`
		UserId         string          `json:"user_id,omitempty"`
	}
)

var SChan chan SendChan

func (this *AwardRecord) TableName() string {
	return "workattence_award_record"
}

func (this *AwardRecord) SendCheckDay(staffId string, bTime, ckt int64) (num int, err error) {
	var count int
	err = db.GDbMgr.Get().Model(AwardRecord{}).Where("user_id = ? and created_at between ? and ? and type = ?", staffId, bTime, ckt, 1).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (this *AwardRecord) AddAuto(ckId int, act AccountMap, staffId string, coin decimal.Decimal, symbol string, atype int) (*AwardRecord, error) {
	awd := &AwardRecord{
		CheckinId:      &ckId,
		AccId:          act.AccId,
		MemberId:       act.MemberId,
		CompanyName:    act.CompanyName,
		CompanyId:      act.CompanyId,
		DepartmentName: act.DepartmentName,
		DepartmentId:   act.DepartmentId,
		Name:           act.Name,
		StaffId:        &staffId,
		Coin:           &coin,
		Symbol:         &symbol,
		Type:           &atype,
	}

	err := db.GDbMgr.Get().Create(awd).Error

	if err != nil {
		return nil, err
	}

	return awd, nil
}

func (this *AwardRecord) ResendRecord() {
	var awds []*AwardRecord
	//get half hour ago unsend successful and attempt less ten times records
	err := db.GDbMgr.Get().
		Where("transfer_flag = ? and created_at <= ? and created_at >= ?", 0, time.Now().Unix()-1800, time.Now().Unix()-12000).
		Limit(100).
		Order("created_at asc").
		Find(&awds).Error

	if err == nil && awds != nil {
		for _, awd := range awds {
			SChan <- SendChan{
				Coin:       *awd.Coin,
				Symbol:     *awd.Symbol,
				MerchantId: config.GConfig.Award.MerchantId,
				Times:      1,
				AwardId:    *awd.Id,
				AccountId:  *awd.AccId,
			}
		}
	}
}

func (this *AwardRecord) SetTransferSuccess(id int) error {
	err := db.GDbMgr.Get().Model(AwardRecord{}).Where("id = ?", id).Update("transfer_flag", 1).Error

	if err != nil {
		return err
	}

	return nil
}

func (this *AwardRecord) GetLatestOvertimeRecordByUserId(userId string) (*AwardRecord, error) {
	awd := new(AwardRecord)
	err := db.GDbMgr.Get().Where("user_id = ? and type = ? and created_at >= ?", userId, 2, common.New().DayBeginTimestamp()).Order("created_at desc").First(&awd).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return awd, nil
}

func (this *AwardRecord) GetOvertimeList(btime, page, size int64) (*common.Result, error) {
	var list []*ResponseOvertimeList
	query := db.GDbMgr.Get().Select("user_id, max(checkin_id) checkin_id,any_value(department_name) as department_name ,any_value(name) as name,any_value(member_id) as member_id,sum(coin) as coin, any_value(created_at) created_at,any_value(symbol) symbol").
		Where("created_at between ? and ? and type = ?", btime, btime+86400, 2).
		Group("user_id").Order("created_at asc")

	return new(common.Result).PageQueryScan(query, &AwardRecord{}, &list, page, size, nil, "")
}
