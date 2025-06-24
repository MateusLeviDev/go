package main

import (
	"fmt"
	"levigo/concurrency"
)

func main() {
	// c := concurrency.Boring("call")
	// j := concurrency.Boring("joe")
	// a := concurrency.Boring("ann")
	c := concurrency.FanIn(concurrency.Boring("x"), concurrency.Boring("y"))

	for i := 0; i < 10; i++ {
		// msg1 := <-j
		// fmt.Println(msg1.Str)
		// msg2 := <-a
		// fmt.Println(msg2.Str)
		// msg1.Wait <- true
		// msg2.Wait <- true
		msg1 := <-c
		fmt.Println(msg1.Str)
		msg1.Wait <- true
	}

	fmt.Println("Leaving")
}
