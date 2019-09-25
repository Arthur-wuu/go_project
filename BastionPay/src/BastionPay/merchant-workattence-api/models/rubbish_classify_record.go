package models

import (
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/db"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type (
	RubbishClassifyRecord struct {
		Id             *int             `json:"id,omitempty"                        gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"`
		TotalNumbers   *int             `json:"total_numbers,omitempty"             gorm:"column:total_numbers;type:mediumint(8)"`
		TotalCoin      *decimal.Decimal `json:"total_coin,omitempty"     gorm:"column:total_coin;type:decimal(38,12)"`
		ScoreDate      *string          `json:"score_date,omitempty"            gorm:"column:score_date;type:varchar(10)"`
		DepartmentInfo *string          `json:"department_info,omitempty"         gorm:"column:department_info;type:varchar(1500)"`      //部门信息
		TransferFlag   *int             `json:"transfer_flag,omitempty"           gorm:"column:transfer_flag;type:tinyint(3);default:0"` //发送标志, 0(未发放) | 1(发放中) | 2(已发放)
		CreatedAt      *int64           `json:"created_at,omitempty"              gorm:"column:created_at;type:int(11)"`
		UpdatedAt      *int64           `json:"updated_at,omitempty"              gorm:"column:updated_at;type:int(11)"`
	}
)

func (this *RubbishClassifyRecord) TableName() string {
	return "rubbish_classify_record"
}

func (this *RubbishClassifyRecord) Add(p *api.ResponseDepartmenInfo) (*RubbishClassifyRecord, error) {
	tInfo := p.Department
	dInfo, err := json.Marshal(tInfo)
	sDinfo := string(dInfo)

	if err != nil {
		return nil, err
	}

	coinConf := config.GConfig.Company.RubbishClassify.Coin
	var tCoin decimal.Decimal

	for _, v := range tInfo {
		if v.Score != 0 {
			nums := decimal.New(v.Numbers, 0)
			tCoin = tCoin.Add(coinConf[v.Score].Mul(nums))
		}
	}

	rcd := &RubbishClassifyRecord{
		TotalNumbers:   &p.TotalNumbers,
		ScoreDate:      &p.Date,
		TotalCoin:      &tCoin,
		DepartmentInfo: &sDinfo,
	}

	err = db.GDbMgr.Get().Create(rcd).Error

	if err != nil {
		return nil, err
	}

	return rcd, nil
}

func (this *RubbishClassifyRecord) Update(p *api.ResponseDepartmenInfo) (*RubbishClassifyRecord, error) {
	tInfo := p.Department
	dInfo, err := json.Marshal(tInfo)
	sDinfo := string(dInfo)

	if err != nil {
		return nil, err
	}

	coinConf := config.GConfig.Company.RubbishClassify.Coin
	var tCoin decimal.Decimal

	for _, v := range tInfo {
		if v.Score != 0 {
			nums := decimal.New(v.Numbers, 0)
			tCoin = tCoin.Add(coinConf[v.Score].Mul(nums))
		}
	}

	rcd := &RubbishClassifyRecord{
		TotalNumbers:   &p.TotalNumbers,
		ScoreDate:      &p.Date,
		TotalCoin:      &tCoin,
		DepartmentInfo: &sDinfo,
	}

	err = db.GDbMgr.Get().Model(RubbishClassifyRecord{}).Where("id=?", p.Id).Update(rcd).Error

	if err != nil {
		return nil, err
	}

	act := new(RubbishClassifyRecord)
	err = db.GDbMgr.Get().Model(act).Where("id = ?", p.Id).Last(act).Error

	if err != nil {
		return nil, err
	}

	return act, nil
}

func (this *RubbishClassifyRecord) BkParseList(p *api.BkRubbishClassifyList) *RubbishClassifyRecord {
	rbr := &RubbishClassifyRecord{
		Id:           p.Id,
		TransferFlag: p.TransferFlag,
		TotalCoin:    p.TotalCoin,
		TotalNumbers: p.TotalNumbers,
		ScoreDate:    p.ScoreDate,
	}

	return rbr
}

func (this *RubbishClassifyRecord) List(page, size int64) (*common.Result, error) {
	var list []*RubbishClassifyRecord
	query := db.GDbMgr.Get().Select("id,total_numbers,total_coin,score_date,transfer_flag,created_at").Where(this).Order("created_at desc")

	return new(common.Result).PageQuery(query, &RubbishClassifyRecord{}, &list, page, size, nil, "")
}

func (this *RubbishClassifyRecord) GetById(id int) (*RubbishClassifyRecord, error) {
	rcd := new(RubbishClassifyRecord)
	err := db.GDbMgr.Get().Where("id = ?", id).Find(&rcd).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return rcd, err
}

func (this *RubbishClassifyRecord) UpdateTransferFlag(id, flag int) error {
	err := db.GDbMgr.Get().Model(RubbishClassifyRecord{}).Where("id = ?", id).Update("transfer_flag = ?", flag).Error

	if err != nil {
		return err
	}

	return nil
}
