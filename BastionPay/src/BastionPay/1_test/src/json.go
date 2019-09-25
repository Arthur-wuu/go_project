package src

import (
	"encoding/json"
	"fmt"
)
func main() {
	type Person struct {
		Name   string
		Age    int
		Gender bool
	}

	//散集数据转换成 结构体，切片结构体，map，切片

	//unmarshal to struct
	var p Person
	var str = `{"Name":"junbin", "Age":21, "Gender":true}`
	json.Unmarshal([]byte(str), &p)
	//result --> junbin : 21 : true
	fmt.Println(p.Name, ":", p.Age, ":", p.Gender)

	// unmarshal to slice-struct
	var ps []Person
	var aJson = `[{"Name":"junbin", "Age":21, "Gender":true},
				{"Name":"junbin", "Age":21, "Gender":true}]`
	json.Unmarshal([]byte(aJson), &ps)
	//result --> [{junbin 21 true} {junbin 21 true}] len is 2
	fmt.Println(ps, "len is", len(ps))

	// unmarshal to map[string]interface{}

	var obj interface{} // var obj map[string]interface{}
	json.Unmarshal([]byte(str), &obj)
	m2 := obj.(map[string]interface{})
	//result --> junbin : 21 : true
	fmt.Println(m2["Name"], ":", m2["Age"], ":", m2["Gender"])

//c := {1,2,"sa"}	//unmarshal to slice
	var strs string = `["Go", "Java", "C", "Php"]`
	//var strs1 string = ["Go", "Java", "C", "Php"]
	var aStr2 []string
	json.Unmarshal([]byte(strs), &aStr2)
	//result --> [Go Java C Php]  len is 4
	fmt.Println(aStr2, " len is", len(aStr2))





	//散集数据转换成json
	//========================================================================

	//tag中的第一个参数是用来指定别名
	//比如Name 指定别名为 username `json:"username"`
	//如果不想指定别名但是想指定其他参数用逗号来分隔
	//omitempty 指定到一个field时
	//如果在赋值时对该属性赋值 或者 对该属性赋值为 zero value
	//那么将Person序列化成json时会忽略该字段
	//- 指定到一个field时
	//无论有没有值将Person序列化成json时都会忽略该字段
	//string 指定到一个field时
	//比如Person中的Count为int类型 如果没有任何指定在序列化
	//到json之后也是int 比如这个样子 "Count":0
	//但是如果指定了string之后序列化之后也是string类型的
	//那么就是这个样子"Count":"0"
	type Person1 struct {
		Name        string `json:"username"`
		Age         int
		Gender      bool `json:",omitempty"`
		Profile     string
		OmitContent string `json:"-"`
		Count       int    `json:",string"`
	}


		var p1  = &Person1{
			Name:        "brainwu",
			Age:         21,
			Gender:      true,
			Profile:     "I am Wujunbin",
			OmitContent: "OmitConent",
		}
		if bs, err := json.Marshal(&p1); err != nil {
			panic(err)
		} else {
			//result --> {"username":"brainwu","Age":21,"Gender":true,"Profile":"I am Wujunbin","Count":"0"}
			fmt.Println(string(bs))
		}

		var p2 *Person1 = &Person1{
			Name:        "brainwu",
			Age:         21,
			Profile:     "I am Wujunbin",
			OmitContent: "OmitConent",
		}
		if bs, err := json.Marshal(&p2); err != nil {
			panic(err)
		} else {
			//result --> {"username":"brainwu","Age":21,"Profile":"I am Wujunbin","Count":"0"}
			fmt.Println(string(bs))
		}

		// slice 序列化为json
		var aStr []string = []string{"Go", "Java", "Python", "Android"}
		if bs, err := json.Marshal(aStr); err != nil {
			panic(err)
		} else {
			//result --> ["Go","Java","Python","Android"]
			fmt.Println(string(bs))
		}

		//map 序列化为json
		var m map[string]string = make(map[string]string)
		m["Go"] = "No.1"
		m["Java"] = "No.2"
		m["C"] = "No.3"
		if bs, err := json.Marshal(m); err != nil {
			panic(err)
		} else {
			//result --> {"C":"No.3","Go":"No.1","Java":"No.2"}
			fmt.Println(string(bs))
		}




}


