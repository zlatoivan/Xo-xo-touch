package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			c1 <- "from 1"
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			c2 <- "from 2"
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			select {
			case msg1 := <-c1:
				fmt.Println(msg1)
			case msg2 := <-c2:
				fmt.Println(msg2)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}

// Unmarshall = json -> struct
// Marshall   = struct -> json

//// Create a Person struct.
//person := Person{
//	Name: "Krunal Lathiya",
//	Age:  30,
//}
//
//// Marshal the Person struct to JSON.
//j, err := json.Marshal(person)
//if err != nil {
//	fmt.Println(err)
//	return
//}
//
//jsonString := string(j)
//
//fmt.Println(jsonString)

//type Person struct {
//	Name string
//	Age  int
//}
//
//func main() {
//	person := Person{}
//
//	jsonString := `{"name":"Krunal Lathiya","age":30}`
//
//	_ = json.Unmarshal([]byte(jsonString), &person)
//
//	fmt.Println(person)
//}
