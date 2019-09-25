package models

import (
	"BastionPay/bas-notify/models/table"
	"BastionPay/bas-notify/db"
	"github.com/jinzhu/gorm"
)

type(
	TemplateAdd struct{
		Id        *int   `valid:"optional" json:"id,omitempty"`
		CreatedAt *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt *int64  `valid:"optional" json:"updatedat,omitempty"`
		Content   *string `valid:"optional" json:"content,omitempty"`
		Name     *string  `valid:"optional" json:"name,omitempty"`
		Title     *string `valid:"optional" json:"title,omitempty"`
		Type      *int    `valid:"optional" json:"type,omitempty" `
		Lang      *string `valid:"required" json:"lang,omitempty" `
		//Sign      *string `valid:"optional" json:"sign,omitempty" `
		GroupId   *int    `valid:"required" json:"groupid,omitempty" `
		//Alive     *int    `valid:"optional" json:"alive,omitempty" `
		//Alias     *string `valid:"optional" json:"alias,omitempty" `
		//DefaultRecipient *string `valid:"optional" json:"default_recipient,omitempty"`
		//SmsPlatform *int  `valid:"optional" json:"sms_platform,omitempty"`
		YuntongxinTempId *string  `valid:"optional" json:"ronglianyun_temp_id,omitempty"`
		DingQunName   *string  `valid:"optional" json:"ding_qun_name,omitempty"`
	}

	Template struct{
		Id        *int   `valid:"optional" json:"id,omitempty"`
		CreatedAt *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt *int64  `valid:"optional" json:"updatedat,omitempty"`
		Content   *string `valid:"optional" json:"content,omitempty"`
		Name     *string  `valid:"optional" json:"name,omitempty"`
		Title     *string `valid:"optional" json:"title,omitempty"`
		Type      *int    `valid:"optional" json:"type,omitempty" `
		Lang      *string `valid:"optional" json:"lang,omitempty" `
		//Sign      *string `valid:"optional" json:"sign,omitempty" `
		GroupId   *int    `valid:"optional" json:"groupid,omitempty" `
		//Alive     *int    `valid:"optional" json:"alive,omitempty" `
		//Alias     *string `valid:"optional" json:"alias,omitempty" `
		//DefaultRecipient *string `valid:"optional" json:"default_recipient,omitempty"`
		//SmsPlatform *int  `valid:"optional" json:"sms_platform,omitempty"`
		YuntongxinTempId *string  `valid:"optional" json:"ronglianyun_temp_id,omitempty"`
		DingQunName   *string  `valid:"optional" json:"ding_qun_name,omitempty"`
	}

	TemplateSave struct{
		Id        *int   `valid:"optional" json:"id,omitempty"`
		CreatedAt *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt *int64  `valid:"optional" json:"updatedat,omitempty"`
		Content   *string `valid:"optional" json:"content,omitempty"`
		Name     *string  `valid:"optional" json:"name,omitempty"`
		Title     *string `valid:"optional" json:"title,omitempty"`
		Type      *int    `valid:"optional" json:"type,omitempty" `
		Lang      *string `valid:"optional" json:"lang,omitempty" `
		//Sign      *string `valid:"optional" json:"sign,omitempty" `
		GroupId   *int    `valid:"optional" json:"groupid,omitempty" `
		//Alive     *int    `valid:"optional" json:"alive,omitempty" `
		//Alias     *string `valid:"optional" json:"alias,omitempty" `
		//DefaultRecipient *string `valid:"optional" json:"default_recipient,omitempty"`
		//SmsPlatform *int  `valid:"optional" json:"sms_platform,omitempty"`
		YuntongxinTempId *string  `valid:"optional" json:"ronglianyun_temp_id,omitempty"`
		DingQunName   *string  `valid:"optional" json:"ding_qun_name,omitempty"`
	}

	TemplateList struct{
		Id        *int   `valid:"optional" json:"id,omitempty"`
		CreatedAt *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt *int64  `valid:"optional" json:"updatedat,omitempty"`
		Content   *string `valid:"optional" json:"content,omitempty"`
		Name     *string  `valid:"optional" json:"name,omitempty"`
		Title     *string `valid:"optional" json:"title,omitempty"`
		Type      *int    `valid:"optional" json:"type,omitempty" `
		Lang      *string `valid:"optional" json:"lang,omitempty" `
		//Sign      *string `valid:"optional" json:"sign,omitempty" `
		GroupId   *int    `valid:"optional" json:"groupid,omitempty" `
		//Alive     *int    `valid:"optional" json:"alive,omitempty" `
		//Alias     *string `valid:"optional" json:"alias,omitempty" `
		//DefaultRecipient *string `valid:"optional" json:"default_recipient,omitempty"`
		//SmsPlatform *int  `valid:"optional" json:"sms_platform,omitempty"`
		YuntongxinTempId *string  `valid:"optional" json:"ronglianyun_temp_id,omitempty"`
		DingQunName   *string  `valid:"optional" json:"ding_qun_name,omitempty"`
	}
)

