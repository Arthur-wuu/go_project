package models

type (
	RoleAccess struct {
		Id       int64 `json:"id" gorm:"column:id"`
		RuleId   int64 `json:"rule_id" gorm:"column:rule_id"`
		AccessId int64 `json:"access_id" gorm:"column:access_id"`
		Model
	}
)

func (this *RoleAccess) TableName() string {
	return "t_role_access"
}
