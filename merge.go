package pipeline

import "sync"

func Merge[O Output](channels ...<-chan O) <-chan O {
	result := make(chan O, len(channels))

	go func() {
		defer close(result)

		var wg sync.WaitGroup
		wg.Add(len(channels))

		defer wg.Wait()

		for _, c := range channels {
			go func(c <-chan O) {
				defer wg.Done()
				for o := range c {
					result <- o
				}
			}(c)
		}
	}()

	return result
}
