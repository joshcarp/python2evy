package main

import "fmt"

func main() {
	person := map[string]interface{}{
		"name": "Bob",
		"age":  30,
		"city": "New York",
	}

	fmt.Println(person["name"])
	person["age"] = 31
	fmt.Println(person)
}
