package table

type (
	AppVersion struct {
		Id            *uint64    `json:"id"              gorm:"AUTO_INCREMENT;primary_key;column:id"`
		Name          *string      `json:"name"          gorm:"column:name";type:varchar(30)`
		LabelId        *uint64   `json:"label_id"        gorm:"column:label_id;type:bigint(20)"`
		Language      *string    `json:"language"        gorm:"column:language";type:varchar(20)`
		UpgradeAt     *int64     `json:"upgrade_at"       gorm:"column:upgrade_at;type:bigint(20)"`
		Version       *string    `json:"version"         gorm:"column:version;type:varchar(20)"`
		Url           *string    `json:"url"             gorm:"column:url;type:varchar(100)"`
		Instructions  *string    `json:"instructions"    gorm:"column:instructions;type:varchar(200)"`
		UpgradeMode   *int       `json:"upgrade_mode"    gorm:"column:upgrade_mode;type:int(11);default:0"`
		ShowMode      *int       `json:"show_mode"       gorm:"column:show_mode;type:int(11);default:0"`
		SysType       *string    `json:"sys_type"        gorm:"column:sys_type;type:varchar(20)"`
		CreatedAt     *int64     `json:"created_at"      gorm:"column:created_at;type:bigint(20)"`
		UpdatedAt     *int64     `json:"updated_at"      gorm:"column:updated_at;type:bigint(20)"`
	}

	AppWhiteLabel struct {
		Id            *uint64    `json:"id"              gorm:"AUTO_INCREMENT;primary_key;column:id"`
		//NameId        *int       `json:"name_id"        gorm:"column:name_id;type:int(11)"`
		Name          *string      `json:"name"          gorm:"column:name";type:varchar(30)`
		CreatedAt     *int64     `json:"created_at"      gorm:"column:created_at;type:bigint(20)"`
		UpdatedAt     *int64     `json:"updated_at"      gorm:"column:updated_at;type:bigint(20)"`
	}
)

func (this *AppVersion) TableName() string {
	return "app_version"
}

func (this *AppWhiteLabel) TableName() string {
	return "app_white_label"
}