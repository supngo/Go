package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"https://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.com",
	}

	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c)
	}

	// loop to all links to print from channel
	// for i := 0; i < len(links); i++ {
	// 	fmt.Println(<-c)
	// }

	// Infinite loop
	// for {
	// 	go checkLink(<-c, c)
	// }

	// Infinite loop but wait for channel
	for l := range c {
		// go checkLink(l, c)

		// functional literal (like lambda)
		go func(link string) {
			time.Sleep(3 * time.Second)
			checkLink(link, c)
		}(l)
	}
}

func checkLink(link string, c chan string) {
	// fmt.Println("accessing", link)
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "migh be down!")
		// c <- "migh be down!"
		c <- link
		return
	}
	fmt.Println(link, "is up!")
	// c <- "is up!"
	c <- link
	// time.Sleep(5 * time.Second)
}
