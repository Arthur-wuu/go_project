package table

type Template struct {
	Id        *int   `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt *int64  `json:"createdat,omitempty" gorm:"type:bigint(20)"`
	UpdatedAt *int64  `json:"updatedat,omitempty" gorm:"type:bigint(20)"`
	Name      *string `json:"name,omitempty" gorm:"type:varchar(50)"`
	Title     *string `json:"title,omitempty" gorm:"type:varchar(50)"`
	Type      *int    `json:"type,omitempty" gorm:"type:int(11)" `
	Content   *string `json:"content,omitempty" gorm:"type:text" `
	Lang      *string `json:"lang,omitempty" gorm:"type:varchar(20)" `
	//Sign      *string           //`json:"sign,omitempty" gorm:"type:varchar(50)" `
	GroupId   *int    `json:"groupid,omitempty" gorm:"type:int(11)" `
	//Alive     *int              //`json:"alive,omitempty" gorm:"type:int(11)" `
	//Alias     *string           //`json:"alias,omitempty" gorm:"type:varchar(30)" `
	//DefaultRecipient *string    //`json:"default_recipient,omitempty" gorm:"type:varchar(100)"`
	SmsPlatform *int            `json:"-" gorm:"-"`
	YuntongxinTempId *string   `json:"ronglianyun_temp_id,omitempty" gorm:"column:ronglianyun_temp_id;type:varchar(30)"`
	DingQunName *string   `json:"ding_qun_name,omitempty" gorm:"column:ding_qun_name;type:varchar(30)"`
}

func (this *Template) TableName() string {
	return "notify_template"
}
