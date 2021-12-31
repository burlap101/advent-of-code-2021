package main

import (
	"fmt"
	"os"
	"strings"
)

type Polymerizator struct {
	Polymer string
	InsertionRules map[string]byte
}

func processFile(filename string) (Polymerizator, error) {
	pztor := Polymerizator{"", make(map[string]byte)}
	dat, err := os.ReadFile(filename)
	if err != nil {
		return pztor, err
	}
	manualSections := strings.Split(string(dat), "\n\n")
	pztor.Polymer = strings.TrimSpace(manualSections[0])
	for _, line := range strings.Split(manualSections[1], "\n") {
		ruleSplit := strings.Split(line, "->")
		pztor.InsertionRules[strings.TrimSpace(ruleSplit[0])] = strings.TrimSpace(ruleSplit[1])[0]
	}
	return pztor, nil
}

func (p Polymerizator) PerformInsertions() (string, error) {
	result := p.Polymer
	insertions := make([]byte, 0)
	
	for i := 0; i < len(p.Polymer)-1; i++ {
		actualPair := p.Polymer[i:i+2]
		insertion, ok := p.InsertionRules[string(actualPair)]
		if !ok {
			return "", fmt.Errorf("no insertion rule found; pair='%s'", string(actualPair))
		}
		insertions = append(insertions, insertion)
	}

	for i := len(p.Polymer)-1; len(insertions) > 0; i-- {
		upper := result[i:]
		result = result[:i] + string(insertions[len(insertions)-1]) + upper
		insertions = insertions[:len(insertions)-1]
	}

	return result, nil
}

func (p Polymerizator) String() (string) {
	result := fmt.Sprintf("Polymer:\n%s\n\n", p.Polymer)
	result += "Insertion Rules:\n"
	for key, val := range p.InsertionRules {
		result += fmt.Sprintf("%s -> %c\n", key, val)
	}
	return result
}

func CountLetterOccurrences(polymer string) (map[rune]int) {
	result := make(map[rune]int)
	for _, el := range polymer {
		if _, ok := result[el]; !ok {
			result[el] = 0
		}
		result[el] += 1
	}
	return result
}

func main() {
	fmt.Println("Advent of code day 14")

	pztor, err := processFile("input.txt")
	if err != nil {
		panic(err)
	}
	iterations := 10
	for i := 0; i < iterations; i++ {
		pztor.Polymer, err = pztor.PerformInsertions()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Progress: %d of %d\r", i, iterations)
	}
	fmt.Println()
	max, min := 0, 0
	letterOccurrence := CountLetterOccurrences(pztor.Polymer)
	for el, count := range letterOccurrence {
		fmt.Println(string(el), count)
		if count > max {
			max = count
		}
		if count < min || min == 0 {
			min = count
		}
	}
	fmt.Println("Part 1 Answer:", max-min)
}