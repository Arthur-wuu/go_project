package rbac

import (
	"errors"
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/services/account"
	"strconv"
)

type (
	UserRole struct {
		UserId int64  `valid:"required" json:"user_id"`
		RoleId string `valid:"required" json:"role_id"`
	}

	UserRoleList struct {
		UserId int64 `valid:"required" json:"user_id"`
		RoleId int64 `valid:"optional" json:"role_id"`
		Page   int64 `valid:"optional" json:"page"`
		Size   int64 `valid:"optional" json:"size"`
	}

	SearchUserRole struct {
		models.UserRole
		Name string `json:"name" gorm:"column:name"`
	}
)

func (this *UserRole) SetUserRule() error {
	user, err := new(account.Account).GetUserById(this.UserId)
	if err != nil || user.Id == 0 {
		return errors.New("user does not exist")
	}

	result, err := new(Role).GetRoleById(this.RoleId)
	if err != nil {
		return err
	}

	if len(result) <= 0 {
		return errors.New("role does not exist")
	}

	metaRoleList, err := this.GetUserRoleList()
	if err != nil {
		return err
	}

	afterRoleIdList := this.ColumnRoleKey(result)
	metaRoleIdList := this.ColumnUserRoleKey(metaRoleList)
	newRoleIdList := Tools.ArrayDiff(afterRoleIdList, metaRoleIdList)

	delRoleIdList := Tools.ArrayDiff(metaRoleIdList, afterRoleIdList)
	err = models.DB.Where("user_id = ? AND role_id in (?)", this.UserId, delRoleIdList).Delete(&models.UserRole{}).Error
	if err != nil {
		return errors.New("remove user role failure")
	}

	var values []string
	for _, v := range newRoleIdList {
		value := fmt.Sprintf("%d, %s, '%s', '%s', %d",
			this.UserId,
			v,
			Tools.GetDateNowString(),
			Tools.GetDateNowString(), 0)
		values = append(values, value)
	}

	if len(values) > 0 {
		err = this.BatchAddUserRole(values)
	}

	if err != nil {
		return errors.New("no role can be set")
	}

	return nil
}

func (this *UserRole) BatchAddUserRole(values []string) error {

	fields := []string{"user_id", "role_id", "created_at", "updated_at", "valid"}

	f := models.UserRole{}
	return f.BatchInsert(f.TableName(), fields, values).Error
}

func (this *UserRole) GetUserRoleList() ([]*models.UserRole, error) {
	var userRoleList []*models.UserRole

	err := models.DB.Where("user_id = ? AND valid = ?", this.UserId, "0").Find(&userRoleList).Error

	return userRoleList, err
}

func (this *UserRole) ColumnUserRoleKey(list []*models.UserRole) []string {
	var roleIdList []string

	for _, v := range list {
		roleIdList = append(roleIdList, strconv.Itoa(int(v.RoleId)))
	}

	return roleIdList
}

func (this *UserRole) ColumnRoleKey(list []*models.Role) []string {
	var roleIdList []string

	for _, v := range list {
		roleIdList = append(roleIdList, strconv.Itoa(int(v.Id)))
	}

	return roleIdList
}

func (this *UserRoleList) SearchUserRole() (*common.Result, error) {
	var list []SearchUserRole
	var tatal int64
	query := models.DB.Table(new(models.UserRole).TableName()+" a").
		Select("a.*, b.name").
		Where("a.valid = ?", "0")

	if this.UserId > 0 {
		query = query.Where("a.user_id = ?", this.UserId)
	}

	if this.RoleId > 0 {
		query = query.Where("a.role_id = ?", this.RoleId)
	}

	query = query.Order("a.id DESC")
	query = query.Joins("left join " + new(models.Role).TableName() + " b on b.id = a.role_id").Count(&tatal)
	err := query.Limit(this.Size).Offset((this.Page - 1) * this.Size).Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return new(common.Result).PageResult(list, tatal, this.Page, this.Size)
}

func (this *UserRole) DeleteRole(RoleId int64) error {
	err := models.DB.Where("role_id = ?", RoleId).Delete(&models.UserRole{}).Error

	return err
}
