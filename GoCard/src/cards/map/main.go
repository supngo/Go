package main

import "fmt"

func main() {
	colors := map[string]string{
		"red":   "#FF0000",
		"green": "#00FF00",
	}
	// fmt.Println(colors)
	printMap(colors)

	// new empty map
	var colors2 map[string]string
	fmt.Println(colors2)

	// new empty map
	colors3 := make(map[int]string)
	colors3[1] = "#FFFFFF"
	fmt.Println(colors3)
	delete(colors3, 1)
	fmt.Println(colors3)
}

func printMap(c map[string]string) {
	for color, hex := range c {
		fmt.Println("hex code for", color, "is", hex)
	}
}
