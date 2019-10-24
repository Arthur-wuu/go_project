package main

import (
	"container/list"
	"fmt"
	"strings"
	"time"
	//"sort"
)

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	    Val int
	    Next *ListNode
	  //  s  chan *error
}

//zijide
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {

	res := &ListNode{0, nil}
	l := res
	temp := 0

	for{  //每个位数相加，
		if l1.Val + l2.Val + temp > 9 {
			l.Val = l1.Val + l2.Val + temp - 10
			temp = 1
		}else {
			l.Val = l1.Val + l2.Val + temp
			temp = 0
		}

		//没有可以加的了
		if l1.Next == nil && l2.Next == nil && temp == 0 {
			break
		}

		//高位
		if l1.Next != nil {
			l1 = l1.Next
		}else {
			l1 = &ListNode{0,nil}
		}
		if l2.Next != nil {
			l2 = l2.Next
		}else {
			l2 = &ListNode{0,nil}
		}
		l.Next = &ListNode{0, nil}
		l = l.Next

	}


return res

}



func addTwoNumber(l1 *ListNode, l2 *ListNode) *ListNode {
	var res = &ListNode{0, nil}
	l := res   // 头节点
	temp := 0
	for {
		if temp + l1.Val + l2.Val > 9 {
			l.Val = temp + l1.Val + l2.Val - 10
			temp = 1
		} else {
			l.Val = l1.Val + l2.Val + temp
			temp = 0
		}


		if l1.Next == nil && l2.Next == nil && temp == 0 {
			break
		}

		l.Next = &ListNode{0, nil}

		if l1.Next != nil {
			l1 = l1.Next
		} else {
			l1 = &ListNode{0, nil}
		}

		if l2.Next != nil {
			l2 = l2.Next
		} else {
			l2 = &ListNode{0, nil}
		}

		l = l.Next
	}
	return res
}


func main(){
	l1 := ListNode{2,nil}
	l2 := ListNode{3,nil}
	s := addTwoNumbers(&l1, &l2)
	fmt.Println(s)

	s1 := singleNumber([]int{1,1,2})
	fmt.Println(s1)

	s2 := majorityElement([]int{1,1,2})
	fmt.Println(s2)


	s3 := maxProfit([]int{1,1,2,4,5})
	fmt.Println(s3)
	Ts()
}

func lengthOfLongestSubstring(s string) int {
	val := []byte(s)
	kvMap := make([]int, 128)
	lens := len(s)
	var max, num int

	for i, j := 0, 0; i < lens && j < lens; j++ {
		if kvMap[val[j]] > i {
			i = kvMap[val[j]]
		}
		num = j - i + 1
		if num > max {
			max = num
		}
		kvMap[val[j]] = j + 1
	}
	return max
}




func isValids(s string) bool {
	sMap := map[string]string{
		")":"(",
		"]":"[",
		"}":"{",
	}

	stack := make([]string,0)

	// 遇到前括号，入栈，  遇到后括号，匹配，出栈

	for _, v := range s {
		if string(v) == "(" || string(v) == "{" || string(v) == "["  {
			stack = append(stack, string(v))
		}else if  string(v) == ")" || string(v) == "}" || string(v) == "]" {
			if sMap[string(v)] == stack[len(stack) - 1] && len(stack) > 0 {
				stack = stack[:len(stack)-1]
			} else {
				return  false
			}
		}

	}

	if len(stack) > 0 {
		return  false
	} else {
		return  true
	}

}


/**
* Definition for singly-linked list.
* type ListNode struct {
*     Val int
*     Next *ListNode
* }
*/

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	// 新建一个链表
	res := &ListNode{0,nil}
	//

	for{
		if l1.Next != nil && l2.Next != nil {
			if l1.Val > l2.Val {
				res.Val = l2.Val
			}else {
				res.Val = l1.Val
			}
		}
	}



}

















func isValid(s string) bool {
	stack := []string{}
	// 后括号映射表
	frontBracket := map[string]string{ ")":"(", "]":"[", "}":"{" }

	for _, x := range s {
		if x=='(' || x=='[' || x=='{' {     // 遇到前括号，入栈
			stack = append(stack,string(x))
		} else if x==')' || x==']' || x=='}' {    // 遇到后括号，判断
			if len(stack)!=0 && stack[len(stack)-1] == frontBracket[string(x)] { // 栈非空，和栈顶元素匹配，匹配成功，出栈
				stack = stack[0:len(stack)-1]
			} else {    // 栈空或者匹配失败，返回错误
				return false
			}
		}
	}
	if (len(stack)==0) {
		return true
	} else {
		return false
	}
}


func singleNumber(nums []int) int {
	// 数组里只出现一次的数字	 放到map里面去，计数
	smap := make(map[int]bool)

	for _, v := range nums {
		if _, ok :=  smap[v] ;ok {   //ok  // 说明有了
			smap[v] = false
		}else {
			smap[v] = true
		}
	}

	for k, v := range smap {
		if v == true {
			return  k
		}
	}
return  0
}



func majorityElement(nums []int) int {

	smap := make(map[int]int)
	for _, v := range nums {
		smap[v] += 1
	}

	max := 0
	keyInt := 0
	for k, v := range smap {
		if v > max {
			max = v
			keyInt = k
		}
	}
	return  keyInt

}


func maxProfit(prices []int) int {
	shouyi := 0
	max := 0
	for i:= 0; i< len(prices)-1 ; i++  {
		for j:=i+1; j < len(prices) ; j++  {
			shouyi = prices[j] - prices[i]
			if shouyi > max {
				max = shouyi
			}
		}
	}

	return max
}

func Ts() time.Time {
	t := time.Now().Format("2006-01-02 15:04:05")
	time, _ :=time.Parse("2006-01-02 15:04:05",t)
	return time
}

func backspaceCompare(S string, T string) bool {
	//"sa#a"  // "sad#d"
	//循环string，一个个添加，遇到#，减去一个
	stack1 := make([]string,0)
	stack2 := make([]string,0)

	for _, v := range S {
		if string(v) != "#" {
			stack1 = append(stack1, string(v))
		}else if  string(v) == "#" && len(stack1) > 0   {
			stack1 = stack1[:len(stack1)-1]
		}
	}

	for _, v := range T {
		if string(v) != "#" {
			stack2 = append(stack2, string(v))
		}else if  string(v) == "#" && len(stack2) > 0  {
			stack2 = stack2[:len(stack1)-1]
		}
	}

	strings1 := ""
	for _,v := range stack1 {
		strings1 = strings1 + string(v)
	}

	strings2 := ""
	for _,v := range stack2 {
		strings2 = strings2 + string(v)
	}

	if strings1 == strings2 {
		return true
	}
	cha := make(chan int)

	return false

}



func preorderTraversal(root *TreeNode) []int {
	sInt := make([]int,0)

	digui(root, sInt)

	return sInt
}

func digui (root *TreeNode, sInt []int) {
	if root != nil {
		sInt = append(sInt, root.Val)
		digui(root.Left, sInt)
		digui(root.Right, sInt)
	}
}






//选择  冒泡   插入   快速   希尔   归并   计数   基数   堆排序     kmp lru

//倒排索引 哈希表(数据结构)

//




















