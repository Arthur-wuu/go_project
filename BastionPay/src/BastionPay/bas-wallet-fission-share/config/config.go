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
	if GConfig.Cache.UserLevelMaxKey < 100 {
		GConfig.Cache.UserLevelMaxKey = 100
	}
	if GConfig.Cache.UserLevelTimeout < 5 {
		GConfig.Cache.UserLevelTimeout = 60
	}
	if GConfig.Cache.LevelAuthMaxKey < 10 {
		GConfig.Cache.LevelAuthMaxKey = 10
	}
	if GConfig.Cache.LevelAuthTimeout < 1800 {
		GConfig.Cache.LevelAuthTimeout = 1800
	}
	if GConfig.Cache.LevelRuleMaxKey < 10 {
		GConfig.Cache.LevelRuleMaxKey = 10
	}
	if GConfig.Cache.LevelRuleTimeout < 1800 {
		GConfig.Cache.LevelRuleTimeout = 1800
	}
}

type Config struct {
	Server System `yaml:"system"`
	Redis  Redis  `yaml:"redis"`
	Db     Mysql  `yaml:"mysql"`
	Cache  Cache  `yaml:"cache"`
}

type System struct {
	Port    string `yaml:"port"`
	Debug   bool   `yaml:"debug"`
	LogPath string `yaml:"log_path"`
	Monitor string `yaml:"monitor"`
}

type Cache struct {
	UserLevelMaxKey  int `yaml:"user_level_max_key"`
	UserLevelTimeout int `yaml:"user_level_timeout"`
	LevelAuthMaxKey  int `yaml:"level_max_key"`
	LevelAuthTimeout int `yaml:"level_timeout"`
	LevelRuleMaxKey  int `yaml:"level_rule_max_key"`
	LevelRuleTimeout int `yaml:"level_rule_timeout"`
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
