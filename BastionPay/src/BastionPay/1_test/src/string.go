package main

import (
	"fmt"
	//"strconv"
	//"strings"
	"time"
)

// Trim 将删除 s 首尾连续的包含在 cutset 中的字符
//func Trim(s string, cutset string) string
//func main() {
//	s := " Hello 世界 ! "
//	ts := strings.Trim(s, "lo!")
//
//	fmt.Printf("%q\n", ts) // "世界"
//}

//
//func main() {
//	ss := []string{"Monday", "Tuesday", "Wednesday"}
//	s := strings.Join(ss, "|")
//	fmt.Println(s)
//}



func main() {
	//x := 20
	//y := 12
	//fmt.Println(x<<1)
	//fmt.Println(y>>1)


//from := "ab,cd,cd,fv,gn"
	//
	//from = strings.TrimRight(from, ",")
	//
	//fmt.Println(
	//	from)
	//fromArr := strings.Split(from, ",")
	//fmt.Print(fromArr)
	//
	//
	//limit := "1"
	////limit, err := strconv.Atoi(count)
	////if err != nil {
	////    ZapLog().Error("count param err")
	//
	////req.Topic = sub+".market."+obj+".kline|{\\\"period\\\":\\\""+period+"\\\"}"
	//
	//obj:="btc_btc"
	//period:="1m"
	//
	//Topic := "market." + obj + ".kline|{\"period\":\"" + period + "\",\"limit\":" + limit + "}"
	//fmt.Println("req.Topic:" + Topic)
	//float64 := 12.34
	//st := strconv.FormatFloat(float64, 'f', -1, 64)
	//
	//str1 := "cdsvds(xxcdscdsxx)"
	//index1 := strings.Index(str1,"(")
	//
	//index2 := strings.Index(str1,")")
	//fmt.Println("req.index1:" , index1)
	//fmt.Println("req.index2:" , index2)
	//
	//str := str1[index1+1 : index2]
	//fmt.Println("str:" + str)
	//fmt.Println("st:" + st)

	fmt.Println("day:" ,GenDay(1541120390))

	//fmt.Println("d" ,min(6,8))
}

func min(a,b int64) int64 {
	if a >= b {
		return b
	}
	return a
}

func  GenDay(t int64) int64{
	t1:=time.Unix(t, 0).Year()        //年
	t2:=time.Unix(t, 0).Month()       //月
	t3:=time.Unix(t, 0).Day()         //日
	lc,_ :=time.LoadLocation("UTC")
	currentTimeData:=time.Date(t1,t2,t3,0,0,0,0, lc)
	return currentTimeData.Unix()


	go func(id int) {

		for  {
			ls :=1
			fmt.Println(ls)
		}
	}(3)

	return 2
}
//output:4 2
