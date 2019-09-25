package models

type (
	Access struct {
		Id       int64  `json:"id" gorm:"column:id"`
		ParentId int64  `json:"parent_id" gorm:"column:parent_id"`
		Name     string `json:"name" gorm:"column:name"`
		Uri      string `json:"uri" gorm:"column:uri"`
		IsAuth   int    `json:"is_auth" gorm:"column:is_auth"`
		IsMenu   int    `json:"is_menu" gorm:"column:is_menu"`
		Sort     int64  `json:"sort" gorm:"column:sort"`
		Model
		Children []*Access `json:"children" orm:"-"`
	}
)

func (this *Access) TableName() string {
	return "t_access"
}
