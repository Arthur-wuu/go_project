//公告功能，包含公告列表查询，特定公告查询
package controllers

import (
	"BastionPay/bas-admin-api/models"
	"BastionPay/bas-api/apibackend"
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-api/admin"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/domac/gosler"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"io/ioutil"
	"sync"
)

func NewNoticeController(db *gorm.DB, config *config.Config) *NoticeController {
	notice := new(NoticeController)
	notice.mConfig = config
	notice.mNoticeModel = models.NewNoticeModel(db)
	notice.addTimeWork()
	notice.mUserMaxNoticeGroup = make(map[uint]uint, 0)
	return notice
}

type NoticeController struct {
	mConfig             *config.Config
	mNoticeModel        *models.NoticeModel
	mUserMaxNoticeGroup map[uint]uint
	mLock               sync.Mutex
}

func (this *NoticeController) GetFromInner(ctx iris.Context) {
	params := new(admin.NoticeIdsParam)
	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}
	infos, err := this.mNoticeModel.GetNoticeInfo(params.Ids, models.NoticeInfoSelectElemsFromInner)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Models GetNoticeInfo err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	resInfos := make([]admin.NoticeInfo, len(infos))
	for i := 0; i < len(infos); i++ {
		resInfos[i].Id = infos[i].ID
		resInfos[i].CreatedAt = infos[i].CreatedAt
		resInfos[i].UpdatedAt = infos[i].UpdatedAt
		resInfos[i].Language = infos[i].Language
		resInfos[i].Title = infos[i].Title
		resInfos[i].Content = infos[i].Content
		resInfos[i].Focus = infos[i].Focus
		resInfos[i].Author = infos[i].Author
		resInfos[i].Race = infos[i].Race
		resInfos[i].OnlinedAt = infos[i].OnlinedAt
		resInfos[i].OfflinedAt = infos[i].OfflinedAt
		resInfos[i].Alive = this.genAliveStatus(0, infos[i].OnlinedAt, infos[i].OfflinedAt)
	}
	dataAck, err := json.Marshal(resInfos)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("json Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}
	ctx.JSON(common.NewSuccessResponse(ctx, string(dataAck)))
}

