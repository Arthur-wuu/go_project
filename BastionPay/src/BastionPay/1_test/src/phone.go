
package main

import (

	"io"
	"os"
	"strconv"

	//"fmt"

	"fmt"
	//"io/ioutil"
)

func main(){
	var start_num int
	var end_num int

	start_num = 250000
	end_num = 999999

	var f *os.File
	var err1 error
	filename := "./test.csv"
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0777) //打开文件
		fmt.Println("文件存在")
		check(err1)
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
		check(err1)
	}
	check(err1)

//尾号三个相同的
	//_, err1 = io.WriteString(f, "\r\n\r\n"+"尾号三个相同 abcddd"+"\r\n")
	for i:=start_num ; i<=end_num ; i++  {

		s1 := i % 10              //6
		s2 := i / 10 % 10         //5
		s3 := i / 100 % 10        //4
		//s4 := i / 1000 % 10     //3
		//s5 := i / 10000 % 10    //2
		//s6 := i / 100000 % 10   //1

		if s1 == s2 && s2 == s3 {

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}

	//250052    250250
	//_, err1 = io.WriteString(f, "\r\n\r\n"+"对称号 abccba"+"\r\n")
	for i:=start_num ; i<=end_num ; i++  {

		s1 := i % 10              //6
		s2 := i / 10 % 10         //5
		s3 := i / 100 % 10        //4
		s4 := i / 1000 % 10     //3
		s5 := i / 10000 % 10    //2
		s6 := i / 100000 % 10   //1

		if s1 == s6 && s2 == s5 && s3 == s4 {

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}

	//250250
	//_, err1 = io.WriteString(f, "\r\n\r\n"+"对称号 abcabc"+"\r\n")
	for i:=start_num ; i<=end_num ; i++  {

		s1 := i % 10              //6
		s2 := i / 10 % 10         //5
		s3 := i / 100 % 10        //4
		s4 := i / 1000 % 10     //3
		s5 := i / 10000 % 10    //2
		s6 := i / 100000 % 10   //1

		if s1 == s4 && s2 == s5 && s3 == s6 {

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}

	//303030
	//_, err1 = io.WriteString(f, "\r\n\r\n"+"对称号 ababab"+"\r\n")
	for i:=start_num ; i<=end_num ; i++  {

		s1 := i % 10              //6
		s2 := i / 10 % 10         //5
		s3 := i / 100 % 10        //4
		s4 := i / 1000 % 10     //3
		s5 := i / 10000 % 10    //2
		s6 := i / 100000 % 10   //1

		if s2 == s4 && s1 == s5 && s4 == s6 && s3 == s5 {

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}

	//223344  221100
	//_, err1 = io.WriteString(f, "\r\n\r\n"+"对顺序 对倒序"+"\r\n")
	for i:=start_num ; i<=end_num ; i++  {

		s1 := i % 10              //6
		s2 := i / 10 % 10         //5
		s3 := i / 100 % 10        //4
		s4 := i / 1000 % 10     //3
		s5 := i / 10000 % 10    //2
		s6 := i / 100000 % 10   //1

		if s5+1 == s4 && s3 +1 == s2 && s1 == s2 && s4 == s3 && s6 == s5{

		stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}

		if s5-1 == s4 && s3 -1 == s2 && s1 == s2 && s4 == s3 && s6 == s5{

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}

		if s5-1 == s3 && s3-1 == s1  && s4 == s2 && s6 == s4{

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}

		if s5+1 == s3 && s3+1 == s1  && s4 == s2 && s6 == s4{

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}

	//_, err1 = io.WriteString(f, "\r\n\r\n"+"4顺序 4倒序"+"\r\n")
	for i:=start_num ; i<=end_num ; i++  {
		t:= IsSuccessive(i,4)
		if t == true {
			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}

	//_, err1 = io.WriteString(f, "\r\n\r\n"+"4相同"+"\r\n")
	for i:=start_num ; i<=end_num ; i++  {
		t:= IsAlike(i,4)
		if t == true {
			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}


	for i:=start_num ; i<=end_num ; i++  {

		s1 := i % 10              //6
		s2 := i / 10 % 10         //5

		if s1 == s2 && s2 == 0 {

			stringI :=strconv.Itoa(i)
			_, err1 = io.WriteString(f, "\r\n"+stringI) //写入文件(字符串)
			check(err1)
		}
	}




}


func check(e error) {
	if e != nil {
		panic(e)
	}
}






func IsSuccessive(n, lens int) bool {
//统计正顺次数 12345
z := 0
//统计反顺次数  654321
f := 0
//判断3个数字是否是顺子，只需要判断2次
lens = lens - 1
for {
// 个位数
g := n % 10
n = n / 10
// 十位数
s := n % 10

if s-g == 1 {
f = f + 1
} else {
f = 0
}

if g-s == 1 {
z = z + 1
} else {
z = 0
}

if f == lens || z == lens {
return true
}

if n == 0 {
return false
}
}
}
//判断是否满足重复次数的号码
//IsAlike 是否相同数字
func IsAlike(n, lens int) bool {
c := 0
lens = lens - 1
var g, s int
for {
g = n % 10
n = n / 10
s = n % 10

if s == g {
c = c + 1
} else {
c = 0
}

if c == lens {
return true
}

if n == 0 {
return false
}
}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}