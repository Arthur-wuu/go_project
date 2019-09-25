package main

import (
	sdk "BastionPay/bas-tools/sdk.aws.ses"
	"fmt"
)

func main() {
	fmt.Println("start")
	ss := sdk.NewSesSdk("us-east-1", "AKIAJEJUHBCEJGTFGNHA", "oflCXIs+8jsbBYffy8lPedKok90NkSSabv8SVY66", "")

	err := ss.SendTemplate("531169454@qq.com", "", "pingzilao@qq.com", "t1", "{ \"name\":\"Alejandro\"}", 60)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("ok")
}
