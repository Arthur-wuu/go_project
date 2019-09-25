package rbac

import (
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/services/access"
	"sort"
	"strconv"
)

type (
	UserPertainAccess struct {
		UserId int64 `valid:"required" json:"user_id"`
	}
)

//获取用户所有权限
func (this *VerifyAccess) SearchUserPertainAccess(userAccess []*models.Access) []*models.Access {
	data := make([]*models.Access, 0)

	var parentIds []string
	for _, v := range userAccess {
		data = append(data, v)
		if v.ParentId > -1 {
			parentIds = append(Tools.Unique(parentIds), strconv.FormatInt(v.ParentId, 10))
		}
	}

	if len(parentIds) > 0 {
		data, _ = this.RecursiveQueryClass(parentIds, data)
	}

	result := new(access.AccessList).Tree(data)
	sort.Sort(access.AccessListSort(result))

	return result
}

func (this *VerifyAccess) RecursiveQueryClass(parentIds []string, data []*models.Access) ([]*models.Access, error) {
	var result []*models.Access

	err := models.DB.Where("id IN (?) AND valid = ?", parentIds, "0").Find(&result).Error
	if err != nil {
		return nil, err
	}

	var parentId []string
	for _, v := range result {
		data = append(data, v)
		if v.ParentId > -1 {
			parentId = append(Tools.Unique(parentId),
				strconv.FormatInt(v.ParentId, 10))
		}
	}

	if len(parentId) > 0 {
		data, err = this.RecursiveQueryClass(parentId, data)
	}

	return data, nil
}
