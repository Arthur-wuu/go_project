package rbac

import (
	"errors"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"strings"
)

type (
	Role struct {
		Id     int64  `valid:"optional" json:"id"`
		Name   string `valid:"required" json:"name"`
		Status int64  `valid:"optional" json:"status"`
	}

	RoleList struct {
		Id     int64  `valid:"optional" json:"id"`
		Name   string `valid:"optional" json:"name"`
		Status int64  `valid:"optional" json:"status"`
		Page   int64  `valid:"optional" json:"page"`
		Size   int64  `valid:"optional" json:"size"`
	}

	RoleDelete struct {
		Id int64 `valid:"required" json:"id"`
	}

	RoleUpdate struct {
		Id     int64  `valid:"required" json:"id"`
		Name   string `valid:"optional" json:"name"`
		Status int64  `valid:"optional" json:"status"`
	}
)

func (this *Role) AddRole() (*models.Role, error) {
	rule := &models.Role{
		Name:   this.Name,
		Status: this.Status,
	}

	result, _ := this.GetRoleByName()

	switch {
	case result.Id == 0:
		rule.Status = 1
		err := models.DB.Save(rule).Error
		return rule, err
	case result.Valid == 1:
		result.Valid = 0
		err := models.DB.Save(&result).Error
		return result, err
	case result.Valid == 0 && result.Id > 0:
		return nil, errors.New("the role already exists")
	}

	return nil, errors.New("add rbac failure")
}

func (this *Role) GetRoleByName() (*models.Role, error) {
	rule := &models.Role{}

	err := models.DB.Where("name = ? AND valid = ?", this.Name, "0").First(rule).Error

	return rule, err
}

func (this *Role) GetRoleById(ids string) ([]*models.Role, error) {
	var rules []*models.Role

	idList := strings.Split(ids, ",")
	err := models.DB.Where("valid = ? AND id IN (?)", "0", idList).Find(&rules).Error

	return rules, err
}

func (this *Role) GetRoleInfoById() (*models.Role, error) {
	var roles models.Role

	err := models.DB.Where("valid = ? AND id = ?", "0", this.Id).Find(&roles).Error

	return &roles, err
}

func (this *RoleList) RoleList() (*common.Result, error) {
	var list []*models.Role

	query := models.DB.Where("valid = ?", "0")
	if this.Name != "" {
		query = query.Where("name = ?", this.Name)
	}

	return new(common.Result).PageQuery(query, &models.Role{}, &list, this.Page, this.Size, nil, "")
}

func (this *RoleDelete) Delete() (*models.Role, error) {
	var role models.Role

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").First(&role).Error
	if err != nil {
		return nil, err
	}

	role.Valid = 1
	err = models.DB.Save(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (this *RoleUpdate) Update() (*models.Role, error) {
	var role models.Role

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").First(&role).Error
	if err != nil {
		return nil, err
	}

	role.Name = this.Name
	err = models.DB.Save(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (this *RoleUpdate) Disabled() (*models.Role, error) {
	var role models.Role

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").First(&role).Error
	if err != nil {
		return nil, err
	}

	role.Status = this.Status
	err = models.DB.Save(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}
