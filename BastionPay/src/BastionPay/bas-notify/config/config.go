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
	return &GConfig
}

type Config struct {
	Server     Server     `yaml:"server"`
	Db         Db         `yaml:"db"`
	Aws        Aws        `yaml:"aws"`
	Sms        Sms        `yaml:"sms"`
	Email      Email      `yaml:"mail"`
	Monitor    Monitor    `yaml:"monitor"`
	ChuangLan  ChuangLan  `yaml:"chuanglan"`
	Twilio     Twilio     `yaml:"twilio"`
	Nexmo      Nexmo      `yaml:"nexmo"`
	YunTongXun YunTongXun `yaml:"yuntongxun"`
	DingDing   []DingDing `yaml:"dingding"`
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
	Debug         bool   `yaml:"debug"`
}

type Aws struct {
	Accesskeyid string `yaml:"accesskeyid"`
	Accesskey   string `yaml:"accesskey"`
	Accesstoken string `yaml:"accesstoken"`
	SesRegion   string `yaml:"sesregion"`
	SesTimeout  int    `yaml:"sestimeout"`
	SnsRegion   string `yaml:"snsregion"`
	SnsTimeout  int    `yaml:"snstimeout"`
	SrcEMail    string `yaml:"srcemail"`
	MetadataCmd string `yaml:"metadata_cmd"`
}

type Sms struct {
	LimitSend         int     `yaml:"limit_send"`
	FailRateThreshold float32 `yaml:"fail_rate_threshold"`
	FailThreshold     int     `yaml:"fail_threshold"`
}

type Email struct {
	LimitSend         int     `yaml:"limit_send"`
	FailRateThreshold float32 `yaml:"fail_rate_threshold"`
	FailThreshold     int     `yaml:"fail_threshold"`
}

type Monitor struct {
	TmpGNameMail string `yaml:"mail_tempgroup_name"`
	TmpLangMail  string `yaml:"mail_tempgroup_lang"`
	TmpGNameSms  string `yaml:"sms_tempgroup_name"`
	TmpLangSms   string `yaml:"sms_tempgroup_lang"`
}

type ChuangLan struct {
	Account string `yaml:"account"`
	Pwd     string `yaml:"pwd"`
	Url     string `yaml:"url"`
}

type Twilio struct {
	Sid       string `yaml:"account_sid"`
	Token     string `yaml:"auth_token"`
	FromPhone string `yaml:"from_phone"`
}

type Nexmo struct {
	ApiKey        string `yaml:"api_key"`
	ApiSecret     string `yaml:"api_secret"`
	Url           string `yaml:"url"`
	DefaultSender string `yaml:"default_sender"`
}

type YunTongXun struct {
	AppId      string `yaml:"app_id"`
	AuthToken  string `yaml:"auth_token"`
	Url        string `yaml:"url"`
	AccountSid string `yaml:"account_sid"`
}

type DingDing struct {
	QunName  string `yaml:"qun_name"`
	RobToken string `yaml:"rob_token"`
}
