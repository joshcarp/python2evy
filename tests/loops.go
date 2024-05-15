package main

import "fmt"

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println("for", i)
	}

	count := 0
	for count < 5 {
		fmt.Println("while", count)
		count++
	}

	for i := 1; i < 4; i++ {
		for j := 1; j < 4; j++ {
			if i != j {
				fmt.Printf("(%v, %v)\n", i, j)
			}
		}
	}
}
