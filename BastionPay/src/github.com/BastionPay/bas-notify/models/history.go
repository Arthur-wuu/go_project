package models

import (
	"BastionPay/bas-notify/models/table"
	"BastionPay/bas-notify/db"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"BastionPay/bas-notify/config"
	. "BastionPay/bas-base/log/zap"
)

type(
	 TemplateHistoryAdd struct {
		Id        *int   `valid:"optional" json:"id,omitempty"`
		CreatedAt *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt *int64  `valid:"optional" json:"updatedat,omitempty"`
		Day       *int64  `valid:"optional" json:"day,omitempty"`
		DaySucc   *int    `valid:"optional" json:"day_succ,omitempty"`
		DayFail   *int    `valid:"optional" json:"day_fail,omitempty"`
		GroupId   *int    `valid:"required" json:"group_id,omitempty"`
		RateFail  *float32 `valid:"optional" json:"rate_fail,omitempty"`
		Inform    *int     `valid:"optional" json:"inform,omitempty"`
		Type      *int    `valid:"optional" json:"type,omitempty"`
	 }

	TemplateHistory struct {
		Id        *int   `valid:"required" json:"id,omitempty"`
		CreatedAt *int64  `valid:"optional" json:"createdat,omitempty"`
		UpdatedAt *int64  `valid:"optional" json:"updatedat,omitempty"`
		Day       *int64  `valid:"optional" json:"day,omitempty"`
		DaySucc   *int    `valid:"optional" json:"day_succ,omitempty"`
		DayFail   *int    `valid:"optional" json:"day_fail,omitempty"`
		GroupId   *int    `valid:"optional" json:"group_id,omitempty"`
		RateFail  *float32 `valid:"optional" json:"rate_fail,omitempty"`
		Inform    *int     `valid:"optional" json:"inform,omitempty"`
		Type      *int    `valid:"optional" json:"type,omitempty"`
	}

	 TemplateHistoryList struct{
	 	GroupId         int             `valid:"required" json:"groupid,omitempty"`
		Total_lines     int              `valid:"optional" json:"total_lines"`
		Page_index      int              `valid:"optional" json:"page_index"`
		Max_disp_lines  int              `valid:"optional" json:"max_disp_lines"`
	 }
)

func (this * TemplateHistoryAdd)Add() error {
	his := &table.History{
		Day: this.Day,
		DaySucc: this.DaySucc,
		DayFail: this.DayFail,
		GroupId: this.GroupId,
		RateFail: this.RateFail,
		Inform: this.Inform,
	}

	return db.GDbMgr.Get().Create(his).Error
}

func (this *TemplateHistory)Update() error {
	his := &table.History{
		Id:  this.Id,
		Day: this.Day,
		DaySucc: this.DaySucc,
		DayFail: this.DayFail,
		GroupId: this.GroupId,
		RateFail: this.RateFail,
		Inform: this.Inform,
	}
	return db.GDbMgr.Get().Model(his).Updates(his).Error
}

func (this *TemplateHistory) Count(groupid int) (int,error) {
	cc := 0
	err := db.GDbMgr.Get().Model(&table.History{}).Where("group_id = ?", groupid).Count(&cc).Error
	return cc, err
}

func (this *TemplateHistoryList) List() ([]*table.History ,error){
	his := &table.History{
		GroupId: &this.GroupId,
	}
	if this.Max_disp_lines < 1 || this.Max_disp_lines > 100 {
		this.Max_disp_lines = 50
	}
	Page_index := this.Max_disp_lines * (this.Page_index - 1)
	var list []*table.History
	err := db.GDbMgr.Get().Model(his).Where(his).Offset(Page_index).Limit(this.Max_disp_lines).Order("day desc").Find(&list).Error
	return list, err
}

func (this *TemplateHistory) IncrTemplateHistoryCountAndRate( genRateFail func(int, int)float32, firstOverRateFailThd func(oldhis, newhis *table.History, tp int)bool ) (bool, error) {
	tx := db.GDbMgr.Get().Begin()
	defer func() {
		if r := recover(); r != nil {
			ZapLog().Error("panic")
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return false, tx.Error
	}
	oldHis := new(table.History)
	err := tx.Model(&table.History{}).Where("group_id = ? and day = ?", this.GroupId, this.Day).Find(oldHis).Error
	if err == gorm.ErrRecordNotFound {
		newHis := &table.History{
			GroupId : this.GroupId,
			Day : this.Day,
			Type : this.Type,
			DaySucc : this.DaySucc,
			DayFail : this.DayFail,
			RateFail: new(float32),
			Inform: new(int),
		}
		*newHis.Inform = 0
		*newHis.RateFail = genRateFail(*this.DaySucc, *this.DayFail)
		err = tx.Create(newHis).Error //创建的时候不需要判断是否错误率超过阈值
		if err != nil {
			tx.Rollback()
			return false, err
		}
		return false, tx.Commit().Error
	}
	if err != nil {
		tx.Rollback()
		return false, err
	}
	newHis := &table.History{
		GroupId : this.GroupId,
		Day : this.Day,
		Type : this.Type,
		DaySucc : this.DaySucc,
		DayFail : this.DayFail,
		RateFail: new(float32),
	}
	*newHis.DaySucc += *oldHis.DaySucc
	*newHis.DayFail +=  *oldHis.DayFail
	*newHis.RateFail = genRateFail(*newHis.DaySucc, *newHis.DayFail)

	flag := firstOverRateFailThd(oldHis, newHis, *this.Type)
	if flag {
		newHis.Inform = new(int)
		*newHis.Inform = 1
	}
	err = tx.Model(&table.History{}).Where("group_id = ? and day = ?", this.GroupId, this.Day).Update(newHis).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}
	return flag, tx.Commit().Error
}



func  RecordHistory(gid int,tp, succ, fail int, flag bool) {
	defer PanicPrint()
	informAdminFlag, err := IncrHistoryCount(gid, tp, succ, fail)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("IncrTemplateHistoryCount err")
		return
	}
	if !informAdminFlag {
		return
	}
	if !flag {
		return
	}
	mon := &config.GConfig.Monitor
	ZapLog().Warn("Monitor start", zap.Int("groupid", gid))
	if len(mon.TmpGNameMail) != 0 {
		req := &EmailMsg{
			GroupName: &mon.TmpGNameMail,
			Params :make(map[string]interface{}),
			Lang: &mon.TmpLangMail,
		}
		req.Params["key1"] = gid
		errCode, errMsg := req.Send( false)
		if errCode != 0 {
			ZapLog().With(zap.Error(errMsg), zap.Int("errcode",errCode )).Error("monitor sendMail err")
		}
	}
	if len(mon.TmpGNameSms) != 0 {
		req := &SmsMsg{
			GroupName: &mon.TmpGNameSms,
			Params :make(map[string]interface{}),
			Lang: &mon.TmpLangSms,
		}
		req.Params["key1"] = gid
		errCode, errMsg := req.Send( false)
		if errCode != 0 {
			ZapLog().With(zap.Error(errMsg), zap.Int("errcode",errCode )).Error("monitor sendSms err")
		}
	}

}