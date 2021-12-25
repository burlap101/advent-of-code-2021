package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"errors"
	"github.com/burlap101/advent-of-code-2021/day8/segmentcodec"
)

/* Each index corresponds to length of segment count while its value is the total count */
func digitLengthCounts(dat string) ([]int) {
	counts := make([]int, 8)
	for _, digit := range strings.Fields(dat) {
		if digit == "|" {
			_ = 0
		} else {
			counts[len(digit)] += 1
		}
	}
	return counts
}

/* Returns all unique digits found in a given string */
func uniqueDigits(str string) ([]string, error) {
	digitStrs := make([]string, 0)

	for _, digitStr := range strings.Fields(str) {
		if !alreadyFound(digitStrs, digitStr) {
			digitStrs = append(digitStrs, digitStr)
		}
		if len(digitStrs) == 10 {
			return digitStrs, nil
		}
	}
	
	return make([]string, 0), errors.New("there was a problem finding unique digits")
}

func alreadyFound(digitStrs []string, str string) bool {
	for _, ds := range digitStrs {
		if segmentcodec.DigitStringMatch(ds, str) {
			return true
		}
	}
	return false
}

func main() {
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	sdat := strings.ReplaceAll(string(dat), " |\n", " | ")
	sdatLines := strings.Split(sdat, "\n")

	outputConcat := ""
	mappedLines := make(map[string]string)
	for _, line := range sdatLines {
		input, output := strings.Split(line, "|")[0], strings.Split(line, "|")[1]
		outputConcat += output + " "
		mappedLines[strings.TrimSpace(input)] = strings.TrimSpace(output)
	}
	
	outputCounts := digitLengthCounts(outputConcat)
	
	fmt.Println(outputCounts)
	p1Counts := outputCounts[2:5]
	p1Counts = append(p1Counts, outputCounts[7])
	fmt.Println(p1Counts)
	totalP1Count := 0
	for _, count := range p1Counts {
		totalP1Count += count
	}

	
	fmt.Println("Part 1 answer:", totalP1Count)
	// outputNums := make([]int, 0)
	// encodedDigitStrs, err := uniqueDigits(outputConcat)
	// if err != nil {
	// 	panic(err)
	// }
	// var codec segmentcodec.Converter
	total := 0
	for input, output := range mappedLines {
		var codec segmentcodec.Converter
		encodedDigitStrs := strings.Fields(input)
		codec.PopulateConverter(encodedDigitStrs)
		decodedOutputStr := ""
		for _, encodedDigit := range strings.Fields(output) {
			digit, err := codec.ConvertToDigit(encodedDigit)
			if err != nil {
				panic(err)
			}
			decodedOutputStr += string(digit)
		}
		fmt.Println(decodedOutputStr)
		val, err := strconv.Atoi(decodedOutputStr)
		if err != nil {
			panic(err)
		}
		total += val
	}
	fmt.Println("Part2 Sum:", total)
	
}