package pipeline

import (
	"sync"
)

func Fork[I Input](ch <-chan I, n int) []<-chan I {
	ins := make([]chan I, n)
	for i := range ins {
		ins[i] = make(chan I)
	}

	go func() {
		defer func() {
			for _, is := range ins {
				close(is)
			}
		}()

		var wg sync.WaitGroup
		defer wg.Wait()

		for p := range ch {
			wg.Add(len(ins))

			for _, is := range ins {
				go func(is chan<- I, p I) {
					defer wg.Done()
					is <- p
				}(is, p)
			}
		}
	}()

	return toOutputChannels(ins)
}

func toOutputChannels[O Output](chans []chan O) []<-chan O {
	ins := make([]<-chan O, len(chans))

	for i, c := range chans {
		ins[i] = c
	}

	return ins
}
