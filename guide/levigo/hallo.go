package main

import "fmt"

const deutschHalloPrefix = "Hallo, "

func Hallo(name string) string {
	if name == "" {
		name = "Welt"
	}
	return deutschHalloPrefix + name
}

func main() {
	fmt.Println(Hallo(""))
}
