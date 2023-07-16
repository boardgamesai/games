package main

// This is a dummy file that does nothing but satisfy Go's need for every main package to have a main()

import "fmt"

func main() {
	ai := &AI{}
	fmt.Printf("ai: %+v\n", ai)
}
