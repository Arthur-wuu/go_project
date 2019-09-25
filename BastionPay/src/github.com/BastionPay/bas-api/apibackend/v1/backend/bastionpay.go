package backend

type (
	SpReqQueryParam struct {
		IsAsyn       int    `json:"is_asyn" doc:"0:同步返回数据;1:异步返回数据"`
		PageIndex    int    `json:"page_index" doc:"页索引,从1开始"`
		MaxDispLines int    `json:"max_disp_lines" doc:"查询返回数据条数"`
		TotalLines   int    `json:"total_lines" doc:"查询数据总条数"`
		UserParam    string `json:"user_param,omitempty" doc:"用户自定义字段"`
	}

	SpAckQueryParam struct {
		IsAsyn       int    `json:"is_asyn" doc:"0:同步返回数据;1:异步返回数据"`
		QueryID      string `json:"query_id" doc:"查询序号"`
		PageIndex    int    `json:"page_index" doc:"页索引,1开始"`
		MaxDispLines int    `json:"max_disp_lines" doc:"页最大数,100以下"`
		TotalLines   int    `json:"total_lines" doc:"总数,0：表示首次查询"`
		UserParam    string `json:"user_param,omitempty" doc:"用户自定义字段"`
	}

	SpReqQueryAsynData struct {
		QueryID string `json:"query_id" doc:"查询序号"`
	}

	SpAckQueryAsynData struct {
		QueryID   string `json:"query_id" doc:"查询序号"`
		Status    int    `json:"status" doc:"0:数据准备中;1:查询数据成功;2:数据查询失败"`
		Data      string `json:"data" doc:"异步返回的数据"`
		UserParam string `json:"user_param,omitempty" doc:"用户自定义字段"`
	}

	SpCollectTransaction struct {
		AssetName string  `json:"asset_name"`
		From      string  `json:"from"`
		Amount    float64 `json:"amount"`
		IsTotal   int     `json:"is_total"`
	}

	SpReqCollectTransaction struct {
		SpCollectTransaction
	}

	SpAckCollectTransaction struct {
		OrderID string `json:"order_id"`
	}

	SpReqBatchCollectTransaction []SpCollectTransaction

	SpAckBatchCollectTransaction struct {
		BatchID  string   `json:"batch_id"`
		OrderIDs []string `json:"order_ids"`
	}

	// 获取币种属性
	SpReqAssetsAttributeList struct {
		AssetNames []string `json:"asset_names" doc:"需要查询属性的币种列表，不空表示精确查找"`
		IsToken    int      `json:"is_token" doc:"是否代币，-1:所有，0：不是代币，非0：代币"`
		HasMemo    int      `json:"has_memo" doc:"是否可以设置Memo"`
		Enabled    int      `json:"enabled" doc:"是否支持服务，-1:所有，0：不支持， 1：支持"`
		SpReqQueryParam
	}

	SpAckAssetsAttribute struct {
		AssetName             string  `json:"asset_name"`
		FullName              string  `json:"full_name"`
		OfficialName          string  `json:"official_name"`
		IsToken               int     `json:"is_token"`
		ParentName            string  `json:"parent_name"`
		ContractAddress       string  `json:"contract_address"`
		Logo                  string  `json:"logo"`
		BrowserUrl            string  `json:"browser_url"`
		Wallet                string  `json:"wallet"`
		DepositMin            float64 `json:"deposit_min"`
		DepositStategy        float64 `json:"deposit_stategy"`
		WithdrawalRate        float64 `json:"withdrawal_rate"`
		WithdrawalValue       float64 `json:"withdrawal_value"`
		WithdrawalReserveRate float64 `json:"withdrawal_reserve_rate"`
		WithdrawalAlertRate   float64 `json:"withdrawal_alert_rate"`
		WithdrawalStategy     float64 `json:"withdrawal_stategy"`
		ConfirmationNum       int     `json:"confirmation_num"`
		Decimals              int     `json:"decimals"`
		GasFactor             float64 `json:"gas_factor"`
		Debt                  float64 `json:"debt"`
		ParkAmount            float64 `json:"park_amount"`
		HasMemo               int     `json:"has_memo"`
		HasAccount            int     `json:"has_account"`
		TxGap                 int     `json:"tx_gap"`
		IsOffline             int     `json:"is_offline"`
		IsMultiTx             int     `json:"is_multi_tx"`
		ExchangeRate          string  `json:"exchange_rate"`
		Enabled               int     `json:"enabled"`
	}

	SpAckAssetsAttributeList struct {
		Data []SpAckAssetsAttribute `json:"data" doc:"币种属性列表"`
		SpAckQueryParam
	}

	SpReqSetAssetAttribute struct {
		AssetName             string  `json:"asset_name"`
		FullName              string  `json:"full_name"`
		IsToken               int     `json:"is_token"`
		ParentName            string  `json:"parent_name"`
		OfficialName          string  `json:"official_name"`
		ContractAddress       string  `json:"contract_address"`
		Logo                  string  `json:"logo"`
		BrowserUrl            string  `json:"browser_url"`
		Wallet                string  `json:"wallet"`
		DepositMin            float64 `json:"deposit_min"`
		DepositStategy        float64 `json:"deposit_stategy"`
		WithdrawalRate        float64 `json:"withdrawal_rate"`
		WithdrawalValue       float64 `json:"withdrawal_value"`
		WithdrawalReserveRate float64 `json:"withdrawal_reserve_rate"`
		WithdrawalAlertRate   float64 `json:"withdrawal_alert_rate"`
		WithdrawalStategy     float64 `json:"withdrawal_stategy"`
		ConfirmationNum       int     `json:"confirmation_num"`
		Decimals              int     `json:"decimals"`
		GasFactor             float64 `json:"gas_factor"`
		Debt                  float64 `json:"debt"`
		ParkAmount            float64 `json:"park_amount"`
		ExchangeRate          string  `json:"exchange_rate"`
		HasMemo               int     `json:"has_memo"`
		HasAccount            int     `json:"has_account"`
		TxGap                 int     `json:"tx_gap"`
		IsOffline             int     `json:"is_offline"`
		IsMultiTx             int     `json:"is_multi_tx"`
		Enabled               int     `json:"enabled"`
	}

	// 获取用户地址
	SpReqUserAddress struct {
		UserName          string   `json:"user_name" doc:"用户名称"`
		UserKey           string   `json:"user_key" doc:"用户Key"`
		UserClass         int      `json:"user_class" doc:"用户类型"`
		AssetNames        []string `json:"asset_names" doc:"币种"`
		MaxAllocationTime int64    `json:"max_allocation_time" doc:"分配地址时间"`
		MinAllocationTime int64    `json:"min_allocation_time" doc:"分配地址时间"`
		Address           string   `json:"address" doc:"地址"`
		Memo              string   `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		SpReqQueryParam
	}

	SpAckUserAddress struct {
		UserKey         string  `json:"user_key"`
		UserName        string  `json:"user_name"`
		UserClass       int     `json:"user_class"`
		AssetName       string  `json:"asset_name"`
		Address         string  `json:"address"`
		Memo            string  `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		PrivateKey      string  `json:"private_key"`
		AvailableAmount float64 `json:"available_amount"`
		FrozenAmount    float64 `json:"frozen_amount"`
		Enabled         int     `json:"enabled"`
		CreateTime      int64   `json:"create_time"`
		AllocationTime  int64   `json:"allocation_time"`
		UpdateTime      int64   `json:"update_time"`
	}

	SpAckUserAddressList struct {
		Data []SpAckUserAddress `json:"data" doc:"用户地址列表"`
		SpAckQueryParam
	}

	// 获取用户余额
	SpReqChainBalance struct {
		AssetName string `json:"asset_name" doc:"需要查询余额的币种列表"`
		Address   string `json:"address"`
		SpReqQueryParam
	}

	SpAckChainBalance struct {
		AssetName string  `json:"asset_name" doc:"币种简称"`
		Address   string  `json:"address"`
		Amount    float64 `json:"available_amount" doc:"可用余额"`
	}

	SpAckChainBalanceList struct {
		Data []SpAckChainBalance `json:"data" doc:"币种余额列表"`
		SpAckQueryParam
	}

	// 历史交易订单
	SpReqTransactionBill struct {
		UserName       string  `json:"user_name" doc:"用户名称"`
		UserKey        string  `json:"user_key"`
		ID             int64   `json:"id" doc:"流水号"`
		OrderID        string  `json:"order_id" doc:"订单号"`
		AssetName      string  `json:"asset_name" doc:"币种"`
		Memo           string  `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Address        string  `json:"address" doc:"地址"`
		TransType      int     `json:"trans_type,omitempty" doc:"交易类型"`
		TransTypes     []int   `json:"trans_types,omitempty" doc:"交易类型"`
		Status         int     `json:"status" doc:"交易状态"`
		Hash           string  `json:"hash" doc:"交易哈希"`
		MaxAmount      float64 `json:"max_amount" doc:"最大金额"`
		MinAmount      float64 `json:"min_amount" doc:"最小金额"`
		MaxConfirmTime int64   `json:"max_confirm_time" doc:"开始时间"`
		MinConfirmTime int64   `json:"min_confirm_time" doc:"结束时间"`
		SpReqQueryParam
	}

	SpAckTransactionBill struct {
		ID               int64   `json:"id" doc:"流水号"`
		UserKey          string  `json:"user_key"`
		UserName         string  `json:"user_name"`
		OrderID          string  `json:"order_id" doc:"交易订单"`
		UserOrderID      string  `json:"user_order_id" doc:"用户订单号"`
		AssetName        string  `json:"asset_name" doc:"币种"`
		Address          string  `json:"address" doc:"地址"`
		Memo             string  `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		TransType        int     `json:"trans_type" doc:"交易类型"`
		FeeType          int     `json:"fee_type" doc:"费用类型"`
		Amount           float64 `json:"amount" doc:"数量"`
		Sign             string  `json:"sign" doc:"正负符号"`
		PayFee           float64 `json:"pay_fee" doc:"交易费用(0:充提金额,1:提币手续费,2:矿工费,3:利润收入,4:矿工费支出)"`
		MinerFee         float64 `json:"miner_fee" doc:"矿工费"`
		Balance          float64 `json:"balance"doc:"当前余额"`
		Hash             string  `json:"hash"`
		Status           int     `json:"status" doc:"交易状态"`
		BlockinHeight    int64   `json:"blockin_height" doc:"入块高度"`
		CreateOrderTime  int64   `json:"create_order_time" doc:"订单创建时间"`
		BlockinTime      int64   `json:"blockin_time" doc:"入块时间"`
		ConfirmTime      int64   `json:"confirm_time" doc:"确认时间"`
		RelationUserKey  string  `json:"relation_user_key,omitempty"`
		RelationUserName string  `json:"relation_user_name,omitempty"`
	}

	SpAckTransactionBillList struct {
		Data []SpAckTransactionBill `json:"data" doc:"历史交易订单列表"`
		SpAckQueryParam
	}

	SpReqTransactionBillDaily struct {
		UserKey   string `json:"user_key"`
		AssetName string `json:"asset_name" doc:"币种"`
		MaxPeriod int    `json:"max_period" doc:"最大周期值"`
		MinPeriod int    `json:"min_period" doc:"最小周期值"`
		SpReqQueryParam
	}

	SpAckTransactionBillDaily struct {
		UserKey     string  `json:"user_key"`
		UserName    string  `json:"user_name"`
		Period      int     `json:"period"`
		AssetName   string  `json:"asset_name"`
		SumDPAmount float64 `json:"sum_dp_amount"`
		SumWDAmount float64 `json:"sum_wd_amount"`
		SumT0Amount float64 `json:"sum_t0_amount"`
		SumT1Amount float64 `json:"sum_t1_amount"`
		SumPayFee   float64 `json:"sum_pay_fee"`
		SumMinerFee float64 `json:"sum_miner_fee"`
		PreBalance  float64 `json:"pre_balance"`
		LastBalance float64 `json:"last_balance"`
	}

	SpAckTransactionBillDailyList struct {
		Data []SpAckTransactionBillDaily `json:"data" doc:"历史日结帐单"`
		SpAckQueryParam
	}

	// 获取用户余额
	SpReqUserBalance struct {
		UserKey    string   `json:"user_key" doc:"用户Key"`
		AssetNames []string `json:"asset_names" doc:"需要查询余额的币种列表"`
		SpReqQueryParam
	}

	SpAckUserBalance struct {
		UserKey         string  `json:"user_key" doc:"用户Key"`
		UserName        string  `json:"user_name" doc:"用户名称"`
		AssetName       string  `json:"asset_name" doc:"币种简称"`
		AvailableAmount float64 `json:"available_amount" doc:"可用余额"`
		FrozenAmount    float64 `json:"frozen_amount" doc:"冻结余额"`
		Time            int64   `json:"time" doc:"刷新时间"`
	}

	SpAckUserBalanceList struct {
		Data []SpAckUserBalance `json:"data" doc:"币种余额列表"`
		SpAckQueryParam
	}

	SpAckSetPayAddress struct {
		AssetName string `json:"asset_name" doc:"币种简称"`
		Address   string `json:"address" doc:"地址"`
	}

	SpReqPayAddress struct {
		AssetNames []string `json:"asset_names" doc:"需要查询属性的币种列表，不空表示精确查找"`
		SpReqQueryParam
	}

	SpAckPayAddress struct {
		AssetName  string  `json:"asset_name" doc:"币种简称"`
		Address    string  `json:"address" doc:"地址"`
		Amount     float64 `json:"amount" doc:"数量"`
		UpdateTime int64   `json:"update_time" doc:"更新时间"`
	}

	SpAckPayAddressList struct {
		Data []SpAckPayAddress `json:"data" doc:"热钱包地址列表"`
		SpAckQueryParam
	}

	SpReqWithdrawalOrder struct {
		UserName       string  `json:"user_name" doc:"用户名称"`
		UserKey        string  `json:"user_key" doc:"用户Key"`
		OrderID        string  `json:"order_id" doc:"订单号"`
		AssetName      string  `json:"asset_name" doc:"币种"`
		Address        string  `json:"address" doc:"地址"`
		Memo           string  `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Status         int     `json:"status" doc:"交易状态"`
		Hash           string  `json:"hash" doc:"交易哈希"`
		MinAmount      float64 `json:"min_amount" doc:"最小金额"`
		MaxAmount      float64 `json:"max_amount" doc:"最大金额"`
		MinOrderTime   int64   `json:"min_order_time" doc:"按订单创建时间查询开始"`
		MaxOrderTime   int64   `json:"max_order_time" doc:"按订单创建时间查询结束"`
		MinConfirmTime int64   `json:"min_confirm_time" doc:"按确认时间查询开始"`
		MaxConfirmTime int64   `json:"max_confirm_time" doc:"按确认时间查询结束"`
		Tag            int     `json:"tag"`
		SpReqQueryParam
	}

	SpAckWithdrawalOrder struct {
		UserKey       string `json:"user_key" doc:"用户Key"`
		UserName      string `json:"user_name" doc:"用户名称"`
		OrderID       string `json:"order_id" doc:"订单号"`
		Hash          string `json:"hash" doc:"交易哈希值"`
		AssetName     string `json:"asset_name" doc:"币种名称"`
		Address       string `json:"address" doc:"提币地址"`
		Memo          string `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Status        int    `json:"status" doc:"交易状态"`
		Amount        string `json:"amount" doc:"交易金额"`
		PayFee        string `json:"pay_fee" doc:"交易手续费"`
		Balance       string `json:"balance" doc:"当前余额"`
		BlockinHeight int64  `json:"blockin_height" doc:"入块高度"`
		OrderTime     int64  `json:"order_time" doc:"订单创建时间"`
		BlockinTime   int64  `json:"blockin_time" doc:"入块时间"`
		ConfirmTime   int64  `json:"confirm_time" doc:"确认时间"`
		UserOrderID   string `json:"user_order_id" doc:"用户自定义ID"`
		RefundInfo    string `json:"refund_info" doc:"退款说明"`
		OrderMemo     string `json:"order_memo"`
		Tag           int    `json:"tag"`
		TagUser       string `json:"tag_user"`
	}

	SpAckWithdrawalOrderList struct {
		Data []SpAckWithdrawalOrder `json:"data" doc:"提币订单列表"`
		SpAckQueryParam
	}

	SpReqDepositOrder struct {
		UserName       string  `json:"user_name" doc:"用户名称"`
		UserKey        string  `json:"user_key" doc:"用户Key"`
		OrderID        string  `json:"order_id" doc:"订单号"`
		AssetName      string  `json:"asset_name" doc:"币种"`
		Address        string  `json:"address" doc:"地址"`
		Memo           string  `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Status         int     `json:"status" doc:"交易状态"`
		Hash           string  `json:"hash" doc:"交易哈希"`
		MaxAmount      float64 `json:"max_amount" doc:"最大金额"`
		MinAmount      float64 `json:"min_amount" doc:"最小金额"`
		MinOrderTime   int64   `json:"min_order_time" doc:"按订单创建时间查询开始"`
		MaxOrderTime   int64   `json:"max_order_time" doc:"按订单创建时间查询结束"`
		MinConfirmTime int64   `json:"min_confirm_time" doc:"按确认时间查询开始"`
		MaxConfirmTime int64   `json:"max_confirm_time" doc:"按确认时间查询结束"`
		Tag            int     `json:"tag" doc:"状态标记"`
		SpReqQueryParam
	}

	SpAckDepositOrder struct {
		UserKey       string `json:"user_key" doc:"用户Key"`
		UserName      string `json:"user_name" doc:"用户名称"`
		OrderID       string `json:"order_id" doc:"订单号"`
		Hash          string `json:"hash" doc:"交易哈希值"`
		AssetName     string `json:"asset_name" doc:"币种名称"`
		Address       string `json:"address" doc:"充值地址"`
		Memo          string `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Status        int    `json:"status" doc:"交易状态"`
		Amount        string `json:"amount" doc:"交易金额"`
		PayFee        string `json:"pay_fee" doc:"交易手续费"`
		Balance       string `json:"balance" doc:"当前余额"`
		BlockinHeight int64  `json:"blockin_height" doc:"入块高度"`
		OrderTime     int64  `json:"order_time" doc:"订单创建时间"`
		BlockinTime   int64  `json:"blockin_time" doc:"入块时间"`
		ConfirmTime   int64  `json:"confirm_time" doc:"确认时间"`
		UserOrderID   string `json:"user_order_id" doc:"用户自定义ID"`
		OrderMemo     string `json:"order_memo"`
		Tag           int    `json:"tag"`
		TagUser       string `json:"tag_user"`
	}

	SpAckDepositOrderList struct {
		Data []SpAckDepositOrder `json:"data" doc:"提币订单列表"`
		SpAckQueryParam
	}

	SpReqTransferOrder struct {
		UserName         string  `json:"user_name" doc:"用户名称"`
		UserKey          string  `json:"user_key" doc:"用户Key"`
		OrderID          string  `json:"order_id" doc:"订单号"`
		AssetName        string  `json:"asset_name" doc:"币种"`
		Address          string  `json:"address" doc:"地址"`
		Memo             string  `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		MinAmount        float64 `json:"min_amount" doc:"最小金额"`
		MaxAmount        float64 `json:"max_amount" doc:"最大金额"`
		Out              int     `json:"out,omitempty"`
		In               int     `json:"in,omitempty"`
		MinOrderTime     int64   `json:"min_order_time" doc:"按订单创建时间查询开始"`
		MaxOrderTime     int64   `json:"max_order_time" doc:"按订单创建时间查询结束"`
		MinConfirmTime   int64   `json:"min_confirm_time" doc:"按确认时间查询开始"`
		MaxConfirmTime   int64   `json:"max_confirm_time" doc:"按确认时间查询结束"`
		RelationUserKey  string  `json:"relation_user_key,omitempty"`
		RelationUserName string  `json:"relation_user_name,omitempty"`
		SpReqQueryParam
	}

	SpAckTransferOrder struct {
		UserKey          string `json:"user_key" doc:"用户Key"`
		UserName         string `json:"user_name" doc:"用户名称"`
		OrderID          string `json:"order_id" doc:"订单号"`
		TransType        int    `json:"trans_type"`
		Hash             string `json:"hash" doc:"交易哈希值"`
		AssetName        string `json:"asset_name" doc:"币种名称"`
		Address          string `json:"address" doc:"提币地址"`
		Memo             string `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Status           int    `json:"status" doc:"交易状态"`
		Amount           string `json:"amount" doc:"交易金额"`
		PayFee           string `json:"pay_fee" doc:"交易手续费"`
		Balance          string `json:"balance" doc:"当前余额"`
		BlockinHeight    int64  `json:"blockin_height" doc:"入块高度"`
		OrderTime        int64  `json:"order_time" doc:"订单创建时间"`
		BlockinTime      int64  `json:"blockin_time" doc:"入块时间"`
		ConfirmTime      int64  `json:"confirm_time" doc:"确认时间"`
		RelationUserKey  string `json:"relation_user_key"`
		RelationUserName string `json:"relation_user_name"`
		UserOrderID      string `json:"user_order_id" doc:"用户自定义ID"`
		Remark           string `json:"remark"`
		OrderMemo        string `json:"order_memo"`
	}

	SpAckTransferOrderList struct {
		Data []SpAckTransferOrder `json:"data" doc:"提币订单列表"`
		SpAckQueryParam
	}

	SpReqCollectOrder struct {
		BatchID        string  `json:"batch_id" doc:"批次号"`
		OrderID        string  `json:"order_id" doc:"订单号"`
		AssetName      string  `json:"asset_name" doc:"币种"`
		Address        string  `json:"address" doc:"地址"`
		Status         int     `json:"status" doc:"交易状态"`
		Hash           string  `json:"hash" doc:"交易哈希"`
		MaxAmount      float64 `json:"max_amount" doc:"最大金额"`
		MinAmount      float64 `json:"min_amount" doc:"最小金额"`
		MinOrderTime   int64   `json:"min_order_time" doc:"按订单创建时间查询开始"`
		MaxOrderTime   int64   `json:"max_order_time" doc:"按订单创建时间查询结束"`
		MinConfirmTime int64   `json:"min_confirm_time" doc:"按确认时间查询开始"`
		MaxConfirmTime int64   `json:"max_confirm_time" doc:"按确认时间查询结束"`
		SpReqQueryParam
	}

	SpAckCollectOrder struct {
		BatchID       string `json:"batch_id" doc:"批次号"`
		OrderID       string `json:"order_id" doc:"订单号"`
		Hash          string `json:"hash" doc:"交易哈希值"`
		TransType     int    `json:"trans_type" doc:"交易类型 {2:加油, 3:归集}"`
		AssetName     string `json:"asset_name" doc:"币种名称"`
		Address       string `json:"address" doc:"充值地址"`
		Status        int    `json:"status" doc:"交易状态"`
		Amount        string `json:"amount" doc:"交易金额"`
		MinerFee      string `json:"miner_fee" doc:"交易手续费"`
		BlockinHeight int64  `json:"blockin_height" doc:"入块高度"`
		OrderTime     int64  `json:"order_time" doc:"订单创建时间"`
		BlockinTime   int64  `json:"blockin_time" doc:"入块时间"`
		ConfirmTime   int64  `json:"confirm_time" doc:"确认时间"`
		OrderMemo     string `json:"order_memo"`
	}

	SpAckCollectOrderList struct {
		Data []SpAckCollectOrder `json:"data" doc:"提币订单列表"`
		SpAckQueryParam
	}

	SpReqAbnormalOrder struct {
		UserName         string  `json:"user_name" doc:"用户名称"`
		UserKey          string  `json:"user_key" doc:"用户Key"`
		OrderID          string  `json:"order_id" doc:"订单号"`
		AssetName        string  `json:"asset_name" doc:"币种"`
		Address          string  `json:"address" doc:"地址"`
		Memo             string  `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Status           int     `json:"status" doc:"交易状态"`
		Hash             string  `json:"hash" doc:"交易哈希"`
		MinAmount        float64 `json:"min_amount" doc:"最小金额"`
		MaxAmount        float64 `json:"max_amount" doc:"最大金额"`
		MinOrderTime     int64   `json:"min_order_time" doc:"按订单创建时间查询开始"`
		MaxOrderTime     int64   `json:"max_order_time" doc:"按订单创建时间查询结束"`
		MinConfirmTime   int64   `json:"min_confirm_time" doc:"按确认时间查询开始"`
		MaxConfirmTime   int64   `json:"max_confirm_time" doc:"按确认时间查询结束"`
		RelationUserKey  string  `json:"relation_user_key,omitempty"`
		RelationUserName string  `json:"relation_user_name,omitempty"`
		SpReqQueryParam
	}

	SpAckAbnormalOrder struct {
		UserKey          string `json:"user_key" doc:"用户Key"`
		UserName         string `json:"user_name" doc:"用户名称"`
		OrderID          string `json:"order_id" doc:"订单号"`
		Hash             string `json:"hash" doc:"交易哈希值"`
		AssetName        string `json:"asset_name" doc:"币种名称"`
		Address          string `json:"address" doc:"提币地址"`
		Memo             string `json:"memo" doc:"特定币种用于充提币时的用户标识"`
		Status           int    `json:"status" doc:"交易状态"`
		Amount           string `json:"amount" doc:"交易金额"`
		PayFee           string `json:"pay_fee" doc:"交易手续费"`
		Balance          string `json:"balance" doc:"当前余额"`
		BlockinHeight    int64  `json:"blockin_height" doc:"入块高度"`
		OrderTime        int64  `json:"order_time" doc:"订单创建时间"`
		BlockinTime      int64  `json:"blockin_time" doc:"入块时间"`
		ConfirmTime      int64  `json:"confirm_time" doc:"确认时间"`
		RelationUserKey  string `json:"relation_user_key"`
		RelationUserName string `json:"relation_user_name"`
		UserOrderID      string `json:"user_order_id" doc:"用户自定义ID"`
		AlarmInfo        string `json:"alarm_info" doc:"异常订单备注"`
		AlarmUserKey     string `json:"alarm_user_key" doc:"异常订单备注用户"`
	}

	SpAckAbnormalOrderList struct {
		Data []SpAckAbnormalOrder `json:"data" doc:"提币订单列表"`
		SpAckQueryParam
	}

	SpReqWithdrawalProfitDaily struct {
		AssetName string `json:"asset_name" doc:"币种"`
		MaxPeriod int    `json:"max_period" doc:"最大周期值"`
		MinPeriod int    `json:"min_period" doc:"最小周期值"`
		SpReqQueryParam
	}

	SpAckWithdrawalProfitDaily struct {
		Period       int    `json:"period" doc:"日期"`
		AssetName    string `json:"asset_name"  doc:"数字货币代码"`
		PreSumPayFee string `json:"pre_sum_pay_fee"  doc:"期初提币手续费总和"`
		PayFee       string `json:"pay_fee"  doc:"当期手续费增加"`
		SumPayFee    string `json:"sum_pay_fee"  doc:"期末提币手续费总和"`
		Time         int64  `json:"time" doc:"更新时间"`
	}

	SpAckWithdrawalProfitDailyList struct {
		Data []SpAckWithdrawalProfitDaily `json:"data" doc:"提币手续费统计列表"`
		SpAckQueryParam
	}

	SpReqRealityAccount struct {
		AssetName             string  `json:"asset_name" doc:"币种"`
		Address               string  `json:"address" doc:"地址"`
		UserClass             int     `json:"user_class" doc:"-1:全部;0:冷钱包;1:热钱包"`
		MaxAccount            float64 `json:"max_account" doc:"帐户最大金额"`
		MinAccount            float64 `json:"min_account" doc:"帐户最小金额"`
		MaxReality            float64 `json:"max_reality" doc:"实际最大金额"`
		MinReality            float64 `json:"min_reality" doc:"实际最小金额"`
		IsOffline             int     `json:"is_offline" doc:"是否离线地址"`
		IsAccountRealityEqual int     `json:"is_account_reality_equal" doc:"帐实是否相等"`
		SpReqQueryParam
	}

	SpAckRealityAccount struct {
		UserKey   string `json:"user_key" doc:"用户Key"`
		UserName  string `json:"user_name" doc:"用户名称"`
		AssetName string `json:"asset_name" doc:"币种"`
		Address   string `json:"address" doc:"地址"`
		UserClass int    `json:"user_class" doc:"0:冷钱包;1:热钱包"`
		Account   string `json:"account" doc:"地址余额"`
		Reality   string `json:"reality" doc:"地址链上余额"`
		IsOffline int    `json:"is_offline" doc:"是否离线地址"`
		Time      int64  `json:"time" doc:"更新时间"`
	}

	SpAckRealityAccountList struct {
		Data []SpAckRealityAccount `json:"data" doc:"币种余额列表"`
		SpAckQueryParam
	}

	SpReqRealityAccountDaily struct {
		AssetName string `json:"asset_name" doc:"币种"`
		MaxPeriod int    `json:"max_period" doc:"最大周期值"`
		MinPeriod int    `json:"min_period" doc:"最小周期值"`
		SpReqQueryParam
	}

	SpAckRealityAccountDaily struct {
		ID          int64  `json:"id" doc:"流水号"`
		Period      string `json:"period" doc:"日期"`
		AssetName   string `json:"asset_name" doc:"币种"`
		Deposit     string `json:"deposit" doc:"当期充值累计金额"`
		Withdrawal  string `json:"withdrawal" doc:"当期提币累计金额"`
		Profit      string `json:"profit" doc:"当期利润累计金额"`
		PayFee      string `json:"pay_fee" doc:"当期手续费累计金额"`
		MinerFee    string `json:"miner_fee" doc:"当期矿工费累计金额"`
		Available   string `json:"available" doc:"当期可用金额"`
		Frozen      string `json:"frozen" doc:"当期冻洁金额"`
		Balance     string `json:"balance" doc:"当期余额"`
		SumCold     string `json:"sum_cold" doc:"当期冷钱包帐上余额"`
		SumHot      string `json:"sum_hot" doc:"当期热钱包帐上余额"`
		SumRealCold string `json:"sum_real_cold" doc:"当期冷钱包链上余额"`
		SumRealHot  string `json:"sum_real_hot" doc:"当期热钱包链上余额"`
		PreBalance  string `json:"pre_balance" doc:"上期余额"`
		Sign        int    `json:"sign" doc:"标记{0:未检查; 1:没有问题; 2:有问题; 3:已标记}"`
		Memo        string `json:"memo" doc:"标记内容"`
		ID2         int64  `json:"id_2" doc:"流水号2"`
		Sign2       int    `json:"sign_2" doc:"标记{0:未检查; 1:没有问题; 2:有问题; 3:已标记}"`
		Memo2       string `json:"memo_2" doc:"标记内容"`
		Time        int64  `json:"time" doc:"最后更新时间"`
	}

	SpAckRealityAccountDailyList struct {
		Data []SpAckRealityAccountDaily `json:"data" doc:"币种余额列表"`
		SpAckQueryParam
	}

	SpSetDailyMemo struct {
		ID   int64  `json:"id" doc:"ID号"`
		Type int    `json:"type" doc:"标记类型{1:订单对帐; 2:帐户对帐}"`
		Sign int    `json:"sign"  doc:"标记{0:未检查; 1:没有问题; 2:有问题; 3:已标记}"`
		Memo string `json:"memo"  doc:"标记内容"`
	}

	SpSetAlarmInfo struct {
		OrderID   string `json:"order_id" doc:"订单号"`
		AlarmInfo string `json:"alarm_info"  doc:"标记内容"`
	}

	SpReqRegistToken struct {
		AssetName string `json:"asset_name" doc:"主链名称"`
		Address   string `json:"address" doc:"合约地址"`
		Register  int    `json:"register" doc:"0:注册; 1:注册 2:取消注册"`
	}

	SpReqActiveToken struct {
		AssetName    string `json:"asset_name" doc:"Coin简称"`
		TokenName    string `json:"token_name" doc:"Token简称"`
		CoinAddress  string `json:"coin_address" doc:"Coin地址"`
		TokenAddress string `json:"token_address" doc:"Token地址"`
	}

	SpSetTxFail struct {
		OrderID     string `json:"order_id" doc:"订单号"`
		HasMinerFee int    `json:"has_miner_fee" doc:"是否需要矿工费"`
		RefundInfo  string `json:"refund_info" doc:"退款说明"`
	}

	SpReqCollectBooking struct {
		AssetName string `json:"asset_name" doc:"数字货币代码"`
		IsTotal   int    `json:"is_total"`
		Data      []struct {
			Address string `json:"address"`
			Amount  string `json:"amount"`
		} `json:"data" doc:"归集地址金额列表"`
	}

	SpAckCollectBooking struct {
		BatchID string `json:"batch_id"`
	}

	SpReqCollectBatch struct {
		BatchID        string `json:"batch_id"`
		AssetName      string `json:"asset_name"`
		TxBuild        int    `json:"tx_build" doc:"创建交易状态"`
		TxSigned       int    `json:"tx_signed" doc:"签名完成状态"`
		MinOrderTime   int64  `json:"min_order_time" doc:"按订单创建时间查询开始"`
		MaxOrderTime   int64  `json:"max_order_time" doc:"按订单创建时间查询结束"`
		MinConfirmTime int64  `json:"min_confirm_time" doc:"按确认时间查询开始"`
		MaxConfirmTime int64  `json:"max_confirm_time" doc:"按确认时间查询结束"`
		SpReqQueryParam
	}

	SpAckCollectBatch struct {
		BatchID     string `json:"batch_id"`
		AssetName   string `json:"asset_name"`
		OrderTime   int64  `json:"order_time" doc:"订单创建时间"`
		ConfirmTime int64  `json:"confirm_time" doc:"订单创建完成时间"`
		TotalAmount string `json:"total_amount" doc:"总归集数量"`
		TxBuild     int    `json:"tx_build" doc:"创建交易状态"`
		TxSigned    int    `json:"tx_signed" doc:"签名完成状态"`
		TxTotal     int    `json:"tx_total" doc:"订单总数"`
		TxSuccess   int    `json:"tx_success" doc:"订单成功数"`
		TxComplete  int    `json:"tx_complete" doc:"订单完成数"`
		Error       int    `json:"error"`
		ErrorMsg    string `json:"error_msg" doc:"错误信息"`
	}

	SpAckCollectBatchList struct {
		Data []SpAckCollectBatch `json:"data" doc:"归集批次列表"`
		SpAckQueryParam
	}

	SpReqCollectBatchTxFile struct {
		BatchID string `json:"batch_id"`
		SpReqQueryParam
	}

	SpAckCollectBatchTxFile struct {
		BatchID   string `json:"batch_id"`
		AssetName string `json:"asset_name"`
		UnSignTx  string `json:"un_sign_tx"`
		SpAckQueryParam
	}

	SpReqUploadFile struct {
		UUID string `json:"uuid"`
	}

	SpAckUploadFile struct {
		UUID   string `json:"uuid"`
		Status int    `json:"status"`
	}

	SpReqAddressBatch struct {
		BatchID       string `json:"batch_id"`
		AssetName     string `json:"asset_name"`
		MinCreateTime int64  `json:"min_create_time"`
		MaxCreateTime int64  `json:"max_create_time"`
		SpReqQueryParam
	}

	SpAckAddressBatch struct {
		BatchID    string `json:"batch_id"`
		AssetName  string `json:"asset_name"`
		Total      int    `json:"total"`
		Complete   int    `json:"complete"`
		CreateTime int64  `json:"create_time"`
	}

	SpAckAddressBatchList struct {
		Data []SpAckAddressBatch `json:"data" doc:"离线地址批次"`
		SpAckQueryParam
	}

	SpReqOfflineAddress struct {
		BatchID           string `json:"batch_id"`
		AssetName         string `json:"asset_name"`
		Allocated         int    `json:"allocated"`
		MinCreateTime     int64  `json:"min_create_time"`
		MaxCreateTime     int64  `json:"max_create_time"`
		MinAllocationTime int64  `json:"min_allocation_time"`
		MaxAllocationTime int64  `json:"max_allocation_time"`
		SpReqQueryParam
	}

	SpAckOfflineAddress struct {
		AssetName      string `json:"asset_name"`
		Address        string `json:"address"`
		Allocated      int    `json:"allocated"`
		BatchID        string `json:"batch_id"`
		CreateTime     int64  `json:"create_time"`
		AllocationTime int64  `json:"allocation_time"`
	}

	SpAckOfflineAddressList struct {
		Data []SpAckOfflineAddress `json:"data" doc:"离线地址列表"`
		SpAckQueryParam
	}

	SpReqMonthlyBalance struct {
		UserName    string `json:"user_name"`
		AssetName   string `json:"asset_name"`
		MaxPeriod   string `json:"max_period" doc:"最大周期值"`
		MinPeriod   string `json:"min_period" doc:"最小周期值"`
		MaxRegTime  int64  `json:"max_reg_time"`
		MinRegTime  int64  `json:"min_reg_time"`
		MaxUSDPrice string `json:"max_usd_price"`
		MinUSDPrice string `json:"min_usd_price"`
		SpReqQueryParam
	}

	SpAckMonthlyBalance struct {
		UserKey      string `json:"user_key"`
		UserName     string `json:"user_name"`
		RegTime      int64  `json:"reg_time"`
		Period       string `json:"period"`
		AssetName    string `json:"asset_name"`
		Available    string `json:"available"`
		Frozen       string `json:"frozen"`
		USDPrice     string `json:"usd_price"`
		InterestRate string `json:"interest_rate"`
		Interest     string `json:"interest"`
		Time         int64  `json:"time"`
	}

	SpAckMonthlyBalanceList struct {
		Data []SpAckMonthlyBalance `json:"data" doc:"用户月资产利息列表"`
		SpAckQueryParam
	}

	SpReqWalletCmd struct {
		CMD string `json:"cmd"`
	}

	SpAckWalletCmd struct {
		Data string `json:"data"`
	}

	SpReqSetOrderMemo struct {
		OrderID   string `json:"order_id"`
		OrderMemo string `json:"order_memo"`
	}

	SpReqSetOrderPass struct {
		OrderID string `json:"order_id"`
		TagUser string `json:"tag_user"`
		TagPass int    `json:"tag_pass"`
	}
)
