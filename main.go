package main

import (
	"fmt"
	"strconv"
)

func main() {
	val, err := strconv.Atoi("1")
	if err != nil {
		fmt.Println("error with 1")
	} else {
		fmt.Println(val)
	}

	val2, err := strconv.Atoi("one")
	if err != nil {
		fmt.Println("error with one")
	} else {
		fmt.Println(val2)
	}
}