func (this *NoticeController) GetListFromInner(ctx iris.Context) {
	param := new(admin.NoticeListParam)
	err := ctx.ReadJSON(&param)
	if err != nil {
		body, _ := ioutil.ReadAll(ctx.Request().Body)
		ZapLog().With(zap.Error(err), zap.String("body", string(body))).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	var (
		pageNum    = uint(param.MaxDispLines)
		totalLine  = uint(param.TotalLines)
		pageIndex  = uint(param.PageIndex)
		beginIndex = uint(0)
	)

	conditionMap, vagueStr, order, dur := this.listParamToSqlCondition(param, false)
	//	dur := models.NewTimeDuration(param.StartCreatedAt, param.EndCreatedAt, param.StartUpdatedAt, param.EndUpdatedAt, param.StartOnlinedAt, param.EndOnlinedAt, param.StartOfflinedAt, param.EndOfflinedAt, param.Alive)
	if totalLine == 0 {
		//totalLine, err = db.ListUserCount()
		totalLine, err = this.mNoticeModel.CountNotice(conditionMap, dur, "", vagueStr)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("models CountNotice err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "CountNotice_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	if pageNum < 1 || pageNum > 100 {
		pageNum = 50
	}

	beginIndex = pageNum * (pageIndex - 1)

	//ackUserList, err := db.ListUsers(beginIndex, pageNum)
	lists, err := this.mNoticeModel.GetNoticeList(beginIndex, pageNum, order, conditionMap, vagueStr, dur, models.NoticeListSelectElemsFromInner)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("models GetNoticeList err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GetNoticeList_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	resList := new(admin.ResNoticeList)

	resList.PageIndex = pageIndex
	resList.MaxDispLines = pageNum
	resList.TotalLines = totalLine

	resList.Notices = make([]admin.NoticeInfo, len(lists))
	for i := 0; i < len(lists); i++ {
		resList.Notices[i].Id = lists[i].ID
		resList.Notices[i].CreatedAt = lists[i].CreatedAt
		resList.Notices[i].UpdatedAt = lists[i].UpdatedAt
		resList.Notices[i].OnlinedAt = lists[i].OnlinedAt
		resList.Notices[i].OfflinedAt = lists[i].OfflinedAt
		resList.Notices[i].Language = lists[i].Language
		resList.Notices[i].Author = lists[i].Author
		resList.Notices[i].Focus = lists[i].Focus
		resList.Notices[i].Title = lists[i].Title
		resList.Notices[i].Abstract = lists[i].Abstract
		resList.Notices[i].Race = lists[i].Race
		resList.Notices[i].Alive = this.genAliveStatus(0, lists[i].OnlinedAt, lists[i].OfflinedAt)
	}

	// to ack
	dataAck, err := json.Marshal(resList)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("jsonMarshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}
	ctx.JSON(common.NewSuccessResponse(ctx, string(dataAck)))
}

//查列表+已读未读标记
func (this *NoticeController) GetList(ctx iris.Context) {
	param := new(admin.NoticeListParam)
	err := ctx.ReadJSON(&param)
	if err != nil {
		body, _ := ioutil.ReadAll(ctx.Request().Body)
		ZapLog().With(zap.Error(err), zap.String("body", string(body))).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	var appClaims *common.AppClaims
	body := ctx.Values().Get("app_claims")
	if body == nil {
		//		ZapLog().With(zap.String("error", "nofind app_claims")).Error("app_claims get err")
		//		glog.Error("nofind app_claims err")
		//		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_APP_claims", common.ResponseErrorParams))
		//		return
		appClaims = nil
	} else {
		appClaims, _ = body.(*common.AppClaims)
	}

	var (
		pageNum    = uint(param.MaxDispLines)
		totalLine  = uint(param.TotalLines)
		pageIndex  = uint(param.PageIndex)
		beginIndex = uint(0)
	)

	if param.Alive == nil {
		param.Alive = new(int)
		*param.Alive = admin.STATUS_ALIVE_Online
	}

	conditionMap, vagueStr, order, dur := this.listParamToSqlCondition(param, true)
	//	dur := models.NewTimeDuration(param.StartCreatedAt, param.EndCreatedAt, param.StartUpdatedAt, param.EndUpdatedAt, param.StartOnlinedAt, param.EndOnlinedAt, param.StartOfflinedAt, param.EndOfflinedAt, param.Alive)

	if totalLine == 0 {
		//totalLine, err = db.ListUserCount()
		totalLine, err = this.mNoticeModel.CountNotice(conditionMap, dur, "", vagueStr)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("models CountNotice err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	if pageNum < 1 || pageNum > 100 {
		pageNum = 50
	}

	beginIndex = pageNum * (pageIndex - 1)

	var lists []models.NoticeInfo
	if appClaims == nil {
		lists, err = this.mNoticeModel.GetNoticeList(beginIndex, pageNum, order, conditionMap, vagueStr, dur, models.NoticeListSelectElems)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("models GetNoticeList err")
			//		glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	} else {
		lists, err = this.mNoticeModel.GetNoticeListWithRead(beginIndex, pageNum, appClaims.UserId, order, conditionMap, vagueStr, dur, models.NoticeListSelectElems)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("models GetNoticeListWithRead err")
			//		glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	resList := new(admin.ResNoticeList)

	resList.PageIndex = pageIndex
	resList.MaxDispLines = pageNum
	resList.TotalLines = totalLine

	maxNoticeId := uint(0)
	resList.Notices = make([]admin.NoticeInfo, len(lists))
	for i := 0; i < len(lists); i++ {
		if (lists[i].ID != nil) && (*lists[i].ID > maxNoticeId) {
			maxNoticeId = *lists[i].ID
		}
		resList.Notices[i].Id = lists[i].ID
		resList.Notices[i].CreatedAt = lists[i].CreatedAt
		resList.Notices[i].UpdatedAt = lists[i].UpdatedAt
		resList.Notices[i].OnlinedAt = lists[i].OnlinedAt
		resList.Notices[i].Language = lists[i].Language
		resList.Notices[i].Focus = lists[i].Focus
		resList.Notices[i].Title = lists[i].Title
		resList.Notices[i].Abstract = lists[i].Abstract
		resList.Notices[i].IsRead = lists[i].IsRead
		resList.Notices[i].Race = lists[i].Race
		resList.Notices[i].Alive = this.genAliveStatus(0, lists[i].OnlinedAt, lists[i].OfflinedAt)
	}

	// to ack
	dataAck, err := json.Marshal(resList)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("jsonMarshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, string(dataAck)))

	if appClaims == nil {
		return
	}

	lastMaxId, err := this.getMaxNoticeId(appClaims.UserId)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Uint("userid", appClaims.UserId)).Error("getMaxNoticeId err")
		return
	}
	if maxNoticeId < lastMaxId {
		return
	}
	err = this.setMaxNoticeId(appClaims.UserId, maxNoticeId)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Uint("userid", appClaims.UserId), zap.Uint("maxnoticeid", maxNoticeId)).Error("UpdateUserMaxReadNotice err")
		return
	}
}

//查公告, 同时插入notice_read表
func (this *NoticeController) Get(ctx iris.Context) {
	param := new(admin.NoticeIdsParam)
	err := ctx.ReadJSON(&param)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	var appClaims *common.AppClaims
	body := ctx.Values().Get("app_claims")
	if body == nil {
		//		ZapLog().With(zap.String("error", "nofind app_claims")).Error("app_claims get err")
		//		glog.Error("nofind app_claims err")
		//		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_APP_claims", common.ResponseErrorParams))
		//		return
		appClaims = nil
	} else {
		appClaims, _ = body.(*common.AppClaims)
	}

	var infos []models.NoticeInfo
	if appClaims == nil {
		infos, err = this.mNoticeModel.GetNoticeInfo(param.Ids, models.NoticeInfoSelectElems)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("models GetNoticeInfo err")
			//		glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	} else {
		infos, err = this.mNoticeModel.GetNoticeInfoWithRead(appClaims.UserId, appClaims.Uuid, param.Ids, models.NoticeInfoSelectElems)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("models GetNoticeInfoWithRead err")
			//		glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	resInfos := make([]admin.NoticeInfo, len(infos))
	for i := 0; i < len(infos); i++ {
		resInfos[i].Id = infos[i].ID
		resInfos[i].CreatedAt = infos[i].CreatedAt
		resInfos[i].UpdatedAt = infos[i].UpdatedAt
		resInfos[i].OnlinedAt = infos[i].OnlinedAt
		resInfos[i].Language = infos[i].Language
		resInfos[i].Title = infos[i].Title
		resInfos[i].Content = infos[i].Content
		resInfos[i].Alive = this.genAliveStatus(0, infos[i].OnlinedAt, infos[i].OfflinedAt)
	}
	dataAck, err := json.Marshal(resInfos)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("json.Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}
	ctx.JSON(common.NewSuccessResponse(ctx, string(dataAck)))
}

//插入公告
func (this *NoticeController) AddFromInner(ctx iris.Context) {
	fmt.Println("AddFromIner")
	params := make([]admin.NoticeInfo, 0)
	err := ctx.ReadJSON(&params)
	if err != nil {
		body, _ := ioutil.ReadAll(ctx.Request().Body)
		ZapLog().With(zap.Error(err), zap.String("body", string(body))).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	nowTime := common.NowTimestamp()
	fmt.Println("param:", params)
	oneFlag := false //至少一次成功
	idStates := make([]admin.ResNoticeIdState, len(params))
	for i := 0; i < len(params); i++ {
		if params[i].Title == nil || params[i].Author == nil || params[i].Content == nil || params[i].Language == nil ||
			len(*params[i].Title) == 0 || len(*params[i].Content) == 0 {
			idStates[i].State = new(bool)
			*idStates[i].State = false
			ZapLog().Error("nil in (Title Author Content Language)")
			idStates[i].ErrMsg = new(string)
			*idStates[i].ErrMsg = "nil in (Title Author Content Language)"
			continue
		}
		if (params[i].OnlinedAt != nil) && (params[i].OfflinedAt != nil) {
			if *params[i].OfflinedAt <= *params[i].OnlinedAt {
				idStates[i].State = new(bool)
				*idStates[i].State = false
				idStates[i].ErrMsg = new(string)
				*idStates[i].ErrMsg = "OfflinedAt <= OnlinedAt"
				ZapLog().Error("Title is nil err")
				continue
			}
		}
		noticeInfo := new(models.NoticeInfo)
		noticeInfo.Language = params[i].Language
		noticeInfo.Author = params[i].Author
		noticeInfo.Focus = params[i].Focus
		noticeInfo.Race = params[i].Race
		noticeInfo.Title = params[i].Title
		noticeInfo.Abstract = params[i].Abstract
		noticeInfo.Content = params[i].Content
		noticeInfo.OnlinedAt = params[i].OnlinedAt
		noticeInfo.OfflinedAt = params[i].OfflinedAt
		if params[i].OfflinedAt == nil {
			noticeInfo.OfflinedAt = new(int64)
			*noticeInfo.OfflinedAt = 999999999999999999
		}
		if params[i].OnlinedAt == nil {
			noticeInfo.OnlinedAt = new(int64)
			*noticeInfo.OnlinedAt = nowTime
		}

		//		idStates[i].Id = noticeInfo.ID
		idStates[i].State = new(bool)

		ok := this.mNoticeModel.AddNotice(noticeInfo)
		if !ok {
			ZapLog().Error("model AddNotice not ok")
			idStates[i].ErrMsg = new(string)
			*idStates[i].ErrMsg = "db add not ok"
		} else {
			oneFlag = true
		}
		*idStates[i].State = ok
		//		fmt.Println("ok %v", ok)
		ZapLog().Sugar().Infof("ok %v", ok)
	}

	if !oneFlag {
		ZapLog().Error("AddNotice err[no one]")
		//		fmt.Println("AddNotice err[no one]")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ackData, err := json.Marshal(idStates)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("json.Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}
	ctx.JSON(common.NewSuccessResponse(ctx, string(ackData)))
}

//更新公告,如果是下架就清理无用信息
func (this *NoticeController) UpdateFromInner(ctx iris.Context) {
	params := make([]admin.NoticeInfo, 0)
	err := ctx.ReadJSON(&params)
	if err != nil {
		body, _ := ioutil.ReadAll(ctx.Request().Body)
		ZapLog().With(zap.Error(err), zap.String("body", string(body))).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	nowTime := common.NowTimestamp()
	noticeInfos := make([]models.NoticeInfo, len(params))
	for i := 0; i < len(params); i++ {
		if params[i].Id == nil {
			ZapLog().Error("id is nil")
			//			glog.Error("id is nil")
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS:id have nil", apibackend.BASERR_INVALID_PARAMETER.Code()))
			return
		}
		noticeInfos[i].ID = params[i].Id
		noticeInfos[i].Language = params[i].Language
		noticeInfos[i].Author = params[i].Author
		noticeInfos[i].Focus = params[i].Focus
		noticeInfos[i].Race = params[i].Race
		noticeInfos[i].Title = params[i].Title
		noticeInfos[i].Abstract = params[i].Abstract
		noticeInfos[i].Content = params[i].Content
		noticeInfos[i].OnlinedAt = params[i].OnlinedAt
		noticeInfos[i].OfflinedAt = params[i].OfflinedAt
	}
	err = this.updateNoticeInfo(ctx, noticeInfos, nowTime)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("updateNoticeInfo err")
		//		glog.Error(err.Error())
	}
}

//删除公告
func (this *NoticeController) DelFromInner(ctx iris.Context) {
	param := new(admin.NoticeIdsParam)
	err := ctx.ReadJSON(&param)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}
	err = this.mNoticeModel.DelNotice(param.Ids)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("models DelNotice err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}
	ctx.JSON(common.NewSuccessResponse(ctx, ""))
}

func (this *NoticeController) FocusFromInner(ctx iris.Context) {
	params := make([]admin.NoticeInfo, 0)
	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	nowTime := common.NowTimestamp()
	noticeInfos := make([]models.NoticeInfo, len(params))
	for i := 0; i < len(params); i++ {
		if params[i].Id == nil || params[i].Focus == nil {
			ZapLog().Error("id is nil or focus is nil")
			//			glog.Error("id is nil or focus is nil")
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS:id or focus have nil", apibackend.BASERR_INVALID_PARAMETER.Code()))
			return
		}
		noticeInfos[i].ID = params[i].Id
		noticeInfos[i].Focus = params[i].Focus
	}
	err = this.updateNoticeInfo(ctx, noticeInfos, nowTime)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("updateNoticeInfo err")
		//		glog.Error(err.Error())
	}

}

func (this *NoticeController) CountUserNotices(ctx iris.Context) {
	param := new(admin.CountUserNoticesParam)
	err := ctx.ReadJSON(&param)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	if param.Alive == nil {
		param.Alive = new(int)
		*param.Alive = admin.STATUS_ALIVE_Online
	}

	var appClaims *common.AppClaims
	body := ctx.Values().Get("app_claims")
	if body == nil {
		//		ZapLog().Error("nofind app_claims err")
		//		glog.Error("nofind app_claims err")
		//		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_APP_claims", common.ResponseErrorParams))
		//		return
		appClaims = nil
	} else {
		appClaims, _ = body.(*common.AppClaims)
	}

	//	dur := models.NewTimeDuration(param.StartCreatedAt, param.EndCreatedAt, param.StartUpdatedAt, param.EndUpdatedAt, param.StartOnlinedAt, param.EndOnlinedAt, param.StartOfflinedAt, param.EndOfflinedAt, param.Alive)

	condition, dur := this.countParamToSqlCondition(param)

	var lastMaxNoticeId uint
	if appClaims != nil {
		lastMaxNoticeId, err = this.getMaxNoticeId(appClaims.UserId)
		//		rcount, err = this.mNoticeModel.CountNoticeRead(appClaims.UserId, param.Condition, dur)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("model CountNoticeRead err")
			//		glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	fmt.Println(lastMaxNoticeId)
	res := new(admin.ResCountNotices)
	if param.ReadFlag != nil && *param.ReadFlag == true {
		res.ReadCount = new(uint)
		*res.ReadCount, err = this.mNoticeModel.CountNotice(condition, dur, fmt.Sprintf(" id <= %d ", lastMaxNoticeId), nil)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("model CountNotice err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	} else {

		res.UnReadCount = new(uint)
		*res.UnReadCount, err = this.mNoticeModel.CountNotice(condition, dur, fmt.Sprintf(" id > %d ", lastMaxNoticeId), nil)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("model CountNotice err")
			//			glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
			return
		}
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("json.Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}
	ctx.JSON(common.NewSuccessResponse(ctx, string(resBytes)))
}

func (this *NoticeController) updateNoticeInfo(ctx iris.Context, noticeInfos []models.NoticeInfo, nowTime int64) error {
	if nowTime == 0 {
		nowTime = common.NowTimestamp()
	}
	oneFlag := false //至少一次成功
	idStates := make([]admin.ResNoticeIdState, len(noticeInfos))
	var err error
	for i := 0; i < len(noticeInfos); i++ {
		idStates[i].Id = noticeInfos[i].ID
		idStates[i].State = new(bool)
		noticeInfos[i].UpdatedAt = new(int64)
		*noticeInfos[i].UpdatedAt = nowTime
		noticeInfos[i].CreatedAt = nil
		err := this.mNoticeModel.UpdateNotice(&noticeInfos[i])
		if err != nil {
			*idStates[i].State = false
			ZapLog().With(zap.Error(err)).Error("model UpdateNotice err")
			//			glog.Error(err.Error())
		} else {
			oneFlag = true
			*idStates[i].State = true
		}
	}

	if !oneFlag {
		ZapLog().Error("UpdateFromIner err[no one]")
		//		fmt.Println("UpdateFromIner err[no one]")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return err
	}

	ackData, err := json.Marshal(idStates)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("json.Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JSON_MUSHAL_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return err
	}
	ctx.JSON(common.NewSuccessResponse(ctx, string(ackData)))
	return nil
}

func (this *NoticeController) addTimeWork() {
	if len(this.mConfig.Notice.ClearTimer) != 5 {
		ZapLog().Error("config ClearTimer wrong:" + this.mConfig.Notice.ClearTimer)
		//		glog.Error("config ClearTimer wrong:", this.mConfig.Notice.ClearTimer)
		return
	}
	gosler.Every(1).Day().At(this.mConfig.Notice.ClearTimer).Do(this.timeWork)
	go gosler.Start()
}

func (this *NoticeController) timeWork() {
	ZapLog().Info("start clear Notice Db")
	//	glog.Info("start clear Notice Db")
	if err := this.mNoticeModel.Clear(); err != nil {
		ZapLog().With(zap.Error(err)).Error("model clear err")
		//		glog.Error(err.Error())
		return
	}
	ZapLog().Info("clear Notice Db ok")
	//	glog.Info("clear Notice Db ok")
}

func (this *NoticeController) getMaxNoticeId(userid uint) (uint, error) {
	this.mLock.Lock()
	lastMaxNoticeId, ok := this.mUserMaxNoticeGroup[userid]
	this.mLock.Unlock()
	if ok {
		return lastMaxNoticeId, nil
	}

	lastMaxNoticeId, err := this.mNoticeModel.GetUserMaxReadNotice(userid)
	if err != nil {
		return 0, err
	}
	this.mLock.Lock()
	this.mUserMaxNoticeGroup[userid] = lastMaxNoticeId
	this.mLock.Unlock()
	return lastMaxNoticeId, nil
}

func (this *NoticeController) setMaxNoticeId(userid, maxUserid uint) error {
	this.mLock.Lock()
	this.mUserMaxNoticeGroup[userid] = maxUserid
	this.mLock.Unlock()
	err := this.mNoticeModel.UpdateUserMaxReadNotice(userid, maxUserid)
	if err != nil {
		//		ZapLog().With(zap.Error(err), zap.Uint("userid", userid), zap.Uint("maxnoticeid", maxUserid)).Error("UpdateUserMaxReadNotice err")
		return err
	}
	return nil
}

func (this *NoticeController) listParamToSqlCondition(param *admin.NoticeListParam, focusOrderFlag bool) (mm map[string]interface{}, likeStr []string, orderArr []string, dur *models.TimeDuration) {
	if orderArr == nil {
		orderArr = make([]string, 0)
	}
	if focusOrderFlag {
		orderArr = append(orderArr, "focus  DESC")
	}

	order := "onlined_at DESC"
	if param.Order != nil {
		order = *param.Order
		if (param.Desc != nil) && (*param.Desc == true) {
			order += "   DESC"
		}
	}
	orderArr = append(orderArr, order)

	if param.Language != nil {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm["language"] = *param.Language
	}
	if param.Focus != nil {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm["focus"] = *param.Focus
	}
	if param.Race != nil {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm["race"] = *param.Race
	}
	if param.Title != nil {
		//if mm == nil {
		//	mm = make(map[string]interface{}, 0)
		//}
		//mm["title"] = *param.Title
	}
	if param.Title != nil {
		if likeStr == nil {
			likeStr = make([]string, 0)
		}
		likeStr = append(likeStr, "title LIKE ?")
		likeStr = append(likeStr, "%"+*param.Title+"%")
	}
	if param.Id != nil {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm["id"] = *param.Id
	}
	for k, v := range param.Condition {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm[k] = v
	}

	dur = models.NewTimeDuration(param.StartCreatedAt, param.EndCreatedAt, param.StartUpdatedAt, param.EndUpdatedAt, param.StartOnlinedAt, param.EndOnlinedAt, param.StartOfflinedAt, param.EndOfflinedAt, param.Alive)
	return
}

func (this *NoticeController) countParamToSqlCondition(param *admin.CountUserNoticesParam) (mm map[string]interface{}, dur *models.TimeDuration) {
	if param.Language != nil {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm["language"] = *param.Language
	}
	if param.Focus != nil {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm["focus"] = *param.Focus
	}
	if param.Race != nil {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm["race"] = *param.Race
	}
	for k, v := range param.Condition {
		if mm == nil {
			mm = make(map[string]interface{}, 0)
		}
		mm[k] = v
	}
	dur = models.NewTimeDuration(param.StartCreatedAt, param.EndCreatedAt, param.StartUpdatedAt, param.EndUpdatedAt, param.StartOnlinedAt, param.EndOnlinedAt, param.StartOfflinedAt, param.EndOfflinedAt, param.Alive)
	return
}

func (this *NoticeController) genAliveStatus(nowtime int64, onlineAt, offlineAt *int64) *int {
	if onlineAt == nil || offlineAt == nil {
		return nil
	}
	if nowtime <= 1000 {
		nowtime = common.NowTimestamp()
	}
	status := new(int)
	if *offlineAt <= nowtime {
		*status = admin.STATUS_Alive_AfterOffline
		return status
	}
	if *onlineAt <= nowtime {
		*status = admin.STATUS_ALIVE_Online
		return status
	}
	if *onlineAt > nowtime {
		*status = admin.STATUS_ALIVE_PreOnline
		return status
	}
	return nil
}
