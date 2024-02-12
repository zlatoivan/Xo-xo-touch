package main

import "fmt"

func greet(c chan string) {
	<-c // for John
	<-c // for Mike
}

func main() {
	fmt.Println("main() started")

	c := make(<-chan int, 10)

	//go greet(c)
	//c <- 123
	val, ok := <-c

	//close(c) // closing channel
	val, ok = <-c
	fmt.Println(val, " ", ok)

	//c <- "Mike"
	fmt.Println("main() stopped")
}
