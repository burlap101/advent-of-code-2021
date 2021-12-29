package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Paper struct {
	points [][]int
	folds []Fold
}

type Fold struct {
	location int
	axis rune
}

func (p *Paper) Fold (lineLocation int, axis rune) (error) {
	return nil
}

func (p *Paper) DotCount() (int) {
	result := 0
	for _, row := range p.points {
		for _, val := range row {
			result += val
		}
	}
	return result
}

func processDotString(dotString string) ([][]int, error) {
	points := make([][]int, 0)
	for _, coords := range strings.Split(dotString, "\n") {
		splitted := strings.Split(coords, ",")
		x, err := strconv.Atoi(splitted[0])
		if err != nil {
			return make([][]int, 0), err
		}
		y, err := strconv.Atoi(splitted[1]) 
		if err != nil {
			return make([][]int, 0), err
		}
		if y >= len(points) {
			for m := len(points); m <= y; m++ {
				if len(points) == 0 {
					points = append(points, make([]int, 0))
				} else {
					points = append(points, make([]int, len(points[0])))
				}
				
			}
		}
		if x >= len(points[y]) {
			for m := 0; m < len(points); m++ {
				points[m] = append(points[m], make([]int, x - len(points[m]) + 1)...)
			}
		}
		points[y][x] = 1
	}
	return points, nil
}

func processInstructionString(instructionString string) ([]Fold, error) {
	folds := make([]Fold, 0)
	for _, instruction := range strings.Split(instructionString, "\n") {
		splitted := strings.Fields(instruction)
		location, err := strconv.Atoi(strings.Split(splitted[2], "=")[1])
		if err != nil {
			return make([]Fold, 0), err
		}
		axis := rune(strings.Split(splitted[2], "=")[0][0])
		folds = append(folds, Fold{location, axis})
	}
	return folds, nil
}

func processFile(filename string)  (Paper, error) {
	paper := Paper{make([][]int, 0), []Fold{}}
	dat, err := os.ReadFile(filename)
	if err != nil {
		return paper, err
	}
	splitted := strings.Split(string(dat), "\n\n")
	dotString, instructionString := splitted[0], splitted[1]
	paper.points, err = processDotString(dotString)
	paper.folds, err = processInstructionString(instructionString)

	return paper, nil
}

func (p Paper) String() (string) {
	result := ""
	for _, row := range p.points {
		for _, val := range row {
			if val == 1 {
				result += "#"
			} else {
				result += "."
			}
		}
		result += "\n"
	}
	// result += "FOLDS:\n"
	// for i, fold := range p.folds {
	// 	result += fmt.Sprintf("%d: {location:%d axis:%c}\n", i, fold.location, fold.axis)
	// }
	return result[:len(result)-1]
}

func (p *Paper) PerformNextFold() (error) {
	if len(p.folds) == 0 {
		return errors.New("no more folds to perform")
	}
	fold := p.folds[0]
	p.folds = p.folds[1:]

	if fold.axis == 'x' {
		p.FoldOnX(fold.location)
	} else if fold.axis == 'y' {
		p.FoldOnY(fold.location)
	} else {
		return fmt.Errorf("invalid fold axis; expected 'x' or 'y', got axis='%c'", fold.axis)
	}
	return nil
}

func (p *Paper) FoldOnX(loc int) {
	left := make([][]int, len(p.points))
	right := make([][]int, len(p.points))
	invRight := make([][]int, len(p.points))
	for m := 0; m < len(p.points); m++ {
		left[m], right[m], invRight[m] = make([]int, loc), make([]int, loc), make([]int, loc)
		copy(left[m], p.points[m][:loc])
		copy(right[m], p.points[m][loc+1:])
		copy(invRight[m], right[m])
	}
	

	for n, invn := len(left[0])-1, 0; n >= 0 && invn < len(left[0]); n-- {
		for m := 0; m < len(left); m++ {
			invRight[m][invn] = right[m][n]
		}
		invn++
	}

	for m, row := range p.points {
		p.points[m] = row[:loc]
	}

	
	for m := 0; m < len(left); m++ {
		for n := 0; n < len(left[m]); n++ {
			p.points[m][n] = left[m][n] | invRight[m][n]
		}
	}
}

func (p *Paper) FoldOnY(loc int) {
	upper := make([][]int, loc)
	lower := make([][]int, loc)
	invlower := make([][]int, loc)
	for m := 0; m < loc; m++ {
		upper[m], lower[m], invlower[m] = make([]int, len(p.points[0])), make([]int, len(p.points[0])), make([]int, len(p.points[0]))
	}
	copy(upper, p.points[:loc])
	copy(lower, p.points[loc+1:])
	copy(invlower, lower)
	
	for i, j := len(lower)-1, 0; i >= 0 && j < len(lower); i-- {
		invlower[j] = lower[i]
		j++
	}

	p.points = p.points[:loc]

	for m := 0; m < len(upper); m++ {
		for n := 0; n < len(upper[m]); n++ {
			p.points[m][n] = upper[m][n] | invlower[m][n]
		}
	}
}

func main() {
	fmt.Println("Advent of code day 13")
	paper, err := processFile("input.txt")
	if err != nil {
		panic(err)
	}
	err = paper.PerformNextFold()
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1 Answer (dots visible):", paper.DotCount())

	for i := 0; len(paper.folds) > 0; i++ {
		paper.PerformNextFold()
		
	}
	fmt.Println("Part 2 answer:")
	fmt.Println(paper)
}