package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
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
	return &GConfig
}

func PreProcess() error {
	GConfig.Parities.Country_code = strings.ToUpper(strings.Replace(GConfig.Parities.Country_code, " ", "", -1))
	GConfig.Parities.Country_name = strings.Replace(GConfig.Parities.Country_name, " ", "", -1)
	GConfig.Coinmarketcap.Coins = strings.ToUpper(strings.Replace(GConfig.Coinmarketcap.Coins, " ", "", -1))
	GConfig.Collect.Coin_num = strings.Replace(GConfig.Collect.Coin_num, " ", "", -1)
	GConfig.Collect.Coin_from = strings.Replace(GConfig.Collect.Coin_from, " ", "", -1)
	GConfig.Collect.Coin_to = strings.Replace(GConfig.Collect.Coin_to, " ", "", -1)
	GConfig.Collect.Coin_pairs = strings.Replace(GConfig.Collect.Coin_pairs, " ", "", -1)
	GConfig.Collect.Coin_exchange = strings.Replace(GConfig.Collect.Coin_exchange, " ", "", -1)
	GConfig.Collect.Coin_entrance = strings.Replace(GConfig.Collect.Coin_entrance, " ", "", -1)
	GConfig.Collect.Open_Flag = strings.Replace(GConfig.Collect.Open_Flag, " ", "", -1)

	cNumArr := strings.Split(GConfig.Collect.Coin_num, ",")
	cSymbolArr := strings.Split(GConfig.Collect.Coin_from, ",")
	cToArr := strings.Split(GConfig.Collect.Coin_to, ",")
	cPairsArr := strings.Split(GConfig.Collect.Coin_pairs, ",")
	cEntrance := strings.Split(GConfig.Collect.Coin_entrance, ",")
	cExchangeArr := strings.Split(GConfig.Collect.Coin_exchange, ",")
	if len(GConfig.Collect.Open_Flag) == 0 {
		tempArr := make([]string, 0)
		for i := 0; i < len(cNumArr); i++ {
			tempArr = append(tempArr, "1")
		}
		GConfig.Collect.Open_Flag = strings.Join(tempArr, ",")
	}
	cOpenFlagArr := strings.Split(GConfig.Collect.Open_Flag, ",")

	GPreConfig.FromCollects = make(map[string]*Collect)
	GPreConfig.IdsCollects = make(map[string]*Collect)
	for i := 0; i < len(cNumArr); i++ {
		collectTmp := *new(Collect)
		collectTmp.SetCoin_num(cNumArr[i])
		collectTmp.SetCoin_from(cSymbolArr[i])
		collectTmp.SetCoin_to(cToArr[i])
		collectTmp.SetCoin_pairs(cPairsArr[i])
		collectTmp.SetCoin_entrance(cEntrance[i])
		collectTmp.SetCoin_exchange(cExchangeArr[i])
		collectTmp.SetOpen_Flag(cOpenFlagArr[i])
		GPreConfig.FromCollects[cSymbolArr[i]] = &collectTmp

		GPreConfig.IdsCollects[cNumArr[i]] = &collectTmp
	}

	GPreConfig.MarketCollects = make(map[string][]*Collect)
	for j := 0; j < len(cEntrance); j++ {
		marketCollectBtcexa := make([]*Collect, 0)
		marketCollectCoimkt := make([]*Collect, 0)
		if cEntrance[j] == "btcexa" {
			collectTmpM := *new(Collect)
			collectTmpM.SetCoin_num(cNumArr[j])
			collectTmpM.SetCoin_from(cSymbolArr[j])
			collectTmpM.SetCoin_to(cToArr[j])
			collectTmpM.SetCoin_pairs(cPairsArr[j])
			collectTmpM.SetCoin_entrance(cEntrance[j])
			collectTmpM.SetCoin_exchange(cExchangeArr[j])
			collectTmpM.SetOpen_Flag(cOpenFlagArr[j])
			marketCollectBtcexa = append(marketCollectBtcexa, &collectTmpM)
		}
		if cEntrance[j] == "coinmerit" {
			collectTmpM := *new(Collect)
			collectTmpM.SetCoin_num(cNumArr[j])
			collectTmpM.SetCoin_from(cSymbolArr[j])
			collectTmpM.SetCoin_to(cToArr[j])
			collectTmpM.SetCoin_pairs(cPairsArr[j])
			collectTmpM.SetCoin_entrance(cEntrance[j])
			collectTmpM.SetCoin_exchange(cExchangeArr[j])
			collectTmpM.SetOpen_Flag(cOpenFlagArr[j])
			marketCollectCoimkt = append(marketCollectCoimkt, &collectTmpM)
		}
		GPreConfig.MarketCollects["coinmerit"] = marketCollectCoimkt
		GPreConfig.MarketCollects["btcexa"] = marketCollectBtcexa
	}
	//配置文件中的181国家和国家码一一加载对应
	cNameArr := make([]string, 0)
	cCodeArr := make([]string, 0)

	for i := 0; i < len(GConfig.SinaParities); i++ {
		cNameTmpArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.SinaParities[i].Contry_name, " ", "", -1)), ","), ",")
		cCodeTmpArr := strings.Split(strings.TrimRight(strings.ToUpper(strings.Replace(GConfig.SinaParities[i].Contry_code, " ", "", -1)), ","), ",")
		for j := 0; j < len(cNameTmpArr); j++ {
			cNameArr = append(cNameArr, cNameTmpArr[j])
			cCodeArr = append(cCodeArr, cCodeTmpArr[j])
		}

		if len(cNameArr) != len(cCodeArr) {
			fmt.Println(cNameArr)
			fmt.Println(cCodeArr)
			return errors.New("SinaParitie config err " + fmt.Sprintf("%v %v", len(cNameArr), len(cCodeArr)))
		}
	}
	GPreConfig.CountryCodeArr = cCodeArr
	GPreConfig.CountryNameArr = cNameArr

	GConfig.Coinmarketcap.Api_key = strings.Replace(GConfig.Coinmarketcap.Api_key, " ", "", -1)
	GConfig.Coinmarketcap.Api_key = strings.TrimRight(GConfig.Coinmarketcap.Api_key, ",")
	GPreConfig.CoinmarketcapApiKeys = strings.Split(GConfig.Coinmarketcap.Api_key, ",")
	return nil
}

