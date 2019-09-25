package table

type (
	//UserDetail struct {
//	//	Id                  *int64         `json:"ID"                   gorm:"AUTO_INCREMENT;column:ID;type:bigint(20)"`
//	//	AccId               *int           `json:"ACC_ID"               gorm:"column:ACCID;type:bigint(20)"`
//	//	UserId              *int64         `json:"USER_ID"              gorm:"column:USERID;type:bigint(20)"`
//	//	Assets              *string        `json:"ASSETS"               gorm:"column:ASSETS;type:varchar(10)"`
//	//	Amount              *float64       `json:"AMOUNT"               gorm:"column:AMOUNT;type:decimal(30,10)"`
//	//	AfterAmount         *float64       `json:"AFTER_AMOUNT"         gorm:"column:AFTER_AMOUNT;type:decimal(30,10)"`
//	//	Direction           *int32         `json:"DIRECTION"            gorm:"column:DIRECTION;type:tinyint(4)"`
//	//	Remark              *string        `json:"REMARK"               gorm:"column:REMARK;type:varchar(1000)"`
//	//	CreateTmie          *string        `json:"CREATE_TIME"          gorm:"column:CREATE_TIME;type:datetime"`
//	//	LastUpdateTmie      *string        `json:"LAST_UPDATE_TIME"     gorm:"column:LAST_UPDATE_TIME;type:datetime"`
//	//	OrderId             *int64         `json:"ORDER_ID"             gorm:"column:ORDER_ID;type:bigint(20)"`
//	//	FrozenAmount        *float64       `json:"FROZEN_AMOUNT"        gorm:"column:FROZEN_AMOUNT;type:decimal(30,10)"`
//	//	AfterFrozenAmount   *float64       `json:"AFTER_FROZEN_AMOUNT"  gorm:"column:AFTER_FROZEN_AMOUNT;type:decimal(30,10)"`
//	//}

	UserDetail struct {
		Id                  *int64         `json:"ID"                   gorm:"AUTO_INCREMENT;column:ID;type:bigint(20)"`
		UserId              *int64         `json:"USER_ID"              gorm:"column:USERID;type:bigint(20)"`
		Assets              *string        `json:"ASSETS"               gorm:"column:ASSETS;type:varchar(10)"`
		Balance              *float64       `json:"BALANCE"               gorm:"column:BALANCE;type:decimal(30,10)"`
		FrozenBalance         *float64       `json:"FROZEN_BALANCE"         gorm:"column:FROZEN_BALANCE;type:decimal(30,10)"`
		CreateTmie          *string        `json:"CREATE_TIME"          gorm:"column:CREATE_TIME;type:datetime"`
		LastUpdateTmie      *string        `json:"LAST_UPDATE_TIME"     gorm:"column:LAST_UPDATE_TIME;type:datetime"`
		}
)

func (this *UserDetail) TableName() string {
	return "USER_ACC"
}
