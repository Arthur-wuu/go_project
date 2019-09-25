package backend

const (
	AUDITE_Status_Default = 0
	AUDITE_Status_Blank   = 0 //空
	AUDITE_Status_Doing   = 1 //审核中
	AUDITE_Status_Pass    = 2 //通过
	AUDITE_Status_Deny    = 3 //拒绝
)

// 账号注册-输入--register
type ReqUserRegister struct {
	UserClass   int    `json:"user_class" doc:"用户类型，0:普通用户 1:热钱包; 2:管理员"`
	Level       int    `json:"level" doc:"级别，0：用户，100：普通管理员，200：创世管理员"`
	IsFrozen    int    `json:"is_frozen" doc:"用户冻结状态，0: 正常；1：冻结状态，默认是0"`
	UserName    string `json:"user_name" doc:"用户名称"`
	UserMobile  string `json:"user_mobile" doc:"用户电话"`
	UserEmail   string `json:"user_email" doc:"用户邮箱"`
	CountryCode string `json:"country_code" doc:"国家码"`
	Language    string `json:"language" doc:"语种"`
}

// 账号注册-输出
type AckUserRegister struct {
	UserKey string `json:"user_key" doc:"用户唯一标示"`
}

// 修改公钥和回调地址-输入--update profile
type ReqUserUpdateProfile struct {
	PublicKey   string `json:"public_key" doc:"用户公钥"`
	SourceIP    string `json:"source_ip" doc:"用户源IP，用逗号(,)隔开"`
	CallbackUrl string `json:"callback_url" doc:"用户回调"`
}

// 修改公钥和回调地址-输出
type AckUserUpdateProfile struct {
	ServerPublicKey string `json:"server_public_key" doc:"BastionPay公钥"`
}

// 获取公钥和回调地址-输入--read profile
type ReqUserReadProfile struct {
}

// 获取公钥和回调地址-输出
type AckUserReadProfile struct {
	UserKey         string `json:"user_key" doc:"用户唯一标示"`
	PublicKey       string `json:"public_key" doc:"用户公钥"`
	SourceIP        string `json:"source_ip" doc:"用户源IP"`
	CallbackUrl     string `json:"callback_url" doc:"用户回调"`
	ServerPublicKey string `json:"server_public_key" doc:"BastionPay公钥"`
}

// 用户列表-输入--list
// 用户基本资料查询
type UserCondition struct {
	Id         int    `json:"id,omitempty" doc:"用户ID"`
	UserName   string `json:"user_name,omitempty" doc:"用户名称"`
	UserMobile string `json:"user_mobile,omitempty" doc:"用户电话"`
	UserEmail  string `json:"user_email,omitempty" doc:"用户邮箱"`
	UserKey    string `json:"user_key,omitempty" doc:"用户唯一标示"`
	UserClass  int    `json:"user_class,omitempty" doc:"用户类型"`
	Level      int    `json:"level,omitempty" doc:"级别"`
	IsFrozen   int    `json:"is_frozen,omitempty" doc:"用户冻结状态，0: 正常；1：冻结状态，默认是0"`
}
type ReqUserList struct {
	TotalLines   int `json:"total_lines" doc:"总数,0：表示首次查询"`
	PageIndex    int `json:"page_index" doc:"页索引,1开始"`
	MaxDispLines int `json:"max_disp_lines" doc:"页最大数，100以下"`

	Condition UserCondition `json:"condition" doc:"条件查询"`
}

