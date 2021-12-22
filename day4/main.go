package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"errors"
)

type board struct {
	str string
	values [][]int
}

func (b *board) winnerCheck() int {
	
	for n := 0; n < len(b.values[0]); n++ {
		colSum := 0
		for m := 0; m < len(b.values); m++ {
			colSum += b.values[m][n]
		}
		if colSum == -1 * len(b.values) {
			return 0
		}
	}
	for m := 0; m < len(b.values); m++ {
		rowSum := 0
		for n := 0; n < len(b.values[0]); n++ {
			rowSum += b.values[m][n]
		}
		if rowSum == -1 * len(b.values[0]) {
			return 0
		}
	}
	
	return -100
}

func (b *board) populateValues() error {
	rows := strings.Split(string(b.str), "\n")
	b.values = make([][]int, len(rows))
	for m, row := range rows {
		vals := strings.Fields(row)
		b.values[m] = make([]int, len(vals))
		for n, val := range vals {
			if num, err := strconv.Atoi(val); err == nil {
				b.values[m][n] = num
			} else {
				panic(err)
			}
		}
	}
	return nil
}

func (b *board) checkMarkDrawn(num int) (int, error) {
	numFound := 0
	for m, row := range b.values {
		for n, val := range row {
			if val == num {
				b.values[m][n] = -1
				numFound += 1
			}
		}
	}

	if win := b.winnerCheck(); win == 0 {
		return 0, errors.New("Bingo!")
	}

	return numFound, nil
}

func (b *board) sumRemainingValues() (int) {
	sum := 0
	for _, row := range b.values {
		for _, val := range row {
			if val >= 0 {
				sum += val
			}
		}
	}
	return sum
}

func removeBoard(boards []board, index int) ([]board) {
	return append(boards[:index], boards[index+1:]...)
}

func main() {
	fmt.Println("Advent of code day 4")
	fmt.Println("Part1:")
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var drawStr string
	var boards []board

	for i, v := range strings.Split(string(dat), "\n\n") {
		if i == 0 {
			drawStr = v
		} else {
			boards = append(boards, board{v, nil})
			err := boards[i-1].populateValues()
			if err != nil {
				panic(err)
			}
		}
	}
	
	fmt.Print("Numbers drawn: ")
	drawNums := strings.Split(drawStr, ",")
	for _, v := range drawNums {
		num, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d ", num)
		winningBoardIndices := make([]int, 0)
		for i, b := range boards {
			if _, err := b.checkMarkDrawn(num); err != nil {
				fmt.Println(err)
				fmt.Println("Winning number:", num)
				fmt.Println("Winning card:", i)
				fmt.Println(b.str)
				remainSum := b.sumRemainingValues()
				fmt.Println("Sum remaining values:", remainSum)
				fmt.Println("Final answer:", num*remainSum)
				winningBoardIndices = append(winningBoardIndices, i)
			}
		}
		for i := len(winningBoardIndices)-1; i >= 0; i-- {
			boards = removeBoard(boards, winningBoardIndices[i])
		}
	}

}