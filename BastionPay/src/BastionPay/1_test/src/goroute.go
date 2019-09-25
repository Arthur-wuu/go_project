package src

import "fmt"


func main() {
	chs := make([]chan int , 10)   //[]chan int 表示是chan int类型的切片 和   chan []int区别  标识通道类型是[]int
	//c := make(chan []int)
	//t := []int{1,2,3}
	//c <- t
	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go Count1(chs[i])         //创建10个协程
		fmt.Println("Countaaaaa",i)
	}
	for i, ch := range chs {
		<-ch
		fmt.Println("Counting------",i)
	}
	str :="select * from t_roles"

	s :=[]byte(str)
	fmt.Println("s  ",s)
}       //调度的顺序不同，可能goroutinue执行的顺序不一样   /v1/filetransfer/export?sql=select * from t_roles&dbname=admin&pagesize=2&expire_at=30&max_op_time=30&pagehalt=1&max_lines=10&order_id=222222222&file_name=newdata

func Count1(ch chan int) {
	ch <- 1
	fmt.Println("Counting`````")
	str :="select * from t_roles"

	  s :=[]byte(str)
	  fmt.Println(s)
}


