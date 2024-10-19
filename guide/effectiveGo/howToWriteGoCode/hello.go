package main

import (
	"fmt"

	"lev.com/morestrings"
	"github.com/google/go-cmp/cmp"
)

func main() {

	fmt.Println("Another, Hello Fucking World")
	fmt.Println(morestrings.ReverseRunes("ivel"))
	fmt.Println(cmp.Diff("Hello", "Go"))

}