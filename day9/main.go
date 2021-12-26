package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"sort"
)

type matrixCoord struct{
	m, n int
}

type heightMap [][]int

type basin struct {
	lowPoint matrixCoord
	includedPoints []matrixCoord
}

type coordStack struct {
	points []matrixCoord
}

func (cs *coordStack) Push(val matrixCoord) int {
	cs.points = append(cs.points, val)
	return len(cs.points)
}

func (cs *coordStack) Pop()  matrixCoord {
	val := cs.points[len(cs.points)-1]
	cs.points = cs.points[:len(cs.points)-1]
	return val
}

func (hm heightMap) findLowSpots() ([]matrixCoord, error) {
	found := make([]matrixCoord, 0)
	for m, row := range hm {
		for n := range row {
			if hm.testForLowSpot(matrixCoord{m, n}) {
				found = append(found, matrixCoord{m, n})
			}
		}
	}
	return found, nil
}

func (hm heightMap) testForLowSpot(coord matrixCoord) bool {
	n := coord.n
	m := coord.m
	testVal := hm[m][n]
	if n - 1 >= 0 && hm[m][n-1] <= testVal {
		return false
	}
	if m - 1 >= 0 && hm[m-1][n] <= testVal {
		return false
	}
	if n + 1 < len(hm[m]) && hm[m][n+1] <= testVal {
		return false
	}
	if m + 1 < len(hm) && hm[m+1][n] <= testVal {
		return false
	}
	return true
}

func (b *basin) findPoints(hm heightMap, startPoint matrixCoord) (error) {
	stack := coordStack{make([]matrixCoord, 0)}
	stack.Push(startPoint)
	
	for len(stack.points) > 0 {
		point := stack.Pop()
		m, n := point.m, point.n
		posM, posN, negM, negN := false, false, false, false
		// positive n direction
		if n+1 < len(hm[m]) && hm[m][n+1] >= hm[m][n] && hm[m][n+1] < 9 {
			fmt.Println("Adding point:", matrixCoord{m, n+1})
			b.addPoint(matrixCoord{m, n+1})
			posN = true
		}
		// positive m direction
		if m+1 < len(hm) && hm[m+1][n] >= hm[m][n] && hm[m+1][n] < 9 {
			fmt.Println("Adding point:", matrixCoord{m+1, n})
			b.addPoint(matrixCoord{m+1, n})
			posM = true
		}

		// negative n direction
		if n-1 >= 0 && hm[m][n-1] > hm[m][n] && hm[m][n-1] < 9 {
			fmt.Println("Adding point:", matrixCoord{m, n-1})
			b.addPoint(matrixCoord{m, n-1})
			negN = true
		}
		// negative m direction
		if m-1 >= 0 && hm[m-1][n] > hm[m][n] && hm[m-1][n] < 9 {
			fmt.Println("Adding point:", matrixCoord{m-1, n})
			b.addPoint(matrixCoord{m-1, n})
			negM = true
		}

		if negM {
			fmt.Println("Progressing to", matrixCoord{m-1, n}, "from", point)
			stack.Push(matrixCoord{m-1, n})
		}

		if negN {
			fmt.Println("Progressing to", matrixCoord{m, n-1}, "from", point)
			stack.Push(matrixCoord{m, n-1})
		}

		if posM {
			fmt.Println("Progressing to", matrixCoord{m+1, n}, "from", point)
			stack.Push(matrixCoord{m+1, n})
		}
		
		if posN {
			fmt.Println("Progressing to", matrixCoord{m, n+1}, "from", point)
			stack.Push(matrixCoord{m, n+1})
		}
	}

	return nil

}

func (b *basin) addPoint(point matrixCoord) (error) {
	if !b.pointExists(point) {
		b.includedPoints = append(b.includedPoints, point)
	}
	return nil
}

func (b *basin) pointExists(point matrixCoord) bool {
	for _, includedPoint := range b.includedPoints {
		if includedPoint == point {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Advent of code day 9")
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	
	hm := make(heightMap, 0)
	for _, rowStr := range strings.Split(string(dat), "\n") {
		row := make([]int, 0)
		for _, numStr := range rowStr {
			num, err := strconv.Atoi(string(numStr))
			if err != nil {
				panic(err)
			}
			row = append(row, num)
		}
		hm = append(hm, row)
	}

	lowSpots, err := hm.findLowSpots()
	if err != nil {
		panic(err)
	}

	lowSpotsRiskSum := 0
	basins := make([]basin, len(lowSpots))
	for i, lowSpot := range lowSpots {
		lowSpotsRiskSum += hm[lowSpot.m][lowSpot.n] + 1
		basins[i] = basin{lowPoint: lowSpot, includedPoints: make([]matrixCoord, 0)}
		basins[i].includedPoints = append(basins[i].includedPoints, lowSpot)
	}

	fmt.Println("Low spots found:", len(lowSpots))
	fmt.Println("Part1 Answer:", lowSpotsRiskSum)

	sizes := make([]int, len(lowSpots))

	for i, basin := range basins {
		basin.findPoints(hm, basin.lowPoint)
		fmt.Printf("Basin %d size: %d\n", i, len(basin.includedPoints))
		fmt.Println("Included points:", basin.includedPoints)
		sizes[i] = len(basin.includedPoints)
	}

	sort.Ints(sizes)
	biggestSizes := sizes[len(sizes)-3:]

	fmt.Println("3 biggest basin sizes:", biggestSizes)
	fmt.Println("Part2 answer:", uint64(biggestSizes[0])*uint64(biggestSizes[1])*uint64(biggestSizes[2]))
	
}