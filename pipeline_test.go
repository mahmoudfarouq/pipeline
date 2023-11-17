package pipeline

import (
	"log"
	"strconv"
	"testing"

	"modernc.org/mathutil"
)

func TestFullFlow(t *testing.T) {
	s := Series(
		factorNumsFilter(),
		greaterThan10Filter(),
		Parallel(
			multiplyBy10(),
			add5(),
			Map(func(i int) int {
				return i
			}),
		),
		Compose(
			Map(func(i int) string {
				return strconv.Itoa(i)
			}),
			FilterMap(func(s string) Option[int] {
				if len(s) <= 2 {
					return Option[int]{}
				}

				i, _ := strconv.Atoi(s)
				return Option[int]{Some: &i}
			}),
		),
	)

	n := 10

	in := Merge(oddNumsGen(n), evenNumsGen(n))

	p := s(in)

	for o := range p {
		log.Printf("%v", o)
	}
}

func add5() Step[int, int] {
	return Map(func(i int) int {
		return i + 5
	})
}

func multiplyBy10() Step[int, int] {
	return Map(func(i int) int {
		return i * 10
	})
}

func oddNumsGen(n int) <-chan int {
	return smartGenerator(n, func(i int) bool {
		return i%2 != 0
	})
}

func evenNumsGen(n int) <-chan int {
	return smartGenerator(n, func(i int) bool {
		return i%2 == 0
	})
}

func factorNumsFilter() Step[int, int] {
	return FilterMap(func(i int) Option[int] {
		if mathutil.IsPrime(uint32(i)) {
			return Option[int]{Some: &i}
		} else {
			return Option[int]{}
		}
	})
}

func greaterThan10Filter() Step[int, int] {
	return Filter(func(i int) bool {
		return i > 10
	})
}

func generator(n int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := 0; i < n; i++ {
			out <- i
		}
	}()

	return out
}

func smartGenerator(n int, predicate Predicate[int]) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		var counter int
		for i := 0; counter < n; i++ {
			if predicate(i) {
				out <- i
				counter++
			}
		}
	}()

	return out
}
