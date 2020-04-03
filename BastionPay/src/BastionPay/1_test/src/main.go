//package src
//
//import (
//	"fmt"
//	"strings"
//)
//
//type Role struct {
//	Id     int64
//	Name   string
//	Status string
//}
//
//func main()  {
//		role1:=Role{6,"sa1","1"}
//		role2:=Role{7,"sa3","2"}
//		role3:=Role{8,"sa4","3"}
//
//	role:=  [3]Role{role1,role2,role3}
//
//	var recipientList []string
//	//recipientList1 := []string{}
//	//recipientList2 := make([]string ,3)
//
//	var index []int
//
//	for k,v:=range role{
//		index = append(index,k)
//		recipientList = append(recipientList, v.Status+v.Name)
//	}
//
//
//
//	fmt.Println(recipientList)
//	fmt.Println(index)
//
//	srcIp := "http://123"
//	srcIp = strings.TrimLeft(srcIp, "http://")
//	fmt.Println(srcIp)
//
//	//fileSuffix := fmt.Sprintf("updatefrozen_%d", 1)
//	//fmt.Println(fileSuffix)
//}

//c := []string{}
//counts := make(map[string]int)
//if err != nil {
//l4g.Error("mailNotify sdk init err[%s]", err.Error())
//}


package main

import (
	"encoding/binary"
	aa "net/url"
	//"github.com/kataras/iris/core/netutil"
	"strconv"

	"bytes"
	//"github.com/bluele/gcache"
	"fmt"
	"math/rand"
	"time"
)


//生成随机字符串
func GenRandomString(l int) int32{
	str := "789"
	bt := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bt[r.Intn(len(bt))])
	}
	bytesBuffer := bytes.NewBuffer(result)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int32(x)
}


func GetRandomString2(l int) int{
	str := "123456789"
	bt := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bt[r.Intn(len(bt))])
	}
fmt.Println(string(result))
	intNum, err := strconv.Atoi(string(result))
	if err != nil {
		//ZapLog().Error( "string to int err", zap.Error(err))
	}
	return  intNum
}

func main() {
	callbackUrl, err := aa.QueryUnescape("__ CALLBACK__%26clktime%3d%7bacttime%7d%26imei%3d%7bIMEI%7d%26idfa%3d% 7bIDFA%7d%26event%3d%7bEVENT%7d")

fmt.Println(callbackUrl)
	fmt.Println(err)


}