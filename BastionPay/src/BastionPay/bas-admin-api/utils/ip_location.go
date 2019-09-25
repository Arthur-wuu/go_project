package utils

import (
	"encoding/json"
	"github.com/BastionPay/bas-admin-api/common"
)

const (
	//ipLocationUrl = "http://ip.taobao.com/service/getIpInfo.php"
	ipLocationUrl = "https://ipfind.co/"
)

var (
	// 修正政治正确
	FixedCountry = map[string]string{"Hong Kong": "China", "Macao": "China", "Taiwan": "China"}
)

//type IpLocationResult struct {
//	Code int
//	Data struct {
//		Ip        string
//		Country   string
//		Area      string
//		Region    string
//		City      string
//		County    string
//		Isp       string
//		CountryId string `json:"country_id"`
//		AreaId    string `json:"area_id"`
//		CityId    string `json:"city_id"`
//		CountyId  string `json:"county_id"`
//		IspId     string `json:"isp_id"`
//	}
//}

//{"ip_address":"60.249.116.191","country":"Taiwan","country_code":"TW","continent":"Asia","continent_code":"AS",
// "city":"Taichung","county":"Taichung City","region":"Taiwan", "region_code":"04","timezone":"Asia\/Taipei",
// "owner":null,"longitude":120.6839,"latitude":24.1469,"currency":"TWD","languages":["zh-TW","zh","nan","hak"]}
type IpLocationResult struct {
	IpAddress     string `json:"ip_address"`
	Country       string
	CountryCode   string `json:"country_code"`
	Continent     string
	ContinentCode string `json:"continent_code"`
	City          string
	County        string
	Region        string
	RegionCode    string `json:"region_code"`
	Timezone      string
	Owner         string
	Longitude     float64
	Latitude      float64
	Currency      string
}

func IpLocation(auth string, ip string) (*IpLocationResult, error) {
	var (
		res = &IpLocationResult{}
		err error
	)
	result, err := common.NewHttp().Get(ipLocationUrl + "?auth=" + auth + "&ip=" + ip)
	if err != nil {
		return nil, err
	}

	res = &IpLocationResult{}

	err = json.Unmarshal(result, res)
	if err != nil {
		return nil, err
	}

	if _, ok := FixedCountry[res.Country]; ok {
		res.Country = FixedCountry[res.Country]
	}

	return res, nil
}
