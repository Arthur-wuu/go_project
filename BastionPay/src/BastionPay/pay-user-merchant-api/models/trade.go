package models

import (
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/db"
	"github.com/jinzhu/gorm"
)

type (
	Trade struct {
		Id              *int    `json:"id,omitempty"                       gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		MerchantTradeNo *string `json:"merchant_trade_no,omitempty"        gorm:"column:merchant_trade_no;type:varchar(64);unique"`
		PayeeId         *string `json:"payee_id,omitempty"                 gorm:"column:payee_id;type:varchar(64)"`
		Assets          *string `json:"assets,omitempty"                   gorm:"column:assets;type:varchar(16)"`
		Amount          *string `json:"amount,omitempty"                   gorm:"column:amount;type:varchar(20)"`

		Legal          *string `json:"legal,omitempty"                    gorm:"column:legal;type:varchar(16)"`
		LegalNum       *string `json:"legal_num,omitempty"                gorm:"column:legal_num;type:varchar(20)"`
		GameCoin       *string `json:"game_coin,omitempty"                gorm:"column:game_coin;type:varchar(20)"`
		ProductName    *string `json:"product_name,omitempty"             gorm:"column:product_name;type:varchar(128)"`
		ProductDetail  *string `json:"product_detail,omitempty"           gorm:"column:product_detail;type:varchar(1000)"`
		ExpireTime     *int64  `json:"expire_time,omitempty"              gorm:"column:expire_time;type:int(20)"`
		Remark         *string `json:"remark,omitempty"                   gorm:"column:remark;type:varchar(128)"`
		TradeNo        *string `json:"trade_no,omitempty"                 gorm:"column:trade_no;type:varchar(64)"`
		DeviceId       *string `json:"device_id,omitempty"                gorm:"column:device_id;type:varchar(64)"`
		OpenStatus     *int    `json:"open_status,omitempty"              gorm:"column:open_status;type:int(11)"`
		TransferStatus *int    `json:"transfer_status,omitempty"          gorm:"column:transfer_status;type:int(11)"`
		Table
	}
)

func (this *Trade) TableName() string {
	return "trade"
}

func (this *Trade) Add() error {
	err := db.GDbMgr.Get().Create(this).Error
	if err != nil {
		return err
	}
	return nil
}

func (this *Trade) UpdateByTradeNo(MerchantTradeNo, tradeNo string, openStatus int) error {

	if err := db.GDbMgr.Get().Model(&Trade{}).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(Trade{TradeNo: &tradeNo, OpenStatus: &openStatus}).Error; err != nil {
		return err
	}
	//mod := new(table.Trade)
	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
	//if err != nil {
	//	return nil, err
	//}
	return nil
}

func (this *Trade) GetByOriginNoList(originNo []string) ([]*Trade, error) {

	mod := make([]*Trade, 0)

	for i := 0; i < len(originNo); i++ {
		model := &Trade{}
		err := db.GDbMgr.Get().Where("merchant_trade_no = ?", originNo[i]).Last(model).Error
		if err != nil {
			return nil, err
		}
		mod = append(mod, model)
	}

	return mod, nil

}

func (this *Trade) UpdateOpenStatusByTradeNo(MerchantTradeNo string, status int) error {

	if err := db.GDbMgr.Get().Model(&Trade{}).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(Trade{OpenStatus: &status}).Error; err != nil {
		return err
	}
	//mod := new(table.Trade)
	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
	//if err != nil {
	//	return nil, err
	//}
	return nil
}

func (this *Trade) UpdateTransferStatusByTradeNo(tx *gorm.DB, MerchantTradeNo string, status int) error {
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	if err := tx.Model(&Trade{}).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(Trade{TransferStatus: &status}).Error; err != nil {
		return err
	}
	//mod := new(table.Trade)
	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
	//if err != nil {
	//	return  err
	//}
	return nil
}

func (this *Trade) UpdateRowsAffected(MerchantTradeNo string, status int) (bool, error) {
	sql := db.GDbMgr.Get().Exec("UPDATE "+new(Trade).TableName()+" SET transfer_status= ? WHERE merchant_trade_no = ? ", status, MerchantTradeNo)
	//if sql.Error != nil {
	//	return false, sql.Error
	//}
	return sql.RowsAffected > 0, nil
}

func (this *Trade) GetStatus(tx *gorm.DB, MerchantTradeNo string) (*Trade, error) {
	t := new(Trade)
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	if err := tx.Model(&Trade{}).Select("open_status, transfer_status, amount").Where("merchant_trade_no = ?", MerchantTradeNo).Last(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (this *Trade) GetDeviceId(tx *gorm.DB, MerchantTradeNo string) (string, string, error) {
	t := new(Trade)
	if tx == nil {
		tx = db.GDbMgr.Get()
	}
	if err := tx.Model(&Trade{}).Select("device_id, game_coin").Where("merchant_trade_no = ?", MerchantTradeNo).Last(t).Error; err != nil {
		return "", "", err
	}
	return *t.DeviceId, *t.GameCoin, nil
}

func (this *Trade) UpdateTransferStatusWithCondition(MerchantTradeNo string, open_cond, transfer_cond interface{}, status int) (bool, error) {
	sql := db.GDbMgr.Get().Exec("UPDATE "+new(Trade).TableName()+" SET transfer_status= (case when open_status = ? and transfer_status in (?) then ? else transfer_status end) WHERE merchant_trade_no = ? ", open_cond, transfer_cond, status, MerchantTradeNo)
	if sql.Error != nil {
		return false, sql.Error
	}
	return sql.RowsAffected > 0, nil
}

//func (this *Trade) UpdateReFundStatusByTradeNo(MerchantTradeNo string, status int) ( error) {
//
//	if err := db.GDbMgr.Get().Model(&Trade{}).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(Trade{ReFundStatus: &status}).Error; err != nil {
//		return err
//	}
//	//mod := new(table.Trade)
//	//err := db.GDbMgr.Get().Where("merchant_trade_no = ?", MerchantTradeNo).Last(mod).Error
//	//if err != nil {
//	//	return nil, err
//	//}
//	return nil
//}

func (this *Trade) List(page, size int64, MerchantTradeNo, PayeeId, TradeNo *string) (*common.Result, error) {
	var list []*Trade
	condition := &Trade{MerchantTradeNo: MerchantTradeNo, PayeeId: PayeeId, TradeNo: TradeNo}

	query := db.GDbMgr.Get().Where(condition)

	return new(common.Result).PageQuery(query, &Trade{}, &list, page, size, nil, "")
}
