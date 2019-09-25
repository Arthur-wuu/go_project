package main

import (
	"crypto/sha1"
	"fmt"
)

func Count(ch chan int) {
	ch <- 1
	fmt.Println("Counting")
}

func main() {
	c :=Sha1("dsa")
	fmt.Println("cc",c)
	//chs := make([] chan int, 10)
	//
	//for i:=0; i<10; i++ {
	//	chs[i] = make(chan int)
	//	go Count(chs[i])
	//}
	//
	//for _, ch := range(chs) {
	//	<-ch
	//}
}
func Sha1(data string) []byte {

	// stage1Hash = SHA1(password)
	crypt := sha1.New()
	crypt.Write([]byte(data))
	stage1 := crypt.Sum(nil)

	return stage1
}