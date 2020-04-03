package api

type BkUserHelpList struct {
	CountryCode *string `valid:"optional" json:"country_code,omitempty" `
	Phone       *string `valid:"optional" json:"phone,omitempty"`
	Email       *string `valid:"optional" json:"email,omitempty"  `
	Name        *string `valid:"optional" json:"name,omitempty" `
	Remark      *string `valid:"optional" json:"remark,omitempty" `
	FeedBack    *string `valid:"optional" json:"feedback,omitempty" `
	Status      *int    `valid:"optional" json:"status,omitempty" `
	Page        int64   `valid:"optional" json:"page,omitempty" `
	Size        int64   `valid:"optional" json:"size,omitempty" `
}

//
type BkUserHelp struct {
	Id          *int    `valid:"required" json:"id,omitempty" `
	CountryCode *string `valid:"optional" json:"country_code,omitempty" `
	Phone       *string `valid:"optional" json:"phone,omitempty"`
	Email       *string `valid:"optional" json:"email,omitempty"  `
	Name        *string `valid:"optional" json:"name,omitempty" `
	Remark      *string `valid:"optional" json:"remark,omitempty" `
	FeedBack    *string `valid:"optional" json:"feedback,omitempty" `
	Status      *int    `valid:"optional" json:"status,omitempty" `
}

type BkListRequest struct {
	Message BkListRequestMessage `valid:"optional" json:"message,omitempty" `
}

type BkListRequestMessage struct {
	Condition BkUserHelpList `valid:"optional" json:"condition,omitempty" `
	Page      int64          `valid:"optional" json:"page,omitempty" `
	Size      int64          `valid:"optional" json:"size,omitempty" `
}

type BkLoginUserInfo struct {
	UserId        *int64  `valid:"required" json:"user_id,omitempty" `
	Thumbnail     *string `valid:"optional" json:"thumbnail,omitempty" `
	UserNick      *string `valid:"optional" json:"user_nick,omitempty" `
	Email         *string `valid:"optional" json:"email,omitempty" `
	Phone         *string `valid:"optional" json:"phone,omitempty" `
	PhoneDistrict *string `valid:"optional" json:"phone_district,omitempty" `
	LoginQrcode   *string `valid:"required" json:"login_qrcode,omitempty" `
	Status        *int    `valid:"-" json:"status,omitempty" `
}
