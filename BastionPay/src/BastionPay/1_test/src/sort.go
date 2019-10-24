package main

import (
	"fmt"
	"strconv"
	"strings"
	//"time"

	//"time"

	//"time"

	//"sort"
)

func main() {
	arr := []int{9, 4, 5, 1, 22, 8, 4, 7, 0, 2}
	//quick(arr,0,len(arr)-1)
	arr = mergeSort(arr)
	//t()
	//S(arr)
	//InsertSort2(arr)
	//	bubbleSort(arr, len(arr))
	//	fmt.Println(arr)
	//
	//fmt.Println(Decimal("1.0000032132"))
	//fmt.Println(Decimal7("0.0243"))
	//amountString := strconv.FormatFloat(0.000014099020627338266, 'f', -1, 64)
	//fmt.Println("a",amountString)
	////选择排序
	////arr := []int{2, 67, 33, 0, 45, 25, 77, 208, -8, -7}
    //selectSort(arr)
	fmt.Print(arr)
	//////c1 := make(chan interface{})
	//////close(c1)
	//////c2 := make(chan interface{})
	//////close(c2)
	//////var c1count, c2count int
	//////for i := 1000; i > 0; i-- {
	//////	select {
	//////	case <-c1:
	//////		fmt.Println("c1", <-c1)
	//////		c1count++
	//////	case <-c2:
	//////		c2count++
	//////	}
	//////}
	//////fmt.Printf("c1count:%d\n c2count:%d\n", c1count, c2count)
	////sort.Ints(arr)
	//fmt.Println(time.Now().UnixNano())
	//
	//channel := make(chan int)
	//go Test1(channel)
	//channel <- 1


	//快排
	//arr1 := []int{2, 1, 3,8,6,9 ,7}
	////quickSort(arr1, 0, len(arr1)-1)
	////fmt.Println(arr1)
	//
	//InsertSort(arr1)
	//t()

}

func Test1(ch chan int)  {
	fmt.Println("hello")
	<- ch
}


//  []int{9, 4, 5, 1, 22, 8, 4, 7, 0, 2}
/*
/ 选择排序：两层循环，一个一个的和后面的比，把小的提到前面，
 */
func selectSort(arr []int) {
	len := len(arr)
	for i := 0; i < len; i++ {
		for j := i + 1; j < len; j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]

				fmt.Println("arr",i,j,arr)
			}
		}
	}
}

//   []int{9, 4, 5, 1, 22, 8, 4, 7, 0, 2}
//         0                           9
//        low                         high

//快速排序
func quick(arr []int, start, end int)  {
	if start >= end {
		return
	}
	low := start
	high := end
	lv := arr[start]

	for low < high {
		for low < high && arr[high] >= lv {
			high --
		}
		//出了循环说明 high 小于等于 mid 则需要交换位置
		arr[high], arr[low] = arr[low], arr[high]

		for low < high && arr[low] < lv{
			low ++
		}
		// 出了循环说明 low 大于等于mid 则需要交换位置
		arr[high], arr[low] = arr[low], arr[high]
	}
	// 大的循环退出以后就找到了 mid的位置,此时mid 左边的都小于mid, 右边都大于mid,以 mid 为界限再排序左右两边
	arr[low] = lv
	quick(arr, start, low - 1) // 排序左边的
	quick(arr, low + 1, end ) // 排序右边的

}


//归并排序
func mergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	i := len(arr) / 2
	left := mergeSort(arr[0:i])
	right := mergeSort(arr[i:])
	result := merge(left, right)
	return result
}

func merge(left, right []int) []int {
	result := make([]int, 0)
	m, n := 0, 0 // left和right的index位置
	l, r := len(left), len(right)
	for m < l && n < r {
		if left[m] > right[n] {
			result = append(result, right[n])
			n++
			continue
		}
		result = append(result, left[m])
		m++
	}
	result = append(result, right[n:]...) // 这里竟然没有报数组越界的异常？
	result = append(result, left[m:]...)
	return result
}


/*      //a[j-1]  a[j]  temp
  插入排序 ： []int {9,     4,      5,   1, 22, 8, 4, 7, 0, 2}
 */
//插入排序 ： []int{9,       4, 5, 1, 22, 8, 4, 7, 0, 2}

func InsertSort2(values []int) {
	length := len(values)
	//if length <= 1 {
	//	return
	//}

	for i := 1; i < length; i++ {
		tmp := values[i] // 从第二个数开始，从左向右依次取数
		key := i - 1     // 下标从0开始，依次从左向右
		// 每次取到的数都跟左侧的数做比较，如果左侧的数比取的数大，就将左侧的数右移一位，直至左侧没有数字比取的数大为止
		for key >= 0 &&  values[key] > tmp {
			values[key+1] = values[key]
			key--
			//fmt.Println(values)
		}

		// 将取到的数插入到不小于左侧数的位置
			values[key+1] = tmp

	}

}




func t(){
	chans := make([]int,10)
	for i, v := range chans{
		fmt.Println(i,v)
	}
	for i := range chans{
		fmt.Println(i)
	}
	s  := cap(chans)
	fmt.Println(s)
}

//冒泡排序
func bubbleSort(arr []int, len int) {
	if len == 1 {
		return
	}
	for i := 0; i < len-1; i++ {
		if arr[i] > arr[i+1] {
			arr[i], arr[i+1] = arr[i+1], arr[i]
		}
	}
	bubbleSort(arr, len-1)
}

func Decimal(value string) string {



	float,err := strconv.ParseFloat(value,64)
	if err != nil {
		return ""
	}
	f64, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", float), 64)
	s2 := strconv.FormatFloat(f64, 'g', -1, 64)//float64
	return s2
}

func Decimal8(value string) string {
	float,err := strconv.ParseFloat(value,64)
	if err != nil {
		return ""
	}
	f64, _ := strconv.ParseFloat(fmt.Sprintf("%.9f", float), 64)
	s2 := strconv.FormatFloat(f64, 'g', -1, 64)//float64
	s2 = s2[:len(s2)-1]
	return s2
}


func Decimal7(value string) string {
	if len(value) < 10 {
		return value
	}
	int := strings.IndexAny(value, ".")
	if len(value) - int <= 8 {
		return  value
	}
	s2 := value[0:int+9]
	return s2
}

//选择排序
func S(arr []int)  {
	for i:= 0; i < len(arr); i++  {
		for j:= i+1; j < len(arr); j++  {
			if arr[j] < arr[i] {
				arr[i] , arr[j] = arr[j], arr[i]
			}
		}
	}
}