package table

type History struct {
	Id        *int     `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt *int64   `json:"createdat,omitempty" gorm:"type:bigint(20)"`
	UpdatedAt *int64   `json:"updatedat,omitempty" gorm:"type:bigint(20)"`
	Day       *int64   `json:"day,omitempty" gorm:"type:bigint(20)" `
	DaySucc   *int     `json:"day_succ,omitempty" gorm:"type:int(11)" `
	DayFail   *int     `json:"day_fail,omitempty" gorm:"type:int(11)" `
	GroupId   *int     `json:"group_id,omitempty" gorm:"type:int(11)" `
	RateFail  *float32 `json:"rate_fail,omitempty" gorm:"type:float(11)"`
	Inform    *int     `json:"inform,omitempty" gorm:"type:int(11)"`
	Type      *int     `json:"type,omitempty" gorm:"type:int(11)"`
}

func (this *History) TableName() string {
	return "notify_template_history"
}
