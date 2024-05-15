package main

import "fmt"

func main() {
	greet("Alice")

	a := "foo"
	b := "bar"
	fmt.Println(concat(a, b))

	result := calculateArea(5, 8)
	fmt.Println("Area of the rectangle:", result)
}

func greet(name string) {
	fmt.Println("Hello,", name)
}

func concat(a string, b string) string {
	return a + b
}

func calculateArea(length int, width int) int {
	area := length * width
	return area
}
