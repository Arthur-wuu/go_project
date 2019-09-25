package config

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GConfig Config

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

type Config struct {
	Server   System   `yaml:"system"`
	Redis    Redis    `yaml:"redis"`
	Db       Mysql    `yaml:"mysql"`
	Award    Award    `yaml:"award"`
	Api      Api      `yaml:"api"`
	Gcache   Gcache   `yaml:"gcache"`
	Dingding Dingding `yaml:"dingding"`
	Company  Company  `yaml:"company"`
}

type System struct {
	Port    string `yaml:"port"`
	Debug   bool   `yaml:"debug"`
	LogPath string `yaml:"log_path"`
	Monitor string `yaml:"monitor"`
}

type Redis struct {
	Network     string `yaml:"network"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Password    string `yaml:"password"`
	Database    int    `yaml:"database"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxActive   int    `yaml:"maxActive"`
	IdleTimeout int    `yaml:"idleTimeout"`
	Prefix      string `yaml:"prefix"`
}

type Mysql struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	User          string `yaml:"user"`
	Pwd           string `yaml:"password"`
	Dbname        string `yaml:"dbname"`
	Charset       string `yaml:"charset"`
	Max_idle_conn int    `yaml:"maxIdle"`
	Max_open_conn int    `yaml:"maxOpen"`
	Debug         bool   `yaml:"debug"`
	ParseTime     bool   `yaml:"parseTime"`
}

type Award struct {
	Extratime struct {
		Coin   decimal.Decimal `yaml:"coin"`
		Symbol string          `yaml:"symbol"`
	} `yaml:"extratime"`
	Checkin struct {
		Coin   decimal.Decimal `yaml:"coin"`
		Symbol string          `yaml:"symbol"`
	} `yaml:"checkin"`
	MerchantId int    `yaml:"merchantId"`
	SendTimes  int    `yaml:"sendTimes"`
	ChanLen    int    `yaml:"chanLen"`
	AwardTime  string `yaml:"awardTime"`
}

type Api struct {
	Account   string `yaml:"account"`
	SecretKey string `yaml:"secretKey"`
}

type Gcache struct {
	SecretKey string `yaml:"secretKey"`
	Expire    int    `yaml:"expire"`
}

type Dingding struct {
	BoJu struct {
		AppKey    string `yaml:"appKey"`
		AppSecret string `yaml:"appSecret"`
	} `yaml:"boJu"`
	Host string `yaml:"host"`
}

type Company struct {
	Id           []int  `yaml:"id"`
	ApiHost      string `yaml:"apiHost"`
	ServiceAward struct {
		CoinBase   decimal.Decimal `yaml:"coinBase"`
		Symbol     string          `yaml:"symbol"`
		MerchantId int             `yaml:"merchantId"`
		SendTimes  int             `yaml:"sendTimes"`
		ChanLen    int             `yaml:"chanLen"`
	} `yaml:"serviceAward"`
	RubbishClassify struct {
		Coin       []decimal.Decimal `yaml:"coin"`
		Symbol     string            `yaml:"symbol"`
		MerchantId int               `yaml:"merchantId"`
		SendTimes  int               `yaml:"sendTimes"`
		ChanLen    int               `yaml:"chanLen"`
	} `yaml:"rubbishClassify"`
}
