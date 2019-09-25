package bastionpay

const (
	TransTypeDeposit    = 0
	TransTypeWithdrawal = 1

	StatusBlockIn      = 0
	StatusBlockConfirm = 1
	StatusFailed       = 2
)

type (
	Source struct {
		ApiKey    string `json:"user_key"`
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}

	Request struct {
		Method string      `json:"method"`
		Params interface{} `json:"params"`
	}

	Response struct {
		Method struct {
			Version  string `json:"version"`
			Srv      string `json:"srv"`
			Function string `json:"function"`
		} `json:"method"`
		Err    int    `json:"err"`
		ErrMsg string `json:"errmsg"`
		Value  Source `json:"value"`
	}

	RequestAddress struct {
		Symbol string `json:"asset_name"`
		Count  int    `json:"count"`
	}

	ResponseAddress struct {
		Symbol    string   `json:"asset_name"`
		Addresses []string `json:"data"`
	}

	RequestHistoryMsg struct {
		MinID int64 `json:"min_msg_id"`
		MaxID int64 `json:"max_msg_id"`
	}

	CallbackMsg struct {
		ID            int64   `json:"msg_id"`
		UserOrderID   string  `json:"user_order_id"`
		TransType     int     `json:"trans_type"`
		Status        int     `json:"status"`
		BlockInHeight int64   `json:"blockin_height"`
		Symbol        string  `json:"asset_name"`
		Address       string  `json:"address"`
		Amount        float64 `json:"amount"`
		PayFee        float64 `json:"pay_fee"`
		TxID          string  `json:"hash"`
		OrderID       string  `json:"order_id"`
		Time          int64   `json:"time"`
	}

	RequestWithdrawal struct {
		UserOrderID string  `json:"user_order_id"`
		Symbol      string  `json:"asset_name"`
		Address     string  `json:"address"`
		Amount      float64 `json:"amount"`
	}

	ResponseWithdrawal struct {
		OrderID string `json:"order_id"`
	}

	ResponseCurrency struct {
		Symbol          string  `json:"asset_name"`
		Name            string  `json:"full_name"`
		IsToken         int     `json:"is_token"`
		ParentName      string  `json:"parent_name"`
		DepositMin      float64 `json:"deposit_min"`
		WithdrawalRate  float64 `json:"withdrawal_rate"`
		WithdrawalValue float64 `json:"withdrawal_value"`
		ConfirmationNum int64   `json:"confirmation_num"`
		Decimal         int     `json:"decimal"`
	}

	ResponseBalance struct {
		Symbol          string  `json:"asset_name"`
		AvailableAmount float64 `json:"available_amount"`
		FrozenAmount    float64 `json:"frozen_amount"`
	}

	RequestOrder struct {
		Symbol        string `json:"asset_name"`
		TransType     int    `json:"trans_type"`
		Status        int    `json:"status"`
		MaxUpdateTime int64  `json:"max_update_time"`
		MinUpdateTime int64  `json:"min_update_time"`
	}

	ResponseOrder struct {
		Symbol    string  `json:"asset_name"`
		TransType int     `json:"trans_type"`
		Status    int     `json:"status"`
		Amount    float64 `json:"amount"`
		PayFee    float64 `json:"pay_fee"`
		TxID      string  `json:"hash"`
		OrderID   string  `json:"order_id"`
		Time      int64   `json:"time"`
	}

	ResponseHeight struct {
		Symbol      string `json:"asset_name"`
		BlockHeight int64  `json:"block_height"`
	}
)
