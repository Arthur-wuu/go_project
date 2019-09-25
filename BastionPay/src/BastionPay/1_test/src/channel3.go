package main

import (
	"fmt"
	"sync"
)

func main() {
	chA := make(chan int)
	chB := make(chan int)
	chBroadcast := make(chan int)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		for v := range chBroadcast {
			chA <- v
			chB <- v
		}
		close(chA)
		close(chB)
		wg.Done()
	}()

	go func() {
		for v := range chA {
			fmt.Println("A", v)
		}
		wg.Done()
	}()

	go func() {
		for v := range chB {
			fmt.Println("B", v)
		}
		wg.Done()
	}()

	for i := 0; i < 10; i++ {
		chBroadcast <- i
	}
	close(chBroadcast)
	wg.Wait()
}