package main

import "fmt"

const (
	filetype = ".txt"
)

func check(err error) {
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func main() {
	fmt.Printf("nothing here yet")
}
