package main

import (
	"fmt"


	//	"os"
)



// 环形链表问题 约瑟夫问题  小孩子循环出圈问题

type Boy struct {
	Id int
	Next * Boy
}

//构建环形链表
func AddCircle (num int) *Boy{
	first := &Boy{
		Id:1,
	}

	//参照Boy
	currencyBoy := &Boy{}

	for i:= 1; i <= num ; i++  {
		boy := &Boy{
			Id:i,
		}

		if i == 1 {
			//只有一个
			first = boy
			currencyBoy = boy
			boy.Next = first
		}else {
			currencyBoy.Next = boy
			currencyBoy = boy
			currencyBoy.Next = first
		}
	}
	return first
}

func PrintLian(first *Boy) {
	curr := first
	for{
		fmt.Println(curr.Id,": -- > ")
		if curr.Next == first {
			break
		}
		curr = curr.Next
	}
}

// 从哪里开始数，数几下
func Play (first *Boy, startNo int, countNo int) {
	//首先，设置一个参照的小孩，first，和 tail 如影相随
	tail := first
	for{
		if tail.Next == first{
			break
		}
		tail = tail.Next
	}

	//从第几个开始数
	for i:= 0; i < startNo - 1; i++  {
		first = first.Next
		tail = tail .Next
	}
	
	//数几下
	for{
		for i:= 0; i < countNo - 1; i++  {
			first = first.Next
			tail = tail .Next
		}
		//删除first小孩哦
		fmt.Println("删除的小孩id", first.Id)
		first = first.Next
		tail.Next = first

		if tail == first {
			break
		}
	}
}

func main() {
	f := AddCircle(5)
	PrintLian(f)

	Play(f,2,3)
}


































/**x
golang实现单链表
单链表是一种链式存取的数据结构，用一组地址任意的存储单元存放线性表中的数据元素。
链表中的数据以节点来表示，每个节点的构成：元素(数据元素的映像)+ 指针(指向后继元素存储位置)。
元素是存储数据的存储单元，指针是连接每个节点的地址数据。
节点结构：
|--data--|--next--|
data域存放节点值的数据域，next域存放节点的直接后继地址(位置)的指针域(链域)
头指针head和终端节点
单链表中每个几点的存储地址是存放在其前趋节点next域中。而开始节点无前趋，所以应该设置头指针head指向开始节点。链表由头指针唯一确定，单链表可以用头指针的名字来命名。
终端节点无后继，所以终端节点的指针域为空，即NULL。
 */
//
//type Node struct {
//	// 值
//	Data interface{}
//	// 后继节点指针
//	Next *Node
//}
//
//// 链表是否为空
//func IsEmpty(node *Node) bool {
//	return node == nil
//}
//
//// 是否是最后一个节点
//func IsLast(node *Node) bool  {
//	return node.Next == nil
//}
//
//// 查找指定节点的前一个节点
//func FindPrevious(data interface{}, node *Node) *Node  {
//	tmp := node
//	for tmp.Next != nil && tmp.Next.Data != data {
//		tmp = tmp.Next
//	}
//	return tmp
//}
//// 查找某个值在哪个节点
//func Find(data interface{}, node *Node) *Node  {
//	tmp := node
//	for tmp.Data != data {
//		tmp = tmp.Next
//	}
//	return tmp
//}
//
//// 插入节点:在指定节点插入节点
///**
//position:指定的节点位置
// */
//func Insert(data interface{}, position *Node)  {
//	// 新建一个节点
//	tmpCell := new(Node)
//	if tmpCell == nil {
//		fmt.Println("err: out of space")
//	}
//	// 给新建的节点的值域赋值为传入的data值
//	tmpCell.Data = data
//	// 新建的节点的next指针指向的是指定节点position的next
//	tmpCell.Next = position.Next
//	// 指定节点position的后继节点变成了新建的节点
//	position.Next = tmpCell
//}
//
//// use myself method
//func Insert2(node , newNode *Node)  {
//	temp := node
//	for{
//		if temp.Next == nil {
//			break
//		}
//		temp = temp.Next
//	}
//	node.Next = temp.Next
//
//	temp.Next = node
//}
//
//// 删除节点
//func Delete(data interface{}, node *Node)  {
//	preview := FindPrevious(data, node)
//	tmp := Find(data, node)
//	preview.Next = tmp.Next
//	tmp = nil
//}
//
//// 删除链表
//func DeleteList(node **Node)  {
//	p := node
//	for *p != nil {
//		tmp := *p
//		*p = nil
//		*p = tmp.Next
//	}
//}
//
////打印列表
//func PrintList(list *Node) {
//	p := list
//	for p != nil {
//		time.Sleep(33333)
//		fmt.Printf("%d------%v------%p ", p.Data, p, p.Next)
//		time.Sleep(33333)
//		p = p.Next
//		fmt.Println()
//	}
//	fmt.Println()
//}
//
//func main() {
//	node4 := &Node{
//		Data: 4,
//		Next: nil,
//	}
//
//	node3 := &Node{
//		Data: 6,
//		Next: node4,
//	}
//
//	node2 := &Node{
//		Data: 6,
//		Next: node3,
//	}
//
//	node1 := &Node{
//		Data: 2,
//		Next: node2,
//	}
//
//	headNode := &Node{
//		Data: 1,
//		Next: node1,
//	}
//	//list := headNode
//	//Insert(2, headNode)
//	//PrintList(list)
//	////fmt.Println(IsEmpty(list))
//	////fmt.Println(IsLast(headNode))
//	////p := Find(0, list)
//	//Insert(2, p)
//	//PrintList(list)
//	//Delete(1, list)
//	//PrintList(list)
//	//DeleteList(&list)
//	//PrintList(list)
//    node  :=Del(6,headNode)
//	PrintList(node)
//
//
//
//
//
//}
//
//func FindPrivious(data  interface{} , node *Node) *Node{
//	tmp := node
//	//找前一个 节点，输入的date 前面的一个节点
//	for  {
//		if tmp.Data != data && tmp.Next.Data == data{
//			fmt.Println("q",tmp)
//			return tmp
//		}
//		tmp = tmp.Next
//	}
//
//	return  tmp
//}
//
//func Finddata(data interface{}, node *Node) *Node{
//	tmp := node
//	for {
//		if tmp.Next.Data == data{
//			fmt.Println("h",tmp.Next.Next)
//			return tmp.Next.Next
//		}
//		tmp = tmp.Next
//	}
//	return  tmp.Next.Next
//}
//
//
////1,2,3,6,4,6,5
//
//func Del(data  interface{}, node *Node ) *Node{
//	tmp := node
//
//
//	 for{
//	 	//if tmp.Data == data {
//	 	//	//第一个就是 删除
//		//
//		//}
//	 	//fmt.Println("qqq",tmp , data  )
//	 	if tmp != nil  &&  tmp.Next.Data == data   {
//	 		for tmp.Next.Next.Data != data{
//					tmp = tmp.Next
//			}
//	 		if tmp.Next == nil {
//	 			return  tmp
//			}
//		}
//
//		 if tmp.Next != nil {
//			 tmp = tmp.Next
//		 }
//		 if tmp.Next == nil {
//			 break
//		 }
//
//
//	 }
//
//
//
//
//return node
//
//
//}
//
//








