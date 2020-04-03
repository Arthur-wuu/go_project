package vipsys

type (
	// /v1/vipsys/user-auth/search  返回值外层包Response
	UserRule struct {
		UserKey string   `valid:"required" json:"user_key"`
		ModName []string `valid:"optional" json:"mod_name,omitempty"` //业务模块按照需要查询
	}

	// /v1/vipsys/user-auth/searchs  返回值外层包Response
	UserRuleBatch struct {
		UserKey []string `valid:"required" json:"user_key"`
		ModName []string `valid:"optional" json:"mod_name,omitempty"` //业务模块按照需要查询
	}

	ResultUserRule struct {
		UserLevel *UserLevel `json:"user_level,omitempty" `
		Level     *Level     `json:"level,omitempty" `
	}

	//返回值是 level或者level数组
)
