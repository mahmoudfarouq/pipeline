package pipeline

import "context"

func OrDone[I Input, O Output](ctx context.Context, step Step[I, O]) Step[I, O] {
	return func(is <-chan I) <-chan O {
		out := make(chan O, 100)

		go func() {
			defer close(out)
			
			s := step(is)

			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-s:
					if ok {
						out <- t
					} else {
						return
					}
				}
			}

		}()

		return out
	}
}
