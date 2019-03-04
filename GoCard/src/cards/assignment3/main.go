package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println(os.Args[1])
	byteSlice, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println(string(byteSlice))
}
