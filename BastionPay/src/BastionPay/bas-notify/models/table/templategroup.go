package table

type TemplateGroup struct {
	Id               *int     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt        *int64   `gorm:"type:bigint(20)" json:"createdat,omitempty"`
	UpdatedAt        *int64   `gorm:"type:bigint(20)" json:"updatedat,omitempty"`
	Name             *string  `gorm:"type:varchar(50)" json:"name,omitempty"`
	SubName          *string  `gorm:"type:varchar(20)" json:"sub_name,omitempty"`
	Detail           *string  `gorm:"type:varchar(50)" json:"detail,omitempty"`
	Alive            *int     `gorm:"type:int(11)" json:"alive,omitempty"`
	Type             *int     `gorm:"type:int(11)" json:"type,omitempty"`
	Author           *string  `gorm:"type:varchar(30)" json:"author,omitempty"`
	Editor           *string  `gorm:"type:varchar(30)" json:"editor,omitempty"`
	DefaultRecipient *string  `gorm:"type:varchar(100)" json:"default_recipient,omitempty"`
	SmsPlatform      *int     `gorm:"type:int(11)" json:"sms_platform,omitempty"`
	Langs            []string `gorm:"-" json:"lang,omitempty"`
}

func (this *TemplateGroup) TableName() string {
	return "notify_template_group"
}
