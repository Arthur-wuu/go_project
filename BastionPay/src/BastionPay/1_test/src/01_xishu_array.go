package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	//"strings"
)

/*
/	稀疏数组的实现

	0 0 0 0
	0 1 0 0
	0 0 0 1
	0 0 0 0

	记录数据在的： 行，列 值
	保存在文件里
	读取出来恢复

*/

type XiShuArr struct {
	row int
	col int
	val int
}

func main() {

	//将上面的数据放到 数组里面
	dataArr := [4][4]int{}
	dataArr[1][1] = 1
	dataArr[2][3] = 1

	fmt.Println(dataArr)

	//遍历数组里的数据，知道数据的行，列，放到 结点里面
	_, err := os.Create("./date.date")
	if err != nil {
		fmt.Println("file 1 err ", err)
	}
	f, err := os.OpenFile("./date.date", os.O_RDWR|os.O_APPEND, 0777) //打开文件
	if err != nil {
		fmt.Println("file 2 err ", err)
	}
	xiShuArr := make([]XiShuArr, 0)
	for i, rows := range dataArr {
		for j, val := range rows {
			if val != 0 {
				arr := XiShuArr{
					row: i,
					col: j,
					val: val,
				}
				xiShuArr = append(xiShuArr, arr)
				s := strconv.Itoa(i) + " " + strconv.Itoa(j) + " " + strconv.Itoa(val) + " " + "\n"
				io.WriteString(f, s)
			}
		}
	}
	fmt.Println(xiShuArr)

	//存文件

	ReturnArr := [4][4]int{}

	//读取文件
	f, err = os.Open("./date.date")
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		//fmt.Println(line)
		s := line
		arr := strings.Split(s, " ")
		if len(arr) < 3 {
			fmt.Println("duqushuju err")
		}
		v,   _ := strconv.Atoi(arr[2])
		row, _ := strconv.Atoi(arr[0])
		col, _ := strconv.Atoi(arr[1])
		ReturnArr[row][col] = v
	}
	fmt.Println(ReturnArr)
	defer f.Close()

}
