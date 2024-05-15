package main

import "fmt"

func main() {
	var fruits []string = []string{"apple", "banana", "orange"}
	fmt.Println(fruits[0])
	fruits = append(fruits, "grape")
	fmt.Println(fruits)
}
