package models

type Table struct {
	Valid     *int   `json:"valid,omitempty" gorm:"column:valid;type:tinyint(1)"`
	CreatedAt *int64 `json:"created_at,omitempty" gorm:"column:created_at;type:bigint(11)"`
	UpdatedAt *int64 `json:"updated_at,omitempty" gorm:"column:updated_at;type:bigint(11)"`
}

const (
	Const_Share_Mode_Random = 1
	Const_Share_Mode_Fixed  = 0
)
