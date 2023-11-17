package pipeline

import (
	"sync"
)

type (
	Input  any
	Output any
	InOut  interface {
		Input
		Output
	}
)

type (
	Step[I Input, O Output]        func(<-chan I) <-chan O
	Transformer[I Input, O Output] func(I) O
	Predicate[I Input]             func(I) bool
)

func Identity[I Input]() Transformer[I, I] {
	return func(i I) I {
		return i
	}
}

func AcceptAll[I Input]() Predicate[I] {
	return func(i I) bool {
		return true
	}
}

func Compose[I Input, IO InOut, O Output](p1 Step[I, IO], p2 Step[IO, O]) Step[I, O] {
	return func(in <-chan I) <-chan O {
		return p2(p1(in))
	}
}

func Series[I InOut](steps ...Step[I, I]) Step[I, I] {
	p := steps[0]

	for _, step := range steps[1:] {
		p = Compose(p, step)
	}

	return p
}

func Parallel[I Input, O Output](steps ...Step[I, O]) Step[I, O] {
	return func(in <-chan I) <-chan O {
		inputs := Fork(in, len(steps))

		outs := make([]<-chan O, len(steps))
		for i, step := range steps {
			outs[i] = step(inputs[i])
		}

		return Merge(outs...)
	}
}

func Map[I Input, O Output](transformer func(I) O) Step[I, O] {
	return func(in <-chan I) <-chan O {
		return Workers(transformer, AcceptAll[O](), 1, in)
	}
}

func Filter[I Input](predicate func(I) bool) Step[I, I] {
	return func(in <-chan I) <-chan I {
		return Workers(Identity[I](), predicate, 1, in)
	}
}

func MapPar[I Input, O Output](transformer func(I) O) Step[I, O] {
	return func(in <-chan I) <-chan O {
		return Workers(transformer, AcceptAll[O](), 4, in)
	}
}

func Workers[I Input, O Output](
	transformer func(I) O,
	predicate func(O) bool,
	n int,
	input <-chan I,
) <-chan O {
	out := make(chan O, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for in := range input {
				if o := transformer(in); predicate(o) {
					out <- o
				}
			}
		}()
	}

	go func() {
		defer close(out)
		defer wg.Wait()
	}()

	return out
}

type Pair[F any, S any] struct {
	First  F
	Second S
}

func MapToChan[K comparable, V any](m map[K]V) <-chan Pair[K, V] {
	ch := make(chan Pair[K, V], len(m))

	go func() {
		defer close(ch)
		for k, v := range m {
			ch <- Pair[K, V]{k, v}
		}
	}()

	return ch
}

func Consume[O Output](ch <-chan O) []O {
	l := make([]O, 0, 100)

	for v := range ch {
		l = append(l, v)
	}

	return l
}
