package api

type (
	//创建订单, 查询
	Trade struct {
		MerchantTradeNo string  `valid:"optional" json:"-"`
		GameCoin        string  `valid:"optional" json:"game_coin"`
		PayeeId         *string `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"required" json:"amount,omitempty"`
		ProductName     *string `valid:"required" json:"product_name,omitempty"`
		ProductDetail   *string `valid:"required" json:"product_detail,omitempty"`
		ExpireTime      int64   `valid:"optional" json:"-"`
		Remark          *string `valid:"optional" json:"remark,omitempty"`
		DeviceId        *string `valid:"optional" json:"device_id,omitempty"`
		ReturnUrl       *string `valid:"optional" json:"return_url,omitempty"`
		ShowUrl         *string `valid:"optional" json:"show_url,omitempty"`
		MerchantId      *string `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl       *string `valid:"optional" json:"notify_url,omitempty"`
	}

	Coffee struct {
		Ver       *string `valid:"optional" json:"ver,omitempty"`
		Orderid   *string `valid:"optional" json:"orderid,omitempty"`
		Machid    *string `valid:"optional" json:"machid,omitempty"`
		Trackno   *string `valid:"optional" json:"trackno,omitempty"`
		Name      *string `valid:"optional" json:"name,omitempty"`
		Price     *int64  `valid:"required" json:"price,omitempty"`
		Channelid *int64  `valid:"required" json:"channelid,omitempty"`
		Randstr   *string `valid:"optional" json:"randstr,omitempty"`
		Timestamp *string `valid:"optional" json:"timestamp,omitempty"`
		Sign      *string `valid:"optional" json:"sign,omitempty"`
	}

	QrTrade struct {
		MerchantTradeNo string  `valid:"optional" json:"-"`
		PayeeId         *string `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"optional" json:"amount,omitempty"`
		Legal           *string `valid:"optional" json:"legal,omitempty"`
		LegalNum        *string `valid:"optional" json:"legal_num,omitempty"`

		ProductName   *string `valid:"optional" json:"product_name,omitempty"`
		ProductDetail *string `valid:"optional" json:"product_detail,omitempty"`
		ExpireTime    int64   `valid:"optional" json:"-"`
		Remark        *string `valid:"optional" json:"remark,omitempty"`
		MerchantId    *string `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl     *string `valid:"optional" json:"notify_url,omitempty"`
	}

	CoffeeQrTrade struct {
		MerchantTradeNo string  `valid:"optional" json:"-"`
		PayeeId         *string `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"optional" json:"amount,omitempty"`
		ProductName     *string `valid:"optional" json:"product_name,omitempty"`
		ProductDetail   *string `valid:"optional" json:"product_detail,omitempty"`
		ExpireTime      int64   `valid:"optional" json:"-"`
		Remark          *string `valid:"optional" json:"remark,omitempty"`
		MerchantId      *string `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl       *string `valid:"optional" json:"notify_url,omitempty"`
	}
	// 创建pos订单
	PosTrade struct {
		PosMachineId  *string `valid:"required" json:"pos_machine_id,omitempty"`
		Legal         *string `valid:"required" json:"legal,omitempty"`
		Assets        *string `valid:"required" json:"assets,omitempty"`
		Amount        *string `valid:"required" json:"amount,omitempty"`
		MerchantId    *string `valid:"required" json:"merchant_id,omitempty"`
		MerchantPosNo string  `valid:"optional" json:"merchant_pos_no"`
		NotifyUrl     *string `valid:"optional" json:"notify_url,omitempty"`
		PayVoucher    *string `valid:"required" json:"pay_voucher,omitempty"`
		PayeeId       *string `valid:"required" json:"payee_id,omitempty"`
		TimeStamp     *string `valid:"optional" json:"timestamp,omitempty"`
		ProductName   *string `valid:"required" json:"product_name,omitempty"`
		ProductDetail *string `valid:"optional" json:"product_detail,omitempty"`
		Remark        *string `valid:"optional" json:"remark,omitempty"`
	}

	// 查询pos订单
	PosOrders struct {
		TimeStamp     *string `valid:"optional" json:"timestamp,omitempty"`
		BeginTime     *string `valid:"optional" json:"begin_time,omitempty"`
		EndTime       *string `valid:"optional" json:"end_time,omitempty"`
		MerchantId    *string `valid:"required" json:"merchant_id,omitempty"`
		MerchantPosNo *string `valid:"optional" json:"merchant_pos_no"`
		NotifyUrl     *string `valid:"optional" json:"notify_url,omitempty"`
		Page          *string `valid:"optional" json:"page,omitempty"`
		PageSize      *string `valid:"optional" json:"page_size,omitempty"`
		PosMachineId  *string `valid:"required" json:"pos_machine_id,omitempty"`
		TradeStatus   *string `valid:"optional" json:"trade_status,omitempty"`
	}

	// 退款订单
	RefundTrade struct {
		MerchantId              *string `valid:"required" json:"merchant_id,omitempty"`
		MerchantRefundNo        *string `valid:"required" json:"merchant_refund_no"`
		NotifyUrl               *string `valid:"optional" json:"notify_url,omitempty"`
		OriginalMerchantTradeNo *string `valid:"required" json:"original_merchant_trade_no,omitempty"`
		Remark                  *string `valid:"optional" json:"remark,omitempty"`
	}
	//供前端查询退单列表
	RefundTradeList struct {
		MerchantId   *string   `valid:"required" json:"merchant_id,omitempty"`
		Page         *int64     `valid:"optional" json:"page,omitempty"`
		Size         *int64     `valid:"optional" json:"size,omitempty"`
	}

	ResTradeAdd struct {
		MerchantTradeNo *string `json:"merchant_trade_no,omitempty"`
		TradeNo         *string `json:"trade_no,omitempty"`
	}

	ResTradeQrAdd struct {
		MerchantTradeNo *string `json:"merchant_trade_no,omitempty"`
		QrCode          *string `json:"qr_code,omitempty"`
	}

	TradeSearch struct {
		PayeeId         *string `valid:"optional" json:"payee_id,omitempty"`
		MerchantTradeNo *string `json:"merchant_trade_no,omitempty"`
		TradeNo         *string `json:"trade_no,omitempty"`
	}

	TradeList struct {
		PayeeId         *string `valid:"required" json:"payee_id,omitempty"`
		MerchantTradeNo *string `valid:"optional" json:"merchant_trade_no,omitempty"`
		TradeNo         *string `valid:"optional" json:"trade_no,omitempty"`
		Page            int64   `valid:"optional" json:"page,omitempty"`
		Size            int64   `valid:"optional" json:"size,omitempty"`
	}

	//退款
	ReFund struct {
		OriginalMerchantTradeNo *string `valid:"optional" json:"original_merchant_trade_no,omitempty"`
		OriginalTradeNo         *string `valid:"optional" json:"original_trade_no,omitempty"`
		MerchantRefundNo        string  `valid:"required" json:"-"`
		Remark                  *string `valid:"optional" json:"remark,omitempty"`
	}

	ResReFund struct {
		MerchantTradeNo *string `valid:"required" json:"merchant_trade_no,omitempty"`
		TradeNo         *string `valid:"required" json:"trade_no,omitempty"`
	}

	//转账 支付
	Pay struct {
		MerchantTradeNo *string `valid:"required" json:"merchant_trade_no,omitempty"`
		PayeeId         *string `valid:"required" json:"payee_id,omitempty"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"required" json:"amount,omitempty"`
		Remark          *string `valid:"optional" json:"remark,omitempty"`
	}

	ResPay struct {
		MerchantTransferNo *string `valid:"required" json:"merchant_transfer_no,omitempty"`
		TransferNo         *string `valid:"required" json:"transfer_no,omitempty"`
	}

	//体现
	FundOut struct {
		MerchantTradeNo string  `valid:"required" json:"-"`
		Assets          *string `valid:"required" json:"assets,omitempty"`
		Amount          *string `valid:"required" json:"amount,omitempty"`
		Address         *string `valid:"required" json:"address,omitempty"`
		Memo            *string `valid:"optional" json:"memo,omitempty"`
		Remark          *string `valid:"optional" json:"remark,omitempty"`
	}

	ResFundOut struct {
		MerchantFundoutNo *string `valid:"required" json:"merchant_fundout_no,omitempty"`
		FundoutNo         *string `valid:"required" json:"fundout_no,omitempty"`
	}

	FundOutList struct {
		PayeeId         *string `valid:"required" json:"payee_id,omitempty"`
		MerchantTradeNo *string `valid:"optional" json:"merchant_trade_no,omitempty"`
		TradeNo         *string `valid:"optional" json:"trade_no,omitempty"`
		Page            int64   `valid:"optional" json:"page,omitempty"`
		Size            int64   `valid:"optional" json:"size,omitempty"`
	}

	TradeInfo struct {
		MerchantTradeNo *string `valid:"required" json:"merchant_trade_no,omitempty"`
		MerchantId      *string `valid:"required" json:"merchant_id,omitempty"`
	}

	AvAssets struct {
		Assets *string `valid:"optional" json:"assets,omitempty"`
		Money  *string `valid:"required" json:"money,omitempty"`
	}

	GetAmount struct {
		Legal  *string `valid:"required" json:"legal,omitempty"`
		Symbol *string `valid:"required" json:"symbol,omitempty"`
		Amount *string `valid:"required" json:"amount,omitempty"`
	}

	Quote struct {
		Symbol *string `valid:"required" json:"symbol,omitempty"`
		Legal  *string `valid:"required" json:"legal,omitempty"`
		//Amount     *string     `valid:"required" json:"amount,omitempty"`
	}
)

const (
	CONST_TRADE_STATUS_DEFAT  = 1
	CONST_TRADE_STATUS_PAY    = 1
	CONST_TRADE_STATUS_REFUND = 2
)

type Merchant struct {
	MerchantTradeNo string `json:"merchant_order_no,omitempty"`
	MerchantId      string `json:"merchant_id,omitempty"`
}
