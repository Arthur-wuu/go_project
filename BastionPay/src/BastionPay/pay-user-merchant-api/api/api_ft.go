package api

type PageParams struct {
	Page  int `params:"page"`
	Limit int `params:"limit"`
}

type ResLoginLog struct {
	Id        *int    `json:"id"`
	Ip        *string `json:"ip"`
	Country   *string `json:"country"`
	City      *string `json:"city"`
	Device    *string `json:"device"`
	CreatedAt *int64  `json:"created_at"`
}

type ResOperationLog struct {
	Id        *int    `json:"id"`
	Operation *string `json:"operation"`
	Ip        *string `json:"ip"`
	Country   *string `json:"country"`
	City      *string `json:"city"`
	CreatedAt *int64  `json:"created_at"`
}

type ResVerificationSend struct {
	Id      *string `json:"id"`
	Captcha *string `json:"captcha,omitempty"`
}

type Verification struct {
	Id        string `json:"id"`
	Value     string `json:"value"`
	Recipient string `json:"recipient"`
}

type Login struct {
	Phone        string `valid:"required" json:"phone"`
	CountryCode  string `valid:"required" json:"country_code"`
	Password     string `valid:"required" json:"password"`
	CaptchaToken string `valid:"required" json:"captcha_token"`
}

type ResLogin struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
	Safe       bool   `json:"safe"`
}

type ResLoginQr struct {
	QrCode     string `json:"qr_code"`
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

type LoginQrCheck struct {
	Token string `valid:"required"  json:"token"`
}

type ResLoginCheck struct {
	Token string `json:"token"`
}

type LoginGa struct {
	GaToken string `json:"ga_token"`
}

type ResLoginGa struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

type Register struct {
	CompanyName string `json:"company_name"`
	Email       string
	Phone       string
	CountryCode string `json:"country_code"`
	Password    string
	Citizenship string
	Language    string
	Timezone    string
	Token       string
}

type ResRegister struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

type Exists struct {
	Username     string
	CaptchaToken string `json:"captcha_token"`
}

type RefreshToken struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

type ResGaGenerate struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
	Image  string `json:"image"`
}

type GaBind struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type GaUnBind struct {
	Value      string `json:"value"`
	EmailToken string `json:"email_token"`
	SmsToken   string `json:"sms_token"`
}

type PasswordModify struct {
	OldPassword string `json:"old_password"`
	Password    string
}

type PasswordInquire struct {
	CaptchaToken string `json:"captcha_token"`
	Username     string
}

type ResPasswordInquire struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	Ga          bool   `json:"ga"`
}

type PasswordReset struct {
	Username   string
	Password   string
	EmailToken string `json:"email_token"`
	SmsToken   string `json:"sms_token"`
	GaValue    string `json:"ga_value"`
}

type ResUserInfo struct {
	Id               *int64 `json:"id"`
	CreatedAt        *int64 `json:"created_at"`
	Phone            string `json:"phone"`
	CountryCode      string `json:"country_code"`
	Email            string `json:"email"`
	BindGa           bool   `json:"bind_ga"`
	VipLevel         uint8  `json:"vip_level"`
	Citizenship      string `json:"citizenship"`
	Language         string `json:"language"`
	Timezone         string `json:"timezone"`
	RegistrationType string `json:"registration_type"`
	CompanyName      string `json:"company_name"`
}

type ResUserInfoNoHide struct {
	Id          *int64 `json:"id"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	Email       string `json:"email"`
	BindGa      bool   `json:"bind_ga"`
}

type UserInfoSet struct {
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

type BindEmail struct {
	EmailToken string `json:"email_token"`
	Email      string
}

type BindPhone struct {
	SmsToken    string `json:"sms_token"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
}

