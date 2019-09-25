package tools

import (
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/bastionpay"
	l4g "github.com/alecthomas/log4go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type (
	Config struct {
		System      System             `yaml:"system"`
		Mysql       Mysql              `yaml:"mysql"`
		Redis       Redis              `yaml:"redis"`
		Wallet      bastionpay.Gateway `yaml:"wallet"`
		WalletPaths []string           `yaml:"wallet_paths"`
		BasAmin     BasAmin            `yaml:"bas_admin"`
		CoinMarket  CoinMarket         `yaml:"coin_market"`
		BasQuote    BasQuote           `yaml:"bas_quote"`
		Monitor     Monitor            `yaml:"monitor"`
		Aws         Aws                `yaml:"aws"`
		Notify      Notify             `yaml:"notify"`
		IpFind      IpFind             `yaml:"ip_find"`
		OperateLog  OperateLog         `yaml:"operate_log"`
	}

	System struct {
		Host        string        `yaml:"host"`
		Port        int           `yaml:"port"`
		Debug       bool          `yaml:"debug"`
		LogPath     string        `yaml:"logPath"`
		CompanyName string        `yaml:"company_name"`
		Expire      time.Duration `yaml:"expire"`
		GaExpire    time.Duration `yaml:"gaExpire"`
	}

	Mysql struct {
		Dialect   string `yaml:"dialect"`
		Host      string `yaml:"host"`
		Port      int64  `yaml:"port"`
		DbName    string `yaml:"dbname"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		Charset   string `yaml:"charset"`
		ParseTime bool   `yaml:"parseTime"`
		MaxIdle   int    `yaml:"maxIdle"`
		MaxOpen   int    `yaml:"maxOpen"`
		Debug     bool   `yaml:"debug"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	}

	CoinMarket struct {
		Url    string `yaml:"url"`
		IdPath string `yaml:"id_path"`
	}

	BasAmin struct {
		Url string `yaml:"url"`
	}

	BasQuote struct {
		Addr string `yaml:"addr"`
	}

	Monitor struct {
		Addr string `yaml:"addr"`
	}

	Aws struct {
		AccessKeyId   string `yaml:"accesskeyid"`
		AccessKey     string `yaml:"accesskey"`
		AccessToken   string `yaml:"accesstoken"`
		LogoRegion    string `yaml:"logoregion"`
		LogoBucket    string `yaml:"logobucket"`
		LogoTimeout   int    `yaml:"logotimeout"`
		SesRegion     string `yaml:"sesregion"`
		SesTimeout    int    `yaml:"sestimeout"`
		SnsRegion     string `yaml:"snsregion"`
		SnsTimeout    int    `yaml:"snstimeout"`
		NoticeRegion  string `yaml:"noticeregion"`
		NoticeBucket  string `yaml:"noticebucket"`
		NoticeTimeout int    `yaml:"noticetimeout"`
		NotifyRegion  string `yaml:"notifyregion"`
		NotifyBucket  string `yaml:"notifybucket"`
		NotifyTimeout int    `yaml:"notifytimeout"`
	}

	Notify struct {
		UserNotify []string `yaml:"usernotify"`
		SysNotify  []string `yaml:"sysnotify"`
		MailId     []string `yaml:"mailid"`
		SmsId      []string `yaml:"smsid"`
		SrcEmail   string   `yaml:"srcemail"`
		TmplateDir string   `yaml:"tmplatedir"`
		Enable     bool     `yaml:"enable"`
		Addr       string   `yaml:"addr"`
	}
	IpFind struct {
		Auth string `yaml:"auth"`
	}
	OperateLog struct {
		RemainDays int64 `yaml:"remain_days"`
		CleanIntvl int   `yaml:"clean_interval"`
	}
)

func Analyze(configPath string) *Config {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		l4g.Crash(err, "analyze config file errors")
		fmt.Printf("analyze config file errors %v\n", err)
		time.Sleep(time.Second * 1)
	}

	var config *Config

	if err = yaml.Unmarshal(content, &config); err != nil {
		l4g.Crash(err, "analyze yaml config unmarshal errors")
		fmt.Printf("analyze yaml config unmarshal errors %v \n", err)
		time.Sleep(time.Second * 1)
	}

	return config
}
