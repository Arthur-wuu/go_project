package main
import (
	"fmt"
	"strings"
)

func main() {

	//for i := 0; i < 10; i++ {
		//fmt.Println(fab(2))
	//}
	s1 := strings.ReplaceAll("tabfreesssss", "tabfree", "rlfree")
	//s := firstUniqChar("aaaaaaavvvc")

	fmt.Println(s1)
}


//斐波那
func fab(n int) int {
	if n <= 1 {
		return 1
	}
	return fab(n-1) + fab(n-2)
}



func FindFirst(s string) int  {
	//第一，把字符串变成 字节数组

	sbytes := []byte(s)

	smap := make(map[byte]int)


	//字节数组 一个个计数
	for _, v  := range sbytes{
		smap[v] ++
	}

	//然后找这个 int = 1的

	for i , v  := range sbytes{
		if smap[v] == 1 {
			return  i
		}
	}

	return -1

}




















