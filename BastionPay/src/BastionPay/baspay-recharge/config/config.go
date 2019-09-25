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
	//PreProcess()
	GAccountMgr.Init()
	return &GConfig
}

func PreProcess() {
	if GConfig.BasNotify.SmsDelay < 1 {
		GConfig.BasNotify.SmsDelay = 7200
	}
	if GConfig.BasNotify.SmsIntvl < 1 {
		GConfig.BasNotify.SmsIntvl = 10
	}
}

var GAccountMgr AccountMgr

type AccountMgr struct {
	mAccountMap map[string]LoginPool
}

func (this *AccountMgr) Init() error {
	this.mAccountMap = make(map[string]LoginPool)
	for i := 0; i < len(GConfig.LoginPool); i++ {
		user := GConfig.LoginPool[i]

		login := new(LoginPool)
		login.Uid = user.Uid
		login.Phone = user.Phone
		login.DeviceId = user.DeviceId
		login.ZfPwd = user.ZfPwd
		login.Pwd = user.Pwd

		this.Add(user.Uid, *login)

	}
	return nil
}

func (this *AccountMgr) Add(id string, d LoginPool) {
	this.mAccountMap[id] = d
}

func (this *AccountMgr) Get(id string) LoginPool {
	return this.mAccountMap[id]
}

type Config struct {
	Server System `yaml:"system"`
	Redis  Redis  `yaml:"redis"`
	Db     Mysql  `yaml:"mysql"`
	//Aws     Aws          `yaml:"aws"`
	BastionpayUrl BastionpayUrl `yaml:"bastionpay_url"`
	RedPackUrl    RedPackUrl    `yaml:"redpack_url"`
	LuckDrawUrl   LuckDrawUrl   `yaml:"luckdraw_url"`
	Login         Login         `yaml:"login"`
	BasNotify     BasNotify     `yaml:"bas_notify"`
	LoginPool     []LoginPool   `yaml:"login_pool"`
}

type BasNotify struct {
	Addr     string `yaml:"addr"`
	SmsDelay int64  `yaml:"fission_rob_delay"`
	SmsIntvl int64  `yaml:"fission_rob_intvl"`
	SmsTemp  string `yaml:"fission_rob_temp"`
}
type BastionpayUrl struct {
	Bastionurl string `yaml:"bastionurl"`
}

type RedPackUrl struct {
	RedPackUrl string `yaml:"redpackurl"`
}

type LuckDrawUrl struct {
	LuckDrawUrl string `yaml:"luckdrawurl"`
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
	Database    string `yaml:"database"`
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

//
//type Aws struct{
//	AccessKeyId string   `yaml:"accesskeyid"`
//	AccessKey  string    `yaml:"accesskey"`
//	AccessToken string   `yaml:"accesstoken"`
//	PicRegion  string    `yaml:"picregion"`
//	PicBucket  string    `yaml:"picbucket"`
//	PicTimeout int    `yaml:"pictimeout"`
//
//}

type Login struct {
	Uid      string `yaml:"pre_uid"`
	Pwd      string `yaml:"password"`
	Phone    string `yaml:"phone"`
	ZfPwd    string `yaml:"zf_pwd"`
	DeviceId string `yaml:"device_id"`
}

type LoginPool struct {
	Uid      string `yaml:"pre_uid"`
	Pwd      string `yaml:"password"`
	Phone    string `yaml:"phone"`
	ZfPwd    string `yaml:"zf_pwd"`
	DeviceId string `yaml:"device_id"`
}
