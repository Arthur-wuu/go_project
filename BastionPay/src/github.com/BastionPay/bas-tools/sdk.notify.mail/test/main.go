package main

import(
	"BastionPay/bas-tools/sdk.notify.mail"
	"encoding/json"
	"fmt"
)

func main(){
	req := new(sdk_notify_mail.ReqNotifyMsg)
	req.SetGroupName("aaa")
	req.SetAppName("bbb")
	req.SetLang("cc")
	req.Recipient = nil
	req.Params = map[string]interface{}{"key1":"aa"}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

