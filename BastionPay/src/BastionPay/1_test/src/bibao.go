package src

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main(){
	var j int = 5
	a:=func() (func()){
		var i int = 10
		return func(){
			fmt.Printf("i,j:%d,%d\n",i,j)
		}
	}()//将一个无需参数，返回值为匿名函数的函数赋值给a()

	sub := "sub"
	obj := "obj"
	period := "1m"
	Topic1 := sub+".market."+obj+".kline|{\\\"period\\\":\\\""+period+"\\\"}"
	Topic2 := sub+".market."+obj+".kline|{\\\"period\\\":\\\""+period+"\\\"}"
	//变成json的字节数组，拼接成json

	json1 := "{\"topic\":\""+Topic1+"\"}"

	fmt.Println(json1)
	fmt.Println(Topic1)
	fmt.Println(Topic2)



	//var ar = []string{"s", "d","x"}
	//s:= ar[1]
	//str := "[\"u\", \"OK\"]"
	//ss := strings.Split(str,",")

	//a()
	//j*=2
	// i*=2这样是错的
	a()


	str := "[\"1539043200000\", \"0.034473\", [\"0.034614\", \"0.03433\"]]"
	byte := []byte(str)
	fmt.Println(str)

	slic := make([]interface{},0)
	json.Unmarshal(byte, &slic)

		s0 := slic[0].(string)
		s1 := slic[1]
		s2 := slic[2]

	c := reflect.TypeOf(s2)


	fmt.Println(s0)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println("*****",c)



}
