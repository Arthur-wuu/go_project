package access

import (
	"errors"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/jinzhu/gorm"
	"sort"
	"strings"
)

type (
	Access struct {
		ParentId int64  `valid:"optional" json:"parent_id"`
		Name     string `valid:"required" json:"name"`
		Uri      string `valid:"optional,requri" json:"uri"`
		IsAuth   int    `valid:"optional,in(1|2)" json:"is_auth"`
		IsMenu   int    `valid:"optional,in(1|2)" json:"is_menu"`
		Sort     int64  `valid:"optional" json:"sort"`
	}

	AccessList struct {
		ParentId int64  `valid:"optional" json:"parent_id"`
		Name     string `valid:"optional" json:"name"`
		Page     int64  `valid:"optional" json:"page"`
		Size     int64  `valid:"optional" json:"size"`
	}

	AccessUpdateList struct {
		Id        int64    `valid:"required" json:"id"`
		ParentId  int64    `valid:"optional" json:"parent_id"`
		Name      string   `valid:"optional" json:"name"`
		Uri       string   `valid:"optional,requri" json:"uri"`
		Valid     int      `valid:"optional" json:"valid"`
		Sort      int64    `valid:"optional" json:"sort"`
		IsAuth    int      `valid:"optional,in(1|2)" json:"is_auth"`
		IsMenu    int      `valid:"optional,in(1|2)" json:"is_menu"`
		ConnBegin *gorm.DB `valid:"-" json:"-"`
	}

	AccessListSort []*models.Access
)

func (this *Access) AddAccess() (*models.Access, error) {
	this.Name = strings.ToLower(this.Name)

	model := &models.Access{
		ParentId: this.ParentId,
		Name:     this.Name,
		Uri:      this.Uri,
		IsAuth:   this.IsAuth,
		IsMenu:   this.IsMenu,
		Sort:     this.Sort,
	}

	err := models.DB.Where("name = ? AND valid = ?", this.Name, "0").First(model).Error
	if model.Id > 0 {
		switch model.Valid {
		case 0:
			return model, errors.New("add permission already exists")
		case 1:
			model.Valid = 0
			err := models.DB.Save(model).Error

			return model, err
		}
	}

	err = models.DB.Create(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (this *AccessList) GetIgnoreUri() []string {
	var ignoreList []*models.Access
	err := models.DB.Where("is_auth = ? AND valid = ?", "2", "0").Find(&ignoreList).Error
	if err != nil {
		return nil
	}

	var list []string
	for _, v := range ignoreList {
		list = append(list, v.Uri)
	}

	return list
}

func (this *AccessList) Search() (*common.Result, error) {
	var list []*models.Access
	var total int64

	query := models.DB.Where("valid = ?", "0")

	if this.Name != "" {
		query = query.Where("name = ?", this.Name)
	}

	if this.ParentId != 0 {
		query = query.Where("parent_id = ?", this.ParentId)
	}

	offset := (this.Page - 1) * this.Size
	err := query.Model(&models.Access{}).Count(&total).Limit(this.Size).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, err
	}

	result := this.Tree(list)
	sort.Sort(AccessListSort(result))

	return new(common.Result).PageResult(result, total, this.Page, this.Size)
}

func (this *AccessList) GetAccessByIdOnSet(ids []string) ([]*models.Access, error) {
	var access []*models.Access

	err := models.DB.Where("valid = ? AND id IN (?) AND is_auth = ? AND is_menu = ?", "0", ids, "1", "2").
		Order("sort ASC").Find(&access).Error

	return access, err
}

func (this *AccessList) GetAccessByIds(ids []string) ([]*models.Access, error) {
	var access []*models.Access

	err := models.DB.Where("valid = ? AND id IN (?)", "0", ids).Find(&access).Error

	return access, err
}

func (this *AccessList) GetAccessById(roleId int64) ([]*models.Access, error) {
	var access []*models.Access

	err := models.DB.Where("role_id = ? AND valid = ?", roleId, "0").Find(&access).Error
	return access, err
}

func (this *AccessUpdateList) DeleteAccess() (*models.Access, error) {
	access := &models.Access{}
	access.Valid = this.Valid

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").Update(&access).Error
	if err != nil || access == nil {
		return nil, err
	}

	return access, nil
}

func (this *AccessUpdateList) UpdateAccess() bool {

	access := &models.Access{
		Name:     this.Name,
		Uri:      this.Uri,
		Sort:     this.Sort,
		IsAuth:   this.IsAuth,
		IsMenu:   this.IsMenu,
		ParentId: this.ParentId,
	}

	access.Valid = this.Valid

	rows := this.ConnBegin.Model(&models.Access{}).
		Where("id = ? AND valid = ?", this.Id, "0").
		Update(&access).RowsAffected

	if rows > 0 {
		return true
	}

	return false
}

func (c AccessListSort) Len() int { return len(c) }

func (c AccessListSort) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c AccessListSort) Less(i, j int) bool { return c[i].Sort < c[j].Sort }
