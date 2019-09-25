package baspay

import (
	"BastionPay/merchant-api/api"
)

type (
	Transfer struct {
		MerchantTradeNo *string `valid:"required" json:"merchant_trade_no,omitempty"`
		PayeeId         *string `valid:"required" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"required" json:"amount,omitempty"`
		Remark          *string `valid:"required" json:"remark,omitempty"`
		Request
	}

	ResTransfer struct {
		MerchantTransferNo *string `valid:"required" json:"merchant_transfer_no,omitempty"`
		TransferNo         *string `valid:"required" json:"transfer_no,omitempty"`
		Response
	}
)

func (this *Transfer) Parse(f *api.Pay) *Transfer {
	return &Transfer{
		MerchantTradeNo: f.MerchantTradeNo,
		PayeeId:         f.PayeeId,
		Assets:          f.Assets,
		Amount:          f.Amount,
		Remark:          f.Remark,
	}
}

func (this *Transfer) Send() (*ResTransfer, error) {
	//	s := this.Amount
	//往baspay发送 付款

	//base.HttpSend("/wallet/api/trade/qr_code_pay", bytes.NewBuffer(),"POST",nil)

	return &ResTransfer{}, nil
}
