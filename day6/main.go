package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/burlap101/advent-of-code-2021/day6/fish"
)

func populateFishes() []fish.LanternFish {
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	var fishes []fish.LanternFish
	for _, inputTime := range strings.Split(string(dat), ",") {
		initialTime, err := strconv.Atoi(inputTime)
		if err != nil {
			panic(err)
		}
		fishes = append(fishes, fish.LanternFish{Timer: initialTime})
	}
	return fishes
}

func serviceFishes(fishes []fish.LanternFish) ([]fish.LanternFish, error) {
	for i := range fishes {
		err := fishes[i].ServiceTimer(false)
		if err != nil {
			newBorn, err := fishes[i].SpawnNew()
			if err != nil {
				panic(err)
			}
			fishes = append(fishes, newBorn)
		}
	}
	return fishes, nil
}

func main() {
	
	fishes := populateFishes()
	
	fmt.Println("Initial count fishes:", len(fishes))

	for i := 0; i < 80; i++ {
		var err error
		fishes, err = serviceFishes(fishes)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Day%d count: %d\n", i+1, len(fishes))
	}


}