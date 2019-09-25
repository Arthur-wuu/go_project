package main

import "fmt"


func main() {
	//firstUniqChar("sssad")
	s := FindUniqe([]int{1,1,2,2,3,3,4,4,5})


	a := []int{1,2,3}
	b := []int{4,5,6}

	c := appends(a, b)



	fmt.Println(s)
}

//leetcode387 : 找第一个不重复的字符的下标

func firstUniqChar(s string) int {
	stmp:=[]byte(s)//字符串转字节数组
	fmt.Println(stmp)
	sMap:=make(map[byte]int)
	for _,v:=range stmp{//利用map统计字频
		sMap[v]++
	}
	for i,v:=range stmp{//找出第一个字频为1的元素
		if sMap[v]==1{
			fmt.Println(i)
			return i
		}
	}
	return -1
}


//leetcode14 :寻找字符串数组的公共前缀

func FindCommonPrev(arr []string) string {

	//输入 ["flads","sdsadsa","dawdf"]

	if len(arr) == 0 {
		return ""
	}

	// 找长度最小的一个单词
	short := arr[0]
	for _, v := range arr {
		if len(v) <= len(short) {
			short = v
		}
	}

	//2。遍历这个单词的每一个 和数组里的比较

	result := ""

	for i:=0; i < len(short); i++ {
		for j:=0; j < len(arr); j++ {
			if short[i] != arr[j][i]{
				return result
			}
		}

		result = result + string(short[i])

		if len(short) == len(result) {
			break
		}


	}

	return result

}


//leetcode 136: 找出数组里的 一个唯一的不成对的数字

func FindUniqe(arr []int) int {

	// 变成一个个的map， 计数
	s := make(map[int]int)  // 数字， 计数

	for _, v := range arr {
		s[v] ++
	}

	for k, v := range s {
		if v ==1 {
			return k
		}

	}

	return -1

}





















