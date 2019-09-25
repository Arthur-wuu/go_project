package main

import (
	"fmt"
	"log"

	"time"
)


//1，获取当前时间  获取当前时间戳
//2，获取当前时间，进行格式化
//3，休眠


func getNow() {
	// 获取当前时间，返回time.Time对象
	fmt.Println(time.Now())
	// output: 2016-07-27 08:57:46.53277327 +0800 CST
	// 其中CST可视为美国，澳大利亚，古巴或中国的标准时间
	// +0800表示比UTC时间快8个小时

	// 获取当前时间戳
	fmt.Println(time.Now().Unix())
	// 精确到纳秒，通过纳秒就可以计算出毫秒和微妙
	fmt.Println(time.Now().UnixNano())
	// output:
	//    1469581066
	//    1469581438172080471
}

func formatUnixTime() {
	// 获取当前时间，进行格式化
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	// output: 2016-07-27 08:57:46

	// 指定的时间进行格式化
	fmt.Println(time.Unix(1469579899, 0).Format("2006-01-02 15:04:05"))
	// output: 2016-07-27 08:38:19
}



func getYear() {
	// 获取指定时间戳的年月日，小时分钟秒
	t := time.Unix(1469579899, 0)
	fmt.Printf("%d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	// output: 2016-7-27 8:38:19
}



// 将2016-07-27 08:46:15这样的时间字符串转换时间戳
func strToUnix() {
	// 先用time.Parse对时间字符串进行分析，如果正确会得到一个time.Time对象
	// 后面就可以用time.Time对象的函数Unix进行获取
	t2, err := time.Parse("2006-01-02 15:04:05", "2016-07-27 08:46:15")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(t2)
	fmt.Println(t2.Unix())
	// output:
	//     2016-07-27 08:46:15 +0000 UTC
	//     1469609175
}



// 根据时间戳获取当日开始的时间戳
// 这个在统计功能中会常常用到
// 方法就是通过时间戳取到2016-01-01 00:00:00这样的时间格式
// 然后再转成时间戳就OK了
// 获取月开始时间和年开始时间类似
func getDayStartUnix() {
	t := time.Unix(1469581066, 0).Format("2006-01-02")
	sts, err := time.Parse("2006-01-02", t)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sts.Unix())
	// output: 1469577600
}




// 休眠
func sleep() {
	// 休眠1秒
	// time.Millisecond    表示1毫秒
	// time.Microsecond    表示1微妙
	// time.Nanosecond    表示1纳秒
	time.Sleep(1 * time.Second)
	// 休眠100毫秒
	time.Sleep(100 * time.Millisecond)

}

func main()  {
	//TimeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	//
	//fmt.Print(TimeStamp)
	datetime := "2018-12-20T16:16:08"  //待转化为时间戳的字符串 2018-12-20 16:16:08

	//日期转化为时间戳
	timeLayout := "2006-01-02T15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("PRC")    //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()    //转化为时间戳 类型是int64
	fmt.Println(timestamp)

	//回调函数
	var f func()

	f = getNow

	f()

	//douhaoIndex:= "var hq_str_fx_susdaed=\"14:27:27,3.672"
	//
	// c:= douhaoIndex[23:31]
	//fmt.Println("c",c)

}

func TimeToTimestamp ( t string ) int64 {
	datetime := t  //待转化为时间戳的字符串

	//日期转化为时间戳
	timeLayout := "2006-01-02T15:04:05.000Z"  //转化所需模板
	loc, _ := time.LoadLocation("GMT")    //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()
	return  timestamp
}

