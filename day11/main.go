package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

type mcoord struct {
	m, n int
}

type cavern struct {
	octopiMap [][]int
	flashCount uint64
	flashedBuffer map[mcoord]struct{}
}

func (c cavern) String() string {
	result := ""
	for _, row := range c.octopiMap {
		result += fmt.Sprintln(row)
	}
	return result
}

func (c *cavern) flashOctopus(point mcoord) error {
	c.flashedBuffer[point] = struct{}{}
	c.flashCount++
	c.octopiMap[point.m][point.n] = 0
	// fmt.Println("FlashPoint", point)
	incrementStack := make([]mcoord, 0)
	startM, startN, endM, endN := point.m, point.n, point.m, point.n
	if point.m-1 >= 0 {
		startM = point.m-1
	}
	if point.n-1 >= 0 {
		startN = point.n-1
	}
	if point.m+1 < len(c.octopiMap) {
		endM = point.m+1
	}
	if point.n+1 < len(c.octopiMap[point.m]) {
		endN = point.n+1
	}

	for m := startM; m <= endM; m++ {
		for n:= startN; n <= endN; n++ {
			_, pointFlashed := c.flashedBuffer[mcoord{m, n}]
			if !(m == point.m && n == point.n) && !pointFlashed {
				incrementStack = append(incrementStack, mcoord{m, n})
			}
		}
	}
	// fmt.Println("stack contents", incrementStack)
	// fmt.Println("flashes", c.flashCount)

	
	for _, octopus := range incrementStack {
		c.incrementOctopus(octopus)
	}
	
	return nil
}

func (c *cavern) incrementOctopus(point mcoord) error {
	_, pointFlashed := c.flashedBuffer[point]
	if pointFlashed {
		return fmt.Errorf("octopus %+v already flashed so shouldn't increment", point)
	}
	c.octopiMap[point.m][point.n]++
	if c.octopiMap[point.m][point.n] > 9 {
		c.flashOctopus(point)
	} 
	return nil
}

func createOctopiMap(filename string) ([][]int, error) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		return make([][]int, 0), err
	}
	octopiMap := make([][]int, len(strings.Split(string(dat), "\n")))
	for m, row := range strings.Split(string(dat), "\n") {
		octopiMap[m] = make([]int, len(row))
		for n, r := range row {
			val, err := strconv.Atoi(string(r))
			if err != nil {
				panic(err)
			}
			octopiMap[m][n] = val
		}
	}
	return octopiMap, nil
}

func (c *cavern) isSynced() bool {
	for m := 0; m < len(c.octopiMap); m++ {
		for n := 0; n < len(c.octopiMap[m]); n++ {
			if c.octopiMap[m][n] != 0 {
				return false
			}
		}
	}
	return true
}

func part1(filename string) error {
	octopiMap, err := createOctopiMap(filename)
	if err != nil {
		panic(err)
	}
	
	cav := cavern{octopiMap, uint64(0), make(map[mcoord]struct{})}

	iterationLimit := 100
	for i := 0; i < iterationLimit; i++ {
		for m := 0; m < len(cav.octopiMap); m++ {
			for n:=0; n < len(cav.octopiMap[m]); n++ {
				cav.incrementOctopus(mcoord{m, n})
			}
		}
		cav.flashedBuffer = make(map[mcoord]struct{})
	}
	fmt.Println("Part 1 answer:", cav.flashCount)
	return nil
}

func part2(filename string) error {
	octopiMap, err := createOctopiMap(filename)
	if err != nil {
		panic(err)
	}
	
	cav := cavern{octopiMap, uint64(0), make(map[mcoord]struct{})}
	i := 0
	for i = 0; !cav.isSynced(); i++ {
		for m := 0; m < len(cav.octopiMap); m++ {
			for n:=0; n < len(cav.octopiMap[m]); n++ {
				cav.incrementOctopus(mcoord{m, n})
			}
		}
		cav.flashedBuffer = make(map[mcoord]struct{})
	}
	fmt.Println("Part 2 answer:", i)
	return nil
}

func main() {
	fmt.Println("Advent of code day 11")
	filename := "input.txt"
	err := part1(filename)
	if err != nil {
		panic(err)
	}
	
	err = part2(filename)
	if err != nil {
		panic(err)
	}

	

}