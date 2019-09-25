package vipsys


type (
	// /v1/vipsys/level/search 返回值外层包Response
	LevelSearch struct {
		Level     int      `valid:"required" json:"level,omitempty"`
		ModName   []string    `valid:"optional" json:"mod_name,omitempty"` //业务模块按照需要查询
	}

	// /v1/vipsys/level/searchs 返回值外层包Response
	LevelBatchSearch struct {
		Level     []int      `valid:"required" json:"level,omitempty"`
		ModName   []string    `valid:"optional" json:"mod_name,omitempty"` //业务模块按照需要查询
	}

	Level struct {
		//Id          *int       `json:"id,omitempty"`
		Level       *int       `json:"level,omitempty"`       //等级
		Name        *string    `json:"name,omitempty"`        //等级名称
		Mod        []*Mod    `json:"mod,omitempty" `       //规则列表
		Valid       *int       `json:"valid,omitempty"`       //1有效
	}

	Mod struct {
		//Id        *int      `json:"id,omitempty"       `
		//LevelId     *int       `json:"level_id,omitempty" `
		//Type       *int       `json:"type,omitempty" `
		Name       *string    `json:"name,omitempty"`        //模块名称
		RuleList    []*Rule   `json:"rule_list,omitempty" `
		Valid       *int       `json:"valid,omitempty"`
	}

	Rule struct {
		//Id        *int      `json:"id,omitempty"  `
		//ModId     *int       `json:"mod_id,omitempty"`
		Type       *int       `json:"type,omitempty"`      //分类
		Key       *string    `json:"key,omitempty"`        //
		Value       *string    `json:"value,omitempty"`
		Valid       *int       `json:"valid,omitempty"`
	}
)
