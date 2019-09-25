package models

import (
	"BastionPay/bas-admin-api/common"
	"fmt"
	"github.com/BastionPay/bas-api/admin"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/bugsnag/bugsnag-go/errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"strings"
)

func NewNoticeModel(conn *gorm.DB) *NoticeModel {
	// 增加company_name字段
	conn.AutoMigrate(&User{})
	return &NoticeModel{mConn: conn}
}

var NoticeListSelectElems = "id, created_at, updated_at,onlined_at, offlined_at, language, focus,race, title, abstract"
var NoticeInfoSelectElems = "id, created_at, updated_at,onlined_at, offlined_at, language, title, content"
var NoticeListSelectElemsFromInner = "id, created_at, updated_at, onlined_at, offlined_at, language,author, focus,race, title, abstract"
var NoticeInfoSelectElemsFromInner = "id, created_at, updated_at, onlined_at, offlined_at, language, author,focus,race, title, abstract, content"

type TimeDuration struct {
	StartCreatedAt  *int64
	EndCreatedAt    *int64
	StartUpdatedAt  *int64
	EndUpdatedAt    *int64
	StartOnlinedAt  *int64
	EndOnlinedAt    *int64
	StartOfflinedAt *int64
	EndOfflinedAt   *int64
	LineAndFlag     bool //上下线 与或关系, 默认与
}

func (this *TimeDuration) Empty() bool {
	if this.StartCreatedAt != nil {
		return false
	}
	if this.EndCreatedAt != nil {
		return false
	}
	if this.StartUpdatedAt != nil {
		return false
	}
	if this.EndUpdatedAt != nil {
		return false
	}
	if this.StartOnlinedAt != nil {
		return false
	}
	if this.EndOnlinedAt != nil {
		return false
	}
	if this.StartOfflinedAt != nil {
		return false
	}
	if this.EndOfflinedAt != nil {
		return false
	}
	return true
}

func (this *TimeDuration) OnOffLineAndFlag() bool {
	return this.LineAndFlag
}

func NewTimeDuration(StartCreatedAt, EndCreatedAt, StartUpdatedAt, EndUpdatedAt, StartOnlinedAt, EndOnlinedAt, StartOfflinedAt, EndOfflinedAt *int64, alive *int) *TimeDuration {
	dur := new(TimeDuration)
	dur.StartCreatedAt = StartCreatedAt
	dur.EndCreatedAt = EndCreatedAt
	dur.StartUpdatedAt = StartUpdatedAt
	dur.EndUpdatedAt = EndUpdatedAt
	dur.StartOnlinedAt = StartOnlinedAt
	dur.EndOnlinedAt = EndOnlinedAt
	dur.StartOfflinedAt = StartOfflinedAt
	dur.EndOfflinedAt = EndOfflinedAt
	dur.LineAndFlag = true

	if alive == nil {
		return dur
	}

	nowTime := common.NowTimestamp()
	switch *alive {
	case admin.STATUS_ALIVE_Online:
		if dur.EndOnlinedAt != nil {
			if *dur.EndOnlinedAt > nowTime {
				*dur.EndOnlinedAt = nowTime
			}
		} else {
			dur.EndOnlinedAt = new(int64)
			*dur.EndOnlinedAt = nowTime
		}
		if dur.StartOfflinedAt != nil {
			if *dur.StartOfflinedAt < nowTime {
				*dur.StartOfflinedAt = nowTime
			}
		} else {
			dur.StartOfflinedAt = new(int64)
			*dur.StartOfflinedAt = nowTime
		}
		break
	case admin.STATUS_ALIVE_PreOnline:
		if dur.StartOnlinedAt != nil {
			if *dur.StartOnlinedAt < nowTime {
				*dur.StartOnlinedAt = nowTime
			}
		} else {
			dur.StartOnlinedAt = new(int64)
			*dur.StartOnlinedAt = nowTime
		}
		break
	case admin.STATUS_Alive_Offline:
		if dur.StartOnlinedAt != nil {
			if *dur.StartOnlinedAt < nowTime {
				*dur.StartOnlinedAt = nowTime
			}
		} else {
			dur.StartOnlinedAt = new(int64)
			*dur.StartOnlinedAt = nowTime
		}
		dur.EndOfflinedAt = new(int64)
		*dur.EndOfflinedAt = nowTime
		dur.LineAndFlag = false
		break
	case admin.STATUS_Alive_AfterOffline:
		if dur.EndOfflinedAt != nil {
			if *dur.EndOfflinedAt > nowTime {
				*dur.EndOfflinedAt = nowTime
			}
		} else {
			dur.EndOfflinedAt = new(int64)
			*dur.EndOfflinedAt = nowTime
		}
	default:

	}

	return dur
}