type RebindPhone struct {
	EmailToken  string `json:"email_token"`
	GaToken     string `json:"ga_token"`
	OldSmsToken string `json:"old_sms_token"`
	NewSmsToken string `json:"new_sms_token"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
}

func (this *ResUserInfo) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = this.Id
	m["created_at"] = this.CreatedAt
	m["phone"] = this.Phone
	m["country_code"] = this.CountryCode
	m["email"] = this.Email
	m["bind_ga"] = this.BindGa
	m["vip_level"] = this.VipLevel
	m["citizenship"] = this.Citizenship
	m["language"] = this.Language
	m["timezone"] = this.Timezone
	m["registration_type"] = this.RegistrationType
	//m["user_key"] = this.UserKey
	m["company_name"] = this.CompanyName

	return m
}

type UserList struct {
	TotalLines   int `json:"total_lines" doc:"总数,0：表示首次查询"`
	PageIndex    int `json:"page_index" doc:"页索引,1开始"`
	MaxDispLines int `json:"max_disp_lines" doc:"页最大数，100以下"`

	Condition map[string]interface{} `json:"condition" doc:"条件查询"`
}

type MerchantAdd struct {
	MerchantId     *string `json:"-" ` //`valid:"required" json:"merchant_id,omitempty" `
	MerchantName   *string `valid:"optional" json:"merchant_name,omitempty" `
	NotifyUrl      *string `valid:"optional,url" json:"notify_url,omitempty"   `
	SignType       *string `valid:"optional" json:"sign_type,omitempty"    `
	SignKey        *string `valid:"optional" json:"sign_key,omitempty"     `
	PayeeId        *int64  `json:"-"   `
	LanguageType   *string `valid:"optional" json:"language_type,omitempty" `
	LegalCurrency  *string `valid:"optional" json:"legal_currency,omitempty" `
	Contact        *string `valid:"optional" json:"contact,omitempty" `
	ContactPhone   *string `valid:"optional" json:"contact_phone,omitempty"  `
	ContactEmail   *string `valid:"optional,email" json:"contact_email,omitempty" `
	Country        *string `valid:"optional" json:"country,omitempty" `
	CreateTime     *int64  `valid:"optional" json:"create_time,omitempty"  `
	LastUpdateTime *int64  `valid:"optional" json:"last_update_time,omitempty"`
}

type Merchant struct {
	ID             *int64  `valid:"required" json:"id,omitempty"`
	MerchantId     *string `valid:"optional" json:"merchant_id,omitempty" `
	MerchantName   *string `valid:"optional" json:"merchant_name,omitempty" `
	NotifyUrl      *string `valid:"optional,url" json:"notify_url,omitempty"   `
	SignType       *string `valid:"optional" json:"sign_type,omitempty"    `
	SignKey        *string `valid:"optional" json:"sign_key,omitempty"     `
	PayeeId        *int64  `json:"-"`
	LanguageType   *string `valid:"optional" json:"language_type,omitempty" `
	LegalCurrency  *string `valid:"optional" json:"legal_currency,omitempty" `
	Contact        *string `valid:"optional" json:"contact,omitempty" `
	ContactPhone   *string `valid:"optional" json:"contact_phone,omitempty"  `
	ContactEmail   *string `valid:"optional,email" json:"contact_email,omitempty" `
	Country        *string `valid:"optional" json:"country,omitempty" `
	CreateTime     *int64  `valid:"optional" json:"create_time,omitempty"  `
	LastUpdateTime *int64  `valid:"optional" json:"last_update_time,omitempty"`
}

type ResMerchant struct {
	ID             *int64  `json:"id,omitempty"`
	MerchantId     *string `json:"merchant_id,omitempty" `
	MerchantName   *string `json:"merchant_name,omitempty" `
	NotifyUrl      *string `json:"notify_url,omitempty"   `
	SignType       *string `json:"sign_type,omitempty"    `
	SignKey        *string `json:"sign_key,omitempty"     `
	PayeeId        *int64  `json:"payee_id,omitempty"   `
	LanguageType   *string `json:"language_type,omitempty" `
	LegalCurrency  *string `json:"legal_currency,omitempty" `
	Contact        *string `json:"contact,omitempty" `
	ContactPhone   *string `json:"contact_phone,omitempty"  `
	ContactEmail   *string `json:"contact_email,omitempty" `
	Country        *string `json:"country,omitempty" `
	CreateTime     *int64  `json:"create_time,omitempty"  `
	LastUpdateTime *int64  `json:"last_update_time,omitempty"`
}

//商户订单 列表，退单
type TradeList struct {
	PayeeId         *string `valid:"required" json:"payee_id,omitempty"`
	MerchantTradeNo *string `valid:"optional" json:"merchant_trade_no,omitempty"`
	TradeNo         *string `valid:"optional" json:"trade_no,omitempty"`
	Page            int64   `valid:"optional" json:"page,omitempty"`
	Size            int64   `valid:"optional" json:"size,omitempty"`
}

type RefundTradeList struct {
	MerchantId *string `valid:"required" json:"merchant_id,omitempty"`
	Page       *int64  `valid:"optional" json:"page,omitempty"`
	Size       *int64  `valid:"optional" json:"size,omitempty"`
}
