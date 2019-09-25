package main

import (
	"go.uber.org/zap"
)

func main() {
	//	file := "zap.conf2"
	ll, _ := zap.NewDevelopment()

	ll.With(zap.String("a1", "a2")).Error("111")

	//log.LoadZapConfig(file)
	//{
	//	for i := 0; true; i++ {
	//		log.GobalZapLog.With(zap.String("a1", "a2")).Info("12345678", zap.String("url", "hahah"))
	//		log.GobalZapLog.Error("err")
	//
	//	}
	//	//	time.Sleep(time.Second * 1)
	//}
}
