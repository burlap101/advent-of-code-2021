package main

import (
	"testing"
)

func TestStateFromExample1(t *testing.T) {
	state, err := processInputFile("example1.txt")
	if err != nil {
		panic(err)
	}

	t.Run("TestAHops", func(t *testing.T) {
		edges := [...]Edge{{"A", "b"},{"A", "c"},{"A", "end"}}
	})
}