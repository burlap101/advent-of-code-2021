package main

import (
	"fmt"
	"os"
	"strings"
	"sort"
)

type bracket struct {
	open rune
	close rune
}

var BRACKETS = []bracket{
	{'(', ')'},
	{'[', ']'},
	{'<', '>'},
	{'{', '}'},
}

var ERROR_SCORES = map[rune]int{
	')':3,
	']':57,
	'}':1197,
	'>':25137,
}

var COMPLETION_SCORES = map[rune]uint64{
	')':1,
	']':2,
	'}':3,
	'>':4,
}

type sourceFile struct {
	lines []string
	bracketStack []bracket
	completions []string
}

func (sf *sourceFile) Push(val bracket) int {
	sf.bracketStack = append(sf.bracketStack, val)
	return len(sf.bracketStack)
}

func (sf *sourceFile) Pop() bracket {
	val := sf.bracketStack[len(sf.bracketStack)-1]
	sf.bracketStack = sf.bracketStack[:len(sf.bracketStack)-1]
	return val
}

func (sf *sourceFile) checkLine(i int) (int, error) {
	sf.bracketStack = make([]bracket, 0)
	for _, r := range sf.lines[i] {
		for _, b := range BRACKETS {
			if r == b.open {
				sf.Push(b)
				break
			} else if r == b.close {
				sb := sf.Pop()
				if sb.close == r {
					break
				} else {
					return ERROR_SCORES[r], fmt.Errorf("expected '%c' but found '%c' instead", sb.close, r)
				}
			}
		}
	}
	return 0, nil
}

func (sf *sourceFile) removeLine(i int) error {
	upper := sf.lines[i+1:]
	sf.lines = sf.lines[:i]
	sf.lines = append(sf.lines, upper...)
	return nil
}

func (sf *sourceFile) completeLine(i int) string {
	result := ""
	for len(sf.bracketStack) > 0 {
		r := sf.Pop()
		result += string(r.close)
	}
	return result
} 

func completionScore(completion string) uint64 {
	result := uint64(0)
	for _, r := range completion {
		result =  5*result + COMPLETION_SCORES[r]
	}
	return result
}


func main() {
	fmt.Println("Advent of code day 10")
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	sf := sourceFile{make([]string, 0), make([]bracket, 0), make([]string, 0)}
	
	sf.lines = append(sf.lines, strings.Split(string(dat), "\n")...)
	totalScore := uint64(0)
	for i := range sf.lines {
		score, err := sf.checkLine(i)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Score:", score)
			totalScore += uint64(score)
		} else {
			sf.completions = append(sf.completions, sf.completeLine(i))
		}
	}
	fmt.Println("Part1 Answer:", totalScore)

	completionScores := make([]uint64, 0)
	for _, completion := range sf.completions {
		cs := completionScore(completion)
		fmt.Println("Completion:", completion, "Score:", cs)
		completionScores = append(completionScores, cs)
	}

	sort.SliceStable(completionScores, func(i, j int) bool { 
		return completionScores[i] < completionScores[j] 
	})
	
	fmt.Println("Part2 Answer:", completionScores[len(completionScores)/2])
	
}