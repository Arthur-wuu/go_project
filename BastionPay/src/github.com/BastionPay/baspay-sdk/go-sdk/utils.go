package go_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	. "BastionPay/bas-base/log/zap"
	"sort"
)


func  GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	return fmt.Sprintf("%s", ud)
}

func RequestBodyToSignStr (body []byte) (string){
	requestParams := make(map[string]string,0)

	err := json.Unmarshal(body, &requestParams)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("requestbody to requestParams err")
		return ""
	}
	//将param的key排序，
	keysSort := make([]string, 0)
	for k, _ := range requestParams{
		keysSort = append(keysSort, k)
	}
	sort.Strings(keysSort)
	//拼接签名字符串
	signH5Str := ""
	for i:=0; i<len(keysSort); i++ {
		signH5Str += keysSort[i]+"="+requestParams[keysSort[i]]+"&"
	}
	signH5Str = signH5Str[0:len(signH5Str)-1]
	return  signH5Str
}