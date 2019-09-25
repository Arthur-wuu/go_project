package db

type CodeTable struct {
	//ID          *uint   `gorm:"type:int(11)"`
	Code        uint    `gorm:"type:int(11)"`
	Symbol      string  `gorm:"primary_key" `
	Name        *string `gorm:"type:varchar(30)" `
	WebsiteSlug *string `gorm:"type:varchar(50)" `
	CreatedAt   *int64  `gorm:"type:bigint(20)"`
	UpdatedAt   *int64  `gorm:"type:bigint(20)"`
	Valid       *int    `gorm:"type:int(11)"`
}

type DbOptions struct {
	Host        string
	Port        string
	User        string
	Pass        string
	DbName      string
	MaxIdleConn int
	MaxOpenConn int
}
