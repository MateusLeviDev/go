package concurrency

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	Str  string
	Wait chan bool
}

func Boring(msg string) <-chan Message {
	c := make(chan Message)
	go func() {
		for i := 0; ; i++ {
			waitForIt := make(chan bool)
			c <- Message{
				Str:  fmt.Sprintf("%s %d", msg, i),
				Wait: waitForIt,
			}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c
}

//	func FanIn(input1, input2 <-chan string) <-chan string {
//		c := make(chan string)
//		go func() {
//			for {
//				c <- <-input1
//			}
//		}()
//		go func() {
//			for {
//				c <- <-input2
//			}
//		}()
//		return c
//	}

func FanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}
