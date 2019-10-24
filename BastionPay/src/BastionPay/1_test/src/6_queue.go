package main

import (
	"errors"
	"fmt"
)

/**
队列：
基本操作是入队(Enqueue)，在表的末端插入一个元素
出队(Dequeue)，删除(或返回)在表头的元素
 */

type Item interface {}

// 队列结构体
type Queue struct {
	Items []Item
	//ITems []interface{}
}

//type IQueue interface {
//	New() Queue
//	Enqueue(t Item)
//	Dequeue(t Item)
//	IsEmpty() bool
//	Size() int
//}

// 新建
func (q *Queue)New() *Queue  {
	//q.Items = []Item{}
	q.Items = make([]Item,0)
	return q
}

// 入队
func (q *Queue) Enqueue(data Item)  {
	q.Items = append(q.Items, data)
}

// 出队
func (q *Queue) Dequeue() *Item {
	// 由于是先进先出，0为队首
	item := q.Items[0]
	q.Items = q.Items[1: len(q.Items)]
	return &item
}

// 队列是否为空
func (q *Queue) IsEmpty() bool  {
	return len(q.Items) == 0
}

// 队列长度
func (q *Queue) Size() int  {
	return len(q.Items)
}

var q Queue

func initQueue() *Queue  {
	if q.Items == nil{
		q = Queue{}
		q.New()
	}
	return &q
}

func main() {
	//q := initQueue()
	//fmt.Println("length:", q.Size())
	//q.Enqueue("a")
	//q.Enqueue(1)
	//fmt.Println("length:", q.Size())
	//fmt.Println(q)
	//tmp := q.Dequeue()



}

/*
循环队列实现思路：
1.循环队列须要几个參数来确定
front，tail，length，capacity
front指向队列的第一个元素，tail指向队列最后一个元素的下一个位置
length表示当前队列的长度，capacity标示队列最多容纳的元素
2.循环队列各个參数的含义
（1）队列初始化时，front和tail值都为零
（2）当队列不为空时，front指向队列的第一个元素，tail指向队列最后一个元素的下一个位置；
（3）当队列为空时，front与tail的值相等，但不一定为零
（4）当（tail+1）% capacity == front ||  （length+1）== capacity 表示队列为满，
因此循环队列默认浪费1个空间
3.循环队列算法实现
（1）把值存在tail所在的位置；
（2）每插入1个元素，length+1，tail=（tail+1）% capacity
（3）每取出1个元素，length-1，front=（front+1）% capacity
（4）扩容功能，当队列容量满，即length+1==capacity时，capacity扩大为原来的2倍
（5）缩容功能，当队列长度小于容量的1/4，即length<=capacity/4时，capacity缩短为原来的一半

*/
// 循环队列实现方法
type loopQueue struct {
	queues   []interface{}
	front    int //队首
	tail     int //队尾
	length   int //队伍长度
	capacity int //队伍容量
}

func NewLoopQueue() *loopQueue {
	loop := &loopQueue{
		queues:   make([]interface{}, 0, 2),
		front:    0,
		tail:     0,
		length:   0,
		capacity: 2,
	}
	//初始化队列
	for i := 0; i < 2; i++ {
		loop.queues = append(loop.queues, "")
	}
	return loop
}

func (q *loopQueue) Len() int {
	return q.length
}

func (q *loopQueue) Cap() int {
	return q.capacity
}

func (q *loopQueue) IsEmpty() bool {
	return q.length == 0
}

func (q *loopQueue) IsFull() bool {
	return (q.length + 1) == q.capacity
}

func (q *loopQueue) GetFront() (interface{}, error) {
	if q.Len() == 0 {
		return nil, errors.New(
			"failed to getFront,queues is empty.")
	}
	return q.queues[q.front], nil
}

func (q *loopQueue) Enqueue(elem interface{}) {
	// 队列扩容
	if q.IsFull() {
		buffer := new(loopQueue)
		//初始化队列
		for i := 0; i < 2*q.capacity; i++ {
			buffer.queues = append(buffer.queues, "")
		}
		for i := 0; i < q.length; i++ {
			buffer.queues[i] = q.queues[q.front]
			q.front = (q.front + 1) % q.capacity
		}
		q.queues = buffer.queues
		q.front = 0
		q.tail = q.length
		q.capacity = 2 * q.capacity
	}

	q.queues[q.tail] = elem
	q.tail = (q.tail + 1) % q.capacity
	q.length++
}

func (q *loopQueue) Dequeue() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New(
			"failed to dequeue,queues is empty.")
	}

	// 当队列长度小于容量1/4时，队列容量缩短一半
	if q.length <= q.capacity/4 {
		buffer := new(loopQueue)
		//初始化队列
		for i := 0; i < q.capacity/2; i++ {
			buffer.queues = append(buffer.queues, "")
		}
		for i := 0; i < q.length; i++ {
			buffer.queues[i] = q.queues[q.front]
			q.front = (q.front + 1) % q.capacity
		}
		q.queues = buffer.queues
		q.front = 0
		q.tail = q.length
		q.capacity = q.capacity / 2
	}

	queue := q.queues[q.front]
	q.front = (q.front + 1) % q.capacity
	q.length--
	return queue, nil
}