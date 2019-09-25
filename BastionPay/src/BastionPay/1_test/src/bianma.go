package main

import (
	//"fmt"
	"fmt"
	"strings"
	"sync"
)

var oSingle sync.Once

func TransString (s string ) string{
	a := "我好爱你啊[开心]粉底dewd刷[开心]ddd[对你撒娇]的积分dsqdwed无法测物权法[的玩家]岑哈["

	if !strings.Contains(a,"[") || !strings.Contains(a,"]") {  //不包含【 或者不包含 】 不处理
		return a
	}
	if !strings.Contains(a,"[") && !strings.Contains(a,"]") {  //不包含【 并且 不包含 】 不处理
		return a
	}
	if strings.Contains(a,"[") && !strings.Contains(a,"]") {  //包含【  并且不包含 】 不处理
		return a
	}
	if !strings.Contains(a,"[") &&  strings.Contains(a,"]") {  //不包含【  并且 包含 】 不处理
		return a
	}
		for i:= 0; i < 1;  {

			if strings.Contains(a, "[") && strings.Contains(a, "]"){
				index1 := strings.Index(a, "[")
				fmt.Println("index1",index1)
				index2 := strings.Index(a, "]")
				fmt.Println("index2",index2)
				a = a[:index1] +" "+ a[index2+1:]
				fmt.Println("a:",a)
			}else {
				break
			}
		}
		return a
	}



func main() {
	//s := "尹正杰到此一游"
	//	//
	//	//
	//	//gbk, err := Utf8ToGbk([]byte(s))
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//} else {
	//	//	fmt.Println("以GBK的编码方式打开:",string(gbk))
	//	//}
	//	//
	//	//utf8, err := GbkToUtf8(gbk)
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//} else {
	//	//	fmt.Println("以UTF8的编码方式打开:",string(utf8))
	//	//}


s :=TransString("dswafnkibaff[fera]dd")

	fmt.Println(s)
}