package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

type movement struct {
	direction string
	magnitude int
}

func (m movement) String() string {
	return fmt.Sprintf("Direction: %s, Magnitude: %d", m.direction, m.magnitude)
}


func main() {
	fmt.Println("Advent of code day 2")
	fmt.Println("Part1:")
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	hPos := 0
	depth := 0

	var movements []movement
	for _, v := range strings.Split(string(dat), "\n") {
	  fields := strings.Fields(v)
		dir := fields[0]
		mag, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(err)
		}
		movements = append(movements, movement{dir, mag})
	}

	for _, m := range movements {
		switch m.direction {
		case "forward":
			hPos += m.magnitude
		case "up":
			depth -= m.magnitude
		case "down":
			depth += m.magnitude
		}
	}
	fmt.Println("Final depth:", depth)
	fmt.Println("Final horizontal position:", hPos)
	fmt.Println("hPos x depth =", hPos*depth)
	aim := 0
	hPos = 0
	depth = 0
	for _, m := range movements {
		switch m.direction {
		case "down":
			aim += m.magnitude
		case "up":
			aim -= m.magnitude
		case "forward":
			hPos += m.magnitude
			depth += aim * m.magnitude
		}
	}

	fmt.Println("\n\nPart2:")
	fmt.Println("Final depth:", depth)
	fmt.Println("Final horizontal position:", hPos)
	fmt.Println("hPos x depth =", hPos*depth)

}