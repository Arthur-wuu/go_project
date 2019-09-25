package baspay

import "BastionPay/merchant-api/api"

type(
	FundOut struct{
		MerchantTradeNo             *string   `valid:"required" json:"merchant_trade_no,omitempty"`
		Assets                      *string   `valid:"required" json:"assets,omitempty"`
		Amount                      *string   `valid:"required" json:"amount,omitempty"`
		Address                     *string   `valid:"required" json:"address,omitempty"`
		Memo                        *string   `valid:"required" json:"memo,omitempty"`
		Remark                      *string   `valid:"required" json:"remark,omitempty"`
		Request
	}

	ResFundOut struct{
		MerchantFundoutNo          *string   `valid:"required" json:"merchant_fundout_no,omitempty"`
		FundoutNo                  *string   `valid:"required" json:"fundout_no,omitempty"`
		Response
	}
)

func (this * FundOut) Parse(f *api.FundOut) *FundOut {
	return &FundOut{
		MerchantTradeNo : &f.MerchantTradeNo,
		Assets: f.Assets,
		Amount: f.Amount,
		Address: f.Address,
		Memo: f.Memo,
		Remark: f.Remark,
	}
}

func (this * FundOut) Send() (*ResFundOut, error){
	return   &ResFundOut{

	},nil
}
