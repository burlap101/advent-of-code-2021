package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func oxygenFilter(vals []string, bitNum int) []string {
	var result []string
	accum := 0
	for _, val := range vals {
		dig, err := strconv.Atoi(string(val[bitNum]))
		if err != nil {
			panic(err)
		}
		accum += dig
	}
	
	ratio := float64(accum)/float64(len(vals))
	// fmt.Println("ratio:", ratio)
	for _, val := range vals {
		if ratio >= 0.5 && val[bitNum] == '1' {
			result = append(result, val)
		} else if ratio < 0.5 && val[bitNum] == '0' {
			result = append(result, val)
		}
	}

	// fmt.Println("Result:", result)
	return result
}

func co2Filter(vals []string, bitNum int) []string {
	var result []string
	accum := 0
	for _, val := range vals {
		dig, err := strconv.Atoi(string(val[bitNum]))
		if err != nil {
			panic(err)
		}
		accum += dig
	}

	ratio := float64(accum)/float64(len(vals))

	for _, val := range vals {
		if ratio >= 0.5 && val[bitNum] == '0' {
			result = append(result, val)
		} else if ratio < 0.5 && val[bitNum] == '1' {
			result = append(result, val)
		}
	}
	// fmt.Println("co2 filter result:", result, ", for bit no.:", bitNum)
	return result

}

func main() {
	fmt.Println("Advent of code day 3")
	fmt.Println("Part1:")

	dat, err := os.ReadFile("input.txt")

	if err != nil {
		panic(err)
	}

	var snums []string
	for _, s := range strings.Split(string(dat), "\n") {
		snums = append(snums, strings.TrimSpace(s))
	}
	
	cols := make([]int, len(snums[0]))
	for _, s := range snums {
		for i, ch := range s {
			chInt, err := strconv.Atoi(string(ch))
			if err != nil {
				panic(err)
			}
			cols[i] += chInt
		}
	}
	// fmt.Println(cols)

	gammaStr := ""
	epsilonStr := ""

	for _, v := range cols {
		ratio := float64(v)/float64(len(snums))
		if ratio >= 0.5 {
			gammaStr += "1"
			epsilonStr += "0"
		} else {
			gammaStr += "0"
			epsilonStr += "1"
		}
	}


	fmt.Println(gammaStr, epsilonStr)
	gamma, err := strconv.ParseInt(gammaStr, 2, 0)
	if err != nil {
		panic(err)
	}
	epsilon, err := strconv.ParseInt(epsilonStr, 2, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("gamma:", gamma, ", epsilon:", epsilon)
	fmt.Println("Power Consumption:", gamma*epsilon)
	// fmt.Println(gammaBin, epsilonBin)

	fmt.Println("\nPart2:")
	
	oxygenStr := ""
	oxygenCandidates := snums[:]
	for i:=0; len(oxygenCandidates) > 1; i++ {
		oxygenCandidates = oxygenFilter(oxygenCandidates, i)
	}
	oxygenStr = oxygenCandidates[0]
	oxygenVal, err := strconv.ParseInt(oxygenStr, 2, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println("Oxygen val:", oxygenVal)
	co2Str := ""
	co2Candidates := snums[:]
	for i := 0; len(co2Candidates) > 1; i++ {
		co2Candidates = co2Filter(co2Candidates, i)
	}
	co2Str = co2Candidates[0]
	co2Val, err := strconv.ParseInt(co2Str, 2, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println("CO2 val:", co2Val)
	fmt.Println("Life support rating:", oxygenVal*co2Val)

}