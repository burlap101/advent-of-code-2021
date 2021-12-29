package main

import (
	"os"
	"testing"
	"fmt"
)

func TestInputDataShape (t *testing.T) {
	paper, err := processFile("input.txt")
	if err != nil {
		t.Error(err)
	}
	length, width := len(paper.points), len(paper.points[0])

	for i, fold := range paper.folds {
		if fold.axis == 'x' {
			width /= 2
			if width != fold.location {
				t.Errorf("Fold #%d location not midpoint x=%d; got x=%d",i+1, width, fold.location)
			}
		} else if fold.axis == 'y' {
			length /= 2
			if length != fold.location {
				t.Errorf("Fold #%d location not midpoint y=%d; got y=%d",i+1, length, fold.location)
			}
		}
	}
}

func TestFolds (t *testing.T) {
	paper, err := processFile("example1.txt")
	if err != nil {
		t.Error(err)
	}
	t.Run("TestFoldOnY", func(t *testing.T) {
		err := paper.PerformNextFold()
		
		if err != nil {
			t.Error(err)
		}
		expectedResult, err := os.ReadFile("firstfold1.txt")
		if err != nil {
			t.Error(err)
		}
		
		if fmt.Sprintf("%v", paper) !=string(expectedResult) {
			t.Error("the result and expected result don't match")
		}
		t.Run("TestDotCount", func(t *testing.T) {
			if paper.DotCount() != 17 {
				t.Errorf("unexpected dot count n=17; got n=%d", paper.DotCount())
			}
		})
	})
	t.Run("TestFoldOnX", func(t *testing.T) {
		err := paper.PerformNextFold()
		
		if err != nil {
			t.Error(err)
		}
		expectedResult, err := os.ReadFile("secondfold1.txt")
		if err != nil {
			t.Error(err)
		}
		
		if fmt.Sprintf("%v", paper) !=string(expectedResult) {
			t.Error("the result and expected result don't match")
		}

	})

}