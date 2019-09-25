package v1

type (
	ReqQueryParam struct {
		IsAsyn       int    `json:"is_asyn"`
		PageIndex    int    `json:"page_index"`
		MaxDispLines int    `json:"max_disp_lines"`
		TotalLines   int    `json:"total_lines"`
		UserParam    string `json:"user_param,omitempty"`
	}

	AckQueryParam struct {
		IsAsyn       int    `json:"is_asyn"`
		QueryID      string `json:"query_id"`
		PageIndex    int    `json:"page_index"`
		MaxDispLines int    `json:"max_disp_lines"`
		TotalLines   int    `json:"total_lines"`
		UserParam    string `json:"user_param,omitempty"`
	}

	ReqQueryAsynData struct {
		QueryID string `json:"query_id" doc:"查询序号"`
	}

	AckQueryAsynData struct {
		QueryID   string `json:"query_id" doc:"查询序号"`
		Status    int    `json:"status" doc:"0:数据准备中;1:查询数据成功;2:数据查询失败"`
		Data      string `json:"data" doc:"异步返回的数据"`
		UserParam string `json:"user_param,omitempty" doc:"用户自定义字段"`
	}

	// 申请地址
	ReqNewAddress struct {
		AssetName string `json:"asset_name"`
	}

	AckNewAddressList struct {
		AssetName string   `json:"asset_name"`
		Address   string   `json:"address"`
		Memo      string   `json:"memo,omitempty"`
		Data      []string `json:"data"`
	}

	// 地址验证
	ReqIsValidAddress struct {
		AssetName string `json:"asset_name"`
		Address   string `json:"address"`
	}

	RspIsValidAddress struct {
		IsValid int `json:"is_valid"`
	}

	// 提币申请
	ReqWithdrawal struct {
		AssetName   string  `json:"asset_name"`
		Amount      float64 `json:"amount"`
		Address     string  `json:"address"`
		Memo        string  `json:"memo,omitempty"`
		UserOrderID string  `json:"user_order_id"`
		IsTransfer  int     `json:"is_transfer,omitempty"`
		Remark      string  `json:"remark,omitempty"`
	}

	AckWithdrawal struct {
		OrderID     string `json:"order_id"`
		UserOrderID string `json:"user_order_id,omitempty"`
	}

	// 获取支持币种
	ReqSupportAssets struct{}

	AckSupportAssetList struct {
		Data []string `json:"data"`
	}

	// 获取币种属性
	ReqAssetsAttributeList struct {
		AssetNames []string `json:"asset_names,omitempty"`
		IsToken    int      `json:"is_token,omitempty"`
		HasMemo    int      `json:"has_memo,omitempty"`
		ReqQueryParam
	}

	AckAssetsAttribute struct {
		AssetName       string  `json:"asset_name" `
		FullName        string  `json:"full_name"`
		IsToken         int     `json:"is_token"`
		ParentName      string  `json:"parent_name"`
		HasMemo         int     `json:"has_memo"`
		Logo            string  `json:"logo"`
		BrowserUrl      string  `json:"browser_url"`
		ContractAddress string  `json:"contract_address"`
		DepositMin      float64 `json:"deposit_min"`
		WithdrawalRate  float64 `json:"withdrawal_rate"`
		WithdrawalValue float64 `json:"withdrawal_value"`
		ConfirmationNum int     `json:"confirmation_num"`
		Decimals        int     `json:"decimals"`
	}

	AckAssetsAttributeList struct {
		Data []AckAssetsAttribute `json:"data"`
		AckQueryParam
	}

	// 获取用户余额
	ReqUserBalance struct {
		AssetNames []string `json:"asset_names,omitempty"`
		ReqQueryParam
	}

	AckUserBalance struct {
		AssetName       string  `json:"asset_name"`
		AvailableAmount float64 `json:"available_amount"`
		FrozenAmount    float64 `json:"frozen_amount"`
		Time            int64   `json:"time"`
	}

	AckUserBalanceList struct {
		Data []AckUserBalance `json:"data"`
		AckQueryParam
	}

	// 获取用户地址
	ReqUserAddress struct {
		AssetNames        []string `json:"asset_names,omitempty"`
		Address           string   `json:"address,omitempty"`
		Memo              string   `json:"memo,omitempty"`
		MaxAllocationTime int64    `json:"max_allocation_time,omitempty"`
		MinAllocationTime int64    `json:"min_allocation_time,omitempty"`
		ReqQueryParam
	}

	AckUserAddress struct {
		AssetName      string `json:"asset_name"`
		Address        string `json:"address"`
		Memo           string `json:"memo"`
		AllocationTime int64  `json:"allocation_time"`
	}

	AckUserAddressList struct {
		Data []AckUserAddress `json:"data"`
		AckQueryParam
	}

	// 历史交易订单
	ReqTransactionBill struct {
		ID             int64   `json:"id,omitempty"`
		OrderID        string  `json:"order_id,omitempty"`
		AssetName      string  `json:"asset_name,omitempty"`
		Address        string  `json:"address,omitempty"`
		Memo           string  `json:"memo,omitempty"`
		TransType      int     `json:"trans_type,omitempty"`
		TransTypes     []int   `json:"trans_types,omitempty"`
		Status         int     `json:"status,omitempty"`
		Hash           string  `json:"hash,omitempty"`
		MinAmount      float64 `json:"min_amount,omitempty"`
		MaxAmount      float64 `json:"max_amount,omitempty"`
		MinConfirmTime int64   `json:"min_confirm_time,omitempty"`
		MaxConfirmTime int64   `json:"max_confirm_time,omitempty"`
		ReqQueryParam
	}

	AckTransactionBill struct {
		ID               int64   `json:"id"`
		OrderID          string  `json:"order_id"`
		UserOrderID      string  `json:"user_order_id"`
		AssetName        string  `json:"asset_name"`
		Address          string  `json:"address"`
		Memo             string  `json:"memo"`
		TransType        int     `json:"trans_type"`
		FeeType          int     `json:"fee_type"`
		Amount           float64 `json:"amount"`
		Sign             string  `json:"sign"`
		PayFee           float64 `json:"pay_fee"`
		Balance          float64 `json:"balance"`
		Hash             string  `json:"hash"`
		Status           int     `json:"status"`
		BlockinHeight    int64   `json:"blockin_height"`
		CreateOrderTime  int64   `json:"create_order_time"`
		BlockinTime      int64   `json:"blockin_time"`
		ConfirmTime      int64   `json:"confirm_time"`
		RelationUserKey  string  `json:"relation_user_key,omitempty"`
		RelationUserName string  `json:"relation_user_name,omitempty"`
	}

	AckTransactionBillList struct {
		Data []AckTransactionBill `json:"data"`
		AckQueryParam
	}

	// 历史交易消息
	ReqTransactionMessage struct {
		Memo         string `json:"memo,omitempty"`
		MaxMessageID int64  `json:"max_msg_id,omitempty"`
		MinMessageID int64  `json:"min_msg_id,omitempty"`
		ReqQueryParam
	}

	AckTransactionMessage struct {
		MsgID            int64   `json:"msg_id"`
		TransType        int     `json:"trans_type"`
		Status           int     `json:"status"`
		BlockinHeight    int64   `json:"blockin_height"`
		AssetName        string  `json:"asset_name"`
		Address          string  `json:"address"`
		Memo             string  `json:"memo,omitempty"`
		Amount           float64 `json:"amount"`
		PayFee           float64 `json:"pay_fee"`
		Balance          float64 `json:"balance"`
		Hash             string  `json:"hash"`
		OrderID          string  `json:"order_id"`
		RelationUserKey  string  `json:"relation_user_key,omitempty"`
		RelationUserName string  `json:"relation_user_name,omitempty"`
		UserOrderID      string  `json:"user_order_id"`
		Remark           string  `json:"remark,omitempty"`
		Time             int64   `json:"time"`
	}

	AckTransactionMessageList struct {
		Data []AckTransactionMessage `json:"data"`
		AckQueryParam
	}

	ReqTransactionBillDaily struct {
		AssetName string `json:"asset_name,omitempty"`
		MaxPeriod int    `json:"max_period,omitempty"`
		MinPeriod int    `json:"min_period,omitempty"`
		ReqQueryParam
	}

	AckTransactionBillDaily struct {
		Period      int     `json:"period"`
		AssetName   string  `json:"asset_name"`
		SumDPAmount float64 `json:"sum_dp_amount"`
		SumWDAmount float64 `json:"sum_wd_amount"`
		SumT0Amount float64 `json:"sum_t0_amount"`
		SumT1Amount float64 `json:"sum_t1_amount"`
		SumPayFee   float64 `json:"sum_pay_fee"`
		PreBalance  float64 `json:"pre_balance"`
		LastBalance float64 `json:"last_balance"`
	}

	AckTransactionBillDailyList struct {
		Data []AckTransactionBillDaily `json:"data"`
		AckQueryParam
	}

	ReqBlockHeight struct {
		AssetNames []string `json:"asset_names,omitempty"`
	}

	AckBlockHeight struct {
		AssetName   string `json:"asset_name"`
		BlockHeight int64  `json:"block_height"`
		UpdateTime  int64  `json:"update_time"`
	}

	AckBlockHeightList struct {
		Data []AckBlockHeight `json:"data"`
	}

	ReqWithdrawalOrder struct {
		OrderID            string  `json:"order_id,omitempty"`
		AssetName          string  `json:"asset_name,omitempty"`
		Address            string  `json:"address,omitempty"`
		Memo               string  `json:"memo,omitempty"`
		Status             int     `json:"status,omitempty"`
		Hash               string  `json:"hash,omitempty"`
		MinAmount          float64 `json:"min_amount,omitempty"`
		MaxAmount          float64 `json:"max_amount,omitempty"`
		MinOrderTime       int64   `json:"min_order_time,omitempty"`
		MaxOrderTime       int64   `json:"max_order_time,omitempty"`
		MinConfirmTime     int64   `json:"min_confirm_time,omitempty"`
		MaxConfirmTime     int64   `json:"max_confirm_time,omitempty"`
		OrderByOrderTime   int     `json:"order_by_order_time,omitempty"`
		OrderByConfirmTime int     `json:"order_by_confirm_time,omitempty"`
		ReqQueryParam
	}

	AckWithdrawalOrder struct {
		UserKey       string `json:"user_key"`
		UserName      string `json:"user_name"`
		OrderID       string `json:"order_id"`
		Hash          string `json:"hash"`
		AssetName     string `json:"asset_name"`
		Address       string `json:"address"`
		Memo          string `json:"memo"`
		Status        int    `json:"status"`
		Amount        string `json:"amount"`
		PayFee        string `json:"pay_fee"`
		Balance       string `json:"balance"`
		BlockinHeight int64  `json:"blockin_height"`
		OrderTime     int64  `json:"order_time"`
		BlockinTime   int64  `json:"blockin_time"`
		ConfirmTime   int64  `json:"confirm_time"`
		UserOrderID   string `json:"user_order_id"`
	}

	AckWithdrawalOrderList struct {
		Data []AckWithdrawalOrder `json:"data"`
		AckQueryParam
	}

	ReqDepositOrder struct {
		OrderID            string  `json:"order_id,omitempty"`
		AssetName          string  `json:"asset_name,omitempty"`
		Address            string  `json:"address,omitempty"`
		Memo               string  `json:"memo,omitempty"`
		Status             int     `json:"status,omitempty"`
		Hash               string  `json:"hash,omitempty"`
		MaxAmount          float64 `json:"max_amount,omitempty"`
		MinAmount          float64 `json:"min_amount,omitempty"`
		MinOrderTime       int64   `json:"min_order_time,omitempty"`
		MaxOrderTime       int64   `json:"max_order_time,omitempty"`
		MinConfirmTime     int64   `json:"min_confirm_time,omitempty"`
		MaxConfirmTime     int64   `json:"max_confirm_time,omitempty"`
		OrderByOrderTime   int     `json:"order_by_order_time,omitempty"`
		OrderByConfirmTime int     `json:"order_by_confirm_time,omitempty"`
		ReqQueryParam
	}

	AckDepositOrder struct {
		UserKey       string `json:"user_key"`
		UserName      string `json:"user_name"`
		OrderID       string `json:"order_id"`
		Hash          string `json:"hash"`
		AssetName     string `json:"asset_name"`
		Address       string `json:"address"`
		Memo          string `json:"memo"`
		Status        int    `json:"status"`
		Amount        string `json:"amount"`
		PayFee        string `json:"pay_fee"`
		Balance       string `json:"balance"`
		BlockinHeight int64  `json:"blockin_height"`
		OrderTime     int64  `json:"order_time"`
		BlockinTime   int64  `json:"blockin_time"`
		ConfirmTime   int64  `json:"confirm_time"`
		UserOrderID   string `json:"user_order_id"`
	}

	AckDepositOrderList struct {
		Data []AckDepositOrder `json:"data"`
		AckQueryParam
	}

	ReqTransferOrder struct {
		OrderID            string   `json:"order_id,omitempty"`
		AssetName          string   `json:"asset_name,omitempty"`
		Address            string   `json:"address,omitempty"`
		Memo               string   `json:"memo,omitempty"`
		MaxAmount          float64  `json:"max_amount,omitempty"`
		MinAmount          float64  `json:"min_amount,omitempty"`
		Out                int      `json:"out,omitempty"`
		In                 int      `json:"in,omitempty"`
		MinOrderTime       int64    `json:"min_order_time,omitempty"`
		MaxOrderTime       int64    `json:"max_order_time,omitempty"`
		MinConfirmTime     int64    `json:"min_confirm_time,omitempty"`
		MaxConfirmTime     int64    `json:"max_confirm_time,omitempty"`
		RelationUserKey    string   `json:"relation_user_key,omitempty"`
		RelationUserName   string   `json:"relation_user_name,omitempty"`
		OrderByOrderTime   int      `json:"order_by_order_time,omitempty"`
		OrderByConfirmTime int      `json:"order_by_confirm_time,omitempty"`
		Remarks            []string `json:"remarks,omitempty"`
		ReqQueryParam
	}

	AckTransferOrder struct {
		UserKey          string `json:"user_key"`
		UserName         string `json:"user_name"`
		OrderID          string `json:"order_id"`
		TransType        int    `json:"trans_type"`
		Hash             string `json:"hash"`
		AssetName        string `json:"asset_name"`
		Address          string `json:"address"`
		Memo             string `json:"memo"`
		Status           int    `json:"status"`
		Amount           string `json:"amount"`
		PayFee           string `json:"pay_fee"`
		Balance          string `json:"balance"`
		BlockinHeight    int64  `json:"blockin_height"`
		OrderTime        int64  `json:"order_time"`
		BlockinTime      int64  `json:"blockin_time"`
		ConfirmTime      int64  `json:"confirm_time"`
		RelationUserKey  string `json:"relation_user_key"`
		RelationUserName string `json:"relation_user_name"`
		UserOrderID      string `json:"user_order_id"`
		Remark           string `json:"remark"`
	}

	AckTransferOrderList struct {
		Data []AckTransferOrder `json:"data"`
		AckQueryParam
	}

	ReqDashboardSummary struct {
		ReqQueryParam
	}

	AckDashboardAssetsUSD struct {
		AssetName string `json:"asset_name"`
		Amount    string `json:"amount"`
	}

	AckDashboardPeriodUSD struct {
		Period   string `json:"period"`
		Amount   string `json:"amount"`
		DPAmount string `json:"dp_amount"`
		WDAmount string `json:"wd_amount"`
		T0Amount string `json:"t0_amount"`
		T1Amount string `json:"t1_amount"`
	}

	AckDashboardBill struct {
		AssetName  string `json:"asset_name"`
		TransType  int    `json:"trans_type"`
		FeeType    int    `json:"fee_type"`
		Amount     string `json:"amount"`
		UpdateTime int64  `json:"update_time"`
	}

	AckDashboardSummary struct {
		Support       int                     `json:"support"`
		Address       int                     `json:"address"`
		Deposit       int                     `json:"deposit"`
		Withdrawal    int                     `json:"withdrawal"`
		Deposit24h    string                  `json:"deposit_24h"`
		Withdarwal24h string                  `json:"withdrawal_24h"`
		AssetsUSD     []AckDashboardAssetsUSD `json:"assets_usd"`
		PeriodUSD     []AckDashboardPeriodUSD `json:"period_usd"`
		BillList      []AckDashboardBill      `json:"bill_list"`
		AckQueryParam
	}

	PushTransactionMessage struct {
		MsgID            int64   `json:"msg_id"`
		TransType        int     `json:"trans_type"`
		Status           int     `json:"status"`
		BlockinHeight    int64   `json:"blockin_height"`
		AssetName        string  `json:"asset_name"`
		Address          string  `json:"address"`
		Memo             string  `json:"memo"`
		Amount           float64 `json:"amount"`
		PayFee           float64 `json:"pay_fee"`
		Balance          float64 `json:"balance"`
		Hash             string  `json:"hash"`
		OrderID          string  `json:"order_id"`
		UserOrderID      string  `json:"user_order_id"`
		RelationUserKey  string  `json:"relation_user_key,omitempty"`
		RelationUserName string  `json:"relation_user_name,omitempty"`
		Remark           string  `json:"remark,omitempty"`
		Time             int64   `json:"time"`
	}
)
