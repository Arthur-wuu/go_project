package main

import "fmt"

func main(){
	//s  := preorderTraversal1(&TreeNode{1,nil,nil})
	//fmt.Println(s)
	s := 4 >> 2
	fmt.Println(s)
}
type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

func preorderTraversal1(root *TreeNode) []int {
	 sInt := make( []int,0 )

	digui1(root, &sInt)

	return sInt
}

func digui1 (root *TreeNode, sInt *[]int) {
	if root != nil {
		*sInt = append(*sInt, root.Val)
		digui1(root.Left, sInt)
		digui1(root.Right, sInt)
	}
}


