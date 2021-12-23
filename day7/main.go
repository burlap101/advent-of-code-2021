package main

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"os"
	"strings"
	"strconv"
	"math"
)

func populatePositions(dat []byte) ([]float64, error) {
	positions := make([]float64, 0)
	for _, spos := range strings.Split(string(dat), ",") {
		ipos, err := strconv.Atoi(spos)
		if err != nil {
			panic(err)
		}
		positions = append(positions, float64(ipos))
	}
	return positions, nil
}

func factorial(num float64) (float64) {
	if num <= 1 {
		return num
	} else {
		return num+factorial(num-1)
	}
}
func main() {
	fmt.Println("Advent of code day 7")
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	positions, err := populatePositions(dat)
	if err != nil {
		panic(err)
	}
	mu, err := stats.Mean(positions)
	if err != nil {
		panic(err)
	}
	fmt.Println("Mean:", mu)
	median, err := stats.Median(positions)
	if err != nil {
		panic(err)
	}
	fmt.Println("Median:", median)
	totalFuel := float64(0)
	for _, fpos := range positions {
		totalFuel += math.Abs(fpos - median)
	}
	fmt.Println("Total fuel (Part1):", totalFuel)

	mu = float64(488)

	totalFuel = 0
	for _, fpos := range positions {
		totalFuel += factorial(math.Abs(fpos - mu))
	}
	fmt.Printf("Total fuel (Part2): %d\n", int64(totalFuel))
}