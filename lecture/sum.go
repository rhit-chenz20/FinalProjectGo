package main

import "fmt"

func sum(c chan int, co chan int) {
	for {
		x := <-c
		y := <-c
		fmt.Println(x, "+", y, "=", x+y)
		co <- x + y
	}
}

// go run sum
func main() {
	c := make(chan int, 6)
	co := make(chan int, 3)
	go sum(c, co)
	c <- 10
	c <- 15
	c <- 99
	c <- 1
	c <- 23
	c <- 7

	r := <-co
	fmt.Println(r)
	r = <-co
	fmt.Println(r)
	r = <-co
	fmt.Println(r)
}