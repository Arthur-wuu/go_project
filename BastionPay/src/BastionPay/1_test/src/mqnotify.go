package main

import (
	"encoding/json"
	"fmt"
)

type (
	Content struct {
		UId           *uint64        `valid:"optional" json:"uid"`
		RegistTime    *int64         `valid:"optional" json:"regist_time"`
		Country       *string        `valid:"optional" json:"country"`
		Phone         *string        `valid:"optional" json:"phone"`
		Channel       *int           `valid:"optional" json:"channel"`
	}
)



func main(){
	reqBody, _ := json.Marshal(map[string]interface{}{
		"uid": 57,
		"regist_time": 1552995187,
		"country": "zh-Cn",
		"phone": "17317351119",
		"channel": 1,
	})
	fmt.Println("body",string(reqBody))

}