package main

import (
	"fmt"
	"time"

	//"sort"

	//	"sort"
	//"unsafe"
)

//对数组进行冒泡排序
func BubbleSort(values []int) []int {
	var arrlen int =len(values)
	for i:=0;i<arrlen-1;i++{   //10shu  0-9

		for j:=0;j<arrlen-1-i;j++{
			if values[j]>values[j+1]{
				values[j],values[j+1]=values[j+1],values[j]
				fmt.Println("arr",values)
			}
		}
	}
	return values
}

//func quickSort(values [] int,left,right int)  {
//
//	temp:=values[left]
//
//	p:=left
//
//	i,j:=left,right
//
//	for i<=j{
//
//		for j>=p&&values[j]>=temp{
//			j--
//		}
//
//		if j>=p{
//			values[p]=values[j]
//			p=j
//		}
//
//		if values[i]<=temp && i<=p{
//			i++
//		}
//
//		if i<=p{
//			values[p]=values[i]
//			p=i
//		}
//	}
//
//	values[p] =temp
//
//	if p-left >1{
//		quickSort(values,left,p-1)
//	}
//
//	if right-p>1{
//		quickSort(values,p+1,right)
//	}
//}

//对数组进行快速排序
//func QuickSort(values [] int)  {
//
//	var arrlen int =len(values)
//	quickSort(values,0,arrlen-1)
//}
//
//type dog interface {
//	name()
//}

func main()  {
	//println(time.Now().Format("2006-01-02 15:04:05"))
	//values:= []int{39, 23, 3 ,2, 5, 93, 45, 208, 27, 67}
	//
	//BubbleSort(values)
	////
	////for i:=0;i<len(values);i++ {
	////	sort.IsSorted(dog)
	//	fmt.Println(values)

	//sort.Sort()
//	GetTimeFromStrLoc()
	//fmt.Println("t",times)
	a :=GetTimeStamp()
	fmt.Println(a)
	}

func GetTime() string {
	const shortForm = "2006-01-01 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	fmt.Println(t)
	return str
}

func GetTimeFromStr() {
	const format = "2006-01-02 15:04:05"
	timeStr := "2018-01-09 20:24:20"
	p, err := time.Parse(format, timeStr)
	if err == nil {
		fmt.Println(p)
	}
}

//带时区匹配，匹配当前时区的时间
func GetTimeFromStrLoc() {
	x := time.Date(2017, 02, 27, 17, 30, 20, 20, time.Local)
	fmt.Println(x.Format("2006-01-02 15:04:05"))
}

func GetTimeStamp() string {
	times := time.Now().Format("2006-01-02 15:04:05")
	return times
}
//	var a uint = 0
//	p := uintptr(unsafe.Pointer(&a))
//	//for i := 0; i < int(unsafe.Sizeof(a)); i++ {
//	//	p += 1
//	//	pb := (*byte)(unsafe.Pointer(p))
//	//	*pb = 1
//	//}
//	fmt.Printf("%x\n", a) //0x1010100
//
//fmt.Println("p",p)


//}


