package modules

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/db"
	"BastionPay/merchant-workattence-api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type Dingtalk struct {
	AppKey        string
	AppSecret     string
	AppAuthAccess string
	Host          string
}

type AuthAccess struct {
	ErrorCode   int    `json:"errcode,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Errmsg      string `json:"errmsg,omitempty"`
	ExpiresIn   *int   `json:"expires_in,omitempty"`
}

type RecordListParams struct {
	UserIds       []string `json:"userIds"`
	CheckDateFrom string   `json:"checkDateFrom"`
	CheckDateTo   string   `json:"checkDateTo"`
	IsI18n        string   `json:"isI18n"`
}

//企业ID先定义成常量，以后会有调整
const CORPID = "dinga9adfe03ff75b84d35c2f4657eb6378f"

func New() *Dingtalk {
	return &Dingtalk{
		AppKey:        config.GConfig.Dingding.BoJu.AppKey,
		AppSecret:     config.GConfig.Dingding.BoJu.AppSecret,
		AppAuthAccess: "",
		Host:          config.GConfig.Dingding.Host,
	}
}

func (this *Dingtalk) SetAuthAccess() {
	url := this.Host + "/gettoken?appkey=" + this.AppKey + "&appsecret=" + this.AppSecret
	resp, err := http.Get(url)

	if err != nil {
		ZapLog().Error("set auth_access error ", zap.Error(err))
		return
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ZapLog().Error("set auth_access read response body error ", zap.Error(err))
		return
	}

	ret := new(AuthAccess)
	err = json.Unmarshal(content, ret)

	if err != nil {
		ZapLog().Error("set auth_access json unmarshal error ", zap.Error(err))
		return
	}

	this.AppAuthAccess = ret.AccessToken

	if ret.ExpiresIn == nil {
		db.GCache.Dingtalk.SetWithExpire(config.GConfig.Gcache.SecretKey, ret.AccessToken, time.Duration(*ret.ExpiresIn)*time.Second)
	} else {
		db.GCache.Dingtalk.SetWithExpire(config.GConfig.Gcache.SecretKey, ret.AccessToken, time.Duration(config.GConfig.Gcache.Expire)*time.Second)
	}
}

func (this *Dingtalk) GetAuthAccess() string {
	token := ""

	if this.AppAuthAccess == "" {
		val, err := db.GCache.Dingtalk.Get(config.GConfig.Gcache.SecretKey)

		if err == nil {
			token = val.(string)
		} else {
			this.SetAuthAccess()
			return this.AppAuthAccess
		}
	} else {
		token = this.AppAuthAccess
	}

	return token
}

func (this *Dingtalk) GetDingdingRecordList() (*[]api.Recordresult, error) {
	//get all account info
	acts := new(models.AccountMap).GetAllAccounts(CORPID)

	if len(acts) == 0 {
		ZapLog().Error("accout map records not exists")
		return nil, fmt.Errorf("accout map records not exists")
	}

	userIds := make([]string, len(acts))

	for k, v := range acts {
		userIds[k] = *v.StaffId
	}

	sTime, eTime := "", ""
	ckRecord, err := new(models.CheckinRecord).GetMaxCheckinAt()

	if err != nil {
		ZapLog().Error("checkin record query err", zap.Error(err))
		return nil, err
	}

	dBeignTimestamp := common.New().DayBeginTimestamp()
	loc := time.FixedZone("UTC", 8*3600)

	if ckRecord == nil || *ckRecord == 0 {
		//sTime = time.Now().Format("2006-01-02 15:04:05")
		sTime = time.Unix(dBeignTimestamp+6*3600, 0).In(loc).Format("2006-01-02 15:04:05")
	} else {
		sTime = time.Unix(*ckRecord+1, 0).In(loc).Format("2006-01-02 15:04:05")
	}

	eTime = time.Unix(dBeignTimestamp+86399, 0).In(loc).Format("2006-01-02 15:04:05")

	//get record list
	resp, err := this.RecordListPost(userIds, sTime, eTime)

	if err != nil {
		ZapLog().Error("record list request err", zap.Error(err))
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ZapLog().Error("response body read err", zap.Error(err))
		return nil, err
	}

	respData := new(api.RecordListResponse)
	err = json.Unmarshal(content, respData)

	if err != nil {
		ZapLog().Error("unmarshal response data err", zap.Error(err))
		return nil, err
	}

	if respData.Errcode != 0 {
		ZapLog().Error("dingding recordlist response err " + respData.Errmsg)
		return nil, err
	}

	return &respData.Recordresult, nil
}

func (this *Dingtalk) AddCheckinRecord() {
	recordResult, err := this.GetDingdingRecordList()

	if err != nil {
		ZapLog().Error("get dingding record list err  ", zap.Error(err))
		return
	}

	ckModel := new(models.CheckinRecord)

	for _, v := range *recordResult {
		if v.TimeResult == "" {
			continue
		}

		err := ckModel.AddNewNoRecord(&v)

		if err != nil {
			ZapLog().Error("dingding checkin record add err ", zap.Error(err))
			return
		}
	}
}

func (this *Dingtalk) AttenRewardSend() {
	recordResult, err := this.GetDingdingRecordList()

	if err != nil {
		ZapLog().Error("get dingding record list err  ", zap.Error(err))
		return
	}

	ckModel := new(models.CheckinRecord)
	uIdMaps := make(map[string]int)
	coin := config.GConfig.Award.Checkin.Coin
	symbol := config.GConfig.Award.Checkin.Symbol
	dayBeginStamp := common.New().DayBeginTimestamp()

	for _, v := range *recordResult {
		if v.TimeResult == "" {
			continue
		}

		rcd, err := ckModel.AddNew(&v)

		if err != nil {
			ZapLog().Error("dingding checkin record add err ", zap.Error(err))
			return
		}

		if v.IsLegal == "N" {
			continue
		}

		uCheckTime := v.UserCheckTime / 1000

		if (uCheckTime >= dayBeginStamp+3600*6 && uCheckTime <= dayBeginStamp+3600*9+1800) ||
			(uCheckTime >= dayBeginStamp+3600*18+1800 && uCheckTime < dayBeginStamp+86400) { //早上 6:00 - 9:30,晚上 18:30 - 23:59:59
			_, ok := uIdMaps[*rcd.StaffId]

			if !ok {
				if uCheckTime <= dayBeginStamp+3600*7 {
					this.SendAward(rcd, dayBeginStamp+3600*6, 3.5, symbol, coin, 1)
				} else if uCheckTime <= dayBeginStamp+3600*9+1800 {
					multTimes := (dayBeginStamp + 3600*9 + 1800 - uCheckTime) / 1800
					coinBase := float32(multTimes)*0.5 + 1
					this.SendAward(rcd, dayBeginStamp+3600*6, coinBase, symbol, coin, 1)
				} else {
					this.SendAward(rcd, dayBeginStamp+3600*6, 1, symbol, coin, 1)
				}
				//go this.SendAward(rcd, common.New().DayBeginTimestamp(), 1)

				uIdMaps[*rcd.StaffId] = 1
			}
		}
	}
}

func (this *Dingtalk) OvertimeAwardSend() {
	ckModel := new(models.CheckinRecord)
	rList, err := ckModel.GetValidWorkOvertimeList(CORPID)
	uIdCrdMap := make(map[string]*models.CheckinRecord, len(rList))
	uIds := make([]string, len(rList))

	if err != nil {
		ZapLog().Error("get overtime checkin record err ", zap.Error(err))
		return
	}

	if len(rList) > 0 {
		for k, v := range rList {
			uIds[k] = *v.StaffId
			uIdCrdMap[*v.StaffId] = v
		}

		eRlist, err := ckModel.GetUsersEarliestCheckinRecordList(uIds)

		if err != nil {
			ZapLog().Error("get onduty checkin record err ", zap.Error(err))
			return
		}

		coin := config.GConfig.Award.Extratime.Coin
		symbol := config.GConfig.Award.Extratime.Symbol
		awdModel := new(models.AwardRecord)

		for _, v := range eRlist {
			crd, ok := uIdCrdMap[*v.StaffId]

			if ok {
				DayBeginTimestamp := common.New().DayBeginTimestamp()

				if (DayBeginTimestamp+13*3600) >= *v.CheckinAt && *crd.CheckinAt-*v.CheckinAt > 37800 { // 13 * 3600 means 13:00
					awdRecord, err := awdModel.GetLatestOvertimeRecordByUserId(*crd.StaffId)

					if err != nil {
						ZapLog().Error("get overtime award send record err ", zap.Error(err))
						return
					}

					if awdRecord == nil {
						var multi int64
						//20:00 timestamp
						eightStmp := common.New().DayBeginTimestamp() + 72000

						if eightStmp-*v.CheckinAt >= 37800 {
							multi = (*crd.CheckinAt - eightStmp) / 1800
						} else {
							multi = (*crd.CheckinAt - *v.CheckinAt - 37800) / 1800
						}

						//8:30之后打卡默认送0.6个币
						this.SendAward(crd, DayBeginTimestamp, float32(multi)+3, symbol, coin, 2)
					} else {
						ckRecord, err := ckModel.GetRecordById(*awdRecord.CheckinId)

						if err != nil || ckRecord == nil {
							ZapLog().Error("get overtime checkin record by award send checkin id err ", zap.Error(err))
							return
						}

						multi := (*crd.CheckinAt-common.New().DayBeginTimestamp()-72000)/1800 -
							(*ckRecord.CheckinAt-common.New().DayBeginTimestamp()-72000)/1800

						if multi > 0 {
							this.SendAward(crd, DayBeginTimestamp, float32(multi), symbol, coin, 2)
						}
					}
				}
			}
		}
	} else {
		ZapLog().Error("overtime checkin record len 0")
	}
}

func (this *Dingtalk) RecordListPost(userIds []string, sTime, eTime string) (resp *http.Response, err error) {
	token := this.GetAuthAccess()
	url := this.Host + "/attendance/listRecord?access_token=" + token
	params := RecordListParams{
		UserIds:       userIds,
		CheckDateFrom: sTime,
		CheckDateTo:   eTime,
		IsI18n:        "false",
	}
	jsonPramas, err := json.Marshal(params)
	b := []byte(jsonPramas)
	b2 := bytes.NewBuffer(b)
	resp, err = http.Post(url, "application/json;charset=utf-8", b2)

	return resp, err
}

/**
 * rcd: CheckinRecord
 * bTime: query begin timestamp
 * coinBase: coin base number
 **/
func (this *Dingtalk) SendAward(rcd *models.CheckinRecord, bTime int64, coinBase float32, symbol string, coin decimal.Decimal, atype int) {
	arm := new(models.AwardRecord)
	//check today award had sended
	dayBeginStamp := common.New().DayBeginTimestamp()
	tHour := time.Now().In(time.FixedZone("UTC", 8*3600)).Hour()
	var num int
	var err error

	if tHour >= 18 && tHour <= 23 {
		num, err = arm.SendCheckDay(*rcd.StaffId, bTime+12*3600+1800, *rcd.CheckinAt)
	} else {
		num, err = arm.SendCheckDay(*rcd.StaffId, bTime, *rcd.CheckinAt)
	}

	if err != nil || (num > 0 && bTime >= dayBeginStamp+6*3600) {
		return
	}

	//get accout info
	act, err := new(models.AccountMap).GetByStaffId(*rcd.StaffId)

	if err != nil {
		return
	}

	//add send record and send award
	dcl := decimal.NewFromFloat32(coinBase)
	tCoin := dcl.Mul(coin)
	awd, err := arm.AddAuto(*rcd.Id, *act, *rcd.StaffId, tCoin, symbol, atype)

	if err != nil {
		return
	}

	models.SChan <- models.SendChan{
		Coin:       tCoin,
		Symbol:     symbol,
		MerchantId: config.GConfig.Award.MerchantId,
		Times:      config.GConfig.Award.SendTimes,
		AwardId:    *awd.Id,
		AccountId:  *awd.AccId,
	}
}
