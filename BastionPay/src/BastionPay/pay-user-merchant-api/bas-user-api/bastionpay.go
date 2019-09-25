package bas_user_api

import(
	"BastionPay/pay-user-merchant-api/common"
	"encoding/json"
	"fmt"
	"BastionPay/pay-user-merchant-api/config"
)

type  ScanLoginQrcode struct{

}

type  ResScanLoginQrcode struct{
	QrCode string `json:"login_qrcode"`
}

func (this * ScanLoginQrcode) Send() (string, error){
	resBytes,err := common.HttpSend(config.GConfig.BasUserApi.Addr + "/wallet/api/user/scan_login_qrcode", nil, "GET", nil)
	if err != nil {
		return "",err
	}
	res := new(ResScanLoginQrcode)
	if err = json.Unmarshal(resBytes, res); err != nil {
		return "",err
	}
	if len(res.QrCode) <= 0 {
		return "", fmt.Errorf("nil QrCode")
	}
	return res.QrCode,nil
}