package config

import (
	"BastionPay/bas-api/utils"
	"encoding/json"
	l4g "github.com/alecthomas/log4go"
	"io/ioutil"
)

func GetBastionPayConfigDir() string {
	appDir, err := utils.GetAppDir()
	if err != nil {
		l4g.Crashf("Get App dir crash %s", err.Error())
	}
	return appDir + "/" + BastionPayConfigDirName
}

func LoadJsonNode(absPath string, name string, value interface{}) error {
	data, err := ioutil.ReadFile(absPath)
	if err == nil {
		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		if err == nil {
			if v, ok := jsonMap[name]; ok {
				data, err = json.Marshal(v)
				if err == nil {
					return json.Unmarshal(data, value)
				}
			}
		}
	}
	return err
}
