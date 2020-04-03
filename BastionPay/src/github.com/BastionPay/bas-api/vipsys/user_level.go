package vipsys

type (
	// /v1/vipsys/user-level/search 返回值外层包Response
	UserLevelSearch struct {
		UserKey string `valid:"optional" json:"user_key"`
	}

	// /v1/vipsys/user-level/searchs 返回值外层包Response
	UserLevelBatchSearch struct {
		UserKey []string `valid:"optional" json:"user_key"`
	}

	UserLevel struct {
		Id       *int    `json:"id,omitempty"`
		UserKey  *string `json:"user_key,omitempty"`
		ExpireAt *int64  `json:"expire_at,omitempty"`
		Level    *int    `json:"level,omitempty"`
		Valid    *int    `json:"valid,omitempty"`
	}
)
