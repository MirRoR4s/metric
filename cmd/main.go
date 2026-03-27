package main

import "fmt"

func f(ch chan int, v int) {
	ch <- v
	ch <- v * 2
	ch <- v * 3
	ch <- v * 7
	close(ch)
}

func main() {

	ch := make(chan int)

	go f(ch, 2)

	for v := range ch {
		fmt.Println(v)
	}
}
