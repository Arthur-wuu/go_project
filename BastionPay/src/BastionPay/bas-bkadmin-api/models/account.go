package models

type (
	Account struct {
		Id           int64  `json:"id" gorm:"column:id"`
		Name         string `json:"name" gorm:"column:name"`
		ActualName   string `json:"actual_name" gorm:"column:actual_name"`
		Password     string `json:"password" gorm:"column:password"`
		Department   string `json:"department" gorm:"column:department"`
		GoogleSecret string `json:"google_secret" gorm:"column:google_secret"`
		Email        string `json:"email" gorm:"column:email"`
		Mobile       string `json:"mobile" gorm:"column:mobile"`
		RoleId       int64  `json:"role_id" gorm:"column:role_id"`
		IsAdmin      int    `json:"is_admin" gorm:"column:is_admin"`
		Status       int    `json:"status" gorm:"column:status"`
		IsGauth      bool   `json:"is_gauth" gorm:"-"`
		Token        string `json:"token" gorm:"-"`
		Model
	}
)

func (this *Account) TableName() string {
	return "t_account"
}
