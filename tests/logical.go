package main

import "fmt"

func main() {
	x := 10
	y := 5
	fmt.Println(x > y)
	fmt.Println(x < y)
	fmt.Println(x == y)
	fmt.Println(x != y)
	fmt.Println(x > 5 && y < 10)
	fmt.Println(x > 5 || y > 10)
	fmt.Println(!(x > 5))
}
