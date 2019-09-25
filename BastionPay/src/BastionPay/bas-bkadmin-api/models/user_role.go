package models

type UserRole struct {
	Id     int64 `json:"id" gorm:"column:id"`
	UserId int64 `json:"user_id" gorm:"column:user_id"`
	RoleId int64 `json:"role_id" gorm:"column:role_id"`
	Model
}

func (c *UserRole) TableName() string {
	return "t_user_role"
}
