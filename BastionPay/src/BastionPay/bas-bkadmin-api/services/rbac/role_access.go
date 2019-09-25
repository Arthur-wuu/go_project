package rbac

import (
	"errors"
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/services/access"
	"strconv"
	"strings"
)

var (
	Tools = common.New()
)

type (
	RoleAccess struct {
		RoleId   int64  `json:"role_id"`
		AccessId string `json:"access_id"`
		Page     int64  `valid:"optional" json:"page"`
		Size     int64  `valid:"optional" json:"size"`
	}

	RoleAccessList struct {
		RoleId int64 `valid:"required" json:"role_id"`
		Page   int64 `valid:"optional" json:"page"`
		Size   int64 `valid:"optional" json:"size"`
	}
)

func (this *RoleAccess) SetRoleAccess() error {
	result, err := new(Role).GetRoleById(strconv.Itoa(int(this.RoleId)))
	if err != nil || result == nil {
		return errors.New("role does not exist")
	}

	a := access.AccessList{}
	accessList, err := a.GetAccessByIdOnSet(strings.Split(this.AccessId, ","))
	if err != nil {
		return err
	}

	//if len(accessList) <= 0 {
	//	return errors.New("access does not exist")
	//}

	roleAccessList, err := this.GetRoleAccessList(strconv.Itoa(int(this.RoleId)))
	if err != nil {
		return err
	}

	aIds := this.ColumnAccessKey(accessList)
	roleAccessIds := this.ColumnRoleAccessKey(roleAccessList)
	newRoleAccessIds := Tools.ArrayDiff(aIds, roleAccessIds)

	//删除原来设置的权限
	delRoleAccessIds := Tools.ArrayDiff(roleAccessIds, aIds)
	err = models.DB.Where("role_id = ? AND access_id in (?)", this.RoleId, delRoleAccessIds).Delete(&models.RoleAccess{}).Error
	if err != nil {
		return errors.New("remove access role failure")
	}

	var values []string
	for _, v := range newRoleAccessIds {
		value := fmt.Sprintf("%d, %s, '%s', '%s', %d",
			this.RoleId,
			v,
			Tools.GetDateNowString(),
			Tools.GetDateNowString(), 0)
		values = append(values, value)
	}

	if len(values) > 0 {
		err = this.BatchAddRoleAccess(values)
	}

	if err != nil {
		return errors.New("no access role can be set")
	}

	return nil
}

func (this *RoleAccess) BatchAddRoleAccess(values []string) error {

	fields := []string{"role_id", "access_id", "created_at", "updated_at", "valid"}

	r := models.RoleAccess{}
	return r.BatchInsert(r.TableName(), fields, values).Error
}

func (this *RoleAccess) GetRoleAccessList(roleIds string) ([]*models.RoleAccess, error) {
	var RoleAccessList []*models.RoleAccess

	ids := strings.Split(roleIds, ",")
	err := models.DB.Where("role_id in (?) AND valid = ?", ids, "0").Find(&RoleAccessList).Error
	return RoleAccessList, err
}

func (this *RoleAccess) ColumnRoleAccessKey(list []*models.RoleAccess) []string {
	var roleIdList []string

	for _, v := range list {
		roleIdList = append(roleIdList, strconv.Itoa(int(v.AccessId)))
	}

	return roleIdList
}

func (this *RoleAccess) ColumnAccessKey(list []*models.Access) []string {
	var accessIdList []string

	for _, v := range list {
		accessIdList = append(accessIdList, strconv.Itoa(int(v.Id)))
	}

	return accessIdList
}

func (this *RoleAccessList) Search() (*common.Result, error) {
	var list []*models.RoleAccess
	query := models.DB.Where("valid = ?", "0")

	if this.RoleId > 0 {
		query = query.Where("role_id = ?", this.RoleId)
	}

	query = query.Order("id DESC")
	return new(common.Result).PageQuery(
		query,
		&models.RoleAccess{},
		&list,
		this.Page,
		this.Size,
		nil, "")
}

func (this *RoleAccess) DeleteAccessById(accessId int64) error {
	err := models.DB.Where("access_id = ?", accessId).Delete(&models.RoleAccess{}).Error

	return err
}

func (this *RoleAccess) DeleteRoleById(roleId int64) error {
	err := models.DB.Where("role_id = ?", roleId).Delete(&models.RoleAccess{}).Error

	return err
}