func (this *TemplateAdd) Add() error {
	temp := &table.Template{
		Id : this.Id,
		Content: this.Content,
		Name: this.Name,
		Title: this.Title,
		Type: this.Type,
		Lang: this.Lang,
		//Sign: this.Sign,
		GroupId: this.GroupId,
		//Alive: this.Alive,
		//Alias: this.Alias,
		//DefaultRecipient: this.DefaultRecipient,
		//SmsPlatform: this.SmsPlatform,
		YuntongxinTempId: this.YuntongxinTempId,
		DingQunName: this.DingQunName,
	}

	//if temp.Alive == nil {
	//	temp.Alive = new(int)
	//	*temp.Alive = 0
	//}
	//if temp.SmsPlatform == nil {
	//	temp.SmsPlatform = new(int)
	//	*temp.SmsPlatform = SMSPlatform_AWS
	//}
	//if temp.Content != nil {
	//	decodeBytes, err := base64.StdEncoding.DecodeString(*temp.Content)
	//	if err != nil {
	//		return err
	//	}
	//	content := string(decodeBytes)
	//	temp.Content = &content
	//}
	//if temp.DefaultRecipient != nil {
	//	*temp.DefaultRecipient = strings.Replace(*temp.DefaultRecipient, " ", "", len(*temp.DefaultRecipient))
	//
	//}
	return db.GDbMgr.Get().Create(temp).Error
}

func (this *TemplateAdd) Unique() (bool, error) {
	temp := &table.Template{
		GroupId:  this.GroupId,
		Lang:     this.Lang,
	}

	count := 0
	err := db.GDbMgr.Get().Model(temp).Where(temp).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, err
}

func (this *Template) Update() error {
	temp := &table.Template{
		Id : this.Id,
		Content: this.Content,
		Name: this.Name,
		Title: this.Title,
		Type: this.Type,
		Lang: this.Lang,
		//Sign: this.Sign,
		GroupId: this.GroupId,
		//Alive: this.Alive,
		//Alias: this.Alias,
		//DefaultRecipient: this.DefaultRecipient,
		//SmsPlatform: this.SmsPlatform,
		YuntongxinTempId: this.YuntongxinTempId,
		DingQunName: this.DingQunName,
	}
	//if temp.Content != nil {
	//	decodeBytes, err := base64.StdEncoding.DecodeString(*temp.Content)
	//	if err != nil {
	//		return err
	//	}
	//	content := string(decodeBytes)
	//	temp.Content = &content
	//}
	//if temp.DefaultRecipient != nil {
	//	*temp.DefaultRecipient = strings.Replace(*temp.DefaultRecipient, " ", "", len(*temp.DefaultRecipient))
	//
	//}

	return db.GDbMgr.Get().Model(temp).Updates(temp).Error
}

func (this *Template) TxSetDefaultRecipient(db * gorm.DB, gid int, recipient string) error {
	return db.Model(&table.Template{}).Where("group_id = ?", gid).Update("default_recipient", recipient).Error
}

func (this *Template) TxSetAllSmsPlatform(db * gorm.DB, smsPlatform int) error {
	return db.Model(&table.Template{}).Update("sms_platform", smsPlatform).Error
}

func (this *Template) TxDelByGid(db *gorm.DB, gid int) error {
	return db.Where("group_id = ?", gid).Delete(&table.Template{}).Error
}

func (this * Template) TxAliveByGid(db *gorm.DB, groupid, alive int) error {
	return db.Model(&table.Template{}).Where("group_id = ?", groupid).Update("alive", alive).Error
}

func (this * Template) GetsByGId(gid int) ([]*table.Template,error) {
	arr := make([]*table.Template, 0)
	err := db.GDbMgr.Get().Model(&table.Template{}).Where(`group_id = ? `, gid).Order("id desc").Order("id desc").Find(&arr).Error
	return arr, err
}

func (this * Template) GetLangByGId(gid int) ([]string,error) {
	arr := make([]string, 0)
	err := db.GDbMgr.Get().Model(&table.Template{}).Select("lang").Where(`group_id = ? and ( content is not null ) and trim(content) !='' `, gid).Order("id desc").Pluck("lang", &arr).Error
	mm := make(map[string]bool)
	newArr := make([]string, 0)
	for i:=0; i < len(arr); i++ {//去重.Table("notify_template")
		_,ok := mm[arr[i]]
		if ok {
			continue
		}
		newArr =append(newArr, arr[i])
		mm[arr[i]] = true
	}
	return newArr, err
}

func (this * Template) GetByGIdAndLang(gid int, lang string) (*table.Template,error) {
	temp := new(table.Template)
	err := db.GDbMgr.Get().Model(&table.Template{}).Where(`group_id = ? and lang = ? `, gid, lang).Last(temp).Error
	if err == gorm.ErrRecordNotFound {
		return nil,nil
	}
	return temp, err
}