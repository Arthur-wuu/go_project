package main

import (
	//"log"
	"fmt"
	"os"
	//"sync"
	"time"
)

func main() {
	word := []string{"AB","CD","EF","GH"}
	//num := []int{1,2,3,4,5,6,7,8}

	///使用两个goroutine交替打印序列，一个goroutinue打印数字，
	// 另外一个goroutine打印字母，最终效果如下12AB34CD56EF78GH910IJ。

	chN := make(chan int, 1)
	chN <- 1
	chC := make(chan int, 1)

//var mu sync.Mutex

	go func() {
		for i:=1; i< 8; i=i+2  {
			//mu.Lock()
			<- chN
			fmt.Print(i)
			fmt.Print(i+1)
			//mu.Unlock()
			chC <- 1
		}
	}()


	go func() {
		for i:=0; i< 8; i++  {
			//mu.Lock()
			<- chC
			fmt.Print(word[i])
			//mu.Unlock()
			chN <- 1
		}
	}()

	time.Sleep(time.Second * 11)

}


func printNum (ch chan int)  {
	ch <- 1
	fmt.Print("12")
}


func printWord (ch chan int)  {
	 <- ch
	fmt.Print("AB")
}













//	os.Create("./xa.txt")
//	os.Create("./xb.txt")
//	os.Create("./xc.txt")
//	os.Create("./xd.txt")
//
//
//	a, err := os.OpenFile("./xa.txt", os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		fmt.Println("err",err)
//	}
//	b, _ := os.OpenFile("./xb.txt", os.O_WRONLY|os.O_APPEND, 0666)
//	c, _ := os.OpenFile("./xc.txt", os.O_WRONLY|os.O_APPEND, 0666)
//	d, _ := os.OpenFile("./xd.txt", os.O_WRONLY|os.O_APPEND, 0666)
//	files := []*os.File{a, b, c, d}
//	i := 0
//	sign := make(chan int, 1)
//	for i < 100 {
//		i++
//		sign <- 1
//		go out1(files[0], sign)
//		sign <- 1
//		go out2(files[1], sign)
//		sign <- 1
//		go out3(files[2], sign)
//		sign <- 1
//		go out4(files[3], sign)
//
//		files = append(files[len(files)-1:], files[:len(files)-1]...)
//	}
//	a.Close()
//	b.Close()
//	c.Close()
//	d.Close()


func out1(f *os.File, c chan int) int {
	f.Write([]byte("1 "))
	// f.Close()
	//log.Println(f.Name() + " write finish...")
	<-c
	return 1
}

func out2(f *os.File, c chan int) int {
	f.Write([]byte("2 "))
	// f.Close()
	//log.Println(f.Name() + " write finish...")
	<-c
	return 2
}
func out3(f *os.File, c chan int) int {
	f.Write([]byte("3 "))
	// f.Close()
	//log.Println(f.Name() + " write finish...")
	<-c
	return 3
}
func out4(f *os.File, c chan int) int {
	f.Write([]byte("4 "))
	// f.Close()
	//log.Println(f.Name() + " write finish...")
	<-c
	return 4
}