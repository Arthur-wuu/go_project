package models

type Table struct {
	//Valid     *int    `json:"valid,omitempty" gorm:"column:valid"`
	CreatedAt *int64 `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *int64 `json:"updated_at,omitempty" gorm:"column:updated_at"`
	//Author    *string `json:"author,omitempty" gorm:"column:author;type:varchar(20)"`
	DeletedAt *int64 `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

const (
	EUM_OPEN_STATUS_OPEN  = 0
	EUM_OPEN_STATUS_CLOSE = 1

	//支付，转账
	EUM_TRANSFER_STATUS_INIT    = 0
	EUM_TRANSFER_STATUS_APPLY   = 0
	EUM_TRANSFER_STATUS_SUCCESS = 1
	EUM_TRANSFER_STATUS_FAIL    = 2
	EUM_TRANSFER_STATUS_CLOSE   = 3

	//退款
	EUM_REFUND_STATUS_APPLY   = 0
	EUM_REFUND_STATUS_SUCCESS = 1
	EUM_REFUND_STATUS_FAIL    = 2
	EUM_REFUND_STATUS_CLOSE   = 3

	//提币
	EUM_FUNDOUT_STATUS_APPLY   = 0
	EUM_FUNDOUT_STATUS_SUCCESS = 1
	EUM_FUNDOUT_STATUS_FAIL    = 2
	EUM_FUNDOUT_STATUS_CLOSE   = 3
)
