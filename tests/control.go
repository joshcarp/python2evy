package main

import "fmt"

func main() {
	age := 25
	if age >= 18 {
		fmt.Println("You are an adult.")
	} else {
		fmt.Println("You are a minor.")
	}
	count := 0
	for count < 5 {
		fmt.Println(count)
		count++
	}
}
