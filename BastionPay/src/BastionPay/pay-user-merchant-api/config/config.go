package config

import (
	"BastionPay/pay-user-merchant-api/common"
	"fmt"
	"github.com/ulule/limiter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
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
	PreProcess()
	return &GConfig
}

func PreProcess() error {
	GPreConfig.PhoneSmsLimits = make([]*limiter.Rate, 0)
	GPreConfig.IpSmsLimits = make([]*limiter.Rate, 0)
	GPreConfig.PhoneEmailLimits = make([]*limiter.Rate, 0)
	GPreConfig.IpEmailLimits = make([]*limiter.Rate, 0)
	for i := 0; i < len(GConfig.BussinessLimits.PhoneSms); i++ {
		rate, err := limiter.NewRateFromFormatted(GConfig.BussinessLimits.PhoneSms[i])
		if err != nil {
			return err
		}
		GPreConfig.PhoneSmsLimits = append(GPreConfig.PhoneSmsLimits, &rate)
	}
	for i := 0; i < len(GConfig.BussinessLimits.IpSms); i++ {
		rate, err := limiter.NewRateFromFormatted(GConfig.BussinessLimits.IpSms[i])
		if err != nil {
			return err
		}
		GPreConfig.IpSmsLimits = append(GPreConfig.IpSmsLimits, &rate)
	}
	for i := 0; i < len(GConfig.BussinessLimits.PhoneMail); i++ {
		rate, err := limiter.NewRateFromFormatted(GConfig.BussinessLimits.PhoneMail[i])
		if err != nil {
			return err
		}
		GPreConfig.PhoneEmailLimits = append(GPreConfig.PhoneEmailLimits, &rate)
	}
	for i := 0; i < len(GConfig.BussinessLimits.IpMail); i++ {
		rate, err := limiter.NewRateFromFormatted(GConfig.BussinessLimits.IpMail[i])
		if err != nil {
			return err
		}
		GPreConfig.IpEmailLimits = append(GPreConfig.IpEmailLimits, &rate)
	}
	GPreConfig.PathWhiteList = make(map[string]struct{})
	for i := 0; i < len(GConfig.PathWhiteList.Paths); i++ {
		GPreConfig.PathWhiteList[GConfig.PathWhiteList.Paths[i]] = struct{}{}
	}
	GPreConfig.PathLimits = make([]*common.PathLimit, 0)
	for i := 0; i < len(GConfig.PathLimits); i++ {
		temp := &common.PathLimit{
			Path:   GConfig.PathLimits[i].Path,
			Method: GConfig.PathLimits[i].Method,
			Limit:  GConfig.PathLimits[i].Limit,
			Time:   GConfig.PathLimits[i].Time,
		}
		GPreConfig.PathLimits = append(GPreConfig.PathLimits, temp)
	}
	GPreConfig.TermBlockLimits = make([]*common.TermBlock, 0)
	for i := 0; i < len(GConfig.TermBlockLimits); i++ {
		temp := &common.TermBlock{
			Name:     GConfig.TermBlockLimits[i].Name,
			Locktime: GConfig.TermBlockLimits[i].Locktime,
			Limit:    GConfig.TermBlockLimits[i].Limit,
			Time:     GConfig.TermBlockLimits[i].Time,
		}
		GPreConfig.TermBlockLimits = append(GPreConfig.TermBlockLimits, temp)
	}
	GConfig.Token.Expiration = strings.Replace(GConfig.Token.Expiration, " ", "", -1)
	if len(GConfig.Token.Expiration) > 0 {
		GConfig.Token.Expirat, _ = strconv.ParseInt(GConfig.Token.Expiration, 10, 64)
		if GConfig.Token.Expirat <= 0 {
			GConfig.Token.Expirat = 3600
		}
	}
	return nil
}

type PreConfig struct {
	PhoneSmsLimits   []*limiter.Rate
	IpSmsLimits      []*limiter.Rate
	PhoneEmailLimits []*limiter.Rate
	IpEmailLimits    []*limiter.Rate
	PathLimits       []*common.PathLimit
	TermBlockLimits  []*common.TermBlock
	PathWhiteList    map[string]struct{}
}

type Config struct {
	Server          System          `yaml:"server"`
	Redis           Redis           `yaml:"redis"`
	Db              Mysql           `yaml:"db"`
	Cache           Cache           `yaml:"cache"`
	Aws             Aws             `yaml:"aws"`
	Token           Token           `yaml:"token"`
	IpFind          IpFind          `yaml:"ip_find"`
	BasMonitor      BasMonitor      `yaml:"monitor"`
	BasNotify       BasNotify       `yaml:"bas_notify"`
	BussinessLimits BussinessLimits `yaml:"bussiness_limits"`
	PathLimits      []*PathLimit    `yaml:"path_limits"`
	TermBlockLimits []*TermBlock    `yaml:"termblock_limits"`
	PathWhiteList   PathWhiteList   `yaml:"path_white_list"`
	BasUserApi      BasUserApi      `yaml:"bas_user_api"`
}

type System struct {
	Port    string `yaml:"port"`
	Debug   bool   `yaml:"debug"`
	LogPath string `yaml:"log_path"`
	Monitor string `yaml:"monitor"`
	BkPort  string `yaml:"bk_port"`
}

type Cache struct {
	//RobberMaxKey          int `yaml:"robber_max_key"`
	//RobberTimeout         int   `yaml:"robber_timeout"`
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
	Dbname        string `yaml:"db"`
	Charset       string `yaml:"charset"`
	Max_idle_conn int    `yaml:"maxIdle"`
	Max_open_conn int    `yaml:"maxOpen"`
	Debug         bool   `yaml:"debug"`
	ParseTime     bool   `yaml:"parseTime"`
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

type BasMonitor struct {
	Version string   `yaml:"version"`
	Name    string   `yaml:"name"`
	Tag     string   `yaml:"tag"`
	RpcAddr string   `yaml:"rpc_addr"`
	Env     []string `yaml:"env"`
}

type IpFind struct {
	Auth string `yaml:"auth"`
}

type Token struct {
	Secret     string `yaml:"secret"`
	Expiration string `yaml:"expiration"`
	Expirat    int64  `yaml:"-"`
}

type PathWhiteList struct {
	Paths []string `yaml:"path"`
}

type PathLimit struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Limit  int    `yaml:"limit"`
	Time   int    `yaml:"time"`
}

type TermBlock struct {
	Name     string `yaml:"name"`
	Limit    int64  `yaml:"limit"`
	Time     int    `yaml:"time"`
	Locktime int    `yaml:"lock_time"`
}

type BussinessLimits struct {
	PhoneSms  []string `yaml:"phone_sms"`
	IpSms     []string `yaml:"ip_sms"`
	PhoneMail []string `yaml:"phone_mail"`
	IpMail    []string `yaml:"ip_mail"`
}

type BasNotify struct {
	Addr                   string `yaml:"addr"`
	VerifyCodeSmsTmp       string `yaml:"verifycode_sms_tmp"`
	VerifyCodeMailTmp      string `yaml:"verifycode_mail_tmp"`
	RegisterSuccessSmsTmp  string `yaml:"register_ok_sms_tmp"`
	RegisterSuccessMailTmp string `yaml:"register_ok_mail_tmp"`
}

type BasUserApi struct {
	Addr string `yaml:"addr"`
}