// 用户列表-输出
// 用户基本资料
type UserBasic struct {
	Id               int    `json:"id" doc:"用户ID"`
	UserName         string `json:"user_name" doc:"用户名称"`
	UserMobile       string `json:"user_mobile" doc:"用户电话"`
	UserEmail        string `json:"user_email" doc:"用户邮箱"`
	UserKey          string `json:"user_key" doc:"用户唯一标示"`
	UserClass        int    `json:"user_class" doc:"用户类型"`
	Level            int    `json:"level" doc:"级别"`
	IsFrozen         int    `json:"is_frozen" doc:"用户冻结状态，0: 正常；1：冻结状态，默认是0"`
	CreateTime       int64  `json:"create_time" doc:"用户注册时间"`
	UpdateTime       int64  `json:"update_time" doc:"用户更新时间"`
	AuditeStatus     uint   `json:"audite_status" doc:"用户审核状态"`
	AuditeInfo       string `json:"audite_info" doc:"用户审核信息"`
	TransferStatus   uint   `json:"transfer_status" doc:"用户转账状态"`
	PubKeyFlag       uint   `json:"pubkey_flag" doc:"用户更新公钥开关"`
	CountryCode      string `json:"country_code" doc:"国家码"`
	NickName         string `json:"nick_name" doc:"用户别名"`
	CustodianFeeRate string `json:"custodian_fee_rate,omitempty" doc:"托管费率"`
	ParentUserKey    string `json:"parent_user_key,omitempty" doc:"用户父唯一标示"`
	IsParent         int    `json:"is_parent" doc:"是否父用户"`
}
type AckUserList struct {
	Data []UserBasic `json:"data" doc:"用户列表"`

	TotalLines   int `json:"total_lines" doc:"总数"`
	PageIndex    int `json:"page_index" doc:"页索引"`
	MaxDispLines int `json:"max_disp_lines" doc:"页最大数"`
}

// 设置冻结开关
type ReqFrozenUser struct {
	UserKey  string `json:"user_key" doc:"用户唯一标示"`
	IsFrozen int    `json:"is_frozen" doc:"用户冻结状态，0: 正常；1：冻结状态，默认是0"`
}

// 返回冻结开关
type AckFrozenUser struct {
	UserKey  string `json:"user_key" doc:"用户唯一标示"`
	IsFrozen int    `json:"is_frozen" doc:"用户冻结状态，0: 正常；1：冻结状态，默认是0"`
}

type ReqUserAuditeStatus struct {
	AuditeStatus uint   `json:"audite_status" doc:"用户审核状态, 0:空（默认）；1：审核中；2：通过；3：拒绝"`
	AuditeInfo   string `json:"audite_info" doc:"用户审核信息"`
}

type ResUserAuditeStatus struct {
	AuditeStatus uint    `json:"audite_status" doc:"用户审核状态, 0:空（默认）；1：审核中；2：通过；3：拒绝"`
	AuditeInfo   *string `json:"audite_info,omitempty" doc:"用户审核信息"`
}

type ReqUserTransferStatus struct {
	TransferStatus uint `json:"transfer_status"`
}

type ResUserAccountStatus struct {
	AuditeStatus   uint    `json:"audite_status" doc:"用户审核状态, 0:空（默认）；1：审核中；2：通过；3：拒绝"`
	AuditeInfo     *string `json:"audite_info,omitempty" doc:"用户审核信息"`
	TransferStatus uint    `json:"transfer_status" doc:"用户是否允许转账"`
	IsFrozen       int     `json:"is_frozen" doc:"用户冻结状态，0: 正常；1：冻结状态，默认是0"`
}

type ReqUserUpdatePublicKey struct {
	PublicKey string `json:"pubkey" doc:"用户公钥"`
}

type AckUserUpdatePublicKey struct {
	PublicKey string `json:"pubkey" doc:"用户公钥"`
}

type ReqUserUpdatePublicKeyFlag struct {
	PublicKeyFlag int `json:"pubkey_flag" doc:"用户公钥状态"`
}

type AckUserUpdatePublicKeyFlag struct {
	PublicKeyFlag int `json:"pubkey_flag" doc:"用户公钥状态"`
}

// 用户基本操作
type (
	UserOpCondition struct {
		Status int `json:"status,omitempty" doc:"状态，1：成功；2：失败"`
	}

	ReqUserOpeation struct {
		TotalLines   int `json:"total_lines" doc:"总数,0：表示首次查询"`
		PageIndex    int `json:"page_index" doc:"页索引,1开始"`
		MaxDispLines int `json:"max_disp_lines" doc:"页最大数，100以下"`

		Condition UserOpCondition `json:"condition" doc:"条件查询"`
	}

	UserOperation struct {
		Id         int    `json:"id" doc:"操作ID"`
		UserKey    string `json:"user_key" doc:"用户唯一标示"`
		CreateTime int64  `json:"create_time" doc:"用户操作时间"`

		Status      int    `json:"status" doc:"状态，1：成功；2：失败"`
		Function    string `json:"function" doc:"调用方法"`
		Description string `json:"description" doc:"描述"`
	}

	AckUserOpeationList struct {
		Data []UserOperation `json:"data" doc:"用户操作列表"`

		TotalLines   int `json:"total_lines" doc:"总数"`
		PageIndex    int `json:"page_index" doc:"页索引"`
		MaxDispLines int `json:"max_disp_lines" doc:"页最大数"`
	}
)

