package main

import (
	"fmt"
)

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contactInfo
}

func main() {
	// alex := person{"Alex", "Anderson", contactInfo{"alex@email.com", 20172}}
	// peter := person{firstName: "Peter", lastName: "Parker", contactInfo: contactInfo{"peter@email.com", 20101}}
	var tony person
	// fmt.Println(alex)
	// fmt.Println(peter)
	// tony.print()

	// tonyPointer := &tony
	// tonyPointer.updateName("Anthony")
	// tony.print()

	tony.updateName("Anthony")
	tony.print()

	// name := "bill"
	// namePointer := &name
	// // will print different address below because everything in Go is pass-by-value
	// fmt.Println(&namePointer)
	// printPointer(namePointer)
}

// func printPointer(namePointer *string) {
// 	fmt.Println(&namePointer)
// }

func (p person) print() {
	fmt.Printf("%+v", p)
}

func (pointer *person) updateName(newFN string) {
	(*pointer).firstName = newFN
	// pointer.firstName = newFN
}
