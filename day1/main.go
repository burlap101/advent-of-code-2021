package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {


	fmt.Println("Advent of code day 1.")
	fmt.Println("Part 1")
	dat, err := os.ReadFile("input.txt")
	
	check(err)
	inputStr := strings.Split(string(dat), "\n")
	inputInts := make([]int, len(inputStr))
	nIncreased := 0
	for i, v := range inputStr {
		if inputInts[i], err = strconv.Atoi(v); err != nil {
			check(err)
		}
		if i > 0 && (inputInts[i] - inputInts[i-1]) > 0 {
			nIncreased += 1
		}
	}

	fmt.Println("No. of increases:", nIncreased)

	fmt.Println("\n\nPart 2:")

	nWinIncreased := 0
	for i := range inputInts {
		if i + 3 >= len(inputInts) {
			break
		}
		if inputInts[i] + inputInts[i+1] + inputInts[i+2] < inputInts[i+1] + inputInts[i+2] + inputInts[i+3] {
			nWinIncreased += 1
		}
	}

	fmt.Println("No. of window increases:", nWinIncreased)
	
}