package main

import (
	"fmt"
	"strconv"
)

func main1() {
	list := []int{}

	for i := 0; i <= 10; i++ {
		list = append(list, i)
	}

	for _, value := range list {
		if value%2 == 0 {
			fmt.Println(strconv.Itoa(value) + " is even")
		} else {
			fmt.Println(strconv.Itoa(value) + " is odd")
		}
	}
}
