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
	//S(arr)
	InsertSort(arr)
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


func quickSort(arr []int, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		fmt.Println("key",key)
		for i <= j {                     // 0  14
			for arr[i] < key {           // 3 < 42
				i++                      // 1
			}
			for arr[j] > key {           //
				j--
			}
			if i <= j {
				fmt.Println(i,j)
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}

		if start < j {
			fmt.Println("start j *******", start, j )
			quickSort(arr, start, j)
		}
		if end > i {
			fmt.Println("end i #######", end ,i )
			quickSort(arr, i, end)
		}
	}
}

/*
  插入排序 ： []int{9, 4, 5, 1, 22, 8, 4, 7, 0, 2}
 */
func InsertSort(a []int) {
	var j int
	for i := 1; i < len(a); i++ {
		temp := a[i]
		for j = i - 1; j >= 0 && a[j] > temp; j-- {
			a[j+1] = a[j]
		}
		fmt.Println("a[j+1]",j+1,a[j+1])
		a[j+1] = temp
		fmt.Println(a)
	}
	fmt.Println(a)
}

func t(){
	chans := make([]int,10)
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