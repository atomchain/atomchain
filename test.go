package main

import (
	"fmt"
	"time"
)

func foo() {
	go func() {
		// time.Sleep(5 * time.Second)
		time.Sleep(1 * time.Second)
		fmt.Print("hello")
	}()
}
func main() {
	foo()
	// time.Sleep(5 * time.Second)
}
