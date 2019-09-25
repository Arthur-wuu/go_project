package table

type (
	UserHelp struct {
		Id            *uint64     `json:"id"              gorm:"AUTO_INCREMENT;primary_key;column:id"`
		//Level         *int       `json:"level"           gorm:"column:level;type:tinyint(4)"`
		CreatedAt     *int64     `json:"created_at"      gorm:"column:created_at;type:bigint(20)"`
		UpdatedAt     *int64     `json:"updated_at"      gorm:"column:updated_at;type:bigint(20)"`
		//DeletedAt     *int64     `json:"deleted_at"      gorm:"column:deleted_at;type:bigint(20)"`
		CountryCode   *string    `json:"country_code"    gorm:"column:country_code;type:varchar(20)"`
		Phone         *string    `json:"phone"           gorm:"column:phone;type:varchar(30)"`
		Email         *string    `json:"email"           gorm:"column:email;type:varchar(50)"`
		Website       *string    `json:"website"         gorm:"column:website;type:varchar(100)"`
		AppName       *string    `json:"app_name"        gorm:"column:app_name;type:varchar(50)"`
		Remark        *string    `json:"remark"          gorm:"column:remark;type:varchar(200)"`
		Status        *int       `json:"status"          gorm:"column:status;type:int(11);default:0"`
		Type          *int       `json:"type"            gorm:"column:type;type:int(11)"`
		Entrance      *string    `json:"entrance"        gorm:"column:entrance;type:varchar(30)"`
		Name          *string    `json:"name"            gorm:"column:name;type:varchar(30)"`
	}
)

func (this *UserHelp) TableName() string {
	return "user_help"
}
