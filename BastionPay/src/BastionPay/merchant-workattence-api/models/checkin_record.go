package models

import (
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/db"
	"github.com/jinzhu/gorm"
)

type CheckinRecord struct {
	Id             *int    `json:"id,omitempty"               gorm:"column:id;primary_key;AUTO_INCREMENT:1;not null"`               //加上type:int(11)后AUTO_INCREMENT无效
	StaffId        *string `json:"user_id,omitempty"          gorm:"column:user_id;type:varchar(50);index:idx_user_id"`             //考勤机员工id
	CheckinAt      *int64  `json:"user_check_time,omitempty"  gorm:"column:user_check_time;type:int(11);index:idx_user_check_time"` //考勤打卡时间
	WorkDate       *int64  `json:"work_date,omitempty"        gorm:"column:work_date;type:int(11)"`                                 //工作日
	CorpId         *string `json:"corp_id,omitempty"          gorm:"column:corp_id;type:varchar(100)"`                              //企业ID
	CheckType      *string `json:"check_type,omitempty"       gorm:"column:check_type;type:varchar(10)"`                            //考勤类型，OnDuty：上班，OffDuty：下班
	SourceType     *string `json:"source_type,omitempty"      gorm:"column:source_type;type:varchar(15)"`                           //数据来源，ATM：考勤机;BEACON：IBeacon;DING_ATM：钉钉考勤机; USER：用户打卡; BOSS：老板改签; APPROVE：审批系统; SYSTEM：考勤系统; AUTO_CHECK：自动打卡
	TimeResult     *string `json:"time_result,omitempty"      gorm:"column:time_result;type:varchar(15)"`                           //时间结果， Normal：正常; Early：早退; Late：迟到; SeriousLate：严重迟到； Absenteeism：旷工迟到； NotSigned：未打卡
	UserAddress    *string `json:"user_address,omitempty"     gorm:"column:user_address;type:varchar(100)"`                         //用户打卡地址
	LocationMethod *string `json:"location_method,omitempty"  gorm:"column:location_method;type:varchar(15)"`                       //定位方法
	DeviceId       *string `json:"device_id,omitempty"        gorm:"column:device_id;type:varchar(100)"`                            //设备id
	IsLegal        *string `json:"is_legal,omitempty"         gorm:"column:is_legal;type:varchar(3)"`                               //是否合法，当timeResult和locationResult都为Normal时，该值为Y；否则为N
	LocationResult *string `json:"location_result,omitempty"  gorm:"column:location_result;type:varchar(15)"`                       //位置结果， Normal：范围内 Outside：范围外，外勤打卡时为这个值
}

func (this *CheckinRecord) TableName() string {
	return "workattence_checkin_record"
}

func (this *CheckinRecord) Add(p *api.PushCheckin) (*CheckinRecord, error) {
	tm := common.New().DatetimeToTimestamp(*p.Time)
	rcd := &CheckinRecord{
		StaffId:   p.Ccid,
		CheckinAt: &tm,
	}

	err := db.GDbMgr.Get().Create(rcd).Error

	if err != nil {
		return nil, err
	}

	return rcd, nil
}

func (this *CheckinRecord) AddNew(p *api.Recordresult) (*CheckinRecord, error) {
	ct := p.UserCheckTime / 1000
	wd := p.WorkDate / 1000
	rcd := &CheckinRecord{
		StaffId:        &p.UserId,
		CheckinAt:      &ct,
		WorkDate:       &wd,
		CorpId:         &p.CorpId,
		CheckType:      &p.CheckType,
		SourceType:     &p.SourceType,
		TimeResult:     &p.TimeResult,
		UserAddress:    &p.UserAddress,
		LocationMethod: &p.LocationMethod,
		LocationResult: &p.LocationResult,
		DeviceId:       &p.DeviceId,
		IsLegal:        &p.IsLegal,
	}

	err := db.GDbMgr.Get().Create(rcd).Error

	if err != nil {
		return nil, err
	}

	return rcd, nil
}

func (this *CheckinRecord) AddNewNoRecord(p *api.Recordresult) error {
	ct := p.UserCheckTime / 1000
	wd := p.WorkDate / 1000
	rcd := &CheckinRecord{
		StaffId:        &p.UserId,
		CheckinAt:      &ct,
		WorkDate:       &wd,
		CorpId:         &p.CorpId,
		CheckType:      &p.CheckType,
		SourceType:     &p.SourceType,
		TimeResult:     &p.TimeResult,
		UserAddress:    &p.UserAddress,
		LocationMethod: &p.LocationMethod,
		LocationResult: &p.LocationResult,
		DeviceId:       &p.DeviceId,
		IsLegal:        &p.IsLegal,
	}

	err := db.GDbMgr.Get().Create(rcd).Error

	if err != nil {
		return err
	}

	return nil
}

func (this *CheckinRecord) GetMaxCheckinAt() (*int64, error) {
	var count int64
	crd := new(CheckinRecord)
	err := db.GDbMgr.Get().Model(crd).Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return &count, nil
	}

	err = db.GDbMgr.Get().Select("max(user_check_time) user_check_time").Find(&crd).Error

	if err != nil {
		return nil, err
	}

	return crd.CheckinAt, nil

}

func (this *CheckinRecord) GetValidWorkOvertimeList(corpId string) ([]*CheckinRecord, error) {
	var crs []*CheckinRecord
	DayBeginTimestamp := common.New().DayBeginTimestamp() + 86400 // (DayBeginTimestamp - 14400) means 20:00

	err := db.GDbMgr.Get().
		Select("max(user_check_time) user_check_time, user_id, max(id) id").
		Where("check_type = ? and user_check_time >= ? and user_check_time <= ? and is_legal=? and corp_id=?", "OffDuty", DayBeginTimestamp-14400, DayBeginTimestamp-1, "Y", corpId).
		Group("user_id").
		Find(&crs).Error

	if err != nil {
		return nil, err
	}

	return crs, nil
}

func (this *CheckinRecord) GetUsersEarliestCheckinRecordList(uIds []string) ([]*CheckinRecord, error) {
	var crs []*CheckinRecord
	DayBeginTimestamp := common.New().DayBeginTimestamp() + 86400

	//39600 means 11 * 3600
	err := db.GDbMgr.Get().
		Select("min(user_check_time) user_check_time, user_id").
		Where("check_type = ? and user_check_time >= ? and user_check_time <= ? and user_id in (?) and is_legal=?", "OnDuty", DayBeginTimestamp-86400, DayBeginTimestamp-39600, uIds, "Y").
		Group("user_id").
		Find(&crs).Error

	if err != nil {
		return nil, err
	}

	return crs, nil
}

func (this *CheckinRecord) GetRecordById(Id int) (*CheckinRecord, error) {
	crd := new(CheckinRecord)
	err := db.GDbMgr.Get().Where("id = ?", Id).Find(&crd).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return crd, nil
}

func (this *CheckinRecord) GetEarliestCheckinRecord(uId string) (*CheckinRecord, error) {
	crd := new(CheckinRecord)
	err := db.GDbMgr.Get().Where("user_check_time >= ? and user_id = ?", common.New().DayBeginTimestamp(), uId).Order("user_check_time asc").First(&crd).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return crd, nil
}
