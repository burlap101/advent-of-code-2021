package main

import (
	"testing"
	"os"
	"strings"
	"sort"
)

func TestPathsFromExample1(t *testing.T) {
	state, err := processInputFile("example1.txt")
	if err != nil {
		t.Error(err)
	}
	
	pathCount, err := state.FindPaths()
	if err != nil {
		t.Error(err)
	}

	t.Run("TestPathCountEqual", func(t *testing.T) {
		if pathCount != 36 {
			t.Errorf("PathCount = %d; want 36", pathCount)
		}
	})

	
	dat, err := os.ReadFile("paths1.txt")
	if err != nil {
		t.Error(err)
	}
	expectedPaths := strings.Split(string(dat), "\n")

	sort.Strings(expectedPaths)


	for pathStr := range state.Paths {
		i := sort.SearchStrings(expectedPaths, pathStr)
		if i >= len(expectedPaths) {
			t.Errorf("PathString = %s; invalid.", pathStr)
		} else {
			upper := expectedPaths[i+1:]
			expectedPaths = expectedPaths[:i]
			expectedPaths = append(expectedPaths, upper...)
		}
	}
	t.Run("TestExpectedNotFound", func(t *testing.T) {
		for _, pathStr := range expectedPaths {
			t.Errorf("PathString = %s; not found", pathStr)
		}
	})
	// t.Run("TestFoundNotExpected", func(t *testing.T) {
	// 	for 
	// })
	
}