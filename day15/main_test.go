package main

import (
	"testing"
)


func TestProcessFile(t *testing.T) {
	result, err := ProcessFile("example.txt")
	if err != nil {
		t.Error(err)
	}
	t.Run("Matrix Sizes Correct", func(t *testing.T) {
		if len(result) != 10 {
			t.Errorf("Number of rows = %d; expected 10", len(result))
		}
		if len(result[0]) != 10 {
			t.Errorf("Number of columns = %d; expected 10", len(result[0]))
		}
	})
}

