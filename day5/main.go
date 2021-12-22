package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"math"
)

type point struct {
	x int
	y int
}

type line struct {
	str string
	start point
	end point
}

func strToPoint(s string) (point, error) {
	coords := strings.Split(s, ",")
	if len(coords) != 2 {
		return point{0,0}, fmt.Errorf("Only 2 coords can be extracted, %d found processing string %s", len(coords), s)
	}
	x, err := strconv.Atoi(strings.TrimSpace(coords[0]))
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(strings.TrimSpace(coords[1]))
	if err != nil {
		panic(err)
	}

	return point{x, y}, nil
}

func (l *line) extractPoints() (error) {
	points := strings.Split(l.str, "->")
	if len(points) != 2 {
		return fmt.Errorf("Only 2 points can be extracted, %d found processing string %s", len(points), l.str)
	}
	start, err := strToPoint(points[0])
	if err != nil {
		panic(err)
	}
	end, err := strToPoint(points[1])
	if err != nil {
		panic(err)
	}
	l.start = start
	l.end = end

	return nil
}

func (l *line) String() string {
	return fmt.Sprint(l.str)
}

func (l *line) gradient() float64 {
	return float64(l.end.x - l.start.x) / float64(l.end.y - l.start.y)
}

func (l *line) testPoint(p point) bool {
	testy := float64(p.x) * l.gradient() + l.yIntercept()
	return testy == float64(p.y)
}

func (l *line) yIntercept() float64 {
	return float64(l.start.y) - (l.gradient() * float64(l.start.x))
}

func (l *line) yGivenX(x int) float64 {
	return l.gradient() * float64(x) + l.yIntercept()
}

type areaMap struct {
	width int
	height int
	coords [][]int
}

func (am *areaMap) createExtendMap() error {
	if len(am.coords) == 0 {
		am.coords = append(am.coords, make([]int, am.width))
	}
	if am.width > len(am.coords[0]) {
		for m := range am.coords {
			am.coords[m] = append(am.coords[m], make([]int, am.width - len(am.coords[m]))...)
		}
	}
	
	if am.height > len(am.coords) {
		prevHeight := len(am.coords) 
		am.coords = append(am.coords, make([][]int, am.height - len(am.coords))...)

		for m := prevHeight; m < am.height; m++ {
			am.coords[m] = make([]int, am.width)
		}
	}

	return nil
}

func (am areaMap) String() string {
	result := make([][]string, am.height)
	for m := range result {
		result[m] = make([]string, am.width)
	}

	for m, row := range am.coords {
		for n, val := range row {
			if val == 0 {
				result[m][n] = "."
			} else {
				result[m][n] = strconv.Itoa(val)
			}
		}
	}
	return fmt.Sprint(result)
}

func main() {
	fmt.Println("Advent of code day 5")
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var lines []line
	var fullArea areaMap

	for i, l := range strings.Split(string(dat), "\n") {
		lines = append(lines, line{l, point{}, point{}})
		err := lines[i].extractPoints()
		if err != nil {
			panic(err)
		}
		if fullArea.height < lines[i].start.y + 1	{
			fullArea.height = lines[i].start.y + 1
			fullArea.createExtendMap()
			// fmt.Println("Map height:", fullArea.height, len(fullArea.coords), lines[i].start.y, " Map width:", fullArea.width, len(fullArea.coords[i]))
		} 
		if fullArea.height < lines[i].end.y + 1{
			fullArea.height = lines[i].end.y + 1
			fullArea.createExtendMap()
			// fmt.Println("Map height:", fullArea.height, " Map width:", fullArea.width)
		}
		if fullArea.width < lines[i].start.x + 1 {
			fullArea.width = lines[i].start.x + 1
			fullArea.createExtendMap()
			// fmt.Println("Map height:", fullArea.height, " Map width:", fullArea.width, lines[i].start.x)
			
		}
		if fullArea.width < lines[i].end.x {
			fullArea.width = lines[i].end.x
			err = fullArea.createExtendMap()
			if err != nil {
				panic(err)
			}
			// fmt.Println("Map height:", fullArea.height, " Map width:", fullArea.width)
		}
	}
	overlaps := 0
	for _, l := range lines {
		if l.start.x - l.end.x == 0 {
			large := int(math.Max(float64(l.start.y), float64(l.end.y)))
			small := int(math.Min(float64(l.start.y), float64(l.end.y)))

			for m := small; m <= large; m++ {
				fullArea.coords[m][l.start.x] += 1
				if fullArea.coords[m][l.start.x] == 2 {
					overlaps += 1
				}
			}
		} else if l.start.y - l.end.y == 0 {
			large := int(math.Max(float64(l.start.x), float64(l.end.x)))
			small := int(math.Min(float64(l.start.x), float64(l.end.x)))

			for n := small; n <= large; n++ {
				fullArea.coords[l.start.y][n] += 1
				if fullArea.coords[l.start.y][n] == 2 {
					overlaps += 1
				}
			}
		} else if gradient := l.gradient(); math.Abs(float64(gradient)) == 1 {
			largex := int(math.Max(float64(l.start.x), float64(l.end.x)))
			smallx := int(math.Min(float64(l.start.x), float64(l.end.x)))

			for n := smallx; n <= largex; n++ {
				m := int(l.yGivenX(n))
				if !l.testPoint(point{n, m}) {
					fmt.Printf("Weird point doesn't fit: (%d, %d)\n", n, m)
				}
				fullArea.coords[m][n] += 1
				if fullArea.coords[m][n] == 2 {
					overlaps += 1
				}
			}
		}
	}



	fmt.Println("Total no. of overlaps:", overlaps)
}