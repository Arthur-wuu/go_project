package config

import (
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/common"
)

type PathLimit struct {
	Path   string
	Method string
	Limit  int
	Time   int
}

type Config struct {
	Server struct {
		Port    string
		Debug   bool
		Logpath string
	}
	Db struct {
		Host        string
		Port        string
		User        string
		Password    string
		DbName      string `yaml:"db"`
		MaxIdleConn int    `yaml:"max_idle_conn"`
		MaxOpenConn int    `yaml:"max_open_conn"`
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		Db       string
	}
	Token struct {
		Secret     string
		Expiration string
	}
	Exchange struct {
		Host string
		Port string
	}
	Ses struct {
		Region      string
		AccessKeyId string `yaml:"access_key_id"`
		SecretKey   string `yaml:"secret_key"`
		Sender      string
	}
	Sns struct {
		Region      string
		AccessKeyId string `yaml:"access_key_id"`
		SecretKey   string `yaml:"secret_key"`
	}
	IpFind struct {
		Auth string
	} `yaml:"ip_find"`
	PathLimits      []*common.PathLimit      `yaml:"path_limits"`
	TermBlockLimits []*common.TermBlock      `yaml:"termblock_limits"`
	LevelPathLimits []*common.LevelPathLimit `yaml:"level_path_limits"`

	Wallet      bastionpay.Gateway `yaml:"wallet"`
	WalletPaths []string           `yaml:"wallet_paths"`

	CoinMarket struct {
		Url    string `yaml:"url"`
		IdPath string `yaml:"id_path"`
	}

	Notice struct {
		ClearTimer string `yaml:"clear_timer"`
	}

	Cache struct {
		AuditeTimeout   int `yaml:"audite_timeout"`
		AuditeMaxKeyNum int `yaml:"audite_max_key"`
	}

	Limits struct {
		IdSms  []int `yaml:"id_sms"`
		IdMail []int `yaml:"id_mail"`
		IpSms  []int `yaml:"ip_sms"`
		IpMail []int `yaml:"ip_mail"`
	}

	Bas_quote struct {
		Addr string `yaml:"addr"`
	}
	Bas_notify struct {
		Addr                   string `yaml:"addr"`
		VerifyCodeSmsTmp       string `yaml:"verifycode_sms_tmp"`
		VerifyCodeMailTmp      string `yaml:"verifycode_mail_tmp"`
		RegisterSuccessSmsTmp  string `yaml:"register_ok_sms_tmp"`
		RegisterSuccessMailTmp string `yaml:"register_ok_mail_tmp"`
	}
}
