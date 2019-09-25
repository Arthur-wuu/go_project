package models

import (
	"BastionPay/bas-notify/db"
	"BastionPay/bas-notify/models/table"
	"github.com/jinzhu/gorm"
	"strings"
)

type (
	TemplateGroupAdd struct {
		Id               *int    `valid:"optional" json:"id,omitempty"`
		Name             *string `valid:"required" json:"name,omitempty"`
		SubName          *string `valid:"optional" json:"sub_name,omitempty"`
		Detail           *string `valid:"optional" json:"detail,omitempty"`
		Alive            *int    `valid:"optional" json:"alive,omitempty"`
		Type             *int    `valid:"optional" json:"type,omitempty"`
		CreatedAt        *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt        *int64  `valid:"optional" json:"updatedat,omitempty"`
		Author           *string `valid:"optional" json:"author,omitempty"`
		Editor           *string `valid:"optional" json:"editor,omitempty"`
		DefaultRecipient *string `valid:"optional" json:"default_recipient,omitempty"`
		SmsPlatform      *int    `valid:"optional" json:"sms_platform,omitempty"`
	}

	TemplateGroup struct {
		Id               *int    `valid:"required" json:"id,omitempty"`
		Name             *string `valid:"optional" json:"name,omitempty"`
		SubName          *string `valid:"optional" json:"sub_name,omitempty"`
		Detail           *string `valid:"optional" json:"detail,omitempty"`
		Alive            *int    `valid:"optional" json:"alive,omitempty"`
		Type             *int    `valid:"optional" json:"type,omitempty"`
		CreatedAt        *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt        *int64  `valid:"optional" json:"updatedat,omitempty"`
		Author           *string `valid:"optional" json:"author,omitempty"`
		Editor           *string `valid:"optional" json:"editor,omitempty"`
		DefaultRecipient *string `valid:"optional" json:"default_recipient,omitempty"`
		SmsPlatform      *int    `valid:"optional" json:"sms_platform,omitempty"`
	}

	TemplateGroupList struct {
		Type           *int    `valid:"optional" json:"type,omitempty"`
		Name           *string `valid:"optional" json:"name,omitempty"`
		SubName        *string `valid:"optional" json:"sub_name,omitempty"`
		Total_lines    int     `valid:"optional" json:"total_lines,omitempty"`
		Page_index     int     `valid:"optional" json:"page_index,omitempty"`
		Max_disp_lines int     `valid:"optional" json:"max_disp_lines,omitempty"`
	}

	TemplateGroupCopy struct {
		Id      *int      `valid:"required" json:"id,omitempty"`
		SubName []*string `valid:"optional" json:"sub_name,omitempty"`
	}
)

func (this *TemplateGroupAdd) Add() (*table.TemplateGroup, error) {
	t := &table.TemplateGroup{
		Name:             this.Name,
		SubName:          this.SubName,
		Detail:           this.Detail,
		Alive:            this.Alive,
		Type:             this.Type,
		Author:           this.Author,
		Editor:           this.Editor,
		DefaultRecipient: this.DefaultRecipient,
		SmsPlatform:      this.SmsPlatform,
	}
	if t.Alive == nil {
		t.Alive = new(int)
		*t.Alive = 0
	}
	if t.SmsPlatform == nil {
		t.SmsPlatform = new(int)
		*t.SmsPlatform = SMSPlatform_AWS
	}
	if t.Name != nil {
		*t.Name = strings.Replace(*t.Name, " ", "", len(*t.Name))
	}
	if t.DefaultRecipient != nil {
		*t.DefaultRecipient = strings.Replace(*t.DefaultRecipient, " ", "", len(*t.DefaultRecipient))

	}
	err := db.GDbMgr.Get().Create(t).Error
	return t, err
}

func (this *TemplateGroupAdd) Unique() (bool, error) {
	mod := &table.TemplateGroup{
		Name:    this.Name,
		SubName: this.SubName,
		Type:    this.Type,
	}

	count := 0
	err := db.GDbMgr.Get().Model(mod).Where(mod).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, err
}

func (this *TemplateGroup) Unique() (bool, error) {
	mod := &table.TemplateGroup{
		Name:    this.Name,
		SubName: this.SubName,
		Type:    this.Type,
	}

	count := 0
	err := db.GDbMgr.Get().Model(mod).Where(mod).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, err
}

