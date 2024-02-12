package main

import (
	"fmt"
)

func main() {
	chanOwner := func() <-chan int {
		out := make(chan int, 5) // 1. create

		go func() {
			defer close(out) // 4. close
			for i := 0; i <= 5; i++ {
				out <- i // 3. write
			}
		}()

		return out // 2. return
	}

	c := chanOwner()
	for i := 0; i < 5; i++ {
		fmt.Println(<-c)
	}
}
