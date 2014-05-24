package digraph

import (
	"testing"
)

func TestStronglyConnectedComponents(t *testing.T) {
	nodes := ParseBasic(`a -> b
a -> c
b -> c
c -> b
c -> d
d -> e`)
	var nlist []Node
	for _, n := range nodes {
		nlist = append(nlist, n)
	}

	sccs := StronglyConnectedComponents(nlist, false)
	if len(sccs) != 4 {
		t.Fatalf("bad: %v", sccs)
	}

	sccs = StronglyConnectedComponents(nlist, true)
	if len(sccs) != 1 {
		t.Fatalf("bad: %v", sccs)
	}

	cycle := sccs[0]
	if len(cycle) != 2 {
		t.Fatalf("bad: %v", sccs)
	}

	if cycle[0].(*BasicNode).Name != "c" {
		t.Fatalf("bad: %v", cycle)
	}
	if cycle[1].(*BasicNode).Name != "b" {
		t.Fatalf("bad: %v", cycle)
	}
}

func TestStronglyConnectedComponents2(t *testing.T) {
	nodes := ParseBasic(`a -> b
a -> c
b -> d
b -> e
c -> f
c -> g
g -> a
`)
	var nlist []Node
	for _, n := range nodes {
		nlist = append(nlist, n)
	}

	sccs := StronglyConnectedComponents(nlist, true)
	if len(sccs) != 1 {
		t.Fatalf("bad: %v", sccs)
	}

	cycle := sccs[0]
	if len(cycle) != 3 {
		t.Fatalf("bad: %v", sccs)
	}

	if cycle[0].(*BasicNode).Name != "g" {
		t.Fatalf("bad: %v", cycle)
	}
	if cycle[1].(*BasicNode).Name != "c" {
		t.Fatalf("bad: %v", cycle)
	}
	if cycle[2].(*BasicNode).Name != "a" {
		t.Fatalf("bad: %v", cycle)
	}
}
