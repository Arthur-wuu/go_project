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
	Server  Server              `yaml:"server"`
	Db      Db                  `yaml:"db"`
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
	Db            string `yaml:"db"`
	Max_idle_conn int    `yaml:"max_idle_conn"`
	Max_open_conn int    `yaml:"max_open_conn"`
}



