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
	pztor.LetterCounts = CountLetterOccurrences(pztor.Polymer) // init counts
	pztor.PairCounts = pztor.CountPairs(pztor.Polymer)

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

		t.Run(fmt.Sprintf("TestPairCountInsertions%d", i+1), func(t *testing.T) {
			err = pztor.PerformPairCountInsertions()
			if err != nil {
				t.Error(err)
			}
			expectedLetterCounts := CountLetterOccurrences(pztor.Polymer)
			for letter, count := range expectedLetterCounts {
				if pztor.LetterCounts[letter] != count {
					t.Errorf("incorrect count for '%c' %d; expected %d", letter, pztor.LetterCounts[letter], count)
				}
			}
		})
	}
}

func TestCountPairs(t *testing.T) {
	pztor, err := processFile("example.txt")
	if err != nil {
		t.Error(err)
	}
	pairCounts := pztor.CountPairs(pztor.Polymer)
	expectedPairCounts := map[string]uint64 {
		"CH" : 0,
		"HH" : 0,
		"CB" : 1,
		"NH" : 0,
		"HB" : 0,
		"HC" : 0,
		"HN" : 0,
		"NN" : 1,
		"BH" : 0,
		"NC" : 1,
		"NB" : 0,
		"BN" : 0,
		"BB" : 0,
		"BC" : 0,
		"CC" : 0,
		"CN" : 0,
	}
	for pair, count := range expectedPairCounts {
		t.Run(fmt.Sprintf("Pair '%s'", pair), func(t *testing.T) {
			if pairCounts[pair] != count {
				t.Errorf(fmt.Sprintf("incorrect count %d; expected %d", pairCounts[pair], count))
			}
		})
	}
}

func TestFinalAnswerPart2(t *testing.T) {
	pztor, err := processFile("example.txt")
	if err != nil {
		t.Error(err)
	}
	pztor.PairCounts = pztor.CountPairs(pztor.Polymer)
	pztor.LetterCounts = CountLetterOccurrences(pztor.Polymer)
	iterations := 40
	for i := 0; i < iterations; i++ {
		err := pztor.PerformPairCountInsertions()
		if err != nil {
			t.Error(err)
		}
	}
	diff := FindMaxMinDifference(pztor.LetterCounts)
	expectedDiff := uint64(2188189693529)
	if diff != expectedDiff{
		t.Errorf("incorrect answer %d; expected %d", diff, expectedDiff)
	}
}