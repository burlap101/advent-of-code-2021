package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"github.com/burlap101/advent-of-code-2021/day15/search"
)

func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func ProcessFile(filename string) (search.RiskMap, error){
	dat, err := os.ReadFile(filename)
	ErrorHandler(err)
	result := make(search.RiskMap, 0)
	for m, row := range strings.Split(string(dat), "\n") {
		result = append(result, make([]uint8, 0))
		for _, val := range row {
			riskLevel, err := strconv.Atoi(string(val))
			ErrorHandler(err)
			result[m] = append(result[m], uint8(riskLevel))
		}
	}
	return result, nil
}

func main() {
	fmt.Println("Advent of code day 15") 
	riskMatrix, err := ProcessFile("input.txt")
	ErrorHandler(err)
	p := search.Problem{StateSpace: riskMatrix}
	pathCost, err := p.Search()
	ErrorHandler(err)
	fmt.Println("\nPart 1 Answer:", pathCost)
	pathCost, _, err = p.SearchPart2()
	ErrorHandler(err)
	fmt.Println("\nPart 2 Answer:", pathCost)
}