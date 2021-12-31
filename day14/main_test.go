package main

import (
	"os"
	"strings"
	"testing"
	"fmt"
	"regexp"
)

func TestSteps(t *testing.T) {
	pztor, err := processFile("example.txt")
	if err != nil {
		t.Error(err)
	}

	dat, err := os.ReadFile("exampleresults.txt")
	if err != nil {
		t.Error(err)
	}

	results := make([]string, 0)
	r, err := regexp.Compile("[1-4]:")
	if err != nil {
		t.Error(err)
	}
	for _, line := range strings.Split(string(dat), "\n") {
		if r.MatchString(line) {
			results = append(results, strings.TrimSpace(strings.Split(line, ":")[1]))
		}
	}

	for i, result := range results {
		pztor.Polymer, err = pztor.PerformInsertions()
		if err != nil {
			t.Error(err)
		}
		t.Run(fmt.Sprintf("TestStep%d", i+1), func(t *testing.T) {
			if result != pztor.Polymer {
				t.Errorf("incorrect polymer.\nreturned \n%s\nexpected \n%s", pztor.Polymer, result)
			}
			for j, el := range pztor.Polymer {
				if el != rune(result[j]) {
					t.Errorf("incorrect element %d of %d, '%c'; expected '%c'", j+1, len(result), el, result[j])
					break
				}
			}
		})
	}

}