package config

import (
	"fmt"
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
	PreProcess()
	return &GConfig
}

func PreProcess() {
	//if GConfig.Cache.VipAuthMaxKey < 5 {
	//	GConfig.Cache.VipAuthMaxKey = 5
	//}
	//if GConfig.Cache.VipAuthTimeout < 5 {
	//	GConfig.Cache.VipAuthTimeout = 60
	//}
	//if GConfig.Cache.VipListMaxKey < 2 {
	//	GConfig.Cache.VipListMaxKey = 2
	//}
	//if GConfig.Cache.VipListTimeout < 5 {
	//	GConfig.Cache.VipListTimeout = 60
	//}
	//if GConfig.Cache.VipDisableMaxKey < 3 {
	//	GConfig.Cache.VipDisableMaxKey = 5
	//}
	//if GConfig.Cache.VipDisableTimeout < 5 {
	//	GConfig.Cache.VipDisableTimeout = 60
	//}
	//if strings.Contains(GConfig.CallBack.ShowUrl, "*") {
	//	//获取本地ip替换掉
	//}
}

type Config struct {
	Server   System      `yaml:"system"`
	Redis    Redis       `yaml:"redis"`
	Db       Mysql       `yaml:"mysql"`
	Cache    Cache       `yaml:"cache"`
	CallBack CallBack    `yaml:"callback"`
	Devices  []Device    `yaml:"device"`
	BastionpayUrl  BastionpayUrl    `yaml:"bastionpay_url"`
	//DeviceMap  map[string]*Device    `yaml:"-"`
	Login   Login        `yaml:"login"`
	PayeeId PayeeId      `yaml:"payeeid"`
	Fee     Fee          `yaml:"fee"`
}

type BastionpayUrl struct{
	Bastionurl   string          `yaml:"bastionurl"`
	BastionUser  string          `yaml:"bastionuser"`
	QuoteUrl     string          `yaml:"quoteurl"`
}

//coffee 收款人 和商户
type PayeeId struct{
	PayId   string          `yaml:"payid"`
	MerchantId string       `yaml:"merchantid"`
}

type Login struct {
	Uid           string `yaml:"pre_uid"`
	Pwd           string `yaml:"password"`
	Phone         string `yaml:"phone"`
	ZfPwd         string `yaml:"zf_pwd"`
	DeviceId      string `yaml:"device_id"`
}

type Device struct{
	Id   string          `yaml:"id"`
	Name string          `yaml:"name"`
	Addr string          `yaml:"addr"`
}

type CallBack struct {
	ReturnUrl    string      `yaml:"return_url"`
	ShowUrl      string      `yaml:"show_url"`
	NotifyUrl    string      `yaml:"notify_url"`
}

type System struct {
	Port    string      `yaml:"port"`
	Debug   bool        `yaml:"debug"`
	LogPath string      `yaml:"log_path"`
	Monitor string      `yaml:"monitor"`
}

type Cache struct {
	VipAuthMaxKey    int      `yaml:"vipauth_max_key"`
	VipAuthTimeout   int      `yaml:"vipauth_timeout"`
	VipListMaxKey    int      `yaml:"viplist_max_key"`
	VipListTimeout   int      `yaml:"viplist_timeout"`
	VipDisableMaxKey    int      `yaml:"vipdisable_max_key"`
	VipDisableTimeout   int      `yaml:"vipdisable_timeout"`
}

type Redis struct {
	Network     string  `yaml:"network"`
	Host        string  `yaml:"host"`
	Port        string  `yaml:"port"`
	Password    string  `yaml:"password"`
	Database    string  `yaml:"database"`
	MaxIdle     int     `yaml:"maxIdle"`
	MaxActive   int     `yaml:"maxActive"`
	IdleTimeout int     `yaml:"idleTimeout"`
	Prefix      string  `yaml:"prefix"`
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


type Fee struct {
	Cash2coin   string      `yaml:"cash2coin"`
	Coin2cash   string      `yaml:"coin2cash"`
}