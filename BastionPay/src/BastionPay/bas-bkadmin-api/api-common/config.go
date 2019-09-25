package common

import (
	l4g "github.com/alecthomas/log4go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Data interface{}
}

func NewConfig(data interface{}) *Config {
	return &Config{Data: data}
}

func (c *Config) Read(file string) {
	var err error

	if err = c.ReadYaml(file); err != nil {
		l4g.Error("Read yml config: %s", err.Error())
	}
	if err = c.ReadEnv(); err != nil {
		l4g.Error("Read env config: %s", err.Error())
	}
}

// 从YAML配置文件读取配置
func (c *Config) ReadYaml(file string) error {
	var (
		err  error
		data []byte
	)

	data, err = ioutil.ReadFile(file)
	if err != nil {
		l4g.Error("Read yml config: %s", err.Error())
		return err
	}

	err = yaml.Unmarshal([]byte(data), c.Data)
	if err != nil {
		l4g.Error("Read yml config: %s", err.Error())
		return err
	}

	return nil
}

// 从环境变量读取配置
func (c *Config) ReadEnv() error {
	var (
		value reflect.Value
		to    reflect.Type
		count int
	)

	value = reflect.ValueOf(c.Data).Elem()
	to = value.Type()
	count = value.NumField()
	for i := 0; i < count; i++ {

		if value.Field(i).Kind() == reflect.Slice {
			return nil
		}

		vv := value.Field(i)
		nn := to.Field(i).Name
		tt := value.Field(i).Type()
		cc := value.Field(i).NumField()

		for j := 0; j < cc; j++ {
			tag := tt.Field(j).Tag.Get("yaml")
			envStr := ""
			if tag != "" {
				envStr = strings.ToUpper(nn) + "_" + strings.ToUpper(tag)
			} else {
				envStr = strings.ToUpper(nn) + "_" + strings.ToUpper(tt.Field(j).Name)
			}

			env := os.Getenv(envStr)

			if env != "" {
				if vv.Field(j).IsValid() && vv.Field(j).CanSet() {
					switch vv.Field(j).Kind() {
					case reflect.String:
						vv.Field(j).SetString(env)
					case reflect.Int:
						d, err := strconv.ParseInt(env, 10, 64)
						if err != nil {
							l4g.Error("Read env config: %s", err.Error())
							return err
						}
						vv.Field(j).SetInt(d)
					case reflect.Bool:
						if strings.ToLower(env) == "true" {
							vv.Field(j).SetBool(true)
						} else {
							vv.Field(j).SetBool(false)
						}
					}

				}
			}
		}
	}

	return nil
}
