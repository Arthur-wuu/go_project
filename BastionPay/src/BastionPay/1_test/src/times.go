package src

import(
	"fmt"
	"strconv"
	//"time"
)

func main() {
	//datetime := "2018-12-14T10:08:04.689Z"  //待转化为时间戳的字符串
	//
	////日期转化为时间戳
	//timeLayout := "2006-01-02T15:04:05.000Z"
	//loc, _ := time.LoadLocation("GMT")    //获取时区
	//tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	//timestamp := tmp.Unix()    //转化为时间戳 类型是int64
	//fmt.Println(timestamp)
	//
	////时间戳转化为日期
	//datetime = time.Unix(timestamp, 0).Format(timeLayout)
	//fmt.Println(datetime)


	idInt := make([]int,4)
	idInt[0] = 1
	idInt[1] = 2
	idInt[2] = 3
	idInt[3] = 5
	fmt.Println(len(idInt))

	s := ArrToString(idInt)
	fmt.Println(s)
}

func ArrToString(idInt []int)  string {
	param := ""
	for i:=0 ; i<len(idInt) ;i++  {
		if i == len(idInt) -1 {
			param = param+strconv.Itoa(idInt[i])
			break
		}
		param = param+strconv.Itoa(idInt[i])+","
	}
	return  param
}