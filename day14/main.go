package main

import (
	"fmt"
	"os"
	"strings"
)

type Polymerizator struct {
	Polymer string
	InsertionRules map[string]byte
	LetterCounts map[byte]uint64
	PairCounts map[string]uint64
}

func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func processFile(filename string) (Polymerizator, error) {
	pztor := Polymerizator{"", make(map[string]byte), make(map[byte]uint64), make(map[string]uint64)}
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

func (p Polymerizator) CountPairs(polymer string) map[string]uint64 {
	pairCounts := make(map[string]uint64)
	for pair := range p.InsertionRules {
		pairCounts[pair] = 0
	}
	for i := 0; i + 1 < len(polymer); i++ {
		pairCounts[polymer[i:i+2]] += 1
	}
	return pairCounts
}

func (p Polymerizator) String() (string) {
	result := fmt.Sprintf("Polymer:\n%s\n\n", p.Polymer)
	result += "Insertion Rules:\n"
	for key, val := range p.InsertionRules {
		result += fmt.Sprintf("%s -> %c\n", key, val)
	}
	return result
}

func CountLetterOccurrences(polymer string) (map[byte]uint64) {
	result := make(map[byte]uint64)
	for _, el := range polymer {
		if _, ok := result[byte(el)]; !ok {
			result[byte(el)] = 0
		}
		result[byte(el)] += 1
	}
	return result
}

func FindMaxMinDifference(letterCounts map[byte]uint64) uint64 {
	var max, min uint64 = 0, 0
	for _, count := range letterCounts {
		if count > max {
			max = count
		}
		if count < min || min == 0 {
			min = count
		}
	}
	return max-min
}

func (pztor Polymerizator) part1(iterations int) Polymerizator {
	for i := 0; i < iterations; i++ {
		err := error(nil)
		pztor.Polymer, err = pztor.PerformInsertions()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Progress: %d of %d\r", i+1, iterations)
	}
	fmt.Println()
	letterOccurrence := CountLetterOccurrences(pztor.Polymer)
	
	fmt.Println("Part 1 Answer:", FindMaxMinDifference(letterOccurrence))
	return pztor
}

func (p *Polymerizator) PerformPairCountInsertions() (error) {
	resultingPairs := make(map[string]uint64)
	for pair, count := range p.PairCounts {
		insertionLetter := p.InsertionRules[pair]
		p.LetterCounts[insertionLetter] += count
		rp := []string{fmt.Sprintf("%c%c", pair[0], insertionLetter), fmt.Sprintf("%c%c", insertionLetter, pair[1])}
		if _, ok := resultingPairs[rp[0]]; !ok {
			resultingPairs[rp[0]] = 0
		}
		resultingPairs[rp[0]] += count
		if _, ok := resultingPairs[rp[1]]; !ok {
			resultingPairs[rp[1]] = 0
		}
		resultingPairs[rp[1]] += count
		if _, ok := resultingPairs[pair]; !ok {
			resultingPairs[pair] = 0
		}
		resultingPairs[pair] -= count
	}
	for pair, count := range resultingPairs {
		p.PairCounts[pair] += count
	}

	return nil
}

func (p Polymerizator) part2(iterations int) {
	p.PairCounts = p.CountPairs(p.Polymer)
	p.LetterCounts = CountLetterOccurrences(p.Polymer)
	fmt.Println("Part 2")
	for i := 0; i < iterations; i++ {
		err := p.PerformPairCountInsertions()
		ErrorHandler(err)
		fmt.Printf("Progress: %d of %d\r", i+1, iterations)
	}
	fmt.Println()
	fmt.Println("Part 2 Answer:", FindMaxMinDifference(p.LetterCounts))
}

func main() {
	fmt.Println("Advent of code day 14")
	pztor, err := processFile("input.txt")
	if err != nil {
		panic(err)
	}
	pztor.part1(10)
	pztor.part2(40)
}