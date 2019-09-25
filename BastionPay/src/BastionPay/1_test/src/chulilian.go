package main


import (
	"fmt"
	//"fmt"
	"strings"
	"time"

	//"time"
)

// 字符串处理函数，传入字符串切片和处理链
func StringProccess(list []string, chain []func(string) string) {

	// 遍历每一个字符串
	for index, str := range list {

		// 第一个需要处理的字符串
		result := str

		// 遍历每一个处理链
		for _, proc := range chain {

			// 输入一个字符串进行处理，返回数据作为下一个处理链的输入。
			result = proc(result)
		}

		// 将结果放回切片
		list[index] = result
	}
}

// 自定义的移除前缀的处理函数
func removePrefix(str string) string {

	return strings.TrimPrefix(str, "go")
}

func main() {

//	// 待处理的字符串列表
//	list := []string{
//		"go scanner",
//		"go parser",
//		"go compiler",
//		"go printer",
//		"go formater",
//	}
//
//	// 处理函数链
//	chain := []func(string) string{
//		removePrefix,
//		strings.TrimSpace,
//		strings.ToUpper,
//	}
//
//	// 处理字符串
//	StringProccess(list, chain)
//
//	// 输出处理好的字符串
//	for _, str := range list {
//		fmt.Println(str)
//	}
//s := time.Now().Local().Format("2006-01-02 15:04:05")
//fmt.Println("s", s)
//	 var UserId  interface{}
//	 UserId = "dw"
//	userId := fmt.Sprintf("%v",UserId)
//
//	fmt.Println("dd", userId)
	//var i interface{} = "TT"
	//var i interface{} = "77"
	//value, ok := i.(int)
	//if ok {
	//	fmt.Printf("类型匹配int:%d\n", value)
	//} else {
	//	fmt.Println("类型不匹配int\n")
	//}
	//if value, ok := i.(int); ok {
	//	fmt.Println("类型匹配整型：%d\n", value)
	//} else if value, ok := i.(string); ok {
	//	fmt.Printf("类型匹配字符串:%s\n", value)
	//}

	const(
		CONST_MIN_RED_DIV_BASE = 10
	)

	a, b := GetTodayUnix()
	fmt.Println(a,b)

}


func GetTodayUnix() (int64, int64) {
	ll := time.FixedZone("UTC", 8*3600)
	temp := time.Now()
	tt := temp.In(ll)
	return time.Date(tt.Year(),   tt.Month(), tt.Day(), 0, 0, 0, 0, ll).Unix(),   temp.Unix()
}