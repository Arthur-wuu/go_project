package rbac

import (
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/services/access"
	"strings"
	"sync"
)

var (
	rw      sync.RWMutex
	UriList = make(map[int64]map[string]bool)
)

type (
	VerifyAccess struct {
		UserId         int64
		IgnoreUri      []string
		UserAccessList []*models.Access
	}

	IUserRole interface {
		GetUserRoleList() ([]*models.UserRole, error)
		ColumnUserRoleKey(list []*models.UserRole) []string
	}

	IRoleAccess interface {
		GetRoleAccessList(roleIds string) ([]*models.RoleAccess, error)
		ColumnRoleAccessKey(list []*models.RoleAccess) []string
	}

	IAccess interface {
		GetAccessByIds(ids []string) ([]*models.Access, error)
	}
)

func (v *VerifyAccess) GetUserAccessList() error {
	var iUserRole IUserRole
	var iRoleAccess IRoleAccess
	var iAccess IAccess

	iUserRole = &UserRole{
		UserId: v.UserId,
	}

	//获取用户所有角色
	userRole, err := iUserRole.GetUserRoleList()
	if err != nil {
		return err
	}

	roleIds := iUserRole.ColumnUserRoleKey(userRole)
	if len(roleIds) <= 0 {
		return fmt.Errorf("No user role found")
	}

	roleIdString := strings.Join(roleIds, ",")
	iRoleAccess = &RoleAccess{
		AccessId: roleIdString,
	}

	//获取角色所有权限
	roleAccess, err := iRoleAccess.GetRoleAccessList(roleIdString)
	if err != nil {
		return err
	}

	roleAccessIdList := iRoleAccess.ColumnRoleAccessKey(roleAccess)
	if len(roleAccessIdList) <= 0 {
		return err
	}

	//获取用户所有权限
	iAccess = new(access.AccessList)
	v.UserAccessList, err = iAccess.GetAccessByIds(roleAccessIdList)
	if err != nil {
		return err
	}

	v.GetColumnUserAccess(v.UserAccessList)

	return nil
}

func (v *VerifyAccess) GetColumnUserAccess(data []*models.Access) {
	list := make(map[string]bool)

	rw.Lock()
	defer rw.Unlock()
	for _, v := range data {
		list[v.Uri] = true
	}

	UriList[v.UserId] = list
}

func (v *VerifyAccess) GetUserList(userId int64, path string) bool {
	rw.RLock()
	defer rw.RUnlock()

	isBool, ok := UriList[userId][path]
	if ok && isBool {
		return true
	}

	return false
}

func (v *VerifyAccess) GetIgnoreUri() *VerifyAccess {
	v.IgnoreUri = new(access.AccessList).GetIgnoreUri()
	v.IgnoreUri = append(v.IgnoreUri, "/")

	return v
}

func (v *VerifyAccess) Ignore(uri string) bool {
	for _, val := range v.IgnoreUri {
		if strings.EqualFold(val, uri) {
			return true
		}
	}

	return false
}
