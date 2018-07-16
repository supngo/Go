package main

import "fmt"

type shape interface {
	getArea() float64
}

type square struct {
	sideLength float64
}

type triangle struct {
	base   float64
	height float64
}

func main() {
	sq := square{10}
	tr := triangle{5, 4}

	printArea(sq)
	printArea(tr)
}

func (sq square) getArea() float64 {
	return sq.sideLength * sq.sideLength
}

func (tr triangle) getArea() float64 {
	return tr.base * tr.height / 2
}

func printArea(sh shape) {
	fmt.Println(sh.getArea())
}
