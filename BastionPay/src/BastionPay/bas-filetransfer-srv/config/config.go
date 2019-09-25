package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
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
	preProcess()
	return &GConfig
}

func preProcess() {
	if GConfig.Task.MaxOpTime < 5 {
		GConfig.Task.MaxOpTime = 5
	}
	if GConfig.Task.FileKeepTime < 5*60 {
		GConfig.Task.FileKeepTime = 300
	}
	if GConfig.Task.StatusKeepTime < 3*60 {
		GConfig.Task.StatusKeepTime = 300
	}
	if GConfig.Task.MaxRecords < 50 {
		GConfig.Task.MaxRecords = 50
	}
	if GConfig.Task.MaxPage < 50 {
		GConfig.Task.MaxPage = 50
	}
	if GConfig.Task.MaxWaitTime < 20 {
		GConfig.Task.MaxWaitTime = 20
	}
	if GConfig.Task.MaxWaitLen < 5 {
		GConfig.Task.MaxWaitLen = 5
	}
	limitStr := strings.ToUpper(GConfig.Task.FileGenLimitStr)
	limitStr = strings.Replace(limitStr, " ", "", len(limitStr))
	limitArr := strings.Split(limitStr, "-")
	if len(limitArr) > 0 {
		GConfig.Task.FileGenlimitCount, _ = strconv.ParseInt(limitArr[0], 10, 64)
	}
	if len(limitArr) > 1 {
		base := int64(3600)
		switch limitArr[1][len(limitArr[1])-1] {
		case 'H':
			break
		case 'M':
			base = 60
			break
		case 'S':
			base = 1
			break
		case 'D':
			base = 2600 * 24
			break
		default:

		}
		GConfig.Task.FileGenlimitTime, _ = strconv.ParseInt(limitArr[1][:len(limitArr[1])-1], 10, 64)
		GConfig.Task.FileGenlimitTime *= base
	}
	if GConfig.Task.FileGenlimitCount < 1 {
		GConfig.Task.FileGenlimitCount = 1
	}
	if GConfig.Task.FileGenlimitTime < 5 {
		GConfig.Task.FileGenlimitTime = 3600
	}
}

type Config struct {
	Server System   `yaml:"system"`
	Redis  Redis    `yaml:"redis"`
	Aws    Aws      `yaml:"aws"`
	Dbs    []*Mysql `yaml:"mysql"`
	Task   Task     `yaml:"task"`
}

type System struct {
	Port        string `yaml:"port"`
	Debug       bool   `yaml:"debug"`
	LogPath     string `yaml:"log_path"`
	MaxWaitTask int64  `yaml:"max_wait_task"`
	TmpPath     string `yaml:"temp_path"`
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

type Aws struct {
	AccessKeyId string `yaml:"accesskeyid"`
	AccessKey   string `yaml:"accesskey"`
	AccessToken string `yaml:"accesstoken"`
	FileRegion  string `yaml:"fileregion"`
	FileBucket  string `yaml:"filebucket"`
	FileTimeout int    `yaml:"filetimeout"`
}

type Task struct {
	MaxWaitLen        int64  `yaml:"max_wait_len"`
	MaxOpTime         int64  `yaml:"max_op_time"`
	MaxWaitTime       int64  `yaml:"max_wait_time"`
	MaxPage           int    `yaml:"max_page"`
	MaxRecords        uint64 `yaml:"max_records"`
	StatusKeepTime    int64  `yaml:"status_keep_time"`
	FileKeepTime      int64  `yaml:"file_keep_time"`
	FileGenLimitStr   string `yaml:"file_gen_limit"` //代码未实现
	FileGenlimitCount int64  `yaml:"-"`
	FileGenlimitTime  int64  `yaml:"-"`
	FileUseExist      bool   `yaml:"file_use_exist"`
}
