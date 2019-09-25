package models

type (
	UserLoginLog struct {
		Id        int64  `json:"id" gorm:"column:id"`
		UserId    int64  `json:"user_id" gorm:"column:user_id"`
		Ip        string `json:"ip" gorm:"column:ip"`
		CreatedAt string `json:"created_at" gorm:"column:created_at"`
	}
)

func (this *UserLoginLog) TableName() string {
	return "t_user_login_log"
}
