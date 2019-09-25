package main

import "fmt"

// Heap 定义堆排序过程中使用的堆结构
type Heap struct {
	arr  []int   // 用来存储堆的数据
	size int     // 用来标识堆的大小
}

// adjustHeap 用于调整堆，保持堆的固有性质
func adjustHeap(h Heap, parentNode int) {
	leftNode := parentNode*2 + 1
	rightNode := parentNode*2 + 2

	maxNode := parentNode
	if leftNode < h.size && h.arr[maxNode] < h.arr[leftNode] {
		maxNode = leftNode
	}
	if rightNode < h.size && h.arr[maxNode] < h.arr[rightNode] {
		maxNode = rightNode
	}

	if maxNode != parentNode {
		h.arr[maxNode], h.arr[parentNode] = h.arr[parentNode], h.arr[maxNode]
		adjustHeap(h, maxNode)
	}
}

// createHeap 用于构造一个堆
func createHeap(arr []int) (h Heap) {
	h.arr = arr
	h.size = len(arr)

	for i := h.size / 2; i >= 0; i-- {
		adjustHeap(h, i)
	}
	return
}

// heapSort 使用堆对数组进行排序
func heapSort(arr []int) {
	h := createHeap(arr)

	for h.size > 0 {
		// 将最大的数值调整到堆的末尾
		h.arr[0], h.arr[h.size-1] = h.arr[h.size-1], h.arr[0]
		// 减少堆的长度
		h.size--
		// 由于堆顶元素改变了，而且堆的大小改变了，需要重新调整堆，维持堆的性质
		adjustHeap(h, 0)
	}
}

func main() {
	// 测试代码
	arr := []int{9, 8, 7, 6, 5, 1, 2, 3, 4, 0}
	fmt.Println(arr)
	heapSort(arr)
	fmt.Println(arr)
}