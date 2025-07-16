package utils

import "sync"

func Split(ch <-chan func(), n int) []<-chan func() {

	cs := make([]chan func(), n)
	for i := 0; i < n; i++ {
		cs[i] = make(chan func())
	}

	distributeToChannels := func(ch <-chan func(), cs []chan func()) {
		defer func(cs []chan func()) {
			for _, c := range cs {
				close(c)
			}
		}(cs)

		for {
			for _, c := range cs {
				select {
				case val, ok := <-ch:
					if !ok {
						return
					}
					c <- val
				}
			}
		}
	}

	go distributeToChannels(ch, cs)

	result := make([]<-chan func(), n)
	for i := 0; i < n; i++ {
		result[i] = cs[i]
	}
	return result
}

func Worker(ch <-chan func(), wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range ch {
		task()
	}
}
