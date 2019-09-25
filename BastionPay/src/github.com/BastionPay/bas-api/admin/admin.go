package admin

//用户详细信息
type UserDetailInfo struct {
	Id               uint   `json:"id"`
	CreatedAt        int64  `json:"created_at"`
	Phone            string `json:"phone"`
	CountryCode      string `json:"country_code"`
	Email            string `json:"email"`
	BindGa           bool   `json:"bind_ga"`
	VipLevel         uint8  `json:"vip_level"`
	Citizenship      string `json:"citizenship"`
	Language         string `json:"language"`
	Timezone         string `json:"timezone"`
	RegistrationType string `json:"registration_type"`
	UserKey          string `json:"user_key"`
	CompanyName      string `json:"company_name"`
}

func (this *UserDetailInfo) ToMap() map[string]interface{} {
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
	m["user_key"] = this.UserKey
	m["company_name"] = this.CompanyName

	return m
}

type AdminResponse struct {
	//	ctx iris.Context
	// response status
	Status struct {
		// response code
		Code int `json:"code"`
		// response msg
		Msg string `json:"msg"`
	} `json:"status"`
	// response result
	Result interface{} `json:"result"`
}
