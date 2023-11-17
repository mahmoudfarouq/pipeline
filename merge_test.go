package pipeline

import (
	"testing"
)

func TestMerge(t *testing.T) {
	g1 := generator(10)
	g2 := generator(10)
	g3 := generator(10)

	chain := Merge(g1, g2, g3)

	out := Consume(chain)

	if len(out) != 30 {
		panic("shit")
	}
}
