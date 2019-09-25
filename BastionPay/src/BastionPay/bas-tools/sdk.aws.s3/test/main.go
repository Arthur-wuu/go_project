package main

import (
	sdks3 "BastionPay/bas-tools/sdk.aws.s3"
	"bytes"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Printf("start")
	bb := []byte("hello,aws!")
	data := bytes.NewReader(bb)
	cc := sdks3.NewS3Sdk("us-east-2", "AKIAJEJUHBCEJGTFGNHA", "oflCXIs+8jsbBYffy8lPedKok90NkSSabv8SVY66", "")

	//// upload file
	addr, err := cc.UpLoad("blockshine-bastionpay-logo", "test25.txt", data, 0)
	if err != nil {
		fmt.Printf("\nerr = %s\n", err.Error())
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Printf("%s\n", addr)

	time.Sleep(3 * time.Second)
	//addr, err = cc.UpLoad("blockshine-bastionpay-logo","test26.txt", data,0)
	//if err != nil {
	//	fmt.Printf("\nerr = %s\n", err.Error())
	//	return
	//}
	//fmt.Print(addr)
	//time.Sleep(time.Second*3)

	// download
	barFile, err := os.Create("./bar.file")
	if err != nil {
		fmt.Printf("\nerr = %s\n", err.Error())
		return
	}
	n, err := cc.Download("blockshine-bastionpay-logo", "test25.txt", barFile, 0)
	if err != nil {
		fmt.Printf("\nerr = %s\n", err.Error())
		return
	}
	fmt.Printf("file size is %d\n", n)
	return
}
