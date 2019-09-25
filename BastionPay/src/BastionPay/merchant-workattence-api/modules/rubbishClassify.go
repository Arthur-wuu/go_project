package modules

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/models"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type RClassify struct {
}

type QueryListParams struct {
	CompanyId    int `json:"company_id"`
	DepartmentId int `json:"department"`
}

func (this *RClassify) SendAward(id int) {
	rcrModel := new(models.RubbishClassifyRecord)
	rcd, err := rcrModel.GetById(id)

	if err != nil {
		ZapLog().Error("get rubbish classify record err", zap.Error(err))
		return
	}

	var depart []*models.DepartmentInfo
	err = json.Unmarshal([]byte(*rcd.DepartmentInfo), &depart)

	if err != nil || len(depart) == 0 {
		ZapLog().Error("rubbish classify department info unmarshal err or department info len(0)", zap.Error(err))
		return
	}

	err = rcrModel.UpdateTransferFlag(*rcd.Id, 1)

	if err != nil {
		ZapLog().Error("update rubbish classify record transfer_flag err", zap.Error(err))
		return
	}

	actModels := new(models.AccountMap)
	rcModl := new(models.RubbishClassify)

	for _, v := range depart {
		acts := actModels.GetAccountsByDepartId(v.DepartmentId)
		symbol := config.GConfig.Company.RubbishClassify.Symbol
		merchantId := config.GConfig.Company.RubbishClassify.MerchantId
		coinArr := config.GConfig.Company.RubbishClassify.Coin
		//calc how much coin to awarding
		coin := coinArr[v.Score]
		//add send record

		for _, act := range acts {
			awd, err := rcModl.AddNewByAccount(*act, coin, symbol, v.Score, *rcd.Id)

			if err != nil {
				return
			}

			models.RcChan <- models.RClassifyChan{
				Coin:       coin,
				Symbol:     symbol,
				MerchantId: merchantId,
				Times:      config.GConfig.Company.RubbishClassify.SendTimes,
				AwardId:    *awd.Id,
				AccountId:  *awd.BpUid,
			}
		}
	}

	rcrModel.UpdateTransferFlag(*rcd.Id, 2)
	//
	//cmpId := config.GConfig.Company.Id
	//
	//for _, v := range cmpId {
	//	//get record list
	//	resp, err := this.GetStaffList(v, departId)
	//
	//	if err != nil {
	//		ZapLog().Error("record list request err", zap.Error(err))
	//		continue
	//	}
	//
	//	defer resp.Body.Close()
	//
	//	content, err := ioutil.ReadAll(resp.Body)
	//
	//	if err != nil {
	//		ZapLog().Error("response body read err", zap.Error(err))
	//		continue
	//	}
	//
	//	respData := new(api.StaffList)
	//	err = json.Unmarshal(content, respData)
	//
	//	if err != nil {
	//		ZapLog().Error("unmarshal response data err", zap.Error(err))
	//		continue
	//	}
	//
	//	if len(respData.Data) == 0 {
	//		continue
	//	}
	//
	//	symbol := config.GConfig.Company.RubbishClassify.Symbol
	//	merchantId := config.GConfig.Company.RubbishClassify.MerchantId
	//	coinArr := config.GConfig.Company.RubbishClassify.Coin
	//	rcModl := new(models.RubbishClassify)
	//
	//	for _, v := range respData.Data {
	//		if *v.Valid != 1 { // not equal one means leaved staff
	//			continue
	//		}
	//
	//		//calc how much coin to awarding
	//		coin := coinArr[score]
	//		//add send record
	//		awd, err := rcModl.AddNew(&v, coin, symbol, score)
	//
	//		if err != nil {
	//			return
	//		}
	//
	//		models.RcChan <- models.RClassifyChan{
	//			Coin:       coin,
	//			Symbol:     symbol,
	//			MerchantId: merchantId,
	//			Times:      config.GConfig.Company.RubbishClassify.SendTimes,
	//			AwardId:    *awd.Id,
	//			AccountId:  *awd.BpUid,
	//		}
	//	}
	//}
}

func (this *RClassify) GetStaffList(companyId, departId int) (resp *http.Response, err error) {
	url := config.GConfig.Company.ApiHost + "/v1/ft/merchant/teammanage/employee/gets"
	params := QueryListParams{
		CompanyId:    companyId,
		DepartmentId: departId,
	}
	jsonPramas, err := json.Marshal(params)
	b2 := bytes.NewBuffer(jsonPramas)
	resp, err = http.Post(url, "application/json;charset=utf-8", b2)

	return resp, err
}
