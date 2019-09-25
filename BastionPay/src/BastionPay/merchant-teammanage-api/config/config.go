package config

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	PreProcess()
	return &GConfig
}

func PreProcess() error {
	//GPreConfig.DayMaxCoins = decimal.NewFromFloat(GConfig.Activity.DayMaxCoins)
	//GPreConfig.EachPlayerCoin = decimal.NewFromFloat(GConfig.Activity.EachPlayerCoin)
	return nil
}

type PreConfig struct {
	DayMaxCoins    decimal.Decimal
	EachPlayerCoin decimal.Decimal
}

type Config struct {
	Server System `yaml:"system"`
	//Redis   Redis       `yaml:"redis"`
	Db Mysql `yaml:"mysql"`
	//Cache   Cache       `yaml:"cache"`
	//User    User        `yaml:"user"`
	//Aws     Aws          `yaml:"aws"`
	//BussinessLimits BussinessLimits   `yaml:"bussiness_limits"`
	//BasNotify       BasNotify         `yaml:"bas_notify"`
	//TestPaper     TestPaper         `yaml:"test_paper"`
	//Activity      Activity           `yaml:"activity"`

}

type System struct {
	Port    string `yaml:"port"`
	Debug   bool   `yaml:"debug"`
	LogPath string `yaml:"log_path"`
	Monitor string `yaml:"monitor"`
	BkPort  string `yaml:"bk_port"`
}

type Cache struct {
	ActivityMaxKey         int `yaml:"activity_max_key"`
	ActivityTimeout        int `yaml:"activity_timeout"`
	SponsorActivityMaxKey  int `yaml:"sponsor_activity_max_key"`
	SponsorActivityTimeout int `yaml:"sponsor_activity_timeout"`
	SponsorMaxKey          int `yaml:"sponsor_max_key"`
	SponsorTimeout         int `yaml:"sponsor_timeout"`
	PageMaxKey             int `yaml:"page_max_key"`
	PageTimeout            int `yaml:"page_timeout"`
	ShareInfoMaxKey        int `yaml:"shareinfo_max_key"`
	ShareInfoTimeout       int `yaml:"shareinfo_timeout"`
	RobberMaxKey           int `yaml:"robber_max_key"`
	RobberTimeout          int `yaml:"robber_timeout"`
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

type User struct {
	Name string `yaml:"name"`
	Pwd  string `yaml:"pwd"`
}

type Aws struct {
	AccessKeyId   string `yaml:"accesskeyid"`
	AccessKey     string `yaml:"accesskey"`
	AccessToken   string `yaml:"accesstoken"`
	PicRegion     string `yaml:"picregion"`
	PicBucket     string `yaml:"picbucket"`
	PicBucketPath string `yaml:"picbucket_path"`
	PicTimeout    int    `yaml:"pictimeout"`
	CdnAddr       string `yaml:"cdn_addr"`
}

type BussinessLimits struct {
	PhoneSms []string `yaml:"phone_sms"`
	IpSms    []string `yaml:"ip_sms"`
}

type BasNotify struct {
	Addr             string `yaml:"addr"`
	VerifyCodeSmsTmp string `yaml:"verifycode_sms_tmp"`
}

type TestPaper struct {
	AllNum    int `yaml:"all_num"`
	PassScore int `yaml:"pass_score"`
}

type Activity struct {
	DayMaxPlayer   int     `yaml:"day_max_player"`
	DayMaxCoins    float64 `yaml:"day_max_coins"`
	EachPlayerCoin float64 `yaml:"each_player_coin"`
	Symbol         string  `yaml:"symbol"`
	NotifyUrl      string  `yaml:"notify_url"`
	SponsorAccount string  `yaml:"sponsor_account"`
	ApiKey         string  `yaml:"api_key"`
	OffAt          int64   `yaml:"off_at"`
	Language       string  `yaml:"language"`
	OnAt           int64   `yaml:"on_at"`
}
