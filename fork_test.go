package pipeline

import (
	"testing"
)

func TestFork(t *testing.T) {
	g := generator(10)

	forks := Fork(g, 5)

	merge := Merge(forks...)

	all := Consume(merge)

	if len(all) != 50 {
		t.Errorf("Expected 50, got %d", len(all))
	}
}
