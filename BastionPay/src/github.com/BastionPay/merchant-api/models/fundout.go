package models

import(
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/db"
	"BastionPay/merchant-api/common"
)

type(
	FundOut struct {
		Id                 *int      `json:"id,omitempty"                       gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"` //加上type:int(11)后AUTO_INCREMENT无效
		MerchantTradeNo    *string   `json:"merchant_trade_no,omitempty"        gorm:"column:merchant_trade_no;type:varchar(64);unique"`
		Assets             *string   `json:"assets,omitempty"                   gorm:"column:assets;type:varchar(16)"`
		Amount             *string   `json:"amount,omitempty"                   gorm:"column:amount;type:varchar(20)"`
		Address            *string   `json:"address,omitempty"                  gorm:"column:address;type:varchar(128)"`
		Memo               *string   `json:"memo,omitempty"                     gorm:"column:memo;type:varchar(128)"`
		Remark             *string   `json:"remark,omitempty"                   gorm:"column:remark;type:varchar(128)"`
		TradeNo            *string   `json:"trade_no,omitempty"                 gorm:"column:trade_no;type:varchar(64)"`
		Status             *int      `json:"status,omitempty"                   gorm:"column:status;type:int(11)"`
		Table
	}
)

func (this *FundOut) TableName() string {
	return "fundout"
}

func (this *FundOut) Parse(p *api.FundOut) *FundOut {
	return &FundOut{
		MerchantTradeNo : &p.MerchantTradeNo,
		Assets: p.Assets,
		Amount: p.Amount,
		Address: p.Address,
		Memo: p.Memo,
		Remark: p.Remark,
	}
}

func (this *FundOut) Add()  (error) {
	err := db.GDbMgr.Get().Create(this).Error
	if err != nil {
		return  err
	}

	return nil
}

func (this *FundOut) UpdateByTradeNo(MerchantTradeNo, tradeNo string, openStatus int) (error) {
	return db.GDbMgr.Get().Model(this).Where("merchant_trade_no = ?", MerchantTradeNo).Updates(Trade{TradeNo: &tradeNo, OpenStatus: &openStatus}).Error
}

func (this *FundOut) List(page, size int64, MerchantTradeNo, PayeeId, TradeNo *string) (*common.Result, error) {
	var list []*FundOut
	condition := &FundOut{MerchantTradeNo:MerchantTradeNo,  TradeNo: TradeNo}

	query := db.GDbMgr.Get().Where(condition)

	return new(common.Result).PageQuery(query, &Trade{}, &list, page, size, nil, "")
}