package pipeline

import (
	"context"
	"log"
	"testing"
)

func TestOrDone(t *testing.T) {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	s := Series(
		factorNumsFilter(),
	)

	s = OrDone(ctx, s)

	n := 10

	in := Merge(oddNumsGen(n), evenNumsGen(n))

	p := s(in)

	for o := range p {
		log.Printf("%v", o)
		cancel()
	}
}
