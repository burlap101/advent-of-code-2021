package main

import "testing"


func TestProcessFile(t *testing.T) {
	expectedString := "38006F45291200"

	if result := ProcessFile("example1.txt"); result != expectedString {
		t.Errorf("Hex string = '%s'; expected '%s'", result, expectedString)
	}
}