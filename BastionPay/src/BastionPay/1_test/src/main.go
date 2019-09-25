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
	"strconv"

	//"github.com/bluele/gcache"
	"fmt"
	"math/rand"
	"time"
	"bytes"
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
	//gc := gcache.New(20).
	//	LRU().
	//	Build()
	//gc.SetWithExpire("ke2", "ok", time.Second*10)
	//value, err := gc.Get("ke")
	//fmt.Println("Get:", value,err)
	//
	//// Wait for value to expire
	//time.Sleep(time.Second*10)
	//
	////value, err := gc.Get("key")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Get:", value)

	//s :=GenRandomString(32)
	////d := GetRandomString(32)
	////
	//fmt.Println("Get:", s)
	//keysSort := make([]string,0)

	//var signH5Str string
	//
	//signH5Str = "="+"&"
	//
	a := "abcdefghij"
	a = a[0:len(a)-1]
	fmt.Println("a:", a)
}