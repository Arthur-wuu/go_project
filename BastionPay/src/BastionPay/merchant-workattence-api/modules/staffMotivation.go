package modules

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/models"
	"bytes"
	"encoding/json"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type Staff struct {
}

type StaffListParams struct {
	CompanyId int `json:"company_id"`
}

var CorpIdMap map[int]string

func (this *Staff) SyncStaffInfo() {
	cmpId := config.GConfig.Company.Id
	amModel := new(models.AccountMap)

	for _, v := range cmpId {
		resp, err := this.GetStaffList(v)

		if err != nil {
			ZapLog().Error("record list request err", zap.Error(err))
			continue
		}

		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			ZapLog().Error("response body read err", zap.Error(err))
			continue
		}

		respData := new(api.StaffList)
		err = json.Unmarshal(content, respData)

		if err != nil {
			ZapLog().Error("unmarshal response data err", zap.Error(err))
			continue
		}

		if len(respData.Data) == 0 {
			continue
		}

		for _, v := range respData.Data {
			if v.EeNo == nil || v.BpUid == nil {
				continue
			}

			err := amModel.AddOrUpdate(&v)

			if err != nil {
				ZapLog().Error("add or update account_map record err", zap.Error(err))
				return
			}
		}
	}

	ZapLog().Sugar().Infof("add or update success")
}

func (this *Staff) SendMotivation() {
	cmpId := config.GConfig.Company.Id
	//get record list

	for _, v := range cmpId {

		resp, err := this.GetStaffList(v)

		if err != nil {
			ZapLog().Error("record list request err", zap.Error(err))
			continue
		}

		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			ZapLog().Error("response body read err", zap.Error(err))
			continue
		}

		respData := new(api.StaffList)
		err = json.Unmarshal(content, respData)

		if err != nil {
			ZapLog().Error("unmarshal response data err", zap.Error(err))
			continue
		}

		if len(respData.Data) == 0 {
			continue
		}

		symbol := config.GConfig.Company.ServiceAward.Symbol
		merchantId := config.GConfig.Company.ServiceAward.MerchantId
		coinBase := config.GConfig.Company.ServiceAward.CoinBase
		smModl := new(models.StaffMotivation)

		for _, v := range respData.Data {
			if *v.Valid != 1 { // not equal one means leaved staff
				continue
			}

			// every whole half a year calc
			mod := (common.New().DayBeginTimestamp() - (*v.HiredAt - (*v.HiredAt % 86400))) % (365 * 43200)

			if mod == 0 {
				div := (time.Now().Unix() - *v.HiredAt) / (365 * 43200)

				if div == 0 {
					continue
				}

				//calc how much coin to awarding
				coin := decimal.NewFromFloat32((1 + float32(div-1)/10)).Mul(coinBase)
				//add send record
				awd, err := smModl.AddNew(&v, coin, symbol)

				if err != nil {
					return
				}

				models.SfChan <- models.StaffChan{
					Coin:       coin,
					Symbol:     symbol,
					MerchantId: merchantId,
					Times:      config.GConfig.Company.ServiceAward.SendTimes,
					AwardId:    *awd.Id,
					AccountId:  *awd.BpUid,
				}
			}

		}
	}

}

func (this *Staff) GetStaffList(companyId int) (resp *http.Response, err error) {
	url := config.GConfig.Company.ApiHost + "/v1/ft/merchant/teammanage/employee/gets"
	params := StaffListParams{
		CompanyId: companyId,
	}
	jsonPramas, err := json.Marshal(params)
	b2 := bytes.NewBuffer(jsonPramas)
	resp, err = http.Post(url, "application/json;charset=utf-8", b2)

	return resp, err
}
