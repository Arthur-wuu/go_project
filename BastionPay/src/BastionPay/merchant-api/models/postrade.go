package models

import (
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/common"
	"BastionPay/merchant-api/db"
	"github.com/jinzhu/gorm"
)

type (
	PosTrade struct {
		Id           *int    `json:"id,omitempty"                   gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		PosMachineId *string `json:"pos_machine_id,omitempty"       gorm:"column:pos_machine_id;type:varchar(64)"`
		PayeeId      *string `json:"payee_id,omitempty"             gorm:"column:payee_id;type:varchar(64)"`
		Assets       *string `json:"assets,omitempty"               gorm:"column:assets;type:varchar(16)"`
		Amount       *string `json:"amount,omitempty"               gorm:"column:amount;type:varchar(20)"`
		Legal        *string `json:"legal,omitempty"                gorm:"column:legal;type:varchar(16)"`
		//LegalNum           *string   `json:"legal_num,omitempty"                gorm:"column:legal_num;type:varchar(20)"`
		//GameCoin           *string   `json:"game_coin,omitempty"                gorm:"column:game_coin;type:varchar(20)"`
		ProductName   *string `json:"product_name,omitempty"         gorm:"column:product_name;type:varchar(128)"`
		ProductDetail *string `json:"product_detail,omitempty"       gorm:"column:product_detail;type:varchar(1000)"`
		//ExpireTime         *int64    `json:"expire_time,omitempty"              gorm:"column:expire_time;type:int(20)"`
		Remark        *string `json:"remark,omitempty"               gorm:"column:remark;type:varchar(128)"`
		MerchantId    *string `json:"merchant_id,omitempty"          gorm:"column:merchant_id;type:varchar(64)"`
		MerchantPosNo *string `json:"merchant_pos_no,omitempty"      gorm:"column:merchant_pos_no;type:varchar(64)"`
		NotifyUrl     *string `json:"notify_url,omitempty"           gorm:"column:notify_url;type:varchar(64)"`
		PayVoucher    *string `json:"pay_voucher,omitempty"          gorm:"column:pay_voucher;type:varchar(64)"`
		TimeStamp     *string `json:"timestamp,omitempty"            gorm:"column:timestamp;type:varchar(64)"`
		Status        *int    `json:"status,omitempty"               gorm:"column:status;type:int(2)"`
		Table
	}

	PosList struct {
		Status        *int    `valid:"optional" json:"status"`
		PosMachineId  *string `valid:"optional" json:"pos_machine_id"`
		MerchantPosNo *string `valid:"optional" json:"merchant_pos_no"`
		TimeStart     *string `valid:"optional" json:"time_start"`
		TimeEnd       *string `valid:"optional" json:"time_end"`
		Page          int64   `valid:"required" json:"page"`
		Size          int64   `valid:"optional" json:"size"`
	}
)

type SqlPairCondition struct {
	Key   interface{}
	Value interface{}
}

func (this *PosTrade) TableName() string {
	return "pos_trade"
}

func (this *PosTrade) Parse(p *api.PosTrade) *PosTrade {
	return &PosTrade{
		PosMachineId:  p.PosMachineId,
		PayeeId:       p.PayeeId,
		Assets:        p.Assets,
		Amount:        p.Amount,
		ProductName:   p.ProductName,
		ProductDetail: p.ProductDetail,
		TimeStamp:     p.TimeStamp,
		Remark:        p.Remark,
		NotifyUrl:     p.NotifyUrl,
		PayVoucher:    p.PayVoucher,
		Legal:         p.Legal,
		MerchantId:    p.MerchantId,
		MerchantPosNo: &p.MerchantPosNo,
	}
}

func (this *PosTrade) Add() error {
	err := db.GDbMgr.Get().Create(this).Error
	if err != nil {
		return err
	}
	return nil
}

func (this *PosList) List(page, size int64, status int, merchantPosNo, timeStart, timeEnd, posMachineId *string) (*common.Result, error) {
	var list []*PosTrade
	condition := &PosTrade{PosMachineId: posMachineId}

	query := db.GDbMgr.Get().Where(condition)

	return new(common.Result).PageQuery(query, &PosTrade{}, &list, page, size, nil, "")
}

func (this *PosList) ListWithConds(page, size int64, start, end *string) (*common.Result, error) {
	var list []*PosTrade
	condition := &PosTrade{PosMachineId: this.PosMachineId, MerchantPosNo: this.MerchantPosNo, Status: this.Status}
	query := db.GDbMgr.Get().Where(condition)

	if start != nil && end != nil {
		//转化一下时间
		condPair := []*SqlPairCondition{
			{"created_at >= ?", *start},
			{"created_at <= ?", *end},
		}

		for i := 0; i < len(condPair); i++ {
			if condPair[i] == nil {
				continue
			}
			query = query.Where(condPair[i].Key, condPair[i].Value)
		}
	}
	//query = query.Order("valid desc").Order("off_at desc").Select(needFields)

	return new(common.Result).PageQuery(query, &PosTrade{}, &list, page, size, nil, "")
}

func (this *PosTrade) UpdateByPosTradeNo(MerchantPosNo string, status int) error {

	if err := db.GDbMgr.Get().Model(&PosTrade{}).Where("merchant_pos_no = ?", MerchantPosNo).Updates(PosTrade{Status: &status}).Error; err != nil {
		return err
	}
	//mod := new(table.Trade)
	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
	//if err != nil {
	//	return nil, err
	//}
	return nil
}

//func (this *PosTrade) UpdateByTradeNo(MerchantTradeNo, tradeNo string, openStatus int) (error) {
//
//	if err := db.GDbMgr.Get().Model(&PosTrade{}).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(PosTrade{TradeNo: &tradeNo, OpenStatus: &openStatus}).Error; err != nil {
//		return err
//	}
//	//mod := new(table.Trade)
//	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
//	//if err != nil {
//	//	return nil, err
//	//}
//	return nil
//}

//func (this *PosTrade) UpdateOpenStatusByTradeNo(MerchantTradeNo string, status int) (error) {
//
//	if err := db.GDbMgr.Get().Model(&Trade{}).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(Trade{OpenStatus: &status}).Error; err != nil {
//		return err
//	}
//	//mod := new(table.Trade)
//	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
//	//if err != nil {
//	//	return nil, err
//	//}
//	return nil
//}

//func (this *PosTrade) UpdateTransferStatusByTradeNo(tx *gorm.DB, MerchantTradeNo string, status int) error {
//	if tx == nil {
//		tx = db.GDbMgr.Get()
//	}
//	if err := tx.Model(&Trade{}).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(Trade{TransferStatus: &status}).Error; err != nil {
//		return err
//	}
//	//mod := new(table.Trade)
//	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
//	//if err != nil {
//	//	return  err
//	//}
//	return nil
//}

func (this *PosTrade) UpdateRowsAffected(MerchantTradeNo string, status int) (bool, error) {
	sql := db.GDbMgr.Get().Exec("UPDATE "+new(Trade).TableName()+" SET transfer_status= ? WHERE merchant_trade_no = ? ", status, MerchantTradeNo)
	//if sql.Error != nil {
	//	return false, sql.Error
	//}
	return sql.RowsAffected > 0, nil
}

func (this *PosTrade) GetStatus(tx *gorm.DB, MerchantTradeNo string) (*Trade, error) {
	t := new(Trade)
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	if err := tx.Model(&Trade{}).Select("open_status, transfer_status, amount").Where("merchant_trade_no = ?", MerchantTradeNo).Last(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (this *PosTrade) GetDeviceId(tx *gorm.DB, MerchantTradeNo string) (string, string, error) {
	t := new(Trade)
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	if err := tx.Model(&Trade{}).Select("device_id, game_coin").Where("merchant_trade_no = ?", MerchantTradeNo).Last(t).Error; err != nil {
		return "", "", err
	}
	return *t.DeviceId, *t.GameCoin, nil
}

func (this *PosTrade) UpdateTransferStatusWithCondition(MerchantTradeNo string, open_cond, transfer_cond interface{}, status int) (bool, error) {
	sql := db.GDbMgr.Get().Exec("UPDATE "+new(Trade).TableName()+" SET transfer_status= (case when open_status = ? and transfer_status in (?) then ? else transfer_status end) WHERE merchant_trade_no = ? ", open_cond, transfer_cond, status, MerchantTradeNo)
	if sql.Error != nil {
		return false, sql.Error
	}
	return sql.RowsAffected > 0, nil
}
