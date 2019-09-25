package main

import "github.com/BastionPay/bas-base/l4gmgr"
import (
	"fmt"
	l4g "github.com/alecthomas/log4go"
)

func main() {
	fmt.Println("start1")

	workPath := "/Users/yy/Desktop/all/BastionPay/src/github.com/BastionPay/bas-base/l4gmgr/test"
	l4gmgr.LoadConfiguration2(workPath+"/a/log.xml", workPath+"/a/b/c/app.log")
	defer l4gmgr.Close()

	l4g.Info("nihao")
	l4g.Info("hahahhahahahah")
	fmt.Println("end")
}