type NoticeModel struct {
	mConn *gorm.DB
}

func (this *NoticeModel) AddNotice(info *NoticeInfo) bool {
	//return this.mConn.NewRecord(*infos)
	if err := this.mConn.Create(info).Error; err != nil {
		fmt.Println("Create err:", err.Error())
		return false
	}
	return true
}

//根据ID更新
func (this *NoticeModel) UpdateNotice(info *NoticeInfo) error {
	tx := this.mConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := this.mConn.Model(&NoticeInfo{}).Updates(*info).Error; err != nil {
		tx.Rollback()
		return err
	}
	if (info.OnlinedAt == nil) && (info.OfflinedAt == nil) {
		return nil
	}
	return nil //联合主键报错了
	readInfo := &NoticeRead{
		Noticeid:   info.ID,
		OnlinedAt:  info.OnlinedAt,
		OfflinedAt: info.OfflinedAt,
	}
	if err := this.mConn.Model(&NoticeRead{}).Updates(*readInfo).Error; err != nil {
		fmt.Println("UpdateNotice not ok", err.Error())
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//查增，NoticeRead设置成联合主键
func (this *NoticeModel) GetNoticeInfo(ids []int, selectStr string) ([]NoticeInfo, error) {
	if len(selectStr) == 0 {
		selectStr = NoticeInfoSelectElemsFromInner
	}
	infos := make([]NoticeInfo, 0)
	err := this.mConn.Model(&NoticeInfo{}).Select(selectStr).Where("id in (?)", ids).Find(&infos).Error
	if err != nil {
		return nil, err
	}

	return infos, nil
}

//查增，NoticeRead设置成联合主键
func (this *NoticeModel) GetNoticeInfoWithRead(userId uint, userKey string, ids []int, selectStr string) ([]NoticeInfo, error) {
	if len(selectStr) == 0 {
		selectStr = NoticeInfoSelectElems
	}
	infos := make([]NoticeInfo, 0)
	err := this.mConn.Model(&NoticeInfo{}).Select(selectStr).Where("id in (?)", ids).Find(&infos).Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(infos); i++ {
		//		fmt.Println("hahahhah:", userId, userKey, *infos[i].ID)
		readInfo := &NoticeRead{
			Noticeid:   infos[i].ID,
			Userid:     &userId,
			Userkey:    &userKey,
			OnlinedAt:  infos[i].OnlinedAt,
			OfflinedAt: infos[i].OfflinedAt,
		}
		if err := this.mConn.Save(readInfo).Error; err != nil {
			fmt.Println("NewRecord not ok", err.Error())
			ZapLog().With(zap.Error(err), zap.Uint("noticeid", *readInfo.Noticeid), zap.Uint("userid", *readInfo.Userid)).Error("Save err")
		} //返回值无法判断是否是网络问题，要么用save
	}
	return infos, nil
}

//查查
func (this *NoticeModel) GetNoticeList(offset, pagesize uint, order []string, condition map[string]interface{}, vagueStr []string, dur *TimeDuration, selectStr string) ([]NoticeInfo, error) {
	if len(selectStr) == 0 {
		selectStr = NoticeListSelectElemsFromInner
	}
	infos := make([]NoticeInfo, 0)
	statment, timeCondicons := this.genTimeCondition(dur)
	//	readInfos := make([]NoticeRead, 0)

	db := this.mConn.Model(&NoticeInfo{}).Select(selectStr).Where(statment, timeCondicons...).Where(condition)

	for i := 0; i < len(vagueStr); i = i + 2 {
		db = db.Where(vagueStr[i], vagueStr[i+1])
	}
	for i := 0; i < len(order); i++ {
		db = db.Order(order[i])
	}
	err := db.Offset(offset).Limit(pagesize).Find(&infos).Error
	if err != nil {
		return nil, err
	}
	//mm := make(map[uint]bool, 0)
	//for i := 0; i < len(readInfos); i++ {
	//	mm[*readInfos[i].Noticeid] = true
	//}
	//
	//for i := 0; i < len(infos); i++ {
	//	infos[i].IsRead = new(bool)
	//	*infos[i].IsRead = false
	//	_, ok := mm[*infos[i].ID]
	//	if ok {
	//		*infos[i].IsRead = true
	//	}
	//}
	return infos, nil
}

//查查
func (this *NoticeModel) GetNoticeListWithRead(offset, pagesize, userId uint, order []string, condition map[string]interface{}, vagueStr []string, dur *TimeDuration, selectStr string) ([]NoticeInfo, error) {
	if len(selectStr) == 0 {
		selectStr = NoticeListSelectElems
	}
	infos := make([]NoticeInfo, 0)
	statment, timeCondicons := this.genTimeCondition(dur)
	readInfos := make([]NoticeRead, 0)
	//	err := this.mConn.Where("userid = ? and noticeid IN (?)", userId, this.mConn.Model(&User{}).Where(statment, condicons...).Where(param.Condition).Order("updated_at").Offset(offset).Limit(pagesize).Find(&infos).Select("noticeid").QueryExpr()).Find(&readInfos).Error

	db := this.mConn.Model(&NoticeInfo{}).Select(selectStr).Where(statment, timeCondicons...).Where(condition)
	for i := 0; i < len(vagueStr); i = i + 2 {
		db = db.Where(vagueStr[i], vagueStr[i+1])
	}
	for i := 0; i < len(order); i++ {
		db = db.Order(order[i])
	}
	err := db.Offset(offset).Limit(pagesize).Find(&infos).Error
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(infos))
	for i := 0; i < len(infos); i++ {
		ids[i] = *infos[i].ID
	}
	err = this.mConn.Where("userid = ? and noticeid IN (?)", userId, ids).Select("noticeid").Find(&readInfos).Error
	if err != nil {
		return nil, err
	}
	mm := make(map[uint]bool, 0)
	for i := 0; i < len(readInfos); i++ {
		mm[*readInfos[i].Noticeid] = true
	}

	for i := 0; i < len(infos); i++ {
		infos[i].IsRead = new(bool)
		*infos[i].IsRead = false
		_, ok := mm[*infos[i].ID]
		if ok {
			*infos[i].IsRead = true
		}
	}
	return infos, nil
}

func (this *NoticeModel) DelNotice(ids []int) error {
	tx := this.mConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := this.mConn.Delete(NoticeInfo{}, "id in (?)", ids).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = this.mConn.Where("noticeid IN (?)", ids).Delete(NoticeRead{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (this *NoticeModel) Clear() error {
	nowTime := common.NowTimestamp() //ms
	for i := 0; i < 30; i++ {
		err := this.mConn.Where("offlined_at <= ?", nowTime).Delete(NoticeRead{}).Limit(1000).Error
		if err == nil {
			continue
		}
		if err == gorm.ErrRecordNotFound {
			fmt.Println("ErrRecordNotFound ")
			return nil
		}
		return err
	}
	return nil
}

func (this *NoticeModel) CountNotice(condition map[string]interface{}, dur *TimeDuration, condStr string, vagueStr []string) (uint, error) {
	count := uint(0)
	statment, timeCondicons := this.genTimeCondition(dur)
	//	fmt.Println("Count", statment, condicons, param.Condition)
	db := this.mConn.Model(&NoticeInfo{}).Where(statment, timeCondicons...).Where(condition).Where(condStr)

	for i := 0; i < len(vagueStr); i = i + 2 {
		db = db.Where(vagueStr[i], vagueStr[i+1])
	}
	err := db.Count(&count).Error
	//	err := this.mConn.Model(&NoticeInfo{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (this *NoticeModel) CountNoticeRead(userid uint, condition map[string]interface{}, dur *TimeDuration) (uint, error) {
	count := uint(0)
	dur.EndUpdatedAt = nil //notice_read 表没有这两个字段
	dur.StartUpdatedAt = nil
	statment, timeCondicons := this.genTimeCondition(dur) //
	fmt.Println(statment, timeCondicons)
	//	fmt.Println("Count", statment, condicons, param.Condition)
	err := this.mConn.Model(&NoticeRead{}).Where("userid = ?", userid).Where(statment, timeCondicons...).Where(condition).Count(&count).Error
	//	err := this.mConn.Model(&NoticeInfo{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (this *NoticeModel) GetUserMaxReadNotice(userid uint) (uint, error) {
	noticeUser := new(NoticeUser)
	err := this.mConn.Model(&NoticeUser{}).Where("userid = ?", userid).Select("max_notice_id").Find(noticeUser).Error
	//	err := this.mConn.Model(&NoticeInfo{}).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if noticeUser.MaxNoticeId == nil {
		return 0, errors.New("MaxNoticeId is nil", 0)
	}
	return *noticeUser.MaxNoticeId, nil
}

func (this *NoticeModel) UpdateUserMaxReadNotice(userid uint, maxNoticeId uint) error {
	noticeUser := new(NoticeUser)
	noticeUser.Userid = new(uint)
	noticeUser.MaxNoticeId = new(uint)
	*noticeUser.Userid = userid
	*noticeUser.MaxNoticeId = maxNoticeId
	err := this.mConn.Save(noticeUser).Error
	//	err := this.mConn.Model(&NoticeInfo{}).Count(&count).Error
	if err != nil {
		return err
	}
	return nil
}

func (this *NoticeModel) existNoticeRead(userId, noticeId uint) (bool, error) {
	count := 0
	err := this.mConn.Model(&NoticeRead{}).Where("userid = ? and noticeid = ?", userId, noticeId).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return true, nil
}

func (this *NoticeModel) genTimeCondition(param *TimeDuration) (string, []interface{}) {
	var statement string
	var condition []interface{}

	havaOnlinedAtFlag := (param.StartOnlinedAt != nil || param.EndOnlinedAt != nil)
	havaOfflinedAtFlag := (param.StartOfflinedAt != nil || param.EndOfflinedAt != nil)
	havaOnAndOffFlag := (havaOnlinedAtFlag && havaOfflinedAtFlag)

	if param.StartCreatedAt != nil {
		statement += " created_at >= ? and"
		condition = append(condition, *param.StartCreatedAt)
	}
	if param.EndCreatedAt != nil {
		statement += " created_at <= ? and"
		condition = append(condition, *param.EndCreatedAt)
	}

	if param.StartUpdatedAt != nil {
		statement += " updated_at >= ? and"
		condition = append(condition, *param.StartUpdatedAt)
	}
	if param.EndUpdatedAt != nil {
		statement += " updated_at <= ? and"
		condition = append(condition, *param.EndUpdatedAt)
	}

	if havaOnAndOffFlag {
		statement += " ( "
	}

	if havaOnlinedAtFlag {
		statement += " ( "
	}
	if param.StartOnlinedAt != nil {
		statement += " onlined_at >= ? and"
		condition = append(condition, *param.StartOnlinedAt)
	}
	if param.EndOnlinedAt != nil {
		statement += "  onlined_at <= ? and"
		condition = append(condition, *param.EndOnlinedAt)
	}

	statement = strings.TrimRight(statement, "and")

	if havaOnlinedAtFlag {
		statement += " ) "
	}

	if havaOnAndOffFlag {
		if param.LineAndFlag {
			statement += " and"
		} else {
			statement += " or"
		}
	}

	if havaOfflinedAtFlag {
		statement += " ( "
	}
	if param.StartOfflinedAt != nil {
		statement += "  offlined_at >= ? and"
		condition = append(condition, *param.StartOfflinedAt)
	}
	if param.EndOfflinedAt != nil {
		statement += "  offlined_at <= ? and"
		condition = append(condition, *param.EndOfflinedAt)
	}
	statement = strings.TrimRight(statement, "and")

	if havaOfflinedAtFlag {
		statement += " ) "
	}

	if havaOnAndOffFlag {
		statement += " ) and"
	}
	statement = strings.TrimRight(statement, "and")
	statement = strings.TrimRight(statement, "or")

	return statement, condition
}