type PreConfig struct {
	MarketCollects       map[string][]*Collect //市场划分
	FromCollects         map[string]*Collect   //from币划分
	CountryNameArr       []string              //sina huilv
	CountryCodeArr       []string              //sina huilv
	CoinmarketcapApiKeys []string
	IdsCollects          map[string]*Collect //from币划分

}

type Config struct {
	Server        Server         `yaml:"server"`
	Db            Db             `yaml:"db"`
	Coinmarketcap Coinmarketcap  `yaml:"coinmarketcap"`
	Redis         Redis          `yaml:"redis"`
	Collect       Collect        `yaml:"collect"`
	CoinMerit     CoinMerit      `yaml:"coinmerit"`
	Parities      Parities       `yaml:"parities"`
	SinaParities  []*SinaParitie `yaml:"sina_parities"`
	Switch        Switch         `yaml:"switch"`
}

//sina汇率
type SinaParitie struct {
	Contry_name string `yaml:"cty_name"`
	Contry_code string `yaml:"cty_code"`
}

//汇率开关
type Switch struct {
	FxSinaFlag  bool `yaml:"fx_sina"`
	FxBaiduFlag bool `yaml:"fx_baidu"`
}

//redis
type Redis struct {
	Network     string `yaml:"network"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxActive   int    `yaml:"maxActive"`
	IdleTimeout int    `yaml:"idleTimeout"`
	Prefix      string `yaml:"prefix"`
}

type Server struct {
	Port    string `yaml:"port"`
	Debug   bool   `yaml:"debug"`
	LogPath string `yaml:"log"`
}

type Db struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	User          string `yaml:"user"`
	Pwd           string `yaml:"password"`
	Quote_db      string `yaml:"db"`
	Max_idle_conn int    `yaml:"max_idle_conn"`
	Max_open_conn int    `yaml:"max_open_conn"`
}

//type Cache struct {
//	Leveldb_path string `yaml:"leveldb_path"`
//}

type Coinmarketcap struct {
	CoinMarket_url     string `yaml:"coinmarket_url"`
	CoinMarket_new_url string `yaml:"coinmarket_new_url"`
	Quote_interval     int    `yaml:"quote_interval"`
	Code_interval      int    `yaml:"code_interval"`
	Coins              string `yaml:"coins"`
	Err_interval       int    `yaml:"err_interval"`
	Api_key            string `yaml:"api_key"`
	Diff_env_interval  int    `yaml:"diff_env_interval"`
	NewApiFlag         bool   `yaml:"new_api_flag"`
}

type Parities struct {
	Country_name string `yaml:"country_name"`
	Country_code string `yaml:"country_code"`
}

type CoinMerit struct {
	Secret_key string `yaml:"secret_key"`
	ApiKey     string `yaml:"api_key"`
	HttpUrl    string `yaml:"http_url"`
}

type Collect struct {
	Coin_num      string `yaml:"coin_from_num"`
	Coin_from     string `yaml:"coin_from"`
	Coin_to       string `yaml:"coin_to"`
	Coin_pairs    string `yaml:"coin_pairs"`
	Coin_exchange string `yaml:"coin_exchange"`
	Coin_entrance string `yaml:"coin_entrance"`
	Open_Flag     string `yaml:"open_flag"`
}

func (this *Collect) SetCoin_num(p string) {
	this.Coin_num = p
}

func (this *Collect) SetCoin_from(p string) {
	this.Coin_from = p
}
func (this *Collect) SetCoin_to(p string) {
	this.Coin_to = p
}
func (this *Collect) SetCoin_pairs(p string) {
	this.Coin_pairs = p
}
func (this *Collect) SetCoin_exchange(p string) {
	this.Coin_exchange = p
}
func (this *Collect) SetCoin_entrance(p string) {
	this.Coin_entrance = p
}
func (this *Collect) SetOpen_Flag(p string) {
	this.Open_Flag = p
}