// 设置用户别名
type ReqUpdateNickName struct {
	UserKey  string `json:"user_key" doc:"用户唯一标示"`
	NickName string `json:"nick_name" doc:"用户新别名"`
}

// 返回别名设置
type AckUpdateNickName struct {
	UserKey  string `json:"user_key" doc:"用户唯一标示"`
	NickName string `json:"nick_name" doc:"返回设置的用户别名"`
}

// 设置用户托管费率
type ReqUpdateCustodianFeeRate struct {
	UserKey          string `json:"user_key" doc:"用户唯一标示"`
	CustodianFeeRate string `json:"custodian_fee_rate" doc:"用户托管费率"`
}

// 获取用户托管费率
type ReqReadCustodianFeeRate struct {
	UserKey string `json:"user_key" doc:"用户唯一标示"`
}
type AckReadCustodianFeeRate struct {
	UserKey          string `json:"user_key" doc:"用户唯一标示"`
	CustodianFeeRate string `json:"custodian_fee_rate" doc:"用户托管费率"`
}

type UserBindInfo struct {
	Id       uint   `json:"id,omitempty" doc:"用户ID"`
	UserName string `json:"user_name,omitempty" doc:"用户名称"`
	UserKey  string `json:"user_key,omitempty" doc:"用户唯一标示"`
}

// 设置子账户
type ReqBindChildUsers struct {
	UserKey    string         `json:"user_key" doc:"用户唯一标示"`
	ChildUsers []UserBindInfo `json:"child_users" doc:"子账户列表"`
}

// 获取子账户
type ReqListChildUsers struct {
	UserKey string `json:"user_key" doc:"用户唯一标示"`
}
type AckListChildUsers struct {
	UserKey    string         `json:"user_key" doc:"用户唯一标示"`
	ChildUsers []UserBindInfo `json:"child_users" doc:"子账户列表"`
}

// 获取绑定信息
type ReqReadBindInfo struct {
	UserKey string `json:"user_key" doc:"用户唯一标示"`
}
type AckReadBindInfo struct {
	UserKey    string         `json:"user_key" doc:"用户唯一标示"`
	ChildUsers []UserBindInfo `json:"child_users" doc:"子账户列表"`
	WildUsers  []UserBindInfo `json:"wild_users" doc:"野子账户列表"`
}

// custodian info
const (
	MaxCustodianFeeRate = 0.05
)

type (
	// info
	UsersCustodianInfo struct {
		ID               uint   `json:"id" doc:"唯一ID"`
		CreateAt         string `json:"create_at" doc:"创建时间"`
		AssetName        string `json:"asset_name" doc:"币种"`
		CustodianFeeRate string `json:"custodian_fee_rate" doc:"费率，不超过0.05"`
	}

	// list
	ReqListUsersCustodianInfo struct {
		UserKey string `json:"user_key" doc:"用户唯一标示"`
	}
	AckListUsersCustodianInfo struct {
		UserKey string               `json:"user_key" doc:"用户唯一标示"`
		List    []UsersCustodianInfo `json:"list" doc:"托管信息列表"`
	}

	// add, update
	ReqAddUpdateUsersCustodianInfo struct {
		UserKey          string `json:"user_key" doc:"用户唯一标示"`
		AssetName        string `json:"asset_name" doc:"币种"`
		CustodianFeeRate string `json:"custodian_fee_rate" doc:"费率，不超过0.05"`
	}
	// del
	ReqDelUsersCustodianInfo struct {
		UserKey   string `json:"user_key" doc:"用户唯一标示"`
		AssetName string `json:"asset_name" doc:"币种"`
	}
)
