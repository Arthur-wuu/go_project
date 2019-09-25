package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"errors"
)

var GConfig Config
var GPreConfig PreConfig

func LoadConfig(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Read yml config[%s] err[%v]", path, err).Error())
	}

	err = yaml.Unmarshal([]byte(data), &GConfig)
	if err != nil {
		panic(fmt.Errorf("yml.Unmarshal config[%s] err[%v]", path, err).Error())
	}
	fmt.Println("data===", GConfig.Server.Port, GConfig.CoinMerit.HttpUrl)
	return &GConfig
}

func PreProcess() error {
	if GPreConfig.MarketMap == nil {
		GPreConfig.MarketMap = make(map[string] *Market)
	}
	if GPreConfig.CoinDetailMap == nil {
		GPreConfig.CoinDetailMap = make(map[string] *CoinDetail)
	}
	if GPreConfig.Markets == nil {
		GPreConfig.Markets = make([]string, 0)
	}
	marketArr := make([]*Market,0)
	for i:=0; i < len(GConfig.Markets); i++ {
		AbbArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.Markets[i].Abb, " ", "", -1)), ","), ",")
		NameArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.Markets[i].Name, " ", "", -1)), ","), ",")
		zhNameArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.Markets[i].ZhName, " ", "", -1)), ","), ",")
		webInfoArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.Markets[i].WebInfo, " ", "", -1)), ","), ",")
		if len(AbbArr) != len(NameArr) || len(NameArr) != len(zhNameArr) || len(zhNameArr) != len(webInfoArr) {
			fmt.Println(AbbArr)
			fmt.Println(NameArr)
			fmt.Println(zhNameArr)
			fmt.Println(webInfoArr)

			return errors.New("Markets config err "+ fmt.Sprintf("%v %v %v %v", len(AbbArr), len(NameArr), len(zhNameArr), len(webInfoArr)))
		}
		for j:=0; j < len(AbbArr); j++ {
			market := new(Market)
			market.Abb = AbbArr[j]
			market.Name = NameArr[j]
			market.ZhName = zhNameArr[j]
			market.WebInfo = webInfoArr[j]
			GPreConfig.MarketMap[market.Abb] = market
			marketArr = append(marketArr, market)
			GPreConfig.Markets = append(GPreConfig.Markets, AbbArr[j])
		}
	}
	GConfig.Markets = marketArr

	coinDetailArr := make([]*CoinDetail,0)
	for i:=0; i < len(GConfig.CoinDetails); i++ {
		NameArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.CoinDetails[i].Name, " ", "", -1)), ","), ",")
		zhNameArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.CoinDetails[i].ZhName, " ", "", -1)), ","), ",")
		pinyinArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.CoinDetails[i].PinYin, " ", "", -1)), ","), ",")
		if  len(NameArr) != len(zhNameArr) || len(zhNameArr) != len(pinyinArr) {
			fmt.Println(NameArr)
			fmt.Println(zhNameArr)
			fmt.Println(pinyinArr)
			return errors.New("CoinDetail config err "+fmt.Sprintf("%v %v %v", len(NameArr) ,len(zhNameArr),len(pinyinArr)))
		}
		for j:=0; j < len(NameArr); j++ {
			deatil := new(CoinDetail)
			deatil.Name = NameArr[j]
			deatil.ZhName = zhNameArr[j]
			deatil.PinYin = pinyinArr[j]
			GPreConfig.CoinDetailMap[deatil.Name] = deatil
			coinDetailArr = append(coinDetailArr, deatil)
		}
	}
	GConfig.CoinDetails = coinDetailArr

	return nil
}

type PreConfig struct{
	CoinDetailMap map[string] *CoinDetail //HBI
	MarketMap     map[string] *Market   //HBI
	Markets       []string             // HBI、BIC
}

type Config struct {
	Server  Server                `yaml:"server"`
	CoinMerit  CoinMerit          `yaml:"coinmerit"`
	BtcExa     BtcExa             `yaml:"btcexa"`
	Markets     []*Market           `yaml:"market"`
	CoinDetails []*CoinDetail       `yaml:"coin_detail"`
}

type CoinMerit struct {
	Secret_key string   `yaml:"secret_key"`
	ApiKey string       `yaml:"api_key"`
	WsUrl   string      `yaml:"ws_url"`
	HttpUrl string      `yaml:"http_url"`
}


type Server struct {
	Port    string `yaml:"port"`
	Debug   bool   `yaml:"debug"`
	LogPath string `yaml:"log"`
}

type BtcExa struct {
	Secret_key string   `yaml:"secret_key"`
	ApiKey string       `yaml:"api_key"`
	WsUrl   string      `yaml:"ws_url"`
	HttpUrl string      `yaml:"http_url"`
}

type Market struct {
	Abb   string         `yaml:"abb"`
	Name  string         `yaml:"name"`
	ZhName string        `yaml:"zh_name"`
	WebInfo string       `yaml:"web_info"`
}

type CoinDetail struct {
	Name   string         `yaml:"name"`
	ZhName string         `yaml:"zh_name"`
	PinYin string         `yaml:"pinyin"`
}