func (this *TemplateGroup) Update() error {
	group := &table.TemplateGroup{
		Id:               this.Id,
		Name:             this.Name,
		SubName:          this.SubName,
		Detail:           this.Detail,
		Alive:            this.Alive,
		Type:             this.Type,
		Author:           this.Author,
		Editor:           this.Editor,
		DefaultRecipient: this.DefaultRecipient,
		SmsPlatform:      this.SmsPlatform,
	}
	if group.DefaultRecipient != nil {
		*group.DefaultRecipient = strings.Replace(*group.DefaultRecipient, " ", "", -1)

	}
	return db.GDbMgr.Get().Model(group).Updates(group).Error
}

func (this *TemplateGroup) TxUpdate(db *gorm.DB) error {
	group := &table.TemplateGroup{
		Id:               this.Id,
		Name:             this.Name,
		SubName:          this.SubName,
		Detail:           this.Detail,
		Alive:            this.Alive,
		Type:             this.Type,
		Author:           this.Author,
		Editor:           this.Editor,
		DefaultRecipient: this.DefaultRecipient,
		SmsPlatform:      this.SmsPlatform,
	}
	if group.DefaultRecipient != nil {
		*group.DefaultRecipient = strings.Replace(*group.DefaultRecipient, " ", "", len(*group.DefaultRecipient))

	}
	return db.Model(group).Updates(group).Error
}

func (this *TemplateGroup) TxDel(db *gorm.DB) error {
	return db.Where("id = ?", *this.Id).Delete(&table.TemplateGroup{}).Error
}

//func (this * TemplateGroup) Alive( id, alive int) error {
//	return db.GDbMgr.Get().Model(&table.TemplateGroup{}).Where("id = ?", id).Update("alive", alive).Error
//}

func (this *TemplateGroup) GetByNameAndType(name string, tp int, subName *string) (*table.TemplateGroup, error) {
	t := new(table.TemplateGroup)
	err := db.GDbMgr.Get().Model(&table.TemplateGroup{}).Where("name = ? and type = ? and sub_name = ?", name, tp, subName).Last(t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return t, err
}

func (this *TemplateGroup) GetByid(id int) (*table.TemplateGroup, error) {
	t := new(table.TemplateGroup)
	err := db.GDbMgr.Get().Model(&table.TemplateGroup{}).Where("id = ?", id).Last(t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return t, err
}

func (this *TemplateGroup) GetAlive() (*table.TemplateGroup, error) {
	t := new(table.TemplateGroup)
	err := db.GDbMgr.Get().Model(&table.TemplateGroup{}).Where("alive = 1").Last(t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return t, err
}

func (this *TemplateGroupList) LikeCount() (int, error) {
	temp := &table.TemplateGroup{
		Type:    this.Type,
		Name:    this.Name,
		SubName: this.SubName,
	}
	newQb := db.GDbMgr.Get().Model(temp).Where(temp)
	count := 0
	if this.Name != nil {
		newQb = newQb.Where("name LIKE ?", "%"+*this.Name+"%")
	}

	err := newQb.Count(&count).Error
	return count, err
}

func (this *TemplateGroupList) LikeList() ([]*table.TemplateGroup, error) {
	temp := &table.TemplateGroup{
		Type: this.Type,
	}
	if this.Max_disp_lines < 1 || this.Max_disp_lines > 100 {
		this.Max_disp_lines = 50
	}
	Page_index := this.Max_disp_lines * (this.Page_index - 1)
	list := make([]*table.TemplateGroup, 0)
	newDb := db.GDbMgr.Get().Model(temp).Where(temp)
	if this.Name != nil {
		newDb = newDb.Where("name LIKE ?", "%"+*this.Name+"%")

	}

	err := newDb.Offset(Page_index).Limit(this.Max_disp_lines).Order("alive desc").Order("id desc").Find(&list).Error
	return list, err
}

func (this *TemplateGroupList) List() ([]*table.TemplateGroup, error) {
	temp := &table.TemplateGroup{
		Type:    this.Type,
		Name:    this.Name,
		SubName: this.SubName,
	}
	if this.Max_disp_lines < 1 || this.Max_disp_lines > 100 {
		this.Max_disp_lines = 50
	}
	Page_index := this.Max_disp_lines * (this.Page_index - 1)
	list := make([]*table.TemplateGroup, 0)
	newDb := db.GDbMgr.Get().Model(temp).Where(temp)
	if this.Name != nil {
		newDb = newDb.Where("name LIKE ?", "%"+*this.Name+"%")
	}

	err := newDb.Offset(Page_index).Limit(this.Max_disp_lines).Order("alive desc").Order("id desc").Find(&list).Error
	return list, err
}

func (this *TemplateGroup) SetAllSmsPlatform(smsPlatform int) error {
	return db.GDbMgr.Get().Model(&table.Template{}).Update("sms_platform", smsPlatform).Error
}
