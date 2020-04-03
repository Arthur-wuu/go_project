package baspay

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/base"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/util"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
)

type (
	Assets struct {
		Assets *string `json:"assets,omitempty"`
		Money  *string `json:"money,omitempty"`
	}

	AssetsResponse struct {
		Code    *int         `json:"code,omitempty"`
		Data    *AssetsInfos `json:"data,omitempty"`
		Message *string      `json:"message,omitempty"`
	}

	AssetsInfos struct {
		Assets []AssetsInfo `json:"assets,omitempty"`
	}

	AssetsInfo struct {
		FullName *string `json:"full_name,omitempty"`
		Assets   *string `json:"assets,omitempty"`
		Logo     *string `json:"logo,omitempty"`
		//Url           *string       `json:"url,omitempty"`
		Price float64 `json:"price,omitempty"`
	}

	ToMoneyInfo struct {
		//Assets        *string       `json:"assets,omitempty"`
		Price float64 `json:"price,omitempty"`
	}

	QuoteResponse struct {
		Code   int               `json:"err,omitempty"`
		Quotes []QuoteDetailInfo `json:"quotes,omitempty" doc:"币简称"`
	}

	QuoteDetailInfo struct {
		Symbol     string      `json:"symbol,omitempty" doc:"币简称"`
		Id         int         `json:"id,omitempty" doc:"币id"`
		MoneyInfos []MoneyInfo `json:"detail,omitempty" doc:"币行情数据"`
	}

	MoneyInfo struct {
		Symbol             *string  `json:"symbol,omitempty" doc:"币简称"`
		Price              float64  `json:"price,omitempty" doc:"最新价"`
		Volume_24h         *float64 `json:"volume_24h,omitempty" doc:"24小时成交量"`
		Market_cap         *float64 `json:"market_cap,omitempty" doc:"总市值"`
		Percent_change_1h  *float64 `json:"percent_change_1h,omitempty" doc:"1小时涨跌幅"`
		Percent_change_24h *float64 `json:"percent_change_24h,omitempty" doc:"24小时涨跌幅"`
		Percent_change_7d  *float64 `json:"percent_change_7d,omitempty" doc:"7天涨跌幅"`
		Last_updated       *int64   `json:"last_updated,omitempty" doc:"最近更新时间"`
	}
)

func (this *Assets) Parse(f *api.AvAssets) *Assets {
	return &Assets{
		Assets: f.Assets,
		Money:  f.Money,
	}
}

func (this *Assets) Send() (interface{}, error) {
	//往baspay 查询订单状态
	signType := "RSA"
	timeStamp := GetTimeStamp()
	notifyUrl := config.GConfig.CallBack.NotifyUrl

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"assets":     this.Assets,
		"timestamp":  timeStamp,
		"notify_url": notifyUrl,
	})

	signStr := RequestBodyToSignStr(reqBodySign)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(CommonPrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return nil, err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"assets":     this.Assets,
		"sign_type":  signType,
		"signature":  finalSign,
		"timestamp":  timeStamp,
		"notify_url": notifyUrl,
	})

	result, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+"/open-api/trade/avail_assets", bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("send message to get assets info err", zap.Error(err))
		return nil, err
	}
	resAssetsInfo := new(AssetsResponse)
	err = json.Unmarshal(result, resAssetsInfo)
	if err != nil {
		ZapLog().Error("unmarshal err", zap.Error(err))
		return nil, err
	}
	ZapLog().Info("**resAssetsInfo**:", zap.Any("res:", resAssetsInfo))

	symbols := make([]string, 0)
	for _, v := range resAssetsInfo.Data.Assets {
		symbols = append(symbols, *v.Assets)
	}

	var fromSybolsStr string
	for i := 0; i < len(symbols); i++ {
		fromSybolsStr = fromSybolsStr + symbols[i] + ","
	}
	fromSybolsStr = fromSybolsStr[:len(fromSybolsStr)-1]

	quotrResult, err := base.HttpSend("http://quote.rkuan.com/api/v1/coin/quote?from="+fromSybolsStr+"&to="+*this.Money, bytes.NewBuffer(reqBody), "GET", nil)
	if err != nil {
		ZapLog().Error("select quote info err", zap.Error(err))
		return nil, err
	}

	quoteResponse := new(QuoteResponse)
	json.Unmarshal(quotrResult, quoteResponse)

	if err != nil {
		ZapLog().Error("unmarshal err", zap.Error(err))
		return nil, err
	}
	assetsInfos := make([]AssetsInfo, 0)
	quoteMap := make(map[string]float64)
	for i := 0; i < len(quoteResponse.Quotes); i++ {
		if len(quoteResponse.Quotes[i].MoneyInfos) == 0 {
			continue
		}

		quoteMap[quoteResponse.Quotes[i].Symbol] = quoteResponse.Quotes[i].MoneyInfos[0].Price
	}
	for i := 0; resAssetsInfo != nil && resAssetsInfo.Data != nil && i < len(resAssetsInfo.Data.Assets); i++ {
		if resAssetsInfo.Data.Assets[i].Assets == nil {
			continue
		}
		price, ok := quoteMap[*resAssetsInfo.Data.Assets[i].Assets]
		if !ok {
			continue
		}
		assetsInfo := new(AssetsInfo)
		assetsInfo.Assets = resAssetsInfo.Data.Assets[i].Assets
		assetsInfo.FullName = resAssetsInfo.Data.Assets[i].FullName
		assetsInfo.Price = price
		assetsInfo.Logo = resAssetsInfo.Data.Assets[i].Logo

		assetsInfos = append(assetsInfos, *assetsInfo)
	}
	return assetsInfos, nil

}
