package go_sdk

type (
	//转账参数
	TransferParam struct{
		Amount           *string   `valid:"required" json:"amount,omitempty"`
		Assets           *string   `valid:"required" json:"assets,omitempty"`
		MerchantId       *string   `valid:"required" json:"merchant_id,omitempty"`
		MerchantTransNo   string   `valid:"optional" json:"-"`
		NotifyUrl        *string   `valid:"optional" json:"notify_url,omitempty"`
		PayeeId          *string   `valid:"required" json:"payee_id,omitempty"`
		ProductName      *string   `valid:"required" json:"product_name,omitempty"`
		Timestamp        *string   `valid:"optional" json:"timestamp,omitempty"`
		SignType         *string   `valid:"optional" json:"sign_type,omitempty"`
		Signature        *string   `valid:"optional" json:"signature,omitempty"`
	}

	TransRes struct{
		Code              int64    `json:"code"`
		Message           string   `json:"message"`
		Data              DataRes  `json:"data"`
	}

	DataRes struct {
		Assets             string  `json:"assets,omitempty"`
		Amount             string  `json:"amount,omitempty"`
		MerchantTransferNo string  `json:"merchant_transfer_no,omitempty"`
		Status			   int     `json:"status,omitempty"`
		TransferNo		   string  `json:"transfer_no,omitempty"`
	}




	//有效币种参数
	AvailAssetsParam struct{
		Assets           *string   `valid:"optional" json:"assets,omitempty"`
	}

	AvailAssetsRes struct{
		Code              int64    `json:"code"`
		Message           string   `json:"message"`
		Data              Assets   `json:"data"`
	}

	Assets struct{
		Array              []Array   `json:"assets"`
	}

	Array struct {
		Assets             string  `json:"assets,omitempty"`
		FullName           string  `json:"full_name,omitempty"`
		Logo               string  `json:"logo,omitempty"`
	}



	//创建二维码订单
	QrTrade struct{
		MerchantTradeNo  string   `valid:"optional" json:"-"`
		PayeeId         *string   `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string   `valid:"required" json:"assets,omitempty"`
		Amount          *string   `valid:"optional" json:"amount,omitempty"`

		ProductName     *string   `valid:"optional" json:"product_name,omitempty"`
		ProductDetail   *string   `valid:"optional" json:"product_detail,omitempty"`
		ExpireTime      int64     `valid:"optional" json:"-"`
		Remark          *string   `valid:"optional" json:"remark,omitempty"`
		MerchantId      *string   `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl       *string   `valid:"optional" json:"notify_url,omitempty"`
		ShowUrl         *string   `valid:"optional" json:"show_url,omitempty"`
		ReturnUrl       *string   `valid:"optional" json:"return_url,omitempty"`
		Timestamp       *string   `valid:"optional" json:"timestamp,omitempty"`
	}

	TradeQrRes struct{
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}






	//创建Sdk订单
	SdkTrade struct{
		MerchantTradeNo  string   `valid:"optional" json:"-"`
		PayeeId         *string   `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string   `valid:"required" json:"assets,omitempty"`
		Amount          *string   `valid:"optional" json:"amount,omitempty"`

		ProductName     *string   `valid:"optional" json:"product_name,omitempty"`
		ProductDetail   *string   `valid:"optional" json:"product_detail,omitempty"`
		ExpireTime      int64     `valid:"optional" json:"-"`
		Remark          *string   `valid:"optional" json:"remark,omitempty"`
		MerchantId      *string   `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl       *string   `valid:"optional" json:"notify_url,omitempty"`
		ShowUrl         *string   `valid:"optional" json:"show_url,omitempty"`
		ReturnUrl       *string   `valid:"optional" json:"return_url,omitempty"`
		Timestamp       *string   `valid:"optional" json:"timestamp,omitempty"`
	}

	SdkRes struct{
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}




	//创建Wap订单
	WapTrade struct{
		MerchantTradeNo  string   `valid:"optional" json:"-"`
		PayeeId         *string   `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string   `valid:"required" json:"assets,omitempty"`
		Amount          *string   `valid:"optional" json:"amount,omitempty"`

		ProductName     *string   `valid:"optional" json:"product_name,omitempty"`
		ProductDetail   *string   `valid:"optional" json:"product_detail,omitempty"`
		ExpireTime      int64     `valid:"optional" json:"-"`
		Remark          *string   `valid:"optional" json:"remark,omitempty"`
		MerchantId      *string   `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl       *string   `valid:"optional" json:"notify_url,omitempty"`
		ShowUrl         *string   `valid:"optional" json:"show_url,omitempty"`
		ReturnUrl       *string   `valid:"optional" json:"return_url,omitempty"`
		Timestamp       *string   `valid:"optional" json:"timestamp,omitempty"`
	}

	WapRes struct{
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}





	//查询交易信息
	TradeInfo struct{
		MerchantTradeNo *string   `valid:"required" json:"merchant_trade_no,omitempty"`
		MerchantId      *string   `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl       *string   `valid:"optional" json:"notify_url,omitempty"`
	}

	InfoRes struct{
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}


    // pos
	PosTrade struct{
		MerchantPosNo    string   `valid:"optional" json:"-"`
		PayeeId         *string   `valid:"optional" json:"payee_id,omitempty"`
		Assets          *string   `valid:"required" json:"assets,omitempty"`
		Amount          *string   `valid:"optional" json:"amount,omitempty"`
		PosMachineId 	*string   `valid:"required" json:"pos_machine_id,omitempty"`
		PayVoucher      *string   `valid:"required" json:"pay_voucher,omitempty"`
		ProductName     *string   `valid:"optional" json:"product_name,omitempty"`
		ProductDetail   *string   `valid:"optional" json:"product_detail,omitempty"`
		Remark          *string   `valid:"optional" json:"remark,omitempty"`
		MerchantId      *string   `valid:"required" json:"merchant_id,omitempty"`
		NotifyUrl       *string   `valid:"optional" json:"notify_url,omitempty"`
		Timestamp       *string   `valid:"optional" json:"timestamp,omitempty"`
	}

	PosRes struct{
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}





	// pos订单列表接口
	PosOrderList struct{
		TimeStamp       *string   `valid:"optional" json:"timestamp,omitempty"`
		BeginTime 		*string   `valid:"optional" json:"begin_time,omitempty"`
		EndTime         *string   `valid:"optional" json:"end_time,omitempty"`
		MerchantId      *string   `valid:"required" json:"merchant_id,omitempty"`
		MerchantPosNo   *string   `valid:"optional" json:"merchant_pos_no"`
		NotifyUrl       *string   `valid:"optional" json:"notify_url,omitempty"`
		Page            *string   `valid:"optional" json:"page,omitempty"`
		PageSize        *string   `valid:"optional" json:"page_size,omitempty"`
		PosMachineId 	*string   `valid:"required" json:"pos_machine_id,omitempty"`
		TradeStatus 	*string   `valid:"optional" json:"trade_status,omitempty"`
	}

	PosListRes struct{
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

)

