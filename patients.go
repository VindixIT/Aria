package main

import (
	"fmt"
)

func WriteName() string {
	const name, age = "Fábio", 24
	fmt.Printf("%s is %d years old.\n", name, age)
}